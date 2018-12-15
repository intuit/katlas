FROM golang:latest
RUN mkdir -p /go/src/github.com/intuit/katlas/controller
WORKDIR /go/src/github.com/intuit/katlas/controller
RUN go get -u -v k8s.io/client-go/...
RUN go get -u -v k8s.io/api/...
RUN go get -u -v github.com/dgraph-io/dgo
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get github.com/Sirupsen/logrus
COPY . .
RUN mkdir -p /root/.kube
RUN touch /root/.kube/config

ENV GOPATH=/go
CMD go run *.go
