FROM golang:1.10
RUN mkdir -p /go/src/github.com/intuit/katlas/rest-service
WORKDIR /go/src/github.com/intuit/katlas/rest-service
RUN go get -u -v github.com/dgraph-io/dgo
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/gorilla/mux
RUN go get github.com/mitchellh/mapstructure
RUN go get github.com/stretchr/testify/mock
RUN apt-get update && apt-get install -y vim openssl
RUN go get -u -v k8s.io/client-go/...
RUN go get -u -v k8s.io/api/...
ENV GOPATH=/go
EXPOSE 8011
COPY . .
CMD go run server.go \
    -envNamespace=$ENV_NAMESPACE \
    -serverType=$SERVER_TYPE \
    -dgraphHost=$DGRAPH_HOST
