package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

var (
	escapeHelper = strings.NewReplacer(
		"\"", "\\\"",
		"\\", "\\\\",
		"\a", "\\a",
		"\b", "\\b",
		"\t", "\\t",
		"\n", "\\n",
		"\f", "\\f",
		"\r", "\\r",
		"\v", "\\v",
	)
)

// JsonBuilder 简单的json构造器，只支持基础类型的键值对
type JsonBuilder struct {
	buf      *bytes.Buffer
	hasField bool
}

func NewJsonBuilder() *JsonBuilder {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("{")
	return &JsonBuilder{buf: buf}
}

func (b *JsonBuilder) Reset() {
	b.buf.Reset()
	b.buf.WriteString("{")
}

func (b *JsonBuilder) Encode() []byte {
	b.buf.WriteString("}")
	defer b.Reset()
	return b.buf.Bytes()
}

func (b *JsonBuilder) String() string {
	return string(b.Encode())
}

// WriteField 添加键值对
// val 只支持基础类型，或者val必须实现String()接口并返回json字符串
func (b *JsonBuilder) WriteField(key string, val any) {
	if b.hasField {
		b.buf.Write([]byte{','})
	}
	b.hasField = true
	fmt.Fprintf(b.buf, "\"%s\":", key)
	switch v := val.(type) {
	case string:
		fmt.Fprintf(b.buf, "\"%s\"", escapeHelper.Replace(v))
	default:
		fmt.Fprintf(b.buf, "%v", val)
	}
}

func (b *JsonBuilder) WriteJSON(key string, json string) {
	if b.hasField {
		b.buf.Write([]byte{','})
	}
	b.hasField = true

	fmt.Fprintf(b.buf, "\"%s\": %s", key, json)
}

func (b *JsonBuilder) WriteJSONAny(key string, v any) {
	data, _ := json.Marshal(v)
	b.WriteJSON(key, string(data))
}

func (b *JsonBuilder) WriteJSONArray(key string, array ...string) {
	if b.hasField {
		b.buf.Write([]byte{','})
	}
	b.hasField = true

	fmt.Fprintf(b.buf, "\"%s\":", key)
	b.buf.WriteString("[")
	b.buf.WriteString(strings.Join(array, ","))
	b.buf.WriteString("]")
}
