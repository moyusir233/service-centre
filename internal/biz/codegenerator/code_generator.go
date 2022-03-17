// Package codegenerator 负责生成数据收集服务和数据处理服务的服务端与客户端代码
package codegenerator

import (
	"bytes"
	v1 "gitee.com/moyusir/util/api/util/v1"
	"strings"
)

type CodeGenerator struct {
	dpRenderer *dataProcessingTmplRenderer
	dcRenderer *dataCollectionTmplRenderer
}

type GeneratedFile struct {
	Name    string
	Content *bytes.Buffer
}

func NewCodeGenerator(dpTmplDir, dcTmplDir string) (*CodeGenerator, error) {
	processingTmplRenderer, err := newDataProcessingTmplRenderer(dpTmplDir)
	if err != nil {
		return nil, err
	}

	collectionTmplRenderer, err := newDataCollectionTmplRenderer(dcTmplDir)
	if err != nil {
		return nil, err
	}

	return &CodeGenerator{
		dpRenderer: processingTmplRenderer,
		dcRenderer: collectionTmplRenderer,
	}, nil
}

func (g *CodeGenerator) GetServiceFiles(
	configInfo []v1.DeviceConfigRegisterInfo, stateInfo []v1.DeviceStateRegisterInfo) (
	dc []GeneratedFile, dp []GeneratedFile, err error) {
	var (
		configs       = make([]Device, len(configInfo))
		states        = make([]Device, len(stateInfo))
		warningStates = make([]Device, len(stateInfo))
	)
	dc = make([]GeneratedFile, 0, 4)
	dp = make([]GeneratedFile, 0, 4)

	// 处理配置注册信息
	for i, info := range configInfo {
		configs[i].DeviceClassID = i
		configs[i].Fields = make([]Field, len(info.Fields))

		for j, f := range info.Fields {
			configs[i].Fields[j].Name = f.Name
			configs[i].Fields[j].Type = strings.ToLower(f.Type.String())
		}
	}

	// 处理状态注册信息
	for i, info := range stateInfo {
		states[i].DeviceClassID = i
		states[i].Fields = make([]Field, len(info.Fields))
		warningStates[i].DeviceClassID = i

		for j, f := range info.Fields {
			states[i].Fields[j].Name = f.Name
			states[i].Fields[j].Type = strings.ToLower(f.Type.String())

			// 预警规则不为空即为预警字段
			if f.WarningRule != nil {
				warningStates[i].Fields = append(warningStates[i].Fields, states[i].Fields[j])
			}
		}
	}

	// 产生数据收集服务相关的代码与服务定义文件
	dcConfigCode, dcConfigProto, err := g.dcRenderer.renderConfigTmpl(configs)
	if err != nil {
		return nil, nil, err
	}
	dc = append(dc, GeneratedFile{Name: "config.go", Content: dcConfigCode})
	dc = append(dc, GeneratedFile{Name: "config.proto", Content: dcConfigProto})

	dcWarningCode, dcWarningProto, err := g.dcRenderer.renderWarningDetectTmpl(states, warningStates)
	if err != nil {
		return nil, nil, err
	}
	dc = append(dc, GeneratedFile{Name: "warning_detect.go", Content: dcWarningCode})
	dc = append(dc, GeneratedFile{Name: "warning_detect.proto", Content: dcWarningProto})

	// 产生数据处理服务相关的代码与服务定义文件
	dpConfigCode, dpConfigProto, err := g.dpRenderer.renderConfigTmpl(configs)
	if err != nil {
		return nil, nil, err
	}
	dp = append(dp, GeneratedFile{Name: "config.go", Content: dpConfigCode})
	dp = append(dp, GeneratedFile{Name: "config.proto", Content: dpConfigProto})

	dpWarningCode, dpWarningProto, err := g.dpRenderer.renderWarningDetectTmpl(states)
	if err != nil {
		return nil, nil, err
	}
	dp = append(dp, GeneratedFile{Name: "warning_detect.go", Content: dpWarningCode})
	dp = append(dp, GeneratedFile{Name: "warning_detect.proto", Content: dpWarningProto})

	return dc, dp, nil
}
