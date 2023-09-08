package elements

const (
	TypeFunc   = 1
	TypeStruct = 2
)

type Node struct {
	Type        int           // 1 func， 2 struct
	Name        string        // func/struct name
	Ptr         any           //
	Path        string        // 完整路径（包括导包路径），如：github.com/example/foo.Hello
	Comments    []string      // 注释
	Annotations []*Annotation // 解析后的注解（以@符号开始的注释）
}

type Annotation struct {
	Raw      string
	Content  string
	Relation []any
}
