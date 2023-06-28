package krt

import (
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/konstellation-io/krt/pkg/parse"
)

func ParseFile(yamlFile string) (*krt.Krt, error) {
	return parse.ParseFile(yamlFile)
}

func ParseString(krtYaml []byte) (*krt.Krt, error) {
	return parse.ParseKrt(krtYaml)
}
