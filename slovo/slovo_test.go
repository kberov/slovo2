package slovo

import (
	"regexp"
	"testing"
)

// TODO
func TestSLOG(t *testing.T) {
	r := regexp.MustCompile(spf("^/%s$", SLOG))
	m := r.FindAllStringSubmatch("/коренъ", -1)
	t.Logf("Match: %#v", m)

}
