package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerSuite struct {
	MiddlewareSuite
}

func (suite *LoggerSuite) setupServer(fn func(w http.ResponseWriter, r *http.Request)) {
	m := http.NewServeMux()
	m.HandleFunc("/", fn)
	log := middleware.NewLogger(middleware.JSONLogger())
	suite.server = httptest.NewServer(log.Handler(m))
}

func (suite *LoggerSuite) TestLogger() {
	suite.setupServer(func(w http.ResponseWriter, r *http.Request) {
		log := middleware.LoggerFromContext(r.Context())
		assert.False(suite.T(), log.Core().Enabled(zapcore.DebugLevel))
		assert.True(suite.T(), log.Core().Enabled(zapcore.InfoLevel))

		w.WriteHeader(http.StatusNoContent)
	})

	req := suite.NewRequest()
	suite.Do(req)
}

func (suite *LoggerSuite) TestLoggerFromContext() {
	req := suite.NewRequest()
	log := middleware.LoggerFromContext(req.Context())
	assert.Equal(suite.T(), zapDiscard(), log)
}

func zapDiscard() *zap.Logger {
	return zap.New(nil)
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerSuite))
}
