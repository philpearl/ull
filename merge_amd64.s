//go:build amd64

#include "textflag.h"

// func mergeRegistersASM(dst, src *byte, length int)
// Merges src into dst by taking element-wise maximum.
// Uses AVX2 instructions to process 32 bytes at a time.
//
// Arguments:
//   dst+0(FP)    = dst pointer
//   src+8(FP)    = src pointer  
//   length+16(FP) = length
TEXT ·mergeRegistersASM(SB), NOSPLIT, $0-24
    MOVQ dst+0(FP), DI      // dst pointer
    MOVQ src+8(FP), SI      // src pointer
    MOVQ length+16(FP), CX  // length

    // Check if length is 0
    TESTQ CX, CX
    JZ    done

    // Process 128 bytes at a time (4 x 32-byte YMM registers)
loop128:
    CMPQ  CX, $128
    JL    loop32

    // Load 128 bytes from dst into YMM0-YMM3
    VMOVDQU (DI), Y0
    VMOVDQU 32(DI), Y1
    VMOVDQU 64(DI), Y2
    VMOVDQU 96(DI), Y3
    
    // Load 128 bytes from src into YMM4-YMM7
    VMOVDQU (SI), Y4
    VMOVDQU 32(SI), Y5
    VMOVDQU 64(SI), Y6
    VMOVDQU 96(SI), Y7
    
    // Compute max for each vector pair (unsigned max for bytes)
    VPMAXUB Y4, Y0, Y0
    VPMAXUB Y5, Y1, Y1
    VPMAXUB Y6, Y2, Y2
    VPMAXUB Y7, Y3, Y3
    
    // Store results back to dst
    VMOVDQU Y0, (DI)
    VMOVDQU Y1, 32(DI)
    VMOVDQU Y2, 64(DI)
    VMOVDQU Y3, 96(DI)
    
    ADDQ  $128, DI
    ADDQ  $128, SI
    SUBQ  $128, CX
    JMP   loop128

    // Process 32 bytes at a time
loop32:
    CMPQ  CX, $32
    JL    loop16

    VMOVDQU (DI), Y0
    VMOVDQU (SI), Y1
    VPMAXUB Y1, Y0, Y0
    VMOVDQU Y0, (DI)
    
    ADDQ  $32, DI
    ADDQ  $32, SI
    SUBQ  $32, CX
    JMP   loop32

    // Process 16 bytes at a time using SSE
loop16:
    CMPQ  CX, $16
    JL    loop1

    VMOVDQU (DI), X0
    VMOVDQU (SI), X1
    VPMAXUB X1, X0, X0
    VMOVDQU X0, (DI)
    
    ADDQ  $16, DI
    ADDQ  $16, SI
    SUBQ  $16, CX
    JMP   loop16

    // Process remaining bytes one at a time
loop1:
    TESTQ CX, CX
    JZ    cleanup

    MOVB  (DI), AL
    MOVB  (SI), BL
    CMPB  BL, AL
    JBE   skip
    MOVB  BL, (DI)
skip:
    INCQ  DI
    INCQ  SI
    DECQ  CX
    JMP   loop1

cleanup:
    // Clear upper bits of YMM registers to avoid SSE/AVX transition penalties
    VZEROUPPER

done:
    RET
