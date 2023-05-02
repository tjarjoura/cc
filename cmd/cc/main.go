package cc

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tjarjoura/cc/pkg/compiler"
	"github.com/tjarjoura/cc/pkg/lexer"
	"github.com/tjarjoura/cc/pkg/parser"
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
	for _, err := range c.Errors() {
		log.Printf("compiler error: %s\n", err.String())
		ret = false
	}

	return ret
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("No input files given.")
		os.Exit(1)
	}

	for _, inputFile := range flag.Args() {
		inp, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		p := parser.New(lexer.New(string(inp)))
		tUnit := p.Parse()

		if checkParserErrors(p) {
			os.Exit(1)
		}

		c := compiler.New(tUnit)
		c.Compile()

		if checkCompilerErrors(c) {
			os.Exit(1)
		}

		// TODO write to file
		//asmFilename := strings.Replace(inputFile, ".c", ".asm", -1)
		c.WriteAssembly(os.Stdout)
	}
}
