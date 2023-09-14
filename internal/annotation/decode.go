package annotation

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type decoder struct {
	root     string
	replaces map[string][]string
}

func (d *decoder) decode() ([]*file, error) {
	m, err := decodeGoMod(d.root)
	if err != nil {
		return nil, err
	}
	var files []*file
	err = filepath.Walk(d.root, func(path string, info fs.FileInfo, err error) error {
		if path == d.root {
			return err
		}
		if info.IsDir() && isFileExist(path+"/"+dotAnnotationFile) {
			return filepath.SkipDir
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		//解析go文件
		fi, err := d.decodeFile(m.module, path)
		if err != nil {
			return err
		}

		//解析注解
		for _, n := range fi.nodes {
			decodeAnnotations(fi.context(), n)
		}
		files = append(files, fi)
		return nil
	})
	return files, err
}

func (d *decoder) decodeFile(module, filename string) (*file, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	re := bufio.NewReader(fi)

	var result = new(file)
	result.filename = strings.ReplaceAll(filename, "\\", "/")

	if n := strings.LastIndex(result.filename, "/"); n >= 0 {
		result.path = result.filename[:n]
	} else {
		result.path = "."
	}

	var comments []string
	for {
		line, err := readLine(re)
		if err != nil {
			break
		}

		if !strings.HasPrefix(line, "/") {
			line = trimInlineComment(line) //过滤行内注释
		}

		if strings.Contains(line, "`") {
			line, err = trimMultiLineText(line, re) //过滤多行文本
			if err != nil {
				break
			}
			//多行文本内可能包含任意内容，对解析会造成干扰，必须去除
		}

		if strings.HasPrefix(line, "package") {
			result.pkg = strings.TrimSpace(strings.TrimPrefix(line, "package")) //包名
			result.packageComments = comments                                   //包注释
			impPath := strings.Trim(strings.TrimPrefix(result.path, d.root), "/")
			result.imp = fmt.Sprintf("%s/%s", module, impPath) //导包列表
			comments = []string{}
			continue
		}

		if strings.HasPrefix(line, "import") {
			result.imports = append(result.imports, d.readImports(line, re)...) //导包列表
			comments = []string{}
			continue
		}

		// 解析节点
		if strings.HasPrefix(line, "func") || strings.HasPrefix(line, "type") {
			e := d.readFuncOrStruct(line, re, comments)
			e.path = fmt.Sprintf("%s.%s", result.imp, e.name)
			result.nodes = append(result.nodes, e)
			comments = []string{}
			continue
		}

		//解析注释
		var isComment = false
		if strings.HasPrefix(line, "/*") {
			isComment = true
			comments = append(comments, d.readComments(line, re)...)
		}
		if strings.HasPrefix(line, "//") {
			isComment = true
			comments = append(comments, strings.TrimSpace(strings.TrimPrefix(line, "//")))
		}

		if !isComment {
			comments = []string{}
		}

	}
	return result, nil
}

func (d *decoder) readComments(lastLine string, re *bufio.Reader) []string {
	if strings.HasSuffix(lastLine, "*/") {
		return []string{strings.Trim(lastLine, "/*")}
	}
	var comments []string
	if theFirst := strings.TrimLeft(lastLine, "/*"); theFirst != "" {
		comments = append(comments, theFirst)
	}
	for {
		line, err := readLine(re)
		if err != nil {
			return comments
		}
		if strings.HasSuffix(line, "*/") {
			if theLast := strings.TrimRight(line, "*/"); theLast != "" {
				comments = append(comments, theLast)
			}
			return comments
		} else {
			comments = append(comments, line)
		}
	}
}

// 解析多行import
func (d *decoder) readImports(lastLine string, re *bufio.Reader) []*imports {
	if !strings.Contains(lastLine, "(") {
		return []*imports{d.readImport(lastLine)}
	}
	var imps []*imports
	for {
		line, err := readLine(re)
		if err != nil {
			return imps
		}
		if line == "" {
			continue
		}
		if line == ")" {
			return imps
		}
		imps = append(imps, d.readImport(line))
	}
}

// 解析单行import
func (d *decoder) readImport(line string) *imports {
	line = strings.TrimSpace(strings.TrimPrefix(line, "import"))
	temp := strings.Fields(line)
	var im = new(imports)
	if len(temp) == 1 {
		im.path = strings.Trim(temp[0], "\"")
	} else if len(temp) == 2 {
		im.path = strings.Trim(temp[1], "\"")
		im.alias = []string{temp[0]}
	}
	if alias := d.replaces[im.path]; len(alias) > 0 {
		im.alias = append(im.alias, alias...)
	}
	return im
}

// 解析func和struct
func (d *decoder) readFuncOrStruct(lastLine string, re *bufio.Reader, comments []string) *node {
	lastLine = trimInlineComment(lastLine)
	funcFields := strings.Fields(lastLine)
	var name = funcFields[1]
	if strings.HasPrefix(lastLine, "func") {
		//if name[0] == '(' {
		//	name = funcFields[3]
		//}
		//name = strings.Split(name, "(")[0]
		name = funcName(lastLine)
	}

	var e = new(node)
	e.name = name
	e.comments = comments
	e.exported = name[0] >= 'A' && name[0] <= 'Z'
	if strings.HasPrefix(lastLine, "func") {
		e.typ = typeFunc
	} else if strings.HasPrefix(lastLine, "type") {
		e.typ = typeStruct
	}

	var left, right int
	left = strings.Count(lastLine, "{")
	right = strings.Count(lastLine, "}")

	for left != right {
		line, err := readLine(re)
		if err != nil {
			break
		}
		left += strings.Count(line, "{")
		right += strings.Count(line, "}")
	}
	return e
}

func funcName(line string) string {
	tmp := strings.TrimSpace(strings.TrimPrefix(line, "func "))
	n := strings.Index(tmp, "(")
	if n == 0 {
		n2 := strings.Index(tmp, ")")
		n3 := strings.Index(tmp[n2+1:], "(")
		return strings.TrimSpace(tmp[n2+1 : n2+1+n3])
	} else if n > 0 {

		return strings.TrimSpace(tmp[:n])
	}
	return ""
}

func readLine(re *bufio.Reader) (string, error) {
	data, err := re.ReadBytes('\n')
	if err != nil || len(data) == 0 {
		return "", io.EOF
	}
	line := strings.TrimSpace(string(data))
	return line, nil
}

// 去除行内注释
func trimInlineComment(line string) string {
	var quotes bool
	var slash bool
	var apostrophe bool
	for i, v := range line {
		if v == '/' {
			if quotes || apostrophe {
				continue
			} else if slash {
				return strings.TrimSpace(line[:i-1])
			} else {
				slash = true
			}
		} else {
			slash = false
		}
		if v == '"' {
			quotes = !quotes
		}
		if v == '`' && !quotes {
			apostrophe = !apostrophe
		}
	}
	return line
}

// 去除多行文本
func trimMultiLineText(line string, re *bufio.Reader) (string, error) {
	count := strings.Count(line, "`")
	if count%2 == 0 {
		return line, nil
	}
	for {
		next, err := readLine(re)
		if err != nil {
			return next, err
		}
		count += strings.Count(next, "`")
		if count%2 == 0 {
			return readLine(re)
		}
	}
}

var reg, _ = regexp.Compile("\\[[\\w.]+\\]|<[\\w.]+>")

// 解析注解
func decodeAnnotations(ctx *context, n *node) {
	for _, line := range n.comments {
		if !strings.HasPrefix(line, "@") {
			continue
		}
		i := &annotation{raw: line}
		i.content = line[1:]
		for _, v := range reg.FindAllString(line, -1) {
			rel := strings.Trim(strings.Trim(v, "[]"), "<>")
			n := strings.LastIndex(rel, ".")
			var relationPath string
			if n > 0 {
				prefix := rel[:n]
				name := rel[n+1:]
				if p := ctx.imports[prefix]; p != "" {
					relationPath = fmt.Sprintf("%s.%s", p, name)
				} else {
					relationPath = rel
				}
			} else {
				relationPath = fmt.Sprintf("%s.%s", ctx.importPath, rel)
			}
			re := &relation{path: relationPath}
			if strings.HasPrefix(v, "[") {
				re.typ = typeStruct
			} else if strings.HasPrefix(v, "<") {
				re.typ = typeFunc
			}
			i.relations = append(i.relations, re)
			i.content = strings.ReplaceAll(i.content, v, v[0:1]+relationPath+v[len(v)-1:])
		}
		n.annotations = append(n.annotations, i)
	}
}
