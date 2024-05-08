	
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
	common a 8:8
	mov r13, 50
	mov [a], r13
	common b 8:8
	common c 8:8
	mov r8, 20
	mov r12, 30
	cmp r8, r12
	setl r8b
	and r8, 255
	mov [b], r8
	mov r14, [a]
	mov r15, 40
	sub r14, r15
	mov [c], r14
	mov r13, [a]
	mov r8, [b]
	mov r11, 3
	imul r8, r11
	add r13, r8
	mov r12, 10
	mov r10, [c]
	mov rax, r12
	cqo
	idiv r10
	mov r12, rax
	sub r13, r12
	mov rdi, r13
	call printint
	mov r11, 20
	mov r9, 30
	cmp r11, r9
	setl r11b
	and r11, 255
	mov rdi, r11
	call printint
	mov r15, 30
	mov r9, 20
	cmp r15, r9
	setl r15b
	and r15, 255
	mov rdi, r15
	call printint
	mov r8, 10
	mov r9, 10
	cmp r8, r9
	setle r8b
	and r8, 255
	mov rdi, r8
	call printint
	mov r14, 10
	mov r15, 10
	cmp r14, r15
	setge r14b
	and r14, 255
	mov rdi, r14
	call printint
	mov r15, 10
	mov r11, 10
	cmp r15, r11
	setl r15b
	and r15, 255
	mov rdi, r15
	call printint
	mov r9, 3
	mov r11, 2
	cmp r9, r11
	setg r9b
	and r9, 255
	mov rdi, r9
	call printint
	mov r15, 2
	mov r11, 2
	cmp r15, r11
	setge r15b
	and r15, 255
	mov rdi, r15
	call printint
	mov r12, 3
	mov r10, 2
	cmp r12, r10
	sete r12b
	and r12, 255
	mov rdi, r12
	call printint
	mov r8, 3
	mov r11, 3
	cmp r8, r11
	sete r8b
	and r8, 255
	mov rdi, r8
	call printint
	mov r15, 3
	mov r8, 2
	cmp r15, r8
	setne r15b
	and r15, 255
	mov rdi, r15
	call printint
	mov r10, 2
	mov r13, 2
	cmp r10, r13
	setne r10b
	and r10, 255
	mov rdi, r10
	call printint

	; Returning with code 0. 60 is code for sys_exit and 0 is return value
	mov rax, 60
	mov rdi, 0
	syscall
