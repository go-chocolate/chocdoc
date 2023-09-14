package annotation

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	dotOptionFile            = ".ann.yaml"
	dotOptionFileDescription = "# the annotation export option.\n#\n" +
		"# root<string>: the annotation exporter running directory.\n#\n" +
		"# output<string>: the directory where the annotation exporter generated code saved.\n#\n" +
		"# replaces<map[string][]string>: the import path alias. \n" +
		"# Some package's name is different from its folder name, may need alias to make sure \n" +
		"# the annotation exporter be able to found the correct import path.\n#\n" +
		"# example: \n#\n" +
		"# root: .\n" +
		"# output: ./annotation\n" +
		"# replace: \n" +
		"#   - github.com/example/example/v1: \n" +
		"#       - example1\n" +
		"#       - example2\n"
)

type option struct {
	Root     string              `yaml:"root"`
	Output   string              `yaml:"output"`
	Replaces map[string][]string `yaml:"replaces"`

	save bool `yaml:"-"`
}

func fromOptionFile() (*option, error) {
	data, err := os.ReadFile(dotOptionFile)
	if err != nil {
		return nil, err
	}
	var opt = new(option)
	err = yaml.Unmarshal(data, opt)
	return opt, err
}

func (o *option) saveToFile() error {
	data, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	data = append([]byte(dotOptionFileDescription+"\n"), data...)
	_ = os.Remove(dotOptionFile)
	err = os.WriteFile(dotOptionFile, data, 0644)
	return err
}

func (o *option) from(src *option) {
	o.Root = src.Root
	o.Output = src.Output
	o.Replaces = src.Replaces
}

func IsDotOptionFileExist() bool {
	return isFileExist(dotOptionFile)
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
