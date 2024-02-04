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
		Domain: "dev.xn--b1arjbl.xn--90ae",
		Box:    MainBox,
		Pub:    2,
		Lang:   "bg",
		Now:    time.Now().Unix(),
	}
}

func TestStranici_FindForDisplay(t *testing.T) {
	str := new(Stranici)
	if err := str.FindForDisplay(args); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	// t.Logf("Stranica: %#v", str)
}

func TestCelini_FindForDisplay(t *testing.T) {
	args.Alias = "ѩꙁыкъ"
	args.Celina = "благодарност"
	t.Logf("StraniciArgs: %#v", args)
	cel := new(Celini)
	if err := cel.FindForDisplay(args); err != nil {
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
	if _, err := SelectMenuItems(args); err != nil {
		t.Fatalf("SelectMenuItems failed: %v", err)
	} // else {//t.Logf("Rows: %#v", items)}
}
