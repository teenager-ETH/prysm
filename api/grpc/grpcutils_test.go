package grpc

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/v5/testing/assert"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/grpc/metadata"
)

type customErrorData struct {
	Message string `json:"message"`
}

func TestAppendHeaders(t *testing.T) {
	t.Run("one_header", func(t *testing.T) {
		ctx := AppendHeaders(context.Background(), []string{"first=value1"})
		md, ok := metadata.FromOutgoingContext(ctx)
		require.Equal(t, true, ok, "Failed to read context metadata")
		require.Equal(t, 1, md.Len(), "MetadataV0 contains wrong number of values")
		assert.Equal(t, "value1", md.Get("first")[0])
	})

	t.Run("multiple_headers", func(t *testing.T) {
		ctx := AppendHeaders(context.Background(), []string{"first=value1", "second=value2"})
		md, ok := metadata.FromOutgoingContext(ctx)
		require.Equal(t, true, ok, "Failed to read context metadata")
		require.Equal(t, 2, md.Len(), "MetadataV0 contains wrong number of values")
		assert.Equal(t, "value1", md.Get("first")[0])
		assert.Equal(t, "value2", md.Get("second")[0])
	})

	t.Run("one_empty_header", func(t *testing.T) {
		ctx := AppendHeaders(context.Background(), []string{"first=value1", ""})
		md, ok := metadata.FromOutgoingContext(ctx)
		require.Equal(t, true, ok, "Failed to read context metadata")
		require.Equal(t, 1, md.Len(), "MetadataV0 contains wrong number of values")
		assert.Equal(t, "value1", md.Get("first")[0])
	})

	t.Run("incorrect_header", func(t *testing.T) {
		logHook := logTest.NewGlobal()
		ctx := AppendHeaders(context.Background(), []string{"first=value1", "second"})
		md, ok := metadata.FromOutgoingContext(ctx)
		require.Equal(t, true, ok, "Failed to read context metadata")
		require.Equal(t, 1, md.Len(), "MetadataV0 contains wrong number of values")
		assert.Equal(t, "value1", md.Get("first")[0])
		assert.LogsContain(t, logHook, "Skipping second")
	})

	t.Run("header_value_with_equal_sign", func(t *testing.T) {
		ctx := AppendHeaders(context.Background(), []string{"first=value=1"})
		md, ok := metadata.FromOutgoingContext(ctx)
		require.Equal(t, true, ok, "Failed to read context metadata")
		require.Equal(t, 1, md.Len(), "MetadataV0 contains wrong number of values")
		assert.Equal(t, "value=1", md.Get("first")[0])
	})
}
