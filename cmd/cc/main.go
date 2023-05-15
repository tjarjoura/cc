package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/tjarjoura/cc/pkg/compiler"
	"github.com/tjarjoura/cc/pkg/lexer"
	"github.com/tjarjoura/cc/pkg/parser"
)

var (
	assembleOnly = flag.Bool("asm", false,
		"If set, will only generate assembly code from each source file")
)

func checkParserErrors(p *parser.Parser) bool {
	ret := true
	for _, err := range p.Errors() {
		log.Printf("parser error: %s\n", err.String())
		ret = false
	}

	return ret
}

func checkCompilerErrors(c *compiler.Compiler) bool {
	ret := true
	errorMap := c.Errors()
	for _, err := range errorMap["global"] {
		log.Printf("compiler error: %s\n", err.String())
		ret = false
	}

	delete(errorMap, "global")
	for name, errors := range errorMap {
		if len(errors) > 0 {
			log.Printf("in function %s:", name)
		}

		for _, err := range errors {
			log.Printf("compiler error: %s", err.String())
			ret = false
		}
	}

	return ret
}

/* Take *.c source files and produce *.asm files */
func compile(outputFile string, sourceFiles ...string) ([]string, error) {
	if len(sourceFiles) > 1 && outputFile != "" {
		return []string{}, fmt.Errorf(
			"cannot specify output file when generating assembly for multiple source files.")
	}

	asmFiles := []string{}
	for _, inputFile := range sourceFiles {
		inp, err := os.ReadFile(inputFile)
		if err != nil {
			return asmFiles,
				fmt.Errorf("error reading %s: %s", inputFile, err)
		}

		p := parser.New(lexer.New(string(inp)))
		tUnit := p.Parse()

		if !checkParserErrors(p) {
			return asmFiles,
				fmt.Errorf("got parser errors for %s", inputFile)
		}

		c := compiler.New(tUnit)
		c.Compile()

		if !checkCompilerErrors(c) {
			return asmFiles,
				fmt.Errorf("got compiler errors for %s", inputFile)
		}

		if outputFile == "" {
			outputFile = strings.ReplaceAll(inputFile, ".c", ".asm")
		}

		f, err := os.OpenFile(outputFile,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return asmFiles, fmt.Errorf(
				"error opening/creating %s: %s", outputFile, err)
		}
		defer f.Close()

		if err := c.WriteAssembly(f); err != nil {
			return asmFiles, fmt.Errorf(
				"error generating assembly: %s", err)
		}

		asmFiles = append(asmFiles, outputFile)
	}

	if outputFile != "" {
		return []string{outputFile}, nil
	}

	return asmFiles, nil
}

/* Take *.asm source files and produces *.o relocatable object files. note: Just calls NASM currently */
func assemble(outputFile string, asmFiles ...string) ([]string, error) {
	if len(asmFiles) > 1 && outputFile != "" {
		return []string{}, fmt.Errorf(
			"cannot specify output file when compiling multiple assembly files.")
	}

	objFiles := []string{}
	for _, asmFile := range asmFiles {
		args := append([]string{"-f", "elf64"}, asmFile)
		cmd := exec.Command("nasm", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return objFiles, fmt.Errorf(
				"got error when trying to assemble using nasm: %s\ncommand output:\n%s",
				err, string(out))
		}

		objFiles = append(objFiles,
			strings.ReplaceAll(asmFile, ".asm", ".o"))
	}

	if outputFile != "" {
		return []string{outputFile}, nil
	}

	return objFiles, nil
}

/* Take *.o source files and links them into a final binary. note: Just calls gcc/ld currently */
func link(outputFile string, objFiles ...string) error {
	if outputFile == "" {
		outputFile = "a.out"
	}

	args := append([]string{"-o", outputFile}, objFiles...)
	cmd := exec.Command("gcc", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"got error when trying to link using gcc: %s", err)
	}

	return nil
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("No input files given.")
		os.Exit(1)
	}

	var asmFiles, objFiles []string
	defer func() {
		for _, f := range append(asmFiles, objFiles...) {
			err := os.Remove(f)
			if err != nil {
				log.Printf("error removing %s: %s", f, err)
			}
		}
	}()

	asmFiles, err := compile("", flag.Args()...)
	if err != nil {
		log.Print(err)
		return
	}

	objFiles, err = assemble("", asmFiles...)
	if err != nil {
		log.Print(err)
		return
	}

	err = link("", objFiles...)
	if err != nil {
		log.Print(err)
		return
	}
}
