package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_FetchReport(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/reports", r.URL.Path)
			require.Equal(t, "2020-01-01", r.URL.Query().Get("from"))
			require.Equal(t, "2020-01-31", r.URL.Query().Get("to"))
			w.Header().Set("Content-Type", "application/json")
			require.NoError(t, json.NewEncoder(w).Encode(map[string]float64{"income": 100.5, "expense": 50.25}))
		}))
		defer server.Close()

		client := New(server.URL, 5*time.Second)
		result, err := client.FetchReport(context.Background(), "2020-01-01", "2020-01-31")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 100.5, result.Income)
		require.Equal(t, 50.25, result.Expense)
	})

	t.Run("bad status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := New(server.URL, 5*time.Second)
		result, err := client.FetchReport(context.Background(), "2020-01-01", "2020-01-31")
		require.Error(t, err)
		require.Nil(t, result)
		require.Contains(t, err.Error(), "500")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("not json"))
			require.NoError(t, err)
		}))
		defer server.Close()

		client := New(server.URL, 5*time.Second)
		result, err := client.FetchReport(context.Background(), "2020-01-01", "2020-01-31")
		require.Error(t, err)
		require.Nil(t, result)
	})
}
