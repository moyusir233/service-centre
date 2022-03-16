package codegenerator

import (
	"bytes"
	"os"
	"testing"
)

func TestRenderServiceTmpl(t *testing.T) {
	configs := []Device{
		{
			DeviceClassID: 0,
			Fields: []Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "status",
					Type: "bool",
				},
			},
		},
		{
			DeviceClassID: 1,
			Fields: []Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "status",
					Type: "bool",
				},
			},
		},
	}
	states := []Device{
		{
			DeviceClassID: 0,
			Fields: []Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "time",
					Type: "google.protobuf.Timestamp",
				},
				{
					Name: "voltage",
					Type: "double",
				},
				{
					Name: "current",
					Type: "double",
				},
				{
					Name: "temperature",
					Type: "double",
				},
			},
		},
		{
			DeviceClassID: 1,
			Fields: []Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "time",
					Type: "google.protobuf.Timestamp",
				},
				{
					Name: "voltage",
					Type: "double",
				},
				{
					Name: "current",
					Type: "double",
				},
				{
					Name: "temperature",
					Type: "double",
				},
			},
		},
	}
	warningDetectInfo := []Device{
		{
			DeviceClassID: 0,
			Fields: []Field{
				{
					Name: "voltage",
					Type: "double",
				},
				{
					Name: "current",
					Type: "double",
				},
				{
					Name: "temperature",
					Type: "double",
				},
			},
		},
		{
			DeviceClassID: 1,
			Fields: []Field{
				{
					Name: "voltage",
					Type: "double",
				},
				{
					Name: "current",
					Type: "double",
				},
				{
					Name: "temperature",
					Type: "double",
				},
			},
		},
	}

	collectionTmplRenderer, err := newDataCollectionTmplRenderer("data-collection-template")
	if err != nil {
		t.Fatal(err)
	}
	processingTmplRenderer, err := newDataProcessingTmplRenderer("data-processing-template")
	if err != nil {
		t.Fatal(err)
	}

	dirs := []string{"output/data-collection", "output/data-processing"}
	for _, dir := range dirs {
		os.MkdirAll(dir, os.ModeDir)
	}

	var buffers []*bytes.Buffer

	configCode1, configProto1, err := collectionTmplRenderer.renderConfigTmpl(configs)
	if err != nil {
		t.Fatal(err)
	}
	buffers = append(buffers, configCode1, configProto1)

	wdCode1, wdProto1, err := collectionTmplRenderer.renderWarningDetectTmpl(states, warningDetectInfo)
	if err != nil {
		t.Fatal(err)
	}
	buffers = append(buffers, wdCode1, wdProto1)

	configCode2, configProto2, err := processingTmplRenderer.renderConfigTmpl(configs)
	if err != nil {
		t.Fatal(err)
	}
	buffers = append(buffers, configCode2, configProto2)

	wdCode2, wdProto2, err := processingTmplRenderer.renderWarningDetectTmpl(states)
	if err != nil {
		t.Fatal(err)
	}
	buffers = append(buffers, wdCode2, wdProto2)

	prefixs := []string{"config", "warning_detect"}
	suffixs := []string{".go", ".proto"}
	for i, b := range buffers {
		dir := dirs[i/4]
		prefix := prefixs[i/2%2]
		suffix := suffixs[i%2]
		file, err := os.Create(dir + "/" + prefix + suffix)
		if err != nil {
			t.Fatal(err)
		}
		_, err = b.WriteTo(file)
		if err != nil {
			return
		}
		file.Close()
	}

}
