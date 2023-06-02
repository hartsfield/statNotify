# syntax=docker/dockerfile:1
FROM fedora:latest

RUN adduser hrtsfld

WORKDIR /home/hrtsfld
ADD statNotify /home/hrtsfld
RUN chown -R hrtsfld /home/hrtsfld
ENV PATH="${PATH}:/home/hrtsfld"
ENV statLogPath="/home/hrtsfld/statlog"
ENV statAdminEmail="johnathanhartsfield@gmail.com"
CMD statNotify
