// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package call_vs_iter_test

import (
	"iter"
	"os"
	"testing"

	"github.com/dr2chase/iter_test"
)

// The benchmarks here compare various ways of iterating.
// The iterations tend to run over 14 things, twice, because reasons.

var t1, t2, t1Small, t2Small *iter_test.T[Int32, sstring]
var t1Len, t2Len int

var m1 map[int]string = make(map[int]string)
var m2 map[int]string = make(map[int]string)

var global, sink int

func TestMain(m *testing.M) {
	t1 = &iter_test.T[Int32, sstring]{}
	t1.Insert(1, sstring{"ant"})
	t1.Insert(2, sstring{"bat"})
	t1.Insert(3, sstring{"cat"})
	t1.Insert(4, sstring{"dog"})
	t1.Insert(5, sstring{"emu"})
	t1.Insert(6, sstring{"fox"})
	t1.Insert(7, sstring{"gnu"})
	t1Small = t1.Copy()
	t1.Insert(8, sstring{"hen"})
	t1.Insert(9, sstring{"imp"})
	t1.Insert(10, sstring{"jay"})
	t1.Insert(11, sstring{"koi"})
	t1.Insert(12, sstring{"loi"})
	t1.Insert(13, sstring{"moi"})
	t1.Insert(14, sstring{"noi"})

	t2 = &iter_test.T[Int32, sstring]{}
	t2.Insert(21, sstring{"auntie"})
	t2.Insert(22, sstring{"batty"})
	t2.Insert(23, sstring{"catty"})
	t2.Insert(24, sstring{"doggie"})
	t2.Insert(25, sstring{"emu-like"})
	t2.Insert(26, sstring{"foxy"})
	t2.Insert(27, sstring{"gnu-like"})
	t2Small = t2.Copy()
	t2.Insert(28, sstring{"henny"})
	t2.Insert(29, sstring{"impish"})
	t2.Insert(30, sstring{"jaylike"})
	t2.Insert(31, sstring{"koioid"})
	t2.Insert(32, sstring{"loioid"})
	t2.Insert(33, sstring{"moioid"})
	t2.Insert(34, sstring{"noioid"})
	// call flag.Parse() here if TestMain uses flags
	t1Len, t2Len = t1.Size(), t2.Size()

	for k, v := range t1.DoAll2 {
		m1[int(k)] = v.s
	}

	for k, v := range t2.DoAll2 {
		m2[int(k)] = v.s
	}

	if t1.Size() != 14 {
		panic("Expected t1 Size == 14")
	}
	if t2.Size() != 14 {
		panic("Expected t2 Size == 14")
	}
	if t1Small.Size() != 7 {
		panic("Expected t1Small Size == 7")
	}
	if t2Small.Size() != 7 {
		panic("Expected t2Small Size == 7")
	}
	if len(m1) != 14 {
		panic("Expected m1 len == 14")
	}
	if len(m2) != 14 {
		panic("Expected m2 len == 14")
	}

	os.Exit(m.Run())
}

/* This is a series of variants on a benchmark intended to measure costs, possibly improve code. */

/* The first three do "exactly the same thing" visiting all the elements of two data structures */

// BenchmarkDoAll2FlatCall measures the cost of the plain call of the iterator for a non-recursive in-order tree visit.
func BenchmarkDoAll2FlatCall(b *testing.B) {
	b.ReportAllocs()
	i := 0
	yield := func(x Int32, y sstring) bool {
		i += int(x) + len(y.s)
		return true
	}
	for range b.N {
		t1.DoAll2Flat(yield)
		t2.DoAll2Flat(yield)
	}
	sink += i
}

// BenchmarkDoAll2FlatFunc measures the cost of iterating a two-value closure that does a non-recursive visit.
func BenchmarkDoAll2FlatFunc(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range t1.DoAll2FlatFunc() {
			i += int(x) + len(y.s)
		}
		for x, y := range t2.DoAll2FlatFunc() {
			i += int(x) + len(y.s)
		}
	}
	sink += i
}

// BenchmarkDoAll2FlatMethod measures the cost of iterating a two-value method value closure that does a non-recursive visit.
func BenchmarkDoAll2FlatMethod(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range t1.DoAll2Flat {
			i += int(x) + len(y.s)
		}
		for x, y := range t2.DoAll2Flat {
			i += int(x) + len(y.s)
		}
	}
	sink += i
}

func True(x Int32, y sstring) bool {
	return true
}

func False(x Int32, y sstring) bool {
	return false
}

// Filter2 returns a Seq2 that yields only the values (v1, v2) of seq for which
// f(v1, v2) returns true.
func Filter2[T, U any](seq iter.Seq2[T, U], f func(T, U) bool) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		seq(func(v T, u U) bool {
			if !f(v, u) {
				return true
			}
			return yield(v, u)
		})
	}
}

/* The next eight benchmarks measure a filtered (actually always true) visit of the elements of two data structures. */

// BenchmarkDoAll2FlatFuncIterIfTrue measures the cost of iterating
// a two-value method value closure that does a non-recursive visit, where the
// body of the loop contains a call to a filter function, that is a constant
// (that can be inlined).
func BenchmarkDoAll2FlatFuncIterIfTrue(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range t1.DoAll2FlatFunc() {
			if True(x, y) {
				i += int(x) + len(y.s)
			}
		}
		for x, y := range t2.DoAll2FlatFunc() {
			if True(x, y) {
				i += int(x) + len(y.s)
			}
		}
	}
	sink += i
}

