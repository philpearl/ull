//go:build amd64

#include "textflag.h"

// func mergeRegistersAVX2(dst, src *byte, length int)
// Merges src into dst by taking element-wise maximum.
// Uses AVX2 instructions to process 128 bytes at a time.
// Length must be a multiple of 128 (which is guaranteed for precision >= 7).
//
// Arguments:
//   dst+0(FP)    = dst pointer
//   src+8(FP)    = src pointer  
//   length+16(FP) = length (must be multiple of 128)
TEXT ·mergeRegistersAVX2(SB), NOSPLIT, $0-24
    MOVQ dst+0(FP), DI      // dst pointer
    MOVQ src+8(FP), SI      // src pointer
    MOVQ length+16(FP), CX  // length

    // Check if length is 0
    TESTQ CX, CX
    JZ    avx2_done

    // Process 128 bytes at a time (4 x 32-byte YMM registers)
avx2_loop128:
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
    JNZ   avx2_loop128

    // Clear upper bits of YMM registers to avoid SSE/AVX transition penalties
    VZEROUPPER

avx2_done:
    RET

// func mergeRegistersAVX512(dst, src *byte, length int)
// Merges src into dst by taking element-wise maximum.
// Uses AVX-512 instructions to process 256 bytes at a time.
// Length must be a multiple of 256 (which is guaranteed for precision >= 8).
//
// Arguments:
//   dst+0(FP)    = dst pointer
//   src+8(FP)    = src pointer  
//   length+16(FP) = length (must be multiple of 256)
TEXT ·mergeRegistersAVX512(SB), NOSPLIT, $0-24
    MOVQ dst+0(FP), DI      // dst pointer
    MOVQ src+8(FP), SI      // src pointer
    MOVQ length+16(FP), CX  // length

    // Check if length is 0
    TESTQ CX, CX
    JZ    avx512_done

    // Process 256 bytes at a time (4 x 64-byte ZMM registers)
avx512_loop256:
    // Load 256 bytes from dst into ZMM0-ZMM3
    VMOVDQU64 (DI), Z0
    VMOVDQU64 64(DI), Z1
    VMOVDQU64 128(DI), Z2
    VMOVDQU64 192(DI), Z3
    
    // Load 256 bytes from src into ZMM4-ZMM7
    VMOVDQU64 (SI), Z4
    VMOVDQU64 64(SI), Z5
    VMOVDQU64 128(SI), Z6
    VMOVDQU64 192(SI), Z7
    
    // Compute max for each vector pair (unsigned max for bytes)
    // VPMAXUB with ZMM registers requires AVX-512BW
    VPMAXUB Z4, Z0, Z0
    VPMAXUB Z5, Z1, Z1
    VPMAXUB Z6, Z2, Z2
    VPMAXUB Z7, Z3, Z3
    
    // Store results back to dst
    VMOVDQU64 Z0, (DI)
    VMOVDQU64 Z1, 64(DI)
    VMOVDQU64 Z2, 128(DI)
    VMOVDQU64 Z3, 192(DI)
    
    ADDQ  $256, DI
    ADDQ  $256, SI
    SUBQ  $256, CX
    JNZ   avx512_loop256

    VZEROUPPER

avx512_done:
    RET
