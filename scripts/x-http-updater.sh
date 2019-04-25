#!/bin/sh

DIR=x/http

curl "https://golang.org/src/net/http/httptest/recorder.go?m=text" | sed "s/httptest/http/g" | sed "s/golang_org/golang\.org/g" > $DIR/recorder.go
curl "https://golang.org/src/net/http/httptest/recorder_test.go?m=text" | sed "s/httptest/http/g" > $DIR/recorder_test.go
curl "https://raw.githubusercontent.com/urfave/negroni/master/response_writer.go" | sed "s/negroni/http/g" > $DIR/response_writer.go
curl "https://raw.githubusercontent.com/urfave/negroni/master/response_writer_test.go" | sed "s/negroni/http/g" > $DIR/response_writer_test.go

# adjust import path
sed -i '' 's|internal/x/net/http/httpguts|golang.org/x/net/http/httpguts|' $DIR/recorder.go
