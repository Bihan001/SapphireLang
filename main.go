package main

import (
	"SLang/core/codegen"
	"SLang/core/parser"
	"SLang/core/scanner"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Welcome to Sapphire Lang!")

	args := os.Args

	if len(args) < 2 {
		panic("No input file")
	}

	inputFilePath := os.Args[1]
	inputFileBytes, inErr := os.ReadFile(inputFilePath)

	if inErr != nil {
		panic(inErr)
	}

	inputFileContents := string(inputFileBytes)

	tokenizer := scanner.NewTokenizer(&inputFileContents)

	tokens := tokenizer.Tokenize()

	syntaxParser := parser.NewParser()

	output := syntaxParser.Parse(tokens)

	llvmIR := codegen.GetGlobals()

	llvmIR += "define i32 @main() {\nentry:\n"

	llvmIR += output

	llvmIR += "ret i32 0\n"

	llvmIR += "}"

	//fmt.Println(llvmIR)

	irCodeOutputPath := "resources/output.ll"
	binaryOutputPath := "resources/output"

	writeIRCodeToFile(llvmIR, irCodeOutputPath)

	buildAndLinkIRCode(irCodeOutputPath, binaryOutputPath)

}

func writeIRCodeToFile(ir string, outputPath string) {
	outErr := os.WriteFile(outputPath, []byte(ir), 0644)

	if outErr != nil {
		panic(outErr)
	}
}

func buildAndLinkIRCode(inputPath string, outputPath string) {
	objectCodeOutputPath := "output.o"
	cmd := exec.Command("/opt/homebrew/Cellar/llvm/17.0.6/bin/llc", inputPath, "-filetype=obj", "-o", objectCodeOutputPath)
	err := cmd.Run()
	if err != nil {
		panic("Failed to compile .ll to .o")
	}

	cmd = exec.Command("clang", objectCodeOutputPath, "-o", outputPath)
	err = cmd.Run()
	if err != nil {
		panic("Failed to link .o to executable")
	}

	err = os.Remove("output.o")
	if err != nil {
		panic("Failed to remove object code file")
	}
}
