package annotation

import "testing"

func TestDecodeGoMod(t *testing.T) {
	m, err := decodeGoMod("testdata")
	if err != nil {
		t.Error(err)
		return
	}
	if m.module != "example" && m.version != "1.19" {
		t.Fail()
	}
	t.Log(m.module, m.version)
}
