FROM alpine:3.8

WORKDIR /opt

# libc6-compat required as the daemon is linked with GNU glibc not musl libc
# https://alexbilbie.com/2017/08/aws-xray-deamon-alpine-linux/
RUN apk add --update --no-cache curl unzip ca-certificates libc6-compat &&\
    update-ca-certificates     &&\
    curl -o /opt/daemon.zip https://s3.dualstack.us-east-2.amazonaws.com/aws-xray-assets.us-east-2/xray-daemon/aws-xray-daemon-linux-3.x.zip &&\
    unzip /opt/daemon.zip      &&\
    rm /opt/daemon.zip         &&\
    mv /opt/xray /usr/bin/xray &&\
    apk del curl unzip

CMD /usr/bin/xray --bind 0.0.0.0:2000
