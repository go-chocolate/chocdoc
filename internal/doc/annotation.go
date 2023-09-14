package doc

import (
	"strings"

	"github.com/go-chocolate/chocdoc/elements"
)

func FromAnnotation(routers []*GINRouter, nodes map[string]*elements.Node) []*document {
	var docs []*document
	for _, router := range routers {
		ele := nodes[router.FuncName]
		var doc = &document{
			path:    router.Path,
			method:  router.Method,
			handler: router.FuncName,
		}
		if ele == nil {
			docs = append(docs, doc)
			continue
		}
		for _, ann := range ele.Annotations {
			n := strings.Index(ann.Content, " ")
			if n <= 0 {
				continue
			}
			var key = ann.Content[:n]
			var content = strings.TrimSpace(ann.Content[n+1:])
			switch key {
			case "req", "request":
				if len(ann.Relation) > 0 {
					doc.request = ann.Relation[0]
				}
			case "rsp", "response":
				if len(ann.Relation) > 0 {
					doc.response = ann.Relation[0]
				}
			case "name":
				doc.name = content
			case "summary":
				doc.summary = content
			case "description":
				doc.description = content
			case "group":
				doc.group = content
			case "header":
				if doc.header == nil {
					doc.header = KV{}
				}
				k, v := splitKV(content)
				doc.header.Add(k, v)
			case "extra":
				if doc.extra == nil {
					doc.extra = KV{}
				}
				k, v := splitKV(content)
				doc.extra.Add(k, v)
			}

			if strings.HasPrefix(ann.Content, "req") {
				if len(ann.Relation) > 0 {
					doc.request = ann.Relation[0]
				}
			} else if strings.HasPrefix(ann.Content, "rsp") {
				if len(ann.Relation) > 0 {
					doc.response = ann.Relation[0]
				}
			}
		}
		docs = append(docs, doc)
	}
	return docs
}

func splitKV(content string) (string, string) {
	n := strings.Index(content, ":")
	if n <= 0 {
		return content, ""
	}
	key := strings.TrimSpace(content[:n])
	val := strings.TrimSpace(content[n+1:])
	return key, val
}
