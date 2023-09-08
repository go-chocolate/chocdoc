package doc

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Array    int    `json:"array"`
	Required bool   `json:"required"`
	Comment  string `json:"comment"`
	Option   string `json:"option"`
	Sub      *Model `json:"sub"`
	Tags     KV     `json:"tags"`
	sub      string
}

func (f *Field) SetName(name string) {
	if f.Name == "" {
		f.Name = name
	}
}

type Model struct {
	Name      string   `json:"name"`
	Fields    []*Field `json:"fields"`
	Array     int      `json:"array"`
	anonymous bool
}

type Example struct {
	Type uint8 `doc:"required;this is field comment;option:1,2,3"`
}

func (m *Model) GetFields() []*Field {
	if m == nil {
		return nil
	}
	return m.Fields
}

func (m *Model) json(prefix string) string {
	var fields []string
	for _, field := range m.GetFields() {
		var text string
		switch strings.TrimLeft(field.Type, "[]") {
		case "string":
			text = text + "\"\""
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float64", "float32":
			text = text + "0"
		case "bool":
			text = text + "false"
		case "Time":
			text = text + "\"2006-01-02 15:04:05\""
		default:
			if field.Sub != nil {
				text = text + field.Sub.json(prefix+"\t")
			} else {
				text = text + "{}"
			}
		}
		if field.Array > 0 {
			text = strings.Repeat("[", field.Array) + text + strings.Repeat("]", field.Array)
		}
		text = fmt.Sprintf("\t\"%s\": ", field.Name) + text
		fields = append(fields, prefix+text)
	}

	if m.Array > 0 {
		return strings.Repeat("[", m.Array) +
			strings.Join(fields, ",\n") +
			strings.Repeat("]", m.Array)
	} else {
		return "{\n" + strings.Join(fields, ",\n") + "\n" + prefix + "}"
	}
}

func (m *Model) JSON() string {
	return m.json("")
}

// 记录结构体链路，避免指针无限递归
// eg:
//
//	type Model struct{
//	  Child *Model
//	}
type tree []string

func (t *tree) join(name string) {
	*t = append(*t, name)
}

func (t tree) contain(name string) bool {
	for _, v := range t {
		if name == v {
			return true
		}
	}
	return false
}

func (t tree) skip(pkg, name string) bool {
	if pkg == "time" && name == "Time" {
		return true
	}
	//and so on...
	return false
}

type decoder struct {
	t tree
}

func newDecoder(t tree) *decoder {
	return &decoder{
		t: t,
	}
}

// 解析结构体
func (d *decoder) decode(model interface{}) *Model {
	if model == nil {
		return &Model{}
	}
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	m := new(Model)
	m.Name = t.Name()
	d.t.join(t.Name())

	rt, n := realType(t, 0)
	m.Array = n

	switch rt.Kind() {
	//case reflect.Map:
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Tag.Get("doc") == "-" {
				continue
			}
			field := d.decodeField(t.Field(i), v.Field(i))
			if field.Sub != nil && field.Sub.anonymous {
				m.Fields = append(m.Fields, field.Sub.Fields...)
			} else {
				m.Fields = append(m.Fields, field)
			}
		}
	default:
	}
	return m
}

// 解析字段
func (d *decoder) decodeField(t reflect.StructField, v reflect.Value) *Field {
	var f = new(Field)
	f.Name = strings.Split(t.Tag.Get("json"), ",")[0]
	f.Tags = StructTag(t.Tag).Lookups()
	f.SetName(t.Name)
	decodeDocTag(f, t.Tag.Get("doc"))
	rt, n := realType(v.Type(), 0)
	f.Array = n
	switch rt.Kind() {
	case reflect.Struct:
		f.Type = rt.Name()
		if d.t.skip(rt.PkgPath(), v.Type().Name()) {
			f.Type = v.Type().Name()
		} else if !d.t.contain(rt.Name()) {
			f.Sub = newDecoder(d.t).decode(reflect.New(rt).Interface())
			f.Sub.anonymous = t.Anonymous
		}
	//case reflect.Map:
	//TODO
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
		f.Type = rt.Kind().String()
	default:
		f.Type = rt.Name()
	}
	if n > 0 {
		f.Type = strings.Repeat("[]", n) + f.Type
	}
	return f
}

// 解析doc标签
func decodeDocTag(f *Field, doc string) {
	fields := strings.Split(doc, ";")
	for _, v := range fields {
		v = strings.TrimSpace(v)
		tmp := strings.Split(v, ":")
		switch len(tmp) {
		case 0:
			continue
		case 1:
			if v == "required" || v == "must" {
				f.Required = true
			} else {
				f.Comment = v
			}
		case 2:
			switch tmp[0] {
			case "option", "opt":
				if len(tmp) > 1 {
					f.Option = tmp[1]
				}
			}
		default:
			continue
		}
	}
}

// 获取指针/切片/数组内的真实类型
func realType(v reflect.Type, i int) (reflect.Type, int) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		i++
		fallthrough
	case reflect.Ptr:
		return realType(v.Elem(), i)
	default:
		return v, i
	}
}

func DecodeModel(v interface{}) *Model {
	return newDecoder(tree{}).decode(v)
}

type StructTag string

func (tag StructTag) Lookups() KV {

	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	kv := KV{}

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if val, err := strconv.Unquote(qvalue); err == nil {
			kv.Add(name, val)
		}
	}
	return kv
}
