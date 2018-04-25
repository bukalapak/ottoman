package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RequestIDSuite struct {
	MiddlewareSuite
}

func (suite *RequestIDSuite) setupServer(fn func(w http.ResponseWriter, r *http.Request)) {
	m := http.NewServeMux()
	m.HandleFunc("/", fn)
	suite.server = httptest.NewServer(middleware.RequestID(m))
}

func (suite *RequestIDSuite) TestRequestID() {
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		id := middleware.RequestIDFromContext(r.Context())
		assert.NotEqual(suite.T(), "", id)

		uid, err := uuid.FromString(id)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), uid.String(), id)

		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	suite.Do(req)
}

func (suite *RequestIDSuite) TestRequestID_fromRequest() {
	uid := "12d36928-06a2-442f-9600-d8c70df6c23c"
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		id := middleware.RequestIDFromContext(r.Context())
		assert.Equal(suite.T(), uid, id)

		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	req.Header.Set("X-Request-Id", uid)
	suite.Do(req)
}

func (suite *RequestIDSuite) TestRequestID_fromRequestSanitized() {
	uid := "12d36928-06a2-442f-9600-d8c70df6c23c"
	rid := "12d36928-06a2-442f-****&&&****9600-d8c70df6c23c"
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		id := middleware.RequestIDFromContext(r.Context())
		assert.Equal(suite.T(), uid, id)

		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	req.Header.Set("X-Request-Id", rid)
	suite.Do(req)
}

func TestRequestIDSuite(t *testing.T) {
	suite.Run(t, new(RequestIDSuite))
}
