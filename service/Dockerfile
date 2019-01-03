FROM alpine:latest

ARG appfolder
ENV appfolder /etc/katlas
RUN mkdir -p $appfolder
ADD bin/katlas $appfolder/katlas
RUN chmod 755 $appfolder/katlas
EXPOSE 8011

CMD $appfolder/katlas \
    -envNamespace=$ENV_NAMESPACE \
    -serverType=$SERVER_TYPE \
    -dgraphHost=$DGRAPH_HOST