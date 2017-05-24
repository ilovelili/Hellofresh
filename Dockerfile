FROM golang:1.7-alpine

# Set apps working directory
WORKDIR /app

# Set an env var that matches $GOPATH
ENV SRC_DIR=/go/src/

# Copy the local package files to the container's workspace
ADD /src/. $SRC_DIR

# get go dependency
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get github.com/lib/pq
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega

# Test it, build it and copy to compiled directory
RUN cd $SRC_DIR/hellofresh; ginkgo; go build -o hellofresh; cp hellofresh /app/

ENTRYPOINT ["./hellofresh"]

EXPOSE 80 8080