FROM golang:1.9

RUN mkdir -p /go/src/github.com/richardalberto/concourse-dashboard
WORKDIR /go/src/github.com/richardalberto/concourse-dashboard
COPY . .

RUN go-wrapper install

EXPOSE 8080

CMD ["go-wrapper", "run"]