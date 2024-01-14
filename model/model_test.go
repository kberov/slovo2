package model

import (
	"testing"
)

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
