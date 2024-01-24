package slovo

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/labstack/echo/v4"
)

// TODO
func TestSLOG(t *testing.T) {
	r := regexp.MustCompile(spf(`^/%s$`, SLOG))
	path := "/коренъ"
	m := r.FindAllStringSubmatch(path, -1)
	if len(m) > 0 && m[0][0] != "" && m[0][1] == path[1:] {
		t.Logf("Match: %#v", m)
	} else {
		t.Fatalf(`SLOG  '%s' did not match : %#v`, SLOG, m)
	}
}

func TestHosts(t *testing.T) {
	dom := "dev.xn--b1arjbl.xn--90ae"
	for _, h := range []string{spf("http://%s:3000", dom), spf("http://%s", dom)} {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, h, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		host := hostName(c)
		if host == dom {
			t.Logf("Expected host name: %s", host)
		} else {
			t.Fatalf("UNExpected host name: %s", host)
		}
		ihost := iHostName(c)
		if ihost == "dev.слово.бг" {
			t.Logf("Expected unicode host name: %s", ihost)
		} else {
			t.Fatalf("UNExpected unicode host name: %s", ihost)
		}
	}
}
