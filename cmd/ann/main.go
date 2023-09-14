package main

import (
	"fmt"
	"os"

	"github.com/go-chocolate/chocdoc/internal/annotation"
)

func main() {
	var err error
	if annotation.IsDotOptionFileExist() {
		err = annotation.Export(annotation.WithDotAnnotationFile())
	} else {
		err = annotation.Export(
			annotation.WithRoot("."),
			annotation.WithOutput("./annotation"),
			annotation.WithSaveDotAnnotationFile(),
		)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
