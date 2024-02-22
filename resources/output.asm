	
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
	mov r9, 50
	mov [a], r9
	common b 8:8
	common c 8:8
	mov r14, 20
	mov r8, 30
	cmp r14, r8
	setl r14b
	and r14, 255
	mov [b], r14
	mov r8, [a]
	mov r14, 40
	sub r8, r14
	mov [c], r8
	mov r8, [a]
	mov r13, [b]
	mov r11, 3
	imul r13, r11
	add r8, r13
	mov r11, 10
	mov r9, [c]
	mov rax, r11
	cqo
	idiv r9
	mov r11, rax
	sub r8, r11
	mov rdi, r8
	call printint
	mov r13, 20
	mov r10, 30
	cmp r13, r10
	setl r13b
	and r13, 255
	mov rdi, r13
	call printint
	mov r10, 30
	mov r15, 20
	cmp r10, r15
	setl r10b
	and r10, 255
	mov rdi, r10
	call printint
	mov r11, 10
	mov r8, 10
	cmp r11, r8
	setle r11b
	and r11, 255
	mov rdi, r11
	call printint
	mov r12, 10
	mov r14, 10
	cmp r12, r14
	setge r12b
	and r12, 255
	mov rdi, r12
	call printint
	mov r10, 10
	mov r14, 10
	cmp r10, r14
	setl r10b
	and r10, 255
	mov rdi, r10
	call printint
	mov r9, 3
	mov r13, 2
	cmp r9, r13
	setg r9b
	and r9, 255
	mov rdi, r9
	call printint
	mov r15, 2
	mov r8, 2
	cmp r15, r8
	setge r15b
	and r15, 255
	mov rdi, r15
	call printint
	mov r8, 3
	mov r14, 2
	cmp r8, r14
	sete r8b
	and r8, 255
	mov rdi, r8
	call printint
	mov r9, 3
	mov r15, 3
	cmp r9, r15
	sete r9b
	and r9, 255
	mov rdi, r9
	call printint
	mov r8, 3
	mov r15, 2
	cmp r8, r15
	setne r8b
	and r8, 255
	mov rdi, r8
	call printint
	mov r9, 2
	mov r15, 2
	cmp r9, r15
	setne r9b
	and r9, 255
	mov rdi, r9
	call printint

	; Returning with code 0. 60 is code for sys_exit and 0 is return value
	mov rax, 60
	mov rdi, 0
	syscall
