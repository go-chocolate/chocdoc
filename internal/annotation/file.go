package annotation

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	typeFunc   = 1
	typeStruct = 2

	nodeTemplate = `{
		Type: %s,
		Name: "%s",
		Ptr: %s,
		Path: "%s",
		Comments: []string{%s},
		Annotations: []*elements.Annotation{
%s
		},
	}`

	annotationTemplate = `        {
				Raw: "%s",
				Content: "%s",
				Relation: []interface{}{%s},
			}`
)

// 包含注解的代码节点
type node struct {
	typ         int           // 1 func， 2 struct
	name        string        // func/struct name
	path        string        // 完整路径（包括导包路径），如：github.com/example/foo.Hello
	comments    []string      // 注释
	exported    bool          // 是否为可导出
	annotations []*annotation // 解析后的注解（以@符号开始的注释）
}

type annotation struct {
	raw       string      // 注解原文（包含@符号）
	content   string      // 注解内容（不包含@符号）
	relations []*relation // 注解内关联对象 以方括号（[]）或尖括号（<>）包裹的内容
}

type relation struct {
	path string //引用地址
	typ  int    //对象类型 1 func， 2 struct
}

type imports struct {
	path  string   // import路径
	alias []string // 别名
}

type context struct {
	imports    map[string]string //map[包名]完整导包路径
	pkg        string
	importPath string
}

type file struct {
	path            string     // 文件路径
	filename        string     // 文件名
	pkg             string     // 所属package
	imp             string     // 本文件导包路径
	imports         []*imports // 本文件内导入的包
	packageComments []string   // 包注释
	nodes           []*node    // 节点
}

func (f *file) context() *context {
	ctx := &context{}
	ctx.pkg = f.pkg
	ctx.importPath = f.imp
	ctx.imports = make(map[string]string)
	for _, v := range f.imports {
		if n := strings.LastIndex(v.path, "/"); n > 0 {
			ctx.imports[v.path[n+1:]] = v.path
		} else {
			ctx.imports[v.path] = v.path
		}
		for _, alias := range v.alias {
			if alias == "." || alias == "_" || alias == "" {
				continue
			}
			ctx.imports[alias] = v.path
		}
	}
	return ctx
}

type mod struct {
	module   string   //module名
	version  string   //go版本
	requires []string //TODO
}

func decodeGoMod(root string) (*mod, error) {
	b, err := os.ReadFile(root + "/go.mod")
	if err != nil {
		return nil, err
	}
	var m = new(mod)
	for _, v := range strings.Split(string(b), "\n") {
		line := strings.TrimSpace(v)
		if strings.HasPrefix(line, "module") {
			m.module = strings.TrimSpace(strings.TrimPrefix(line, "module"))
		} else if strings.HasPrefix(line, "go") {
			m.version = strings.TrimSpace(strings.TrimPrefix(line, "go"))
		}
	}
	return m, nil
}

func (n *node) format() (string, []string) {
	var imports []string
	var anns []string
	for _, v := range n.annotations {
		annName, annImports := v.format()
		imports = append(imports, annImports...)
		anns = append(anns, annName)
	}
	var comments string
	if len(n.comments) > 0 {
		comments = fmt.Sprintf("\"%s\"", strings.Join(n.comments, "\", \""))
	}

	var impt, name = splitImport(n.path)
	imports = append(imports, impt)

	var ptr string
	switch n.typ {
	case typeFunc:
		ptr = name
	case typeStruct:
		ptr = fmt.Sprintf("new(%s)", name)
	}

	text := fmt.Sprintf(nodeTemplate,
		n.typeName(),
		n.name,
		ptr,
		n.path,
		comments,
		strings.Join(anns, ",\n")+",",
	)
	return text, imports
}

func (n *node) typeName() string {
	switch n.typ {
	case typeFunc:
		return "elements.TypeFunc"
	case typeStruct:
		return "elements.TypeStruct"
	}
	return strconv.Itoa(n.typ)
}

func (i *annotation) format() (string, []string) {
	var rel string
	var imports []string

	for _, re := range i.relations {
		imp, name := splitImport(re.path)
		imports = append(imports, imp)
		switch re.typ {
		case typeFunc:
			rel += fmt.Sprintf("%s, ", name)
		case typeStruct:
			rel += fmt.Sprintf("new(%s), ", name)
		}
	}
	text := fmt.Sprintf("    "+annotationTemplate,
		i.raw,
		i.content,
		rel,
	)
	return text, imports
}

func splitImport(path string) (string, string) {
	var name string
	if n := strings.LastIndex(path, "/"); n > 0 {
		name = path[n+1:]
	} else {
		name = path
	}
	var imp string
	if n := strings.LastIndex(path, "."); n > 0 {
		imp = path[:n]
	} else {
		imp = path
	}
	return imp, name
}
