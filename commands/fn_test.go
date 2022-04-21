package commands

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funlessdev/funless-cli/client"
	"github.com/funlessdev/funless-cli/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFn(t *testing.T) {
	testfn := "test-fn"

	t.Run("should use FnService.Invoke to invoke functions", func(t *testing.T) {
		cmd := fn{Name: testfn}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Invoke", testfn).Return(&http.Response{}, nil)

		err := cmd.Run(mockInvoker)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Invoke", testfn)
		mockInvoker.AssertNumberOfCalls(t, "Invoke", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should return error if invalid invoke request", func(t *testing.T) {
		cmd := fn{Name: testfn}
		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Invoke", testfn).Return(nil, fmt.Errorf("some error in FnService.Invoke"))

		err := cmd.Run(mockInvoker)
		require.Error(t, err)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/_/fn/"+testfn, r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	c, _ := client.NewClient(http.DefaultClient, client.Config{Host: server.URL})
	svc := &client.FnService{Client: c}
	t.Run("should send invoke request to server", func(t *testing.T) {
		cmd := fn{Name: testfn}
		err := cmd.Run(svc)
		require.NoError(t, err)
	})

}
