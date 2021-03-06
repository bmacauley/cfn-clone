package main

import (
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	Attributes []string `short:"a" long:"attributes" description:"'=' separated attribute and value"`
	NewName    string   `short:"n" long:"new-name" description:"Name for new stack" required:"true"`
	SourceName string   `short:"s" long:"source-name" description:"Name of source stack to clone" required:"true"`
	Template   string   `short:"t" long:"template" description:"Path to a new template file"`
	Version    func()   `short:"v" long:"version" description:"Display the version of cfn-clone"`
}

func paramsFromCli(attribs []string) map[string]string {
	parameters := map[string]string{}
	for _, a := range attribs {
		p := strings.SplitN(a, "=", 2)
		parameters[p[0]] = p[1]
	}

	return parameters
}

func parseCliArgs() *options {
	opts := &options{}

	opts.Version = func() {
		fmt.Println(version)
		os.Exit(0)
	}

	parser := flags.NewParser(opts, flags.Default)

	args, err := parser.Parse()
	if err != nil {
		helpDisplayed := false

		for _, i := range args {
			if i == "-h" || i == "--help" {
				helpDisplayed = true
				break
			}
		}

		if !helpDisplayed {
			parser.WriteHelp(os.Stderr)
		}
		os.Exit(1)
	}

	if err = validateCliExists("aws"); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err = validateCliParameters(opts.Attributes); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err = validateTemplateExists(opts.Template); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err = validateSourceStackExists(opts.SourceName); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	return opts
}
