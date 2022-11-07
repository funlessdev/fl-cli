package fn

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

func TestFnDelete(t *testing.T) {
	testResult := "test-fn"
	testFn := "test-fn"
	testNs := "test-ns"
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use FnService.Delete to delete functions", func(t *testing.T) {
		cmd := Delete{
			Name:      testFn,
			Namespace: testNs,
		}

		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Delete", testCtx, testFn, testNs).Return(openapi.FunctionDeletionSuccess{Result: &testResult}, nil)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Delete", testCtx, testFn, testNs)
		mockInvoker.AssertNumberOfCalls(t, "Delete", 1)
		mockInvoker.AssertExpectations(t)
	})
	t.Run("should correctly print result", func(t *testing.T) {
		cmd := Delete{
			Name:      testFn,
			Namespace: testNs,
		}

		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Delete", testCtx, testFn, testNs).Return(openapi.FunctionDeletionSuccess{Result: &testResult}, nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockInvoker, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult+"\n", (&outbuf).String())
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should return error if invalid delete request", func(t *testing.T) {
		cmd := Delete{
			Name:      testFn,
			Namespace: testNs,
		}

		mockInvoker := mocks.NewFnHandler(t)

		e := &openapi.GenericOpenAPIError{}
		mockInvoker.On("Delete", testCtx, testFn, testNs).Return(openapi.FunctionDeletionSuccess{}, e)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.Error(t, err)
	})
}
