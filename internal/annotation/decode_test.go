package annotation

import "testing"

func TestDecodeGoFile(t *testing.T) {
	d := &decoder{
		root:     "testdata",
		replaces: map[string][]string{},
	}
	files, err := d.decode()
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range files {
		t.Log(v.path, v.imp, len(v.nodes))
		for _, n := range v.nodes {
			t.Log(n.format())
		}
	}
}
