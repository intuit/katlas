FROM alpine:latest

ARG appfolder
ENV appfolder /etc/katlas
RUN mkdir -p $appfolder
WORKDIR $appfolder

ADD bin/katlas ./katlas
RUN chmod 755 ./katlas
EXPOSE 8011

CMD ./katlas \
    -envNamespace=$ENV_NAMESPACE \
    -serverType=$SERVER_TYPE \
    -dgraphHost=$DGRAPH_HOST