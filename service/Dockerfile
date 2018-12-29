FROM alpine:latest

ADD bin/katlas /katlas
RUN chmod 755 /katlas
EXPOSE 8011

CMD /katlas \
    -envNamespace=$ENV_NAMESPACE \
    -serverType=$SERVER_TYPE \
    -dgraphHost=$DGRAPH_HOST
