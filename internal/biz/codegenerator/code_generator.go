// Package codegenerator 负责生成数据收集服务和数据处理服务的服务端与客户端代码
package codegenerator

import (
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"strings"
	"text/template"
)

var (
	funcMap = template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"indent": func(count int, s string) string {
			return strings.Repeat(" ", count) + s
		},
		"camelCase": func(s string) string {
			return generator.CamelCase(s)
		},
	}
)

type CodeGenerator struct {
}

type Device struct {
	DeviceClassID int
	Fields        []Field
}

type Field struct {
	Name, Type string
}

func NewCodeGenerator() *CodeGenerator {
	return nil
}
