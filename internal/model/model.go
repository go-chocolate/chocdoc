package model

type Model struct {
	Name      string   `json:"name"`
	Fields    []*Field `json:"fields"`
	Array     int      `json:"array"`
	anonymous bool
}
