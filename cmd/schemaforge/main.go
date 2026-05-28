package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/why19970628/schemaforge/internal/convert"
	"github.com/why19970628/schemaforge/internal/server"
)

var version = "dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "schemaforge:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}
	switch args[0] {
	case "convert":
		return runConvert(args[1:])
	case "ui":
		return runUI(args[1:])
	case "version":
		fmt.Println(version)
		return nil
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		if isVersionArg(args[0]) {
			fmt.Println(version)
			return nil
		}
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func isVersionArg(arg string) bool {
	return arg == "-v" || arg == "--version" || arg == "-version"
}

func runConvert(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing convert mode")
	}
	mode := convert.Mode(args[0])
	fs := flag.NewFlagSet("convert "+string(mode), flag.ContinueOnError)
	inputPath := fs.String("i", "", "input file path, defaults to stdin")
	rightPath := fs.String("right", "", "right input file path for json-diff")
	format := fs.String("format", "", "output format for json-diff: structured or unified")
	outputPath := fs.String("o", "", "output file path, defaults to stdout")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	data, err := readInput(*inputPath)
	if err != nil {
		return err
	}
	req := convert.Request{Mode: mode, Input: string(data), Format: convert.Format(*format)}
	if mode == convert.ModeJSONDiff {
		if *rightPath == "" {
			return fmt.Errorf("missing --right for json-diff")
		}
		rightData, readErr := readInput(*rightPath)
		if readErr != nil {
			return readErr
		}
		req.Right = string(rightData)
	}
	resp, err := convert.Convert(req)
	if err != nil {
		return err
	}
	if *outputPath == "" {
		fmt.Print(resp.Output)
		if !strings.HasSuffix(resp.Output, "\n") {
			fmt.Println()
		}
		return nil
	}
	return os.WriteFile(*outputPath, []byte(resp.Output), 0o644)
}

func runUI(args []string) error {
	fs := flag.NewFlagSet("ui", flag.ContinueOnError)
	addr := fs.String("addr", "127.0.0.1:8989", "listen address")
	port := fs.String("port", "", "listen port, shorthand for 127.0.0.1:<port>")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *port != "" {
		*addr = "127.0.0.1:" + *port
	}
	fmt.Printf("SchemaForge UI listening on http://%s\n", *addr)
	return server.ListenAndServe(*addr)
}

func readInput(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return os.ReadFile("/dev/stdin")
	}
	return os.ReadFile(path)
}

func printUsage() {
	fmt.Println(`SchemaForge turns structured text into schemas and code.

Usage:
  schemaforge convert <mode> [-i input] [-o output]
  schemaforge convert json-diff -i left.json --right right.json [--format structured|unified] [-o output]
  schemaforge ui [--port 8989]
  schemaforge version | -v | --version | -version

Modes:
  json-go    JSON sample to Go struct
  json-diff  Compare two JSON documents
  yaml-go    YAML sample to Go struct
  xml-json   XML document to JSON
  sql-ent    MySQL CREATE TABLE to Ent schema
  sql-gorm   MySQL CREATE TABLE to GORM model
  sql-es     MySQL CREATE TABLE to Elasticsearch mapping
  sql-mongo  MySQL CREATE TABLE to MongoDB JSON Schema`)
}
