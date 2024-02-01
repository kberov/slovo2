package model

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/gommon/log"
)

func init() {
	Logger = log.New("DB")
	Logger.SetOutput(os.Stderr)
	Logger.SetHeader(defaultLogHeader)
	Logger.SetLevel(log.DEBUG)
	DSN = "../data/slovo.dev.sqlite"
}

func TestStranici_FindForDisplay(t *testing.T) {
	user := new(Users)
	GetByID(user, 2) // guest
	str := new(Stranici)
	if err := str.FindForDisplay("вѣра", user, 2, "dev.xn--b1arjbl.xn--90ae", "bg"); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	t.Logf("Stranica: %#v", str)
}

func TestCelini_FindForDisplay(t *testing.T) {
	user := new(Users)
	GetByID(user, 2) // guest
	page := new(Stranici)
	GetByID(page, 11)
	cel := new(Celini)
	if err := cel.FindForDisplay(page, "благодарност", user, 2, "bg", "main"); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if cel.ID != 54 {
		t.Fatalf("record not found. expected ID was 54, but it is %d", cel.ID)
	}

	//t.Logf("Stranica: %#v", cel)
}

func TestSQLFor(t *testing.T) {
	table := Record2Table(&Stranici{})
	SQL := SQLFor("GET_PAGE_FOR_DISPLAY", table)

	if strings.Contains(SQL, "${") {
		t.Fatalf("SQL contains placeholders:\n%s", SQL)
	}
	t.Log(SQL)
}

func TestSelectMenuItems(t *testing.T) {
	user := new(Users)
	GetByID(user, 2) // guest
	domain, lang, pub := "dev.xn--b1arjbl.xn--90ae", "bg", 2
	if items, err := SelectMenuItems(QArgs{
		"now":     time.Now().Unix(),
		"user_id": user.ID,
		"domain":  domain,
		"pub":     pub,
		"lang":    lang,
	}); err != nil {
		t.Fatalf("SelectMenuItems failed: %v", err)
	} else {
		t.Logf("Rows: %#v", items)
	}
}
