package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RecoverySuite struct {
	MiddlewareSuite
}

func (suite *RecoverySuite) setupServer(agent middleware.Notifier, fn func(w http.ResponseWriter, r *http.Request)) {
	m := http.NewServeMux()
	m.HandleFunc("/", fn)
	v := middleware.NewRecovery(agent)
	suite.server = httptest.NewServer(v.Handler(m))
}

func (suite *RecoverySuite) Do(r *http.Request) {
	hc := http.DefaultClient

	resp, err := hc.Do(r)
	assert.Nil(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *RecoverySuite) TestRecovery() {
	agent := &SampleAgent{t: suite.T(), Enabled: true}
	suite.setupServer(agent, func(w http.ResponseWriter, r *http.Request) {
		panic("!!!")
	})

	req := suite.NewRequest()
	suite.Do(req)
}

type SampleAgent struct {
	t       *testing.T
	Enabled bool
}

func (a *SampleAgent) Notify(err interface{}, stack []byte) {
	if a.Enabled {
		assert.Equal(a.t, "!!!", err)
		assert.NotEmpty(a.t, stack)
	}
}

func TestRecoverySuite(t *testing.T) {
	suite.Run(t, new(RecoverySuite))
}
