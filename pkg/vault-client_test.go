package pkg

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestReturnsErrorForNonSuccessStatus(t *testing.T) {
	for _, tc := range []struct {
		name       string
		status     int
		body       string
		wantErr    bool
		wantSubstr string
	}{
		{name: "200 with a body succeeds", status: http.StatusOK, body: `{"id":"a"}`},
		{name: "200 with an empty body succeeds", status: http.StatusOK, body: ""},
		{
			name: "401 reports an authorization failure", status: http.StatusUnauthorized,
			wantErr: true, wantSubstr: "not authorized",
		},
		{
			name: "403 with an empty body is an error", status: http.StatusForbidden, body: "",
			wantErr: true, wantSubstr: "403",
		},
		{
			name: "500 carries the response body", status: http.StatusInternalServerError, body: "backend down",
			wantErr: true, wantSubstr: "backend down",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			client := &VaultClient{baseUri: srv.URL, token: "token"}
			var out map[string]string
			err := client.request("GET", "/secret", nil, &out)

			if !tc.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantSubstr) {
				t.Errorf("error %q does not contain %q", err, tc.wantSubstr)
			}
		})
	}
}
