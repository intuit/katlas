FROM alpine:latest

ARG appfolder
ENV appfolder /etc/katlas
RUN mkdir -p $appfolder/data
WORKDIR $appfolder

ADD bin/katlas ./katlas
ADD data/*.json ./data/
RUN chmod 755 ./katlas
EXPOSE 8011

CMD ./katlas \
    -serverType=$SERVER_TYPE \
    -dgraphHost=$DGRAPH_HOST
