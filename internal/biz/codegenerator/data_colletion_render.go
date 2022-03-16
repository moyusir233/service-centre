package codegenerator

import (
	"bytes"
	"path/filepath"
	"text/template"
)

type dataCollectionTmplRenderer struct {
	tmpl *template.Template
}

func newDataCollectionTmplRenderer(root string) (*dataCollectionTmplRenderer, error) {
	// 实例化模板，并注入函数，然后解析指定目录下所有模板文件
	tmpl, err := template.New("config.go.template").
		Funcs(funcMap).
		ParseGlob(filepath.Join(root, "*.template"))
	if err != nil {
		return nil, err
	}
	return &dataCollectionTmplRenderer{tmpl: tmpl}, nil
}

// 渲染配置管理相关的go源码模板和protobuf服务定义模板
func (r *dataCollectionTmplRenderer) renderConfigTmpl(configs []Device) (
	Code *bytes.Buffer, Proto *bytes.Buffer, err error) {
	Code = bytes.NewBuffer(make([]byte, 0, 1024))
	err = r.tmpl.ExecuteTemplate(Code, "config.go.template", configs)
	if err != nil {
		return nil, nil, err
	}

	Proto = bytes.NewBuffer(make([]byte, 0, 1024))
	err = r.tmpl.ExecuteTemplate(Proto, "config.proto.template", configs)
	if err != nil {
		return nil, nil, err
	}

	return
}

// 渲染故障预警相关的go源码模板和protobuf服务定义模板
func (r *dataCollectionTmplRenderer) renderWarningDetectTmpl(states []Device, warningDetectStates []Device) (
	Code *bytes.Buffer, Proto *bytes.Buffer, err error) {
	Code = bytes.NewBuffer(make([]byte, 0, 1024))
	// go源码部分的生成需要注入所有设备的预警字段信息，因此传入warningDetectStates
	err = r.tmpl.ExecuteTemplate(Code, "warning_detect.go.template", warningDetectStates)
	if err != nil {
		return nil, nil, err
	}

	Proto = bytes.NewBuffer(make([]byte, 0, 1024))
	err = r.tmpl.ExecuteTemplate(Proto, "warningDetect.proto.template", states)
	if err != nil {
		return nil, nil, err
	}

	return
}
