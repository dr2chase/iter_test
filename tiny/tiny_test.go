// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter_test

// This is a subset of iter_test.go designed to simplify the search for some performance problems.

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type Int32 int32

var i int

var gslice = []Int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

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

func BenchmarkSliceSpecialized(b *testing.B) {
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

func BenchmarkSliceSpecializedCheck(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		CheckSeq(OfSlice(slice))(fi)
	}
	if i != Int32(b.N*15*7) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*7, i))
	}
}

func BenchmarkSliceSpecializedInlineCallEE(b *testing.B) {
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

func BenchmarkSliceSpecializedInlineCallER(b *testing.B) {
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

func CheckSeq(forall Seq) Seq {
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

// local copy, specialized
type Seq func(yield func(Int32) bool)
type Seq2 func(yield func(int, Int32) bool)
