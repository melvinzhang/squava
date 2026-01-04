// +build amd64

#include "textflag.h"

// Lane order: 0: Horizontal (s=1), 1: Vertical (s=8), 2: Diagonal (s=9), 3: Anti-diagonal (s=7)

DATA ·shifts1+0(SB)/8, $1
DATA ·shifts1+8(SB)/8, $8
DATA ·shifts1+16(SB)/8, $9
DATA ·shifts1+24(SB)/8, $7
GLOBL ·shifts1(SB), RODATA, $32

DATA ·shifts2+0(SB)/8, $2
DATA ·shifts2+8(SB)/8, $16
DATA ·shifts2+16(SB)/8, $18
DATA ·shifts2+24(SB)/8, $14
GLOBL ·shifts2(SB), RODATA, $32

DATA ·shifts3+0(SB)/8, $3
DATA ·shifts3+8(SB)/8, $24
DATA ·shifts3+16(SB)/8, $27
DATA ·shifts3+24(SB)/8, $21
GLOBL ·shifts3(SB), RODATA, $32

// MaskR1: MaskNotH, ALL, MaskNotH, MaskNotA
DATA ·maskR1+0(SB)/8, $0x7F7F7F7F7F7F7F7F
DATA ·maskR1+8(SB)/8, $0xFFFFFFFFFFFFFFFF
DATA ·maskR1+16(SB)/8, $0x7F7F7F7F7F7F7F7F
DATA ·maskR1+24(SB)/8, $0xFEFEFEFEFEFEFEFE
GLOBL ·maskR1(SB), RODATA, $32

// MaskL1: MaskNotA, ALL, MaskNotA, MaskNotH
DATA ·maskL1+0(SB)/8, $0xFEFEFEFEFEFEFEFE
DATA ·maskL1+8(SB)/8, $0xFFFFFFFFFFFFFFFF
DATA ·maskL1+16(SB)/8, $0xFEFEFEFEFEFEFEFE
DATA ·maskL1+24(SB)/8, $0x7F7F7F7F7F7F7F7F
GLOBL ·maskL1(SB), RODATA, $32

// MaskR2: MaskNotGH, ALL, MaskNotGH, MaskNotAB
DATA ·maskR2+0(SB)/8, $0x3F3F3F3F3F3F3F3F
DATA ·maskR2+8(SB)/8, $0xFFFFFFFFFFFFFFFF
DATA ·maskR2+16(SB)/8, $0x3F3F3F3F3F3F3F3F
DATA ·maskR2+24(SB)/8, $0xFCFCFCFCFCFCFCFC
GLOBL ·maskR2(SB), RODATA, $32

// MaskL2: MaskNotAB, ALL, MaskNotAB, MaskNotGH
DATA ·maskL2+0(SB)/8, $0xFCFCFCFCFCFCFCFC
DATA ·maskL2+8(SB)/8, $0xFFFFFFFFFFFFFFFF
DATA ·maskL2+16(SB)/8, $0xFCFCFCFCFCFCFCFC
DATA ·maskL2+24(SB)/8, $0x3F3F3F3F3F3F3F3F
GLOBL ·maskL2(SB), RODATA, $32

// MaskR3: MaskNotFGH, ALL, MaskNotFGH, MaskNotABC
DATA ·maskR3+0(SB)/8, $0x1F1F1F1F1F1F1F1F
DATA ·maskR3+8(SB)/8, $0xFFFFFFFFFFFFFFFF
DATA ·maskR3+16(SB)/8, $0x1F1F1F1F1F1F1F1F
DATA ·maskR3+24(SB)/8, $0xF8F8F8F8F8F8F8F8
GLOBL ·maskR3(SB), RODATA, $32

// MaskL3: MaskNotABC, ALL, MaskNotABC, MaskNotFGH
DATA ·maskL3+0(SB)/8, $0xF8F8F8F8F8F8F8F8
DATA ·maskL3+8(SB)/8, $0xFFFFFFFFFFFFFFFF
DATA ·maskL3+16(SB)/8, $0xF8F8F8F8F8F8F8F8
DATA ·maskL3+24(SB)/8, $0x1F1F1F1F1F1F1F1F
GLOBL ·maskL3(SB), RODATA, $32

