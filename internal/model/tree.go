package model

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
