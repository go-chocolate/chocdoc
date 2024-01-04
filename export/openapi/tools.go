package openapi

import (
	"strings"

	"github.com/go-chocolate/chocdoc/internal/doc"
)

func splitPathParameter(path string) []string {
	var parameters []string
	var parameter []byte
	var mark bool
	for _, v := range path {
		if v == ':' {
			mark = true
			continue
		}
		if v == '/' {
			mark = false
			if len(parameter) > 0 {
				parameters = append(parameters, string(parameter))
				parameter = []byte{}
			}
			continue
		}
		if !mark {
			continue
		}
		parameter = append(parameter, byte(v))
	}

	if len(parameter) > 0 {
		parameters = append(parameters, string(parameter))
	}
	return parameters
}

func in[T comparable](array []T, item T) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}

func arrayDepth(field *doc.Field) (bool, int) {
	if strings.HasPrefix(field.Type, "[]") {
		return true, strings.Count(field.Type, "[]")
	}
	return false, 0
}

func formatGoType(field *doc.Field) Type {
	switch strings.TrimLeft(field.Type, "[]") {
	case "string":
		return TypeString
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return TypeInteger
	case "bool":
		return TypeBoolean
	case "float32", "float64":
		return TypeNumber
	case "any", "interface{}":
		return TypeUnknown
	default:
		return TypeObject
	}
}
