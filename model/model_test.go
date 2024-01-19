package model

import (
	"os"
	"strings"
	"testing"

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
	if err := str.FindForDisplay("вѣра", user, "dev.xn--b1arjbl.xn--90ae"); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	t.Logf("Stranica: %#v", str)
}

func TestSQLFor(t *testing.T) {
	table := Record2Table(&Stranici{})
	SQL := SQLFor("GET_PAGE_FOR_DISPLAY", table)

	if strings.Contains(SQL, "${") {
		t.Fatalf("SQL contains placeholders:\n%s", SQL)
	}
	t.Log(SQL)
}

func TestRecord2Table(t *testing.T) {

	toTable1 := Record2Table(&FirstLogin{})
	if table := "first_login"; table != toTable1 {
		t.Fatalf("%s != %s for %T", table, toTable1, &FirstLogin{})
	} else {
		t.Logf("%s == %s for %T", table, toTable1, &FirstLogin{})
	}
	toTable2 := Record2Table(&Users{})
	if table := "users"; table != toTable2 {
		t.Fatalf("%s != %s for %T", table, toTable2, &Users{})
	} else {
		t.Logf("%s == %s for %T", table, toTable2, &Users{})
	}
	toTable := Record2Table(&UsersInvoicesLastID{})
	if table := "users_invoices_last_id"; table != toTable {
		t.Fatalf("%s != %s for %T", table, toTable, &UsersInvoicesLastID{})
	} else {
		t.Logf("%s == %s for %T", table, toTable, &UsersInvoicesLastID{})
	}
}

func TestField2Column(t *testing.T) {
	fields := map[string]string{
		"ID":                  "id",
		"PageType":            "page_type",
		"GroupID":             "group_id",
		"UsersInvoicesLastID": "users_invoices_last_id",
	}
	for k, v := range fields {
		t.Run(k, func(t *testing.T) {
			if Field2Column(k) != v {
				t.Fatalf("%s => %s", k, v)
			}
			t.Logf("%s => %s", k, v)
		})
	}
}
