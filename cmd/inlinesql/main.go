package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

func main() {
	var (
		in, out, pkg string
	)

	flag.StringVar(&in, "f", "", "SQL file to parse")
	flag.StringVar(&out, "o", "", "File to output Go code to")
	flag.StringVar(&pkg, "p", "main", "Package name to use")
	flag.Parse()

	if in == "" {
		log.Fatalln("Please define a SQL file to parse")
	}

	sqlFile, err := ioutil.ReadFile(in)
	if err != nil {
		log.Fatalln(err)
	}

	// Good enough, small .sql files anyways...
	removeComments := regexp.MustCompile("(?s)--.*?\n|/\\*.*?\\*/")
	queries := strings.Split(strings.Replace(strings.Replace(string(removeComments.ReplaceAll(sqlFile, nil)), "\n", "", -1), "\t", " ", -1), ";")

	outFile, err := os.Create(out)
	if err != nil {
		log.Fatalln(err)
	}
	defer outFile.Close()

	pkgTmpl.Execute(outFile, struct {
		Timestamp time.Time
		Package   string
		Queries   []string
	}{
		Timestamp: time.Now(),
		Package:   pkg,
		Queries:   queries,
	})
}

var pkgTmpl = template.Must(template.New("").Parse(
	`// This file is generated automatically by inlinesql at {{ .Timestamp }}.
package {{ .Package }}

func getQueries() []string {
	return []string{
		{{- range .Queries }}
		{{ printf "%q" . }},
		{{- end }}
	}
}`))
