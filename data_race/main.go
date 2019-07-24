package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-xray-sdk-go/strategy/ctxmissing"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func setUpDDBTestTbl(ctx context.Context, client dynamodbiface.DynamoDBAPI, tableName string) {
	_, err := client.CreateTableWithContext(ctx,
		&dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{AttributeName: aws.String("foo"), AttributeType: aws.String("S")},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{AttributeName: aws.String("foo"), KeyType: aws.String("HASH")},
			},
			TableName: &tableName,
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(100),
				WriteCapacityUnits: aws.Int64(100),
			},
		},
	)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create ddb table")
	}
	logrus.Infof("ddb table [%s] created", tableName)
}

func tearDDBTestTbl(ctx context.Context, client dynamodbiface.DynamoDBAPI, tableName string) {
	_, err := client.DeleteTableWithContext(ctx, &dynamodb.DeleteTableInput{
		TableName: &tableName,
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to delete ddb table")
	}
	logrus.Infof("ddb table [%s] deleted", tableName)
}

type controller struct {
	ddbClient dynamodbiface.DynamoDBAPI
}

func (c *controller) handlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ddbTblName := uuid.New().String()
	setUpDDBTestTbl(ctx, c.ddbClient, ddbTblName)
	defer tearDDBTestTbl(ctx, c.ddbClient, ddbTblName)

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int){
			defer wg.Done()

			_, err := c.ddbClient.PutItemWithContext(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(ddbTblName),
				Item: map[string]*dynamodb.AttributeValue{
					"foo": {
						S: aws.String(uuid.New().String()),
					},
					"bar": {
						S: aws.String(uuid.New().String()),
					},
				},
			})
			if err != nil {
				logrus.WithError(err).Error("error writing ddb item")
			} else {
				logrus.Infof("item [%d] written", i)
			}
		}(i)
	}

	wg.Wait()
	w.WriteHeader(http.StatusOK)
}

func main() {
	ddbEndpoint := "http://dynamodb-local:8000"
	awsSession, _ := session.NewSession(&aws.Config{
		Endpoint: aws.String(ddbEndpoint),
		Region:   aws.String(uuid.New().String()),
	})

	ddbClient := dynamodb.New(awsSession)
	_ = xray.Configure(xray.Config{LogLevel: "trace"})
	xray.AWS(ddbClient.Client)
	_ = xray.Configure(xray.Config{
		ContextMissingStrategy: ctxmissing.NewDefaultLogErrorStrategy(),
		LogLevel:               "debug",
	})

	c := controller{ddbClient: ddbClient}
	httpHandler := xray.Handler(
		xray.NewFixedSegmentNamer("foo"),
		http.HandlerFunc(c.handlerFunc),
	)

	srv := http.Server{
		Handler: httpHandler,
		Addr:    ":8080",
	}
	_ = srv.ListenAndServe()
}
