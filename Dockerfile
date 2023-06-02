# syntax=docker/dockerfile:1
FROM fedora:latest

RUN adduser hrtsfld
WORKDIR /home/hrtsfld

RUN curl https://dl.google.com/go/go1.20.4.linux-amd64.tar.gz --output go1.20.4.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.4.tar.gz
RUN rm go1.20.4.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin:/home/hrtsfld"

ADD . /home/hrtsfld
RUN cd /home/hrtsfld
RUN go mod tidy
RUN go build -o statNotify
RUN chown -R hrtsfld /home/hrtsfld
ENV statLogPath="/home/hrtsfld/statlog"
ENV statAdminEmail="johnathanhartsfield@gmail.com"
CMD statNotify
