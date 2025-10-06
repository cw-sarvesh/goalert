package gupshup

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type requestInfo struct {
	method string
	header http.Header
	form   url.Values
}

func TestClientSendSMS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		status      int
		response    string
		apiKey      string
		wantID      string
		wantErr     bool
		validateReq func(t *testing.T, info requestInfo)
	}{
		{
			name:     "TopLevelMessageID",
			status:   http.StatusOK,
			response: `{"messageId":"msg-123"}`,
			apiKey:   "secret",
			wantID:   "msg-123",
			validateReq: func(t *testing.T, info requestInfo) {
				require.Equal(t, http.MethodPost, info.method)
				require.Equal(t, "application/x-www-form-urlencoded", info.header.Get("Content-Type"))
				require.Equal(t, "secret", info.header.Get("apikey"))
				require.Equal(t, "SMS", info.form.Get("channel"))
				require.Equal(t, "source", info.form.Get("source"))
				require.Equal(t, "+15555550123", info.form.Get("destination"))
				require.Equal(t, "hello world", info.form.Get("message"))
			},
		},
		{
			name:     "NestedResponse",
			status:   http.StatusOK,
			response: `{"response":{"msgId":"msg-999"}}`,
			wantID:   "msg-999",
			validateReq: func(t *testing.T, info requestInfo) {
				require.Equal(t, "", info.header.Get("apikey"))
			},
		},
		{
			name:     "InvalidJSON",
			status:   http.StatusOK,
			response: `not-json`,
		},
		{
			name:     "HTTPError",
			status:   http.StatusBadGateway,
			response: `{"code":"123","message":"error"}`,
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var captured requestInfo

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				vals, err := url.ParseQuery(string(body))
				require.NoError(t, err)

				captured = requestInfo{
					method: r.Method,
					header: r.Header.Clone(),
					form:   vals,
				}

				if tc.status != 0 {
					w.WriteHeader(tc.status)
				}
				_, _ = w.Write([]byte(tc.response))
			}))
			defer srv.Close()

			client := NewClient(Config{
				BaseURL:    srv.URL,
				APIKey:     tc.apiKey,
				Source:     "source",
				HTTPClient: srv.Client(),
			})

			id, err := client.SendSMS(context.Background(), "+15555550123", "hello world")

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantID, id)

			if tc.validateReq != nil {
				tc.validateReq(t, captured)
			}
		})
	}
}
