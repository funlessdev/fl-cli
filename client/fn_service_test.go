package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoke(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
	fnService := &FnService{c}

	t.Run("should send a GET request", func(t *testing.T) {
		res, _ := fnService.Invoke("test")
		require.Equal(t, http.MethodGet, res.Request.Method)
	})

	t.Run("should have /_/fn/test as url path when invoking fn named test", func(t *testing.T) {
		res, _ := fnService.Invoke("test")
		require.Equal(t, "/_/fn/test", res.Request.URL.Path)
	})

	// TODO what does invoke should return?
}
