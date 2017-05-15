package middleware_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RealIPSuite struct {
	MiddlewareSuite
	sampleIP  net.IP
	anotherIP net.IP
	localIP   net.IP
}

func (suite *RealIPSuite) SetupTest() {
	suite.sampleIP = net.ParseIP("202.212.212.202")
	suite.anotherIP = net.ParseIP("202.222.222.202")
	suite.localIP = net.ParseIP("127.0.0.1")
}

func (suite *RealIPSuite) setupServer(fn func(w http.ResponseWriter, r *http.Request)) {
	m := http.NewServeMux()
	m.HandleFunc("/", fn)
	suite.server = httptest.NewServer(middleware.RealIP(m))
}

func (suite *RealIPSuite) TestRealIP() {
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		rip, ok := middleware.IPFromContext(r.Context())
		assert.True(suite.T(), ok)
		assert.Equal(suite.T(), suite.sampleIP.String(), rip)
		assert.Equal(suite.T(), suite.sampleIP.String(), r.Header.Get("X-Forwarded-For"))
		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	req.Header.Add("X-Forwarded-For", suite.sampleIP.String())
	req.Header.Add("X-Real-IP", suite.anotherIP.String())

	suite.Do(req)
}

func (suite *RealIPSuite) TestRealIP_xForwardedFor() {
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		rip, ok := middleware.IPFromContext(r.Context())
		assert.True(suite.T(), ok)
		assert.Equal(suite.T(), suite.sampleIP.String(), rip)
		assert.Equal(suite.T(), suite.sampleIP.String(), r.Header.Get("X-Forwarded-For"))
		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	req.Header.Add("X-Forwarded-For", suite.sampleIP.String())

	suite.Do(req)
}

func (suite *RealIPSuite) TestRealIP_xRealIP() {
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		rip, ok := middleware.IPFromContext(r.Context())
		assert.True(suite.T(), ok)
		assert.Equal(suite.T(), suite.sampleIP.String(), rip)
		assert.Equal(suite.T(), suite.sampleIP.String(), r.Header.Get("X-Real-IP"))
		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	req.Header.Add("X-Real-IP", suite.sampleIP.String())

	suite.Do(req)
}

func (suite *RealIPSuite) TestRealIP_remoteAddr() {
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		rip, ok := middleware.IPFromContext(r.Context())
		assert.True(suite.T(), ok)
		assert.Equal(suite.T(), suite.localIP.String(), rip)

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), suite.localIP.String(), host)
		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	suite.Do(req)
}

func TestRealIPSuite(t *testing.T) {
	suite.Run(t, new(RealIPSuite))
}