// func getWinsAndLossesAVX2(b, e uint64) (w, l uint64)
TEXT ·getWinsAndLossesAVX2(SB), NOSPLIT, $0-32
    MOVQ b+0(FP), AX
    VMOVQ AX, X0
    VPBROADCASTQ X0, Y0      // Y0 = [b, b, b, b]
    MOVQ e+8(FP), BX
    VMOVQ BX, X1
    VPBROADCASTQ X1, Y1      // Y1 = [e, e, e, e]
    
    // Shifts 1
    VMOVDQU ·shifts1(SB), Y2
    VPSRLVQ Y2, Y0, Y3       // Y3 = b >> s
    VPAND ·maskR1(SB), Y3, Y3 // Y3 = r1
    VPSLLVQ Y2, Y0, Y4       // Y4 = b << s
    VPAND ·maskL1(SB), Y4, Y4 // Y4 = l1
    
    // Shifts 2
    VMOVDQU ·shifts2(SB), Y2
    VPSRLVQ Y2, Y0, Y5       // Y5 = b >> 2s
    VPAND ·maskR2(SB), Y5, Y5 // Y5 = r2
    VPSLLVQ Y2, Y0, Y6       // Y6 = b << 2s
    VPAND ·maskL2(SB), Y6, Y6 // Y6 = l2
    
    // Shifts 3
    VMOVDQU ·shifts3(SB), Y2
    VPSRLVQ Y2, Y0, Y7       // Y7 = b >> 3s
    VPAND ·maskR3(SB), Y7, Y7 // Y7 = r3
    VPSLLVQ Y2, Y0, Y8       // Y8 = b << 3s
    VPAND ·maskL3(SB), Y8, Y8 // Y8 = l3
    
    // Calculate L lanes: e & (r1&r2 | r1&l1 | l1&l2)
    VPAND Y3, Y5, Y9         // r1 & r2
    VPAND Y3, Y4, Y10        // r1 & l1
    VPOR Y9, Y10, Y9         // r1&r2 | r1&l1
    VPAND Y4, Y6, Y10        // l1 & l2
    VPOR Y10, Y9, Y9         // r1&r2 | r1&l1 | l1&l2
    VPAND Y1, Y9, Y9         // Y9 = L lanes
    
    // Calculate W lanes: e & (r1&r2&(r3|l1) | l1&l2&(r1|l3))
    VPOR Y7, Y4, Y10         // r3 | l1
    VPAND Y3, Y5, Y11        // r1 & r2
    VPAND Y11, Y10, Y10      // r1&r2&(r3|l1)
    
    VPOR Y3, Y8, Y11         // r1 | l3
    VPAND Y4, Y6, Y12        // l1 & l2
    VPAND Y12, Y11, Y11      // l1&l2&(r1|l3)
    
    VPOR Y10, Y11, Y10       // r1&r2&(r3|l1) | l1&l2&(r1|l3)
    VPAND Y1, Y10, Y10       // Y10 = W lanes
    
    // Horizontal OR for W
    VEXTRACTI128 $1, Y10, X11 // X11 = upper 128 bits of Y10
    VPOR X11, X10, X10        // X10 = [W0|W2, W1|W3]
    VPSHUFD $0x4E, X10, X11   // swap 64-bit halves
    VPOR X11, X10, X10        // X10 = [W0|W1|W2|W3, ...]
    MOVQ X10, AX
    MOVQ AX, w+16(FP)
    
    // Horizontal OR for L
    VEXTRACTI128 $1, Y9, X11
    VPOR X11, X9, X9
    VPSHUFD $0x4E, X9, X11
    VPOR X11, X9, X9
    MOVQ X9, AX
    MOVQ AX, l+24(FP)
    
    VZEROUPPER
    RET
