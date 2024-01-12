// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter_test

// This is a subset of iter_test.go designed to simplify the search for some performance problems.
// None of this code uses either generics or function range iteration, for purposes of simplifying
// diagnosis (and also allows a study back to at least 1.15 to see how how it performs with older
// compilers).

// pkg: github.com/dr2chase/iterbench/tiny
// cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//                     │  tiny.log   │
//                     │   sec/op    │
// SliceOld-8            5.885n ± 1%  <--- the goal
// Slice-8               40.78n ± 1%  <--- didn't fully inline; SliceOf(s) == V2(OfSliceIndex(s))
// SliceCheck-8          31.18n ± 0%  <--- adding check wrapper led to full inlining (???)
// SliceOf-8             22.53n ± 4%  <--- full inlining, but one inline mark and i Addrtaken
// SliceOfCheck-8        92.08n ± 3%  <--- method value closure passed to Check forces allocation
// SliceInlineCallEE-8   29.80n ± 2%  <--- Slice, but callee has V2 inlined.
// SliceInlineCallER-8   27.51n ± 3%  <--- Slice, but inline SliceOf.
//
// SliceInlineCallER and SliceInlineCallEE have same inner loop as SliceOf, plus one more inline mark.

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type Int32 int32
type Seq func(yield func(Int32) bool)
type Seq2 func(yield func(int, Int32) bool)
type Slice []Int32

func (s Slice) Of(yield func(k Int32) bool) {
	for _, x := range s {
		if !yield(x) {
			return
		}
	}
}

var gslice = Slice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

// BenchmarkSliceOld is the baseline; the goal is that versions
// written using function range iterators should apply sufficient
// inlining and optimization to achieve the same performance.
//
// Inner loop:
// v34	00050 (+48) MOVL (CX)(DI*4), R8
// v39	00051 (48) INCQ DI
// v36	00052 (+49) ADDL R8, SI
// v108	00053 (+48) CMPQ DI, DX
// b7	00054 (48) JLT 50
func BenchmarkSliceOld(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	for j := 0; j < b.N; j++ {
		for _, x := range slice {
			i += x
		}
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

// BenchmarkSlice does not manage to fully inline all the closures (as of 1.22)
func BenchmarkSlice(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		OfSlice(slice)(fi)
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

// BenchmarkSliceCheck DOES manage to fully inline all the closures.
// Note that the only change is to replace OfSlice(slice)(fi) with Check(OfSlice(slice)(fi))
//
// Inner loop:
// v185	00068 (+204) INCQ R8
// v169	00069 (?) NOP
// v174	00070 (+86) ADDL R9, github.com/dr2chase/iterbench/tiny_test.i-84(SP)
// v178	00071 (+220) MOVB $1, github.com/dr2chase/iterbench/tiny_test.ret-85(SP)
// v98	00072 (+204) CMPQ R8, DX
// b11	00073 (204) JGE 20
// v156	00074 (204) MOVL (CX)(R8*4), R9
// v157	00075 (+205) XCHGL AX, AX
// v158	00076 (+184) XCHGL AX, AX
// v279	00077 (+217) CMPB github.com/dr2chase/iterbench/tiny_test.ret-85(SP), $0
// b12	00078 (217) JNE 68
func BenchmarkSliceCheck(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		Check(OfSlice(slice))(fi)
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

// BenchmarkSliceOf shows that the iteration variable i is marked Addrtaken
// (which is costly) but otherwise inlines properly.  Compare with later
// BenchmarkSliceInlineCallEE and BenchmarkSliceInlineCallER
//
// Inner loop:
// v58	00050 (+113) MOVL github.com/dr2chase/iterbench/tiny_test.i-36(SP), DI
// v59	00051 (113) ADDL (CX)(SI*4), DI
// v65	00052 (+29) INCQ SI
// v56	00053 (+30) XCHGL AX, AX
// v61	00054 (113) MOVL DI, github.com/dr2chase/iterbench/tiny_test.i-36(SP)
// v27	00055 (+29) CMPQ SI, DX
// b7	00056 (29) JLT 50
func BenchmarkSliceOf(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		slice.Of(fi)
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

// BenchmarkSliceOfCheck shows how a method closure passed to an iterator
// thwarts inlining (and results in a pair of allocations)
func BenchmarkSliceOfCheck(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		Check(slice.Of)(fi)
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

// BenchmarkSliceInlineCallEE is the same as BenchmarkSlice except that the
// callee has V2 manually inlined into it.  Full inlining results, i is marked Addrtaken.
//
// Inner loop:
// v119	00052 (+162) MOVL github.com/dr2chase/iterbench/tiny_test.i-36(SP), DI
// v120	00053 (162) ADDL (CX)(SI*4), DI
// v127	00054 (+220) INCQ SI
// v116	00055 (+221) XCHGL AX, AX
// v117	00056 (+211) XCHGL AX, AX
// v122	00057 (162) MOVL DI, github.com/dr2chase/iterbench/tiny_test.i-36(SP)
// v91	00058 (+220) CMPQ SI, DX
// b9	00059 (220) JLT 52
func BenchmarkSliceInlineCallEE(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		OfSliceInlineV2(slice)(fi)
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

// BenchmarkSliceInlineCallER is the same as BenchmarkSlice except that the call to
// OfSlice is manually inlined into the loop.  Full inlining results, i is marked Addrtaken.
//
// Inner loop:
// v119	00051 (+187) MOVL github.com/dr2chase/iterbench/tiny_test.i-36(SP), DI
// v120	00052 (187) ADDL (CX)(SI*4), DI
// v127	00053 (+228) INCQ SI
// v116	00054 (+229) XCHGL AX, AX
// v117	00055 (+208) XCHGL AX, AX
// v122	00056 (187) MOVL DI, github.com/dr2chase/iterbench/tiny_test.i-36(SP)
// v91	00057 (+228) CMPQ SI, DX
// b9	00058 (228) JLT 51
func BenchmarkSliceInlineCallER(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		V2(OfSliceIndex(slice))(fi)
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

// OfSlice returns a Seq over the elements of s. It is equivalent to
// range s with the index ignored.
func OfSlice(s []Int32) Seq {
	return V2(OfSliceIndex(s))
}

// V2 converts a Seq of pairs into a Seq of the second element.
func V2(seq Seq2) Seq {
	return func(yield func(Int32) bool) {
		seq(func(v1 int, v2 Int32) bool {
			return yield(v2)
		})
	}
}

// OfSliceInlineV2 returns a Seq over the elements of s. It is equivalent to
// range s with the index ignored.  Call to V2 has been hand-inlined.
func OfSliceInlineV2(s []Int32) Seq {
	seq := OfSliceIndex(s)
	return func(yield func(Int32) bool) {
		seq(func(v1 int, v2 Int32) bool {
			return yield(v2)
		})
	}
}

// OfSliceIndex returns a Seq over the elements of s. It is equivalent
// to range s.
func OfSliceIndex(s []Int32) Seq2 {
	return func(yield func(int, Int32) bool) {
		for i, v := range s {
			if !yield(i, v) {
				return
			}
		}
		return
	}
}

func Check(forall Seq) Seq {
	return func(body func(Int32) bool) {
		ret := true
		forall(func(v Int32) bool {
			if !ret {
				panic("Iterator access after exit")
			}
			ret = body(v)
			return ret
		})
	}
}
