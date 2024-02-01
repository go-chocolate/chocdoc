package model

import (
	"reflect"
	"strings"
)

type decoder struct {
	tree tree
}

// 解析结构体
func (d *decoder) decode(model interface{}) *Model {
	if model == nil {
		return &Model{}
	}
	//t := reflect.TypeOf(model)
	//v := reflect.ValueOf(model)
	//if t.Kind() == reflect.Ptr {
	//	t = t.Elem()
	//	v = v.Elem()
	//}

	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := new(Model)
	m.Name = t.Name()
	d.tree.join(t.Name())

	rt, n := realType(t, 0)
	m.Array = n

	switch rt.Kind() {
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
	d.decodeDocTag(f, t.Tag.Get("doc"))
	rt, n := realType(v.Type(), 0)
	f.Array = n
	switch rt.Kind() {
	case reflect.Struct:
		f.Type = rt.Name()
		if d.tree.skip(rt.PkgPath(), v.Type().Name()) {
			f.Type = v.Type().Name()
		} else if !d.tree.contain(rt.Name()) {
			f.Sub = (&decoder{tree: d.tree}).decode(reflect.New(rt).Interface())
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
func (d *decoder) decodeDocTag(f *Field, doc string) {
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
