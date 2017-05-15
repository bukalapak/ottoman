package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestContextKey(t *testing.T) {
	key := middleware.ContextKey("Hello")
	assert.Equal(t, "ottoman/middleware context: Hello", key.String())
}

type MiddlewareSuite struct {
	suite.Suite
	server *httptest.Server
}

func (suite *MiddlewareSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *MiddlewareSuite) URL(s string) string {
	return suite.server.URL + s
}

func (suite *MiddlewareSuite) NewRequest() *http.Request {
	req, _ := http.NewRequest("GET", suite.URL("/"), nil)
	return req
}

func (suite *MiddlewareSuite) Do(r *http.Request) {
	hc := http.DefaultClient

	resp, err := hc.Do(r)
	assert.Nil(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)
}
