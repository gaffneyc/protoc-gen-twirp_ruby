package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

// Version is the release version of the ruby twirp generator
const Version = "1.0.0"

func main() {
	var err error

	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	// protoc pipes the information about what needs generating into os.Stdin.
	// Read the request and return it.
	req, err := readCodeGeneratorRequest(os.Stdin)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// Validate the request
	if len(req.FileToGenerate) == 0 {
		fmt.Println("no files to generate")
		os.Exit(1)
	}

	// Generate the code
	gen := &generator{}
	resp, err := gen.Generate(req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// Write response to os.Stdout
	err = writeCodeGeneratorResponse(os.Stdout, resp)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func readCodeGeneratorRequest(in io.Reader) (*plugin.CodeGeneratorRequest, error) {
	var err error

	// Read the full request before trying to parse it
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, errors.Wrap(err, "reading CodeGeneratorRequest")
	}

	// Unmarshal the request
	req := &plugin.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, req)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling CodeGeneratorRequest")
	}

	return req, nil
}

func writeCodeGeneratorResponse(out io.Writer, resp *plugin.CodeGeneratorResponse) error {
	var err error

	data, err := proto.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "marshaling CodeGeneratorResponse")
	}

	_, err = out.Write(data)
	if err != nil {
		return errors.Wrap(err, "writing CodeGeneratorResponse")
	}

	return nil
}
