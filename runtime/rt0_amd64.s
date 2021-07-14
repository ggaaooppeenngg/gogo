/*
 *  go-routine
 *  https://github.com/teh-cmc/go-internals/blob/master/chapter1_assembly_primer/README.md
 */

#define NOSPLIT 4
#define SYS_mmap                9

TEXT ·gogogo(SB), NOSPLIT, $0-8
	MOVQ	8(SP), AX		// gogobuf
	MOVQ	0(AX), SP		// restore SP
	MOVQ	8(AX), AX
	MOVQ	AX, 0(SP)		// put PC on the stack
	MOVL	$1, AX			// return 1
	RET


TEXT ·gosave(SB), NOSPLIT, $0-16
	MOVQ	8(SP), AX		// gogobuf
	MOVQ	SP, 0(AX)		// save SP
	MOVQ	0(SP), BX
	MOVQ	BX, 8(AX)		// save PC
	MOVB	$0, 16(SP)		// return false
	RET

TEXT ·mmap(SB),NOSPLIT,$0
        MOVQ    addr+0(FP), DI
        MOVQ    n+8(FP), SI
        MOVL    prot+16(FP), DX
        MOVL    flags+20(FP), R10
        MOVL    fd+24(FP), R8
        MOVL    off+28(FP), R9

        MOVL    $SYS_mmap, AX
        SYSCALL
        CMPQ    AX, $0xfffffffffffff001
        JLS     ok
        NOTQ    AX
        INCQ    AX
        MOVQ    $0, p+32(FP)
        MOVQ    AX, err+40(FP)
        RET
ok:
        MOVQ    AX, p+32(FP)
        MOVQ    $0, err+40(FP)
        RET
