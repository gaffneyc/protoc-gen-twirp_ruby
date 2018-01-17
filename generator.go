package main

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/template"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

type generator struct{}

func (g *generator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	resp := &plugin.CodeGeneratorResponse{}

	for _, name := range req.FileToGenerate {
		file, err := getFileDescriptor(req, name)
		if err != nil {
			return nil, err
		}

		genFile, err := g.generateFile(file)
		if err != nil {
			return nil, errors.Wrapf(err, "generating %q", name)
		}

		// Add the generated file to the response
		resp.File = append(resp.File, genFile)
	}

	return resp, nil
}

func (g *generator) generateFile(file *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	var err error

	buffer := &bytes.Buffer{}

	tmpl := template.New("ruby_file")
	tmpl, err = tmpl.Parse(rubyTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}

	// TODO: modulize
	// TODO: camelize
	presenter := &rubyTemplatePresenter{
		proto: file,

		Version:        "1.0.0",
		SourceFilename: file.GetName(),
		Package:        file.GetPackage(),
	}

	err = tmpl.Execute(buffer, presenter)
	if err != nil {
		return nil, errors.Wrapf(err, "rendering template for %q", file.GetName())
	}

	resp := &plugin.CodeGeneratorResponse_File{}
	resp.Name = proto.String("service_twirp.rb")
	resp.Content = proto.String(buffer.String())

	return resp, nil
}

// getFileDescriptor finds the FileDescriptorProto for the given filename.
// Returns a error if the descriptor could not be found.
func getFileDescriptor(req *plugin.CodeGeneratorRequest, name string) (*descriptor.FileDescriptorProto, error) {
	for _, descriptor := range req.ProtoFile {
		if descriptor.GetName() == name {
			return descriptor, nil
		}
	}

	return nil, fmt.Errorf("could not find descriptor for %q", name)
}
