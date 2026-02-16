package hcloud_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeJSON(t *testing.T, writer http.ResponseWriter, data string) {
	t.Helper()

	_, writeErr := writer.Write([]byte(data))
	require.NoError(t, writeErr)
}
