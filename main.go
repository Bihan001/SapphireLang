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

	llvmIR := ""

	llvmIR += codegen.GenPrefixCode()

	llvmIR += output

	llvmIR += codegen.GenPostfixCode()
	//fmt.Println(llvmIR)

	asmCodeOutputPath := "resources/output.asm"
	binaryOutputPath := "resources/output"

	writeAsmCodeToFile(llvmIR, asmCodeOutputPath)

	buildAndLinkAssemblyCode(asmCodeOutputPath, binaryOutputPath)

}

func writeAsmCodeToFile(ir string, outputPath string) {
	outErr := os.WriteFile(outputPath, []byte(ir), 0644)

	if outErr != nil {
		panic(outErr)
	}
}

func buildAndLinkAssemblyCode(inputPath string, outputPath string) {
	/*
		apadmin@apdc1n-dev-mayank-vm:~/nasm$ vi a.asm
		apadmin@apdc1n-dev-mayank-vm:~/nasm$ nasm -felf64 a.asm -o a.o
		apadmin@apdc1n-dev-mayank-vm:~/nasm$ gcc -no-pie a.o -o a
		apadmin@apdc1n-dev-mayank-vm:~/nasm$ ./a
		16
		apadmin@apdc1n-dev-mayank-vm:~/nasm$ $?
		69: command not found
		apadmin@apdc1n-dev-mayank-vm:~/nasm$
	*/
	objectCodeOutputPath := "output.o"
	cmd := exec.Command("nasm", "-felf64", inputPath, "-o", objectCodeOutputPath)
	err := cmd.Run()
	if err != nil {
		panic("Failed to compile .asm to .o")
	}

	cmd = exec.Command("gcc", "-no-pie", objectCodeOutputPath, "-o", outputPath)
	err = cmd.Run()
	if err != nil {
		panic("Failed to link .o to executable")
	}

	err = os.Remove("output.o")
	if err != nil {
		panic("Failed to remove object code file")
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
