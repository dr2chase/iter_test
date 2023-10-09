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

func BenchmarkSliceOld(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for j := 0; j < b.N; j++ {
		for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} {
			i += x
		}
		for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} {
			i += x
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}

}

func BenchmarkSliceOldIGlobal(b *testing.B) {
	b.ReportAllocs()
	for j := 0; j < b.N; j++ {
		for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} {
			i += x
		}
		for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} {
			i += x
		}
	}
}

var gslice = []Int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

func BenchmarkSliceSpecializedOptX(b *testing.B) {
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		OfSliceSpecializedOpt(slice)(fi)
		OfSliceSpecializedOpt(slice)(fi)
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

func BenchmarkSliceSpecializedX(b *testing.B) {
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		OfSlice(slice)(fi)
		OfSlice(slice)(fi)
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

func BenchmarkSliceSpecializedY(b *testing.B) {
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for j := 0; j < b.N; j++ {
		V2(OfSliceIndex(slice))(fi)
		V2(OfSliceIndex(slice))(fi)
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// OfSliceSpecializedOpt returns a Seq over the elements of s.
// It is equivalent to range s with the index ignored.
func OfSliceSpecializedOpt(s []Int32) Seq {
	return func(yield func(Int32) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
		return
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

// local copy, specialized
type Seq func(yield func(Int32) bool)
type Seq2 func(yield func(int, Int32) bool)
