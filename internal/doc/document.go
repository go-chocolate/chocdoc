package doc

import (
	"strings"
)

type Document struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Req         *Model `json:"req"`
	Rsp         *Model `json:"rsp"`
	Group       string `json:"group"`
	Header      KV     `json:"header"`
	Extra       KV     `json:"extra"`
}

type Documents []*Document

type DocumentGroup struct {
	Name      string                    `json:"name"`
	Documents Documents                 `json:"documents"`
	Children  map[string]*DocumentGroup `json:"children"`
	Root      bool                      `json:"root"`
}

func (a Documents) Group(seps ...string) *DocumentGroup {
	var sep = "/"
	if len(seps) > 0 {
		sep = seps[0]
	}
	var groups = &DocumentGroup{Children: make(map[string]*DocumentGroup), Documents: make([]*Document, 0), Root: true}
	for _, doc := range a {
		var ptr = groups
		if doc.Group == "" {
			ptr.Documents = append(ptr.Documents, doc)
			continue
		}
		names := strings.Split(doc.Group, sep)
		for n, name := range names {
			if ptr.Children[name] == nil {
				ptr.Children[name] = &DocumentGroup{Name: name, Children: make(map[string]*DocumentGroup), Documents: make([]*Document, 0)}
			}
			ptr = ptr.Children[name]
			if n == len(names)-1 {
				ptr.Documents = append(ptr.Documents, doc)
			}
		}
	}
	return groups
}
