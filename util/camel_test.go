package util_test

import (
	"testing"

	"github.com/xylonx/config2go/util"
)

func TestCamel(t *testing.T) {
	testCase := []struct {
		str        string
		upperCamel string
	}{
		{"", ""},
		{"snake", "Snake"},
		{"snake__case", "SnakeCase"},
		{"kabbel-case", "KabbelCase"},
		{"upperCamel", "UpperCamel"},
		{"UpperCamel", "UpperCamel"},
	}
	for i := range testCase {
		if testCase[i].upperCamel != util.ConvertString2UpperCamel(testCase[i].str) {
			t.FailNow()
		}
	}
}
