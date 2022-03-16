package codegenerator

import (
	"bytes"
	"os"
	"testing"
)

func TestDataCollection_RenderConfigTmpl(t *testing.T) {
	renderer, err := newDataCollectionTmplRenderer("data-collection-template")
	if err != nil {
		t.Fatal(err)
	}

	configs := []Device{
		{
			DeviceClassID: 0,
			Fields: []Field{
				{
					Name: "id",
					Type: "string",
				},
				{
					Name: "test0_1",
					Type: "int64",
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
					Name: "test1_1",
					Type: "int64",
				},
			},
		},
	}
	code, proto, err := renderer.renderConfigTmpl(configs)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n%s", code.String(), proto.String())

	files := []string{"output/config.go", "output/config.proto"}
	buffers := []*bytes.Buffer{code, proto}
	for i, f := range files {
		file, err := os.Create(f)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		_, err = buffers[i].WriteTo(file)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDataCollection_RenderWarningDetectTmpl(t *testing.T) {
	renderer, err := newDataCollectionTmplRenderer("data-collection-template")
	if err != nil {
		t.Fatal(err)
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
					Name: "Current",
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
					Name: "Current",
					Type: "double",
				},
			},
		},
	}
	warningDetectStates := []Device{
		{
			DeviceClassID: 0,
			Fields: []Field{
				{
					Name: "Current",
					Type: "double",
				},
			},
		},
		{
			DeviceClassID: 1,
			Fields: []Field{
				{
					Name: "Current",
					Type: "double",
				},
			},
		},
	}
	code, proto, err := renderer.renderWarningDetectTmpl(states, warningDetectStates)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n%s", code.String(), proto.String())

	files := []string{"output/warning_detect.go", "output/warningDetect.proto"}
	buffers := []*bytes.Buffer{code, proto}
	for i, f := range files {
		file, err := os.Create(f)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		_, err = buffers[i].WriteTo(file)
		if err != nil {
			t.Fatal(err)
		}
	}
}
