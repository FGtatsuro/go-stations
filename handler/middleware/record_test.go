package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/mileusna/useragent"
)

func TestUserAgentRecord(t *testing.T) {
	testcases := map[string]struct {
		userAgent string
		wantOS    string
	}{
		useragent.MacOS: {
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
			wantOS:    useragent.MacOS,
		},
		useragent.Windows: {
			userAgent: "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
			wantOS:    useragent.Windows,
		},
		useragent.IOS: {
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
			wantOS:    useragent.IOS,
		},
		useragent.Android: {
			userAgent: "Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
			wantOS:    useragent.Android,
		},
	}

	for _, tc := range testcases {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(
			"User-Agent",
			tc.userAgent,
		)
		w := httptest.NewRecorder()
		h := middleware.UserAgentRecord(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Context().Value(middleware.UAContextKeyOS)
			want := tc.wantOS
			if got != want {
				t.Errorf("Contextに正しいOS情報がセットされていません, got = %v, want = %v", got, want)
			}
		}))
		h.ServeHTTP(w, r)
	}
}
