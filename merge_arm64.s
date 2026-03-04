//go:build arm64

#include "textflag.h"

// func mergeRegistersASM(dst, src *byte, length int)
// Merges src into dst by taking element-wise maximum.
// Uses ARM64 NEON instructions to process 16 bytes at a time.
//
// Arguments:
//   R0 = dst pointer
//   R1 = src pointer  
//   R2 = length (must be multiple of 16, minimum 16)
TEXT ·mergeRegistersASM(SB), NOSPLIT, $0-24
    MOVD dst+0(FP), R0      // dst pointer
    MOVD src+8(FP), R1      // src pointer
    MOVD length+16(FP), R2  // length

    // Check if length is 0
    CBZ  R2, done

    // Process 64 bytes at a time (4 x 16-byte vectors) for better throughput
loop64:
    CMP  $64, R2
    BLT  loop16

    // Load 64 bytes from dst
    VLD1.P 64(R0), [V0.B16, V1.B16, V2.B16, V3.B16]
    
    // Load 64 bytes from src
    VLD1.P 64(R1), [V4.B16, V5.B16, V6.B16, V7.B16]
    
    // Compute max for each vector pair (unsigned max for bytes)
    VUMAX V0.B16, V4.B16, V0.B16
    VUMAX V1.B16, V5.B16, V1.B16
    VUMAX V2.B16, V6.B16, V2.B16
    VUMAX V3.B16, V7.B16, V3.B16
    
    // Store results back to dst (need to go back 64 bytes since we post-incremented)
    SUB  $64, R0, R3
    VST1 [V0.B16, V1.B16, V2.B16, V3.B16], (R3)
    
    SUB  $64, R2
    B    loop64

    // Process remaining 16 bytes at a time
loop16:
    CMP  $16, R2
    BLT  done

    // Load 16 bytes from dst
    VLD1 (R0), [V0.B16]
    
    // Load 16 bytes from src
    VLD1 (R1), [V1.B16]
    
    // Compute unsigned max
    VUMAX V0.B16, V1.B16, V0.B16
    
    // Store result back to dst
    VST1 [V0.B16], (R0)
    
    ADD  $16, R0
    ADD  $16, R1
    SUB  $16, R2
    B    loop16

done:
    RET
