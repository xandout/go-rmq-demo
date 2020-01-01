FROM golang

WORKDIR /go/src/app
COPY adder.go ./main.go
COPY million.csv .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]