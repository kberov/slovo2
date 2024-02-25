package slovo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

var cfgFile = Cfg.ConfigFile
var Logger = log.New("slovo2")

const defaultLogHeader = `${prefix}:${level}:${short_file}:${line}`

func init() {
	Logger.SetOutput(os.Stderr)
	Logger.SetHeader(defaultLogHeader)
}

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

func TestRoutes(t *testing.T) {
	/*
		Cfg.Renderer.TemplateRoots = []string{`../templates`}
		e := initEcho(Logger)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/коренъ.bg.html", nil)
		c := e.NewContext(req, rec)
		e.ServeHTTP(rec, req)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, echo.MIMETextHTMLCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
			assert.Equal(t, "Hello, <strong>World!</strong>", rec.Body.String())
		}
	*/
	var testCases = []struct {
		name         string
		whenURL      string
		expectStatus int
		bodyContains string
	}{
		{
			name:         "/",
			whenURL:      "/",
			expectStatus: http.StatusOK,
			bodyContains: `Знакът Ⱄ в горния ляв ъгъл е буква.`,
		},
		{
			name:         "index.html",
			whenURL:      "/index.html",
			expectStatus: http.StatusOK,
			bodyContains: `Знакът Ⱄ в горния ляв ъгъл е буква.`,
		},
		{
			name:         "коренъ.html",
			whenURL:      "/коренъ.html",
			expectStatus: http.StatusOK,
			bodyContains: `Знакът Ⱄ в горния ляв ъгъл е буква.`,
		},
		{
			name:         "коренъ.bg.html",
			whenURL:      "/коренъ.bg.html",
			expectStatus: http.StatusOK,
			bodyContains: `Знакът Ⱄ в горния ляв ъгъл е буква.`,
		},
		{
			name:         "кънигꙑ.bg.html",
			whenURL:      "/кънигꙑ.bg.html",
			expectStatus: http.StatusOK,
			bodyContains: `Тук предлагаме книги, които ни се`,
		},
		{
			name:         "кънигꙑ.bg.html 2",
			whenURL:      "/кънигꙑ.bg.html",
			expectStatus: http.StatusOK,
			bodyContains: `Пропаднал бар в предградията на един тропически град `,
		},
		{
			name:         "кънигꙑ.bg.html ⮊",
			whenURL:      "/кънигꙑ.bg.html",
			expectStatus: http.StatusOK,
			bodyContains: `⮊`,
		},
		{
			name:         "кънигꙑ.bg.html?limit=10&offset=10",
			whenURL:      "/кънигꙑ.bg.html?limit=10&offset=10",
			expectStatus: http.StatusOK,
			bodyContains: `Черно на черно`,
		},
		{
			name:         "кънигꙑ/матере-нашѧ-параскеви.bg.html",
			whenURL:      "/кънигꙑ/матере-нашѧ-параскеви.bg.html",
			expectStatus: http.StatusOK,
			bodyContains: `Електронно издание за свободно изтегляне`,
		},
		{
			name:         "новолѣпьно/малък-свят-на-български.html",
			whenURL:      "/новолѣпьно/малък-свят-на-български.html",
			expectStatus: http.StatusOK,
			// canonical url
			bodyContains: `https://xn--b1arjbl.xn--90ae/новолѣпьно/малък-свят-на-български.bg.html`,
		},
		{
			name:         "новолѣпьно/малък-свят-на-български.bg.html",
			whenURL:      "/новолѣпьно/малък-свят-на-български.bg.html",
			expectStatus: http.StatusOK,
			bodyContains: `<h1>Гуарески за първи път на български</h1>`,
		},
		{
			name:         "notfound.bg.html",
			whenURL:      "/новолѣпьно/unknown.html",
			expectStatus: http.StatusNotFound,
			bodyContains: `<h1>Няма такава страница!</h1>`,
		},
	}

	Cfg.Renderer.TemplateRoots = []string{`../` + Cfg.Renderer.TemplateRoots[0]}
	Cfg.DB.DSN = "../" + Cfg.DB.DSN
	e := initEcho(Logger)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := `http://` + Cfg.Serve.Location + tc.whenURL
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectStatus, rec.Code)
			assert.True(t, strings.Contains(rec.Body.String(), tc.bodyContains))
		})
	}

}
