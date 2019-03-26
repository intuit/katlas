FROM alpine:latest

ARG appfolder
ENV appfolder /etc/katlas
RUN mkdir -p $appfolder
WORKDIR $appfolder

ADD bin/controller ./controller
RUN chmod 755 ./controller

CMD ./controller
