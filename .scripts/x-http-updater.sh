#!/bin/sh

DIR=x/http

curl "https://golang.org/src/net/http/httptest/recorder.go?m=text" | sed "s/httptest/http/g" > $DIR/recorder.go
curl "https://golang.org/src/net/http/httptest/recorder_test.go?m=text" | sed "s/httptest/http/g" > $DIR/recorder_test.go
curl "https://raw.githubusercontent.com/urfave/negroni/master/response_writer.go" | sed "s/negroni/http/g" > $DIR/response_writer.go
curl "https://raw.githubusercontent.com/urfave/negroni/master/response_writer_test.go" | sed "s/negroni/http/g" > $DIR/response_writer_test.go
curl "https://raw.githubusercontent.com/golang/gddo/master/httputil/header/header.go" | sed "/Package header/d" | sed "s/package header/package http/ig" > $DIR/header.go
curl "https://raw.githubusercontent.com/golang/gddo/master/httputil/header/header_test.go" | sed "/Package header/d" | sed "s/package header/package http/ig" > $DIR/header_test.go
