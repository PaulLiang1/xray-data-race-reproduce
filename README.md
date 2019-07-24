here are bunch of test file to help reproduce this github [issue](https://github.com/aws/aws-xray-sdk-go/issues/118).

to run example, it requires
- docker
- docker-compose
- make

## data_race
Contains setup similar in original ticket, using env
```
go version: 1.11.x
goos: linux
goarch: amd64
aws-xray-sdk-go@v1.0.0-rc.11
aws-sdk-go@v1.17.12
```

to run, do
```
make run-data-race
```

this should re-produce data-race error, a copy of the log had been attached [here](./run-data-race.log)
Note the `xray-daemon` generate log line about `The security token included in the request is invalid.` which is expected as this test does not require any real aws cred and we believed this error is irrelevant to data race as we could observe this with our online system where valid cred is present.  

## safe
same env as `data_race`, with xray excluded, to _not_ produce data race to rule out `aws-go-sdk` issue
```
go version: 1.11.x
goos: linux
goarch: amd64
aws-sdk-go@v1.17.12
```

to run, do
```
make run-safe
```

log can be seen [here](./run-safe.log), no data race warning can be observed.
