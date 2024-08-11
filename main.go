package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"regexp"

	"golang.org/x/tools/go/packages"
)

var rgx = `internal|vendor`

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		log.Println(err)
	}
	var b bytes.Buffer
	b.WriteString("package main\n")
	b.WriteString("import (\n")
	for _, p := range pkgs {
		if !regexp.MustCompile(rgx).MatchString(p.PkgPath) {
			b.WriteString(`_ "`)
			b.WriteString(p.PkgPath)
			b.WriteString(`"`)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n)\n")
	b.WriteString("func main() {}")
	fmt, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	tmp, err := os.CreateTemp("", "*.go")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmp.Write(fmt); err != nil {
		log.Fatal(err)
	}
	log.Println(tmp.Name())
}
