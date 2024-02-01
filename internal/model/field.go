package model

import "github.com/go-chocolate/chocdoc/internal/utils/kv"

type Field struct {
	Name     string
	Type     string
	Array    int
	Required bool
	Comment  string
	Option   string
	Sub      *Model
	Tags     kv.KV
	sub      string
}

func (f *Field) SetName(name string) {
	if f.Name == "" {
		f.Name = name
	}
}
