package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	uid := "12d36928-06a2-442f-9600-d8c70df6c23c"
	bid := "12d36928-06a2-442f-****&&&****9600-d8c70df6c23c"

	fn := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, middleware.RequestIDFromContext(r.Context()))
	}

	x := func(req *http.Request, id string) {
		rec := httptest.NewRecorder()

		middleware.RequestID(http.HandlerFunc(fn)).ServeHTTP(rec, req)

		uid, err := uuid.FromString(rec.Body.String())
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		if id != "" {
			assert.Equal(t, uid.String(), id)
		} else {
			assert.Equal(t, uuid.V4, uid.Version())
		}

	}

	t.Run("RequestID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)

		x(req, "")
	})

	t.Run("From-Request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-Request-Id", uid)

		x(req, uid)
	})

	t.Run("Sanitized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-Request-Id", bid)

		x(req, uid)
	})

}
