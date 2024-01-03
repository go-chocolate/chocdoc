package main

import (
	"flag"
	"os"

	"github.com/go-chocolate/chocdoc/internal/annotation"
)

var (
	root   string
	output string
)

func init() {
	flag.StringVar(&root, "root", ".", "chocdoc执行路径，默认当前目录")
	flag.StringVar(&output, "output", "./chocdoc", "chocdoc生成代码路径，默认 ./chocdoc")
	flag.Parse()
}

func main() {
	var err error
	if annotation.IsDotOptionFileExist() {
		err = annotation.Export(annotation.WithDotAnnotationFile())
	} else {
		err = annotation.Export(
			annotation.WithRoot(root),
			annotation.WithOutput(output),
			annotation.WithSaveDotAnnotationFile(),
		)
	}
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	} else {
		os.Stdout.WriteString("chocdoc generated successfully")
	}
}
