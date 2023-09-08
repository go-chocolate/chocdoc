package annotation

import "testing"

func TestExport(t *testing.T) {
	if err := Export(
		WithRoot("testdata"),
		WithOutput("testdata/annotation"),
		WithSaveDotAnnotationFile(),
	); err != nil {
		t.Error(err)
	}
}
