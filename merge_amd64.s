//go:build amd64

#include "textflag.h"

// func mergeRegistersASM(dst, src *byte, length int)
// Merges src into dst by taking element-wise maximum.
// Uses AVX2 instructions to process 128 bytes at a time.
// Length must be a multiple of 128 (which is guaranteed for precision >= 7).
//
// Arguments:
//   dst+0(FP)    = dst pointer
//   src+8(FP)    = src pointer  
//   length+16(FP) = length (must be multiple of 128)
TEXT ·mergeRegistersASM(SB), NOSPLIT, $0-24
    MOVQ dst+0(FP), DI      // dst pointer
    MOVQ src+8(FP), SI      // src pointer
    MOVQ length+16(FP), CX  // length

    // Check if length is 0
    TESTQ CX, CX
    JZ    done

    // Process 128 bytes at a time (4 x 32-byte YMM registers)
loop128:
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
    JNZ   loop128

    // Clear upper bits of YMM registers to avoid SSE/AVX transition penalties
    VZEROUPPER

done:
    RET
