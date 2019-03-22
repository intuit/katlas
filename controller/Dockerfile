FROM alpine:latest

ARG appfolder
ENV appfolder /etc/katlas
RUN mkdir -p $appfolder
WORKDIR $appfolder

ADD bin/controller ./controller
RUN chmod 755 ./controller
RUN mkdir -p /root/.kube
RUN touch /root/.kube/config

CMD ./controller