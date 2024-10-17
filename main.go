package main

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/cloudfoundry/bosh-template-renderer/renderer"
	"io"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	templateStream, err := os.Open(argsWithoutProg[0])
	if err != nil {
		panic(err)
	}
	template, err := renderer.Parse(templateStream)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic("no std in found")
	}
	jsonData, err := gabs.ParseJSON(data)
	if err != nil {
		panic("unable to parse property data")
	}
	fmt.Println(template.Render(jsonData))
}
