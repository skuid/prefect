package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/skuid/prefect/render"
	flag "github.com/ogier/pflag"
	"gopkg.in/yaml.v2"
)

// flags
var (
	context    string
	selector   selectorSet
	targetFile string
)

func dieIfError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	flag.StringVarP(&context, "context", "c", "", "The context values file to use")
	flag.VarP(&selector, "selector", "s", `The selectors to use. Each selector should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list`)
	flag.StringVarP(&targetFile, "target", "t", "", "A target file to use. This is ignored if selectors are specified")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "prefect takes a template and injects context from a given context file")
		fmt.Fprint(os.Stderr, ".\n\n")
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "%s [options] <template>\n\n", os.Args[0])
		fmt.Fprint(os.Stderr, "  <template>\n")
		fmt.Fprint(os.Stderr, "    	The template to read in\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Incorrect number of arguments! Only one argument for the template is needed")
		os.Exit(1)
	}
	template := flag.Args()[0]

	templateData, err := ioutil.ReadFile(template)
	dieIfError(err)

	contextData, err := ioutil.ReadFile(context)
	dieIfError(err)

	if len(selector) > 0 {

		target := render.Target{
			Name:     selector.String(),
			Selector: selector.ToMap(),
		}

		output, err := execute(template, templateData, contextData, target)
		if err != nil {
			fmt.Printf("Error rendering template: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Print(output)
	} else {
		targetData, err := ioutil.ReadFile(targetFile)
		dieIfError(err)

		targets := render.Targets{}
		yaml.Unmarshal(targetData, &targets)

		// Render, then create the path and file for each target
		for _, target := range targets {

			rendered, err := execute(template, templateData, contextData, target)
			dieIfError(err)

			outputPath := filepath.Join(target.Directory, template)

			if len(target.Directory) == 0 {
				fmt.Printf("# Rendered for target %s\n", target.Name)
				fmt.Print(rendered)

			} else {
				err = target.MakedirAll()
				dieIfError(err)

				f, err := os.Create(outputPath)
				dieIfError(err)

				_, err = f.WriteString(rendered)
				f.Close()
				dieIfError(err)

				fmt.Printf("Rendered file at %s\n", outputPath)
			}
		}

	}
}

func execute(templateName string, templateData, contextData []byte, target render.Target) (string, error) {
	doc := render.Document{Name: templateName, Content: string(templateData)}

	values := render.ContextValues{}
	err := yaml.Unmarshal(contextData, &values)
	if err != nil {
		return "", err
	}

	return doc.Render(target, values)
}
