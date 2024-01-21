package util_test

import (
	"testing"

	"github.com/kberov/slovo2/util"
)

func TestSlogifyStripPunctEmptyConnector(t *testing.T) {
	slovoplet := map[string]string{
		"ALA bala Nica, Turska panica": "alabalanicaturskapanica",
		"Кънигы":                       "кънигы",
		"OwnerID":                      "ownerid",
		"PageType":                     "pagetype",
		"UsersInvoicesLastID":          "usersinvoiceslastid",
		"ПРОСТРАННО ЖИТИЕ НА РОМИЛ ВИДИНСКИ": "пространножитиенаромилвидински",
		"Пространно Житие НА Ромил Видински": "пространножитиенаромилвидински",
	}

	for k, v := range slovoplet {
		t.Run(k, func(t *testing.T) {
			slog := util.Slogify(k, "", true)
			t.Logf("%s => %s|%s", k, slog, v)
		})
	}
}

func TestSlogifyStripPunctWithConnector(t *testing.T) {
	slovoplet := map[string]string{
		"ALA bala Nica, Turska panica?": "ala-bala-nica-turska-panica",
		"Кънигы":                        "кънигы",
		"OwnerID":                       "ownerid",
		"PageType":                      "pagetype",
		"UsersInvoicesLastID":           "usersinvoiceslastid",
		"ПРОСТРАННО ЖИТИЕ НА РОМИЛ ВИДИНСКИ": "пространно-житие-на-ромил-видински",
		"Пространно Житие НА Ромил Видински": "пространно-житие-на-ромил-видински",
		"Users Invoices LastID": "users-invoices-lastid",
		"Ѩꙁꙑкъ!?$":              "ѩꙁꙑкъ$",
	}

	for k, v := range slovoplet {
		t.Run(k, func(t *testing.T) {
			slova := util.Slogify(k, "-", true)
			t.Logf("%s =>\n%s\n%s", k, slova, v)
			if slova != v {
				t.Fail()
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	slovo := map[string]string{
		"Кънигы":                "кънигы",
		"НашитеКънигы":          "нашите_кънигы",
		"OwnerID":               "owner_id",
		"PageType":              "page_type",
		"UsersInvoicesLastID":   "users_invoices_last_id",
		"Users Invoices LastID": "users _invoices _last_id",
	}

	for k, v := range slovo {
		t.Run(k, func(t *testing.T) {
			slova := util.ToSnakeCase(k)
			t.Logf("%s => %s|%s", k, slova, v)
			if slova != v {
				t.Fail()
			}
		})
	}
}
