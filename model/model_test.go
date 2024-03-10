package model

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/gommon/log"
)

var args *StraniciArgs

func init() {
	Logger = log.New("DB")
	Logger.SetOutput(os.Stderr)
	Logger.SetHeader(defaultLogHeader)
	Logger.SetLevel(log.DEBUG)
	DSN = "../data/slovo.dev.sqlite"
	args = &StraniciArgs{
		Alias:  "вѣра",
		UserID: 2,
		Domain: "xn--b1arjbl.xn--90ae",
		Box:    MainBox,
		Pub:    2,
		Lang:   "bg",
		Now:    time.Now().Unix(),
	}
}

func TestStranici_FindForDisplay(t *testing.T) {
	str := new(Stranici)
	if err := str.FindForDisplay(*args); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	// t.Logf("Stranica: %#v", str)
}

func TestCelini_FindForDisplay(t *testing.T) {
	args.Alias = "ѩꙁыкъ"
	args.Celina = "благодарност"
	t.Logf("StraniciArgs: %#v", args)
	cel := new(Celini)
	if err := cel.FindForDisplay(*args); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if cel.ID != 54 {
		t.Fatalf("record not found. expected ID was 54, but it is %d", cel.ID)
	}
	// t.Logf("Celina: %#v", cel)
}

func TestSQLFor(t *testing.T) {
	table := Record2Table(&Stranici{})
	SQL := SQLFor("GET_PAGE_FOR_DISPLAY", table)

	if strings.Contains(SQL, "${") {
		t.Fatalf("SQL contains placeholders:\n%s", SQL)
	}
	//t.Log(SQL)
}

func TestSelectMenuItems(t *testing.T) {
	_ = SelectMenuItems(*args)
	errargs := *args
	errargs.Pub = 1
	if items := SelectMenuItems(errargs); len(items) == 0 {
		t.Logf("expected no menuitems")
	} else {
		t.Fatalf("something went terribly wrong (Unexpected items): %#v", items)
	}
}

func TestListStranici(t *testing.T) {
	myArgs := *args
	myArgs.Alias = "коренъ"
	stranici := ListStranici(myArgs)
	if stranici[0].ID != 21 {
		t.Fatalf("ListStranici failed: %v", "Unexpected page at index 0")
	}
}

func TestGetDomain(t *testing.T) {
	domains := map[string]string{
		`xn--b1arjbl.xn--90ae`:     `xn--b1arjbl.xn--90ae`,
		`dev.xn--b1arjbl.xn--90ae`: `xn--b1arjbl.xn--90ae`,
		`www.xn--b1arjbl.xn--90ae`: `xn--b1arjbl.xn--90ae`,
		`localhost`:                `localhost`,
		`qa.localhost`:             `localhost`,
		`dev.i-can.eu`:             `i-can.eu`,
		`127.0.0.1`:                `localhost`,
	}

	for k, v := range domains {
		t.Run(k, func(t *testing.T) {
			dom := new(Domove)
			dom.GetByName(k)
			if dom.Domain != v {
				t.Fatalf("Unexpectedly got domain %s. Expected: %s", dom.Domain, v)
			}
			// t.Logf("%s => %#v", k, dom)
		})
	}
}

func expectPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MISSING PANIC")
		} else {
			t.Log(r)
		}
	}()
	f()
}
