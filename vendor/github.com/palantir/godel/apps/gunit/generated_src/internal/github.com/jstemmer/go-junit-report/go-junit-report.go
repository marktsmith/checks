package amalgomated

import (
	"github.com/palantir/godel/apps/gunit/generated_src/internal/github.com/jstemmer/go-junit-report/amalgomated_flag"
	"fmt"
	"os"

	"github.com/palantir/godel/apps/gunit/generated_src/internal/github.com/jstemmer/go-junit-report/parser"
)

var (
	noXMLHeader	bool
	packageName	string
	goVersionFlag	string
	setExitCode	bool
)

func init() {
	flag.BoolVar(&noXMLHeader, "no-xml-header", false, "do not print xml header")
	flag.StringVar(&packageName, "package-name", "", "specify a package name (compiled test have no package name in output)")
	flag.StringVar(&goVersionFlag, "go-version", "", "specify the value to use for the go.version property in the generated XML")
	flag.BoolVar(&setExitCode, "set-exit-code", false, "set exit code to 1 if tests failed")
}

func AmalgomatedMain() {
	flag.Parse()

	// Read input
	report, err := parser.Parse(os.Stdin, packageName)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		os.Exit(1)
	}

	// Write xml
	err = JUnitReportXML(report, noXMLHeader, goVersionFlag, os.Stdout)
	if err != nil {
		fmt.Printf("Error writing XML: %s\n", err)
		os.Exit(1)
	}

	if setExitCode && report.Failures() > 0 {
		os.Exit(1)
	}
}