// BenchmarkDoAll2FlatFilterCallTrue calls a visit method that also takes a filter parameter.
func BenchmarkDoAll2FlatFilterCallTrue(b *testing.B) {
	b.ReportAllocs()
	i := 0
	yield := func(x Int32, y sstring) bool {
		i += int(x) + len(y.s)
		return true
	}
	for range b.N {
		t1.DoAll2FlatFilter(yield, True)
		t2.DoAll2FlatFilter(yield, True)
	}
	sink += i
}

// BenchmarkDoAll2FlatFilterFuncTrueCall calls a closure from method that also takes a filter parameter.
func BenchmarkDoAll2FlatFilterFuncTrueCall(b *testing.B) {
	b.ReportAllocs()
	i := 0
	yield := func(x Int32, y sstring) bool {
		i += int(x) + len(y.s)
		return true
	}
	for range b.N {
		t1.DoAll2FlatFilterFunc(True)(yield)
		t2.DoAll2FlatFilterFunc(True)(yield)
	}
	sink += i
}

// BenchmarkDoAll2FlatFilterFuncTrueIter calls a method to generate a filtered iterator.
// The filter parameter is a constant function.
func BenchmarkDoAll2FlatFilterFuncTrueIter(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range t1.DoAll2FlatFilterFunc(True) {
			i += int(x) + len(y.s)
		}
		for x, y := range t2.DoAll2FlatFilterFunc(True) {
			i += int(x) + len(y.s)
		}
	}
	sink += i
}

var sel = True

// BenchmarkDoAll2FlatFuncIterIfLocTrue implements "filter" as an inline test in the iterated body, where the filter is a local copy of a global function variable.
func BenchmarkDoAll2FlatFuncIterIfLocTrue(b *testing.B) {
	b.ReportAllocs()
	i := 0
	sel := sel
	for range b.N {
		for x, y := range t1.DoAll2FlatFunc() {
			if sel(x, y) {
				i += int(x) + len(y.s)
			}
		}
		for x, y := range t2.DoAll2FlatFunc() {
			if sel(x, y) {
				i += int(x) + len(y.s)
			}
		}
	}
	sink += i
}

// BenchmarkDoAll2FlatFuncIterIfGlobTrue implements "filter" as an inline test in the iterated body, where the filter is a global function variable.
func BenchmarkDoAll2FlatFuncIterIfGlobTrue(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range t1.DoAll2FlatFunc() {
			if sel(x, y) {
				i += int(x) + len(y.s)
			}
		}
		for x, y := range t2.DoAll2FlatFunc() {
			if sel(x, y) {
				i += int(x) + len(y.s)
			}
		}
	}
	sink += i
}

// BenchmarkDoAll2FlatFuncCombinatorFilterTrueIter iterates a filter combinator applied to a closure returned by a method.
func BenchmarkDoAll2FlatFuncCombinatorFilterTrueIter(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range Filter2(t1.DoAll2FlatFunc(), True) {
			i += int(x) + len(y.s)
		}
		for x, y := range Filter2(t2.DoAll2FlatFunc(), True) {
			i += int(x) + len(y.s)
		}
	}
	sink += i
}

// BenchmarkDoAll2FlatMethodCombinatorFilterTrueIter iterates a filter combinator applied to a method closure
func BenchmarkDoAll2FlatMethodCombinatorFilterTrueIter(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range Filter2(t1.DoAll2Flat, True) {
			i += int(x) + len(y.s)
		}
		for x, y := range Filter2(t2.DoAll2Flat, True) {
			i += int(x) + len(y.s)
		}
	}
	sink += i
}

/* The next three compare different ways of filtering when the filter function is false, i.e., the "body" func is never called. */

// BenchmarkDoAll2FlatFilterCallFalse calls a visit method that also takes a filter parameter, which is false.
// Compare with BenchmarkDoAll2FlatFilterCallTrue
func BenchmarkDoAll2FlatFilterCallFalse(b *testing.B) {
	b.ReportAllocs()
	i := 0
	yield := func(x Int32, y sstring) bool {
		i += int(x) + len(y.s)
		return true
	}
	for range b.N {
		t1.DoAll2FlatFilter(yield, False)
		t2.DoAll2FlatFilter(yield, False)
	}
	sink += i
}

// BenchmarkDoAll2FlatFilterFuncFalseCall measures the cost of a CALL of a two-value method value closure that does a non-recursive filtered visit..
// Compare with BenchmarkDoAll2FlatFilterFuncTrueCall
func BenchmarkDoAll2FlatFilterFuncFalseCall(b *testing.B) {
	b.ReportAllocs()
	i := 0
	yield := func(x Int32, y sstring) bool {
		i += int(x) + len(y.s)
		return true
	}
	for range b.N {
		t1.DoAll2FlatFilterFunc(False)(yield)
		t2.DoAll2FlatFilterFunc(False)(yield)
	}
	sink += i
}

// BenchmarkDoAll2FlatFilterFuncFalseIter measures the cost of a ITERATION of a two-value method value closure that does a non-recursive filtered visit.
// Compare with BenchmarkDoAll2FlatFilterFuncTrueIter
func BenchmarkDoAll2FlatFilterFuncFalseIter(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range t1.DoAll2FlatFilterFunc(False) {
			i += int(x) + len(y.s)
		}
		for x, y := range t2.DoAll2FlatFilterFunc(False) {
			i += int(x) + len(y.s)
		}
	}
	sink += i
}


type sstring struct {
	s string
}

func (s sstring) String() string {
	return s.s
}

var z sstring

func stringer(s string) sstring {
	return sstring{s}
}

type Int32 int32

func (x Int32) Compare(y Int32) int {
	if x < y {
		return -1
	}
	if x > y {
		return 1
	}
	return 0
}

func compare(x, y Int32) int {
	return x.Compare(y)
}
