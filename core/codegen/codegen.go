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
	sub	rsp, 16
	mov	[rbp-4], edi
	mov	eax, [rbp-4]
	mov	esi, eax
	lea	rdi, [rel LC0]
	mov	eax, 0
	call	printf
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
	; Returning with code 69. 60 is code for sys_exit and 69 is return value
	mov rax, 60
	mov rdi, 69
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

func GetAllocateInstruction(a int) (string, string) {
	allc := registerService.GetNewRegister()
	instruction := fmt.Sprintf("\tmov %s, %d\n", allc, a)
	return instruction, allc
}
