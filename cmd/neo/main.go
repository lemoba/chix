package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		newCmd()
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("neo usage:")
	fmt.Println("  neo new <project_name>")
}

func newCmd() {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	module := fs.String("module", "", "go module name, default <name>")
	_ = fs.Parse(os.Args[2:])
	args := fs.Args()
	if len(args) < 1 {
		log.Fatal("project name required")
	}
	name := args[0]
	if *module == "" {
		*module = name
	}

	root := name
	must(os.MkdirAll(filepath.Join(root, "cmd", name), 0o755))
	must(os.MkdirAll(filepath.Join(root, "internal", "handler"), 0o755))
	must(os.MkdirAll(filepath.Join(root, "pkg", "middleware"), 0o755))

	writeFile(filepath.Join(root, "go.mod"), "module "+*module+"\n\ngo 1.21\n\nrequire github.com/lemoba/chix v0.0.0-00010101000000-000000000000\n")
	writeFile(filepath.Join(root, "cmd", name, "main.go"), mainTemplate(name))
	writeFile(filepath.Join(root, "internal", "handler", "health.go"), handlerTemplate())

	fmt.Println("project created:", root)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeFile(path string, content string) {
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	f, err := os.Create(path)
	must(err)
	defer f.Close()
	_, err = f.WriteString(content)
	must(err)
}

func mainTemplate(name string) string {
	return `package main

import (
	"github.com/lemoba/chix"
)

func main() {
	r := chix.New()
	r.Get("/health", func(c chix.Context) { c.Success("ok") })
	r.Run(":3000")
}
`
}

func handlerTemplate() string {
	return `package handler

import (
	"github.com/lemoba/chix"
)

func Health(c chix.Context) {
	c.Success("ok")
}
`
}
