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

var dom = "dev.xn--b1arjbl.xn--90ae"

func TestHosts(t *testing.T) {
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

func Test_publishedStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, spf("http://%s/alabala?preview=bla", dom), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if publishedStatus(c) == 1 {
		t.Log("right guess for published")
	} else {
		t.Fatalf("publishedStatus(c) returned wrong status: %d", publishedStatus(c))
	}
}
