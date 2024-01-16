package codegen

import (
	"SLang/util"
	"fmt"
)

var registerService *util.RegisterService = util.GetNewRegisterService()

var globalStr = ``

func GetGlobals() string {
	return globalStr
}

func AppendToGlobals(str string) {
	globalStr += str
}

func GenPrefixCode() string {
	return `	
global	main
extern	printf
section	.text
LC0: db "%d",10,0
printint:
	push	rbp
	mov	rbp, rsp
    ; Subtracting 16 bytes i.e. 2 stack entries (1 stack entry = 8 bytes in x86_64 systems). The second one is already occupied by rsp, so first one is available for local variables.
	sub	rsp, 16
	mov	[rbp-8], rdi
	mov	rax, [rbp-8]
	lea	rdi, [rel LC0] ; First parameter for printf
	mov	rsi, rax ; Second parameter for printf
	mov	rax, 0 ; Number of vector registers used. In this case, it's 0.
	call printf
	nop
	leave
	ret

main:
	push	rbp
	mov	rbp, rsp
`
}

func GenPostfixCode() string {
	return `
	; Returning with code 0. 60 is code for sys_exit and 0 is return value
	mov rax, 60
	mov rdi, 0
	syscall
`
}

func GenPrintInt(s string) string {
	return fmt.Sprintf("\tmov rdi, %s\n\tcall printint\n", s)
}

func GetAddInstruction(a string, b string) string {
	str := fmt.Sprintf("\tadd %s, %s\n", a, b)
	registerService.FreeRegister(b)
	return str
}

func GetSubtractInstruction(a string, b string) string {
	str := fmt.Sprintf("\tsub %s, %s\n", a, b)
	registerService.FreeRegister(b)
	return str
}

func GetMultiplyInstruction(a string, b string) string {
	str := fmt.Sprintf("\timul %s, %s\n", a, b)
	registerService.FreeRegister(b)
	return str
}

func GetDivideInstruction(a string, b string) string {
	str := fmt.Sprintf("\tmov rax, %s\n", a)
	str += "\tcqo\n"
	str += fmt.Sprintf("\tidiv %s\n", b)
	str += fmt.Sprintf("\tmov %s, rax\n", a)
	registerService.FreeRegister(b)
	return str
}

func GetRegisterAllocateInstruction(a string) (string, string) {
	allc := registerService.GetNewRegister()
	instruction := fmt.Sprintf("\tmov %s, %s\n", allc, a)
	return instruction, allc
}

func GetVariableAllocateInstruction(symbol string) (string, string) {
	str := fmt.Sprintf("\tcommon %s 8:8\n", symbol)
	return str, GetVariableInstructionFromSymbol(symbol)
}

func GetVariableInstructionFromSymbol(symbol string) string {
	return fmt.Sprintf("[%s]", symbol)
}

func GetVariableAssignInstruction(left string, right string) string {
	str := fmt.Sprintf("\tmov %s, %s\n", left, right)

	// Free only if its a register and not a variable
	if len(right) < 2 || (right[0] != '[' && right[len(right)-1] != ']') {
		registerService.FreeRegister(right)
	}

	return str
}
