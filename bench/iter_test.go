// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter_test

import (
	"fmt"
	"iter"
	"math"
	"os"
	_ "runtime"
	"slices"
	"sync"
	"testing"

	"github.com/dr2chase/iter_test"
	"github.com/dr2chase/xiter"
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

// BenchmarkCountOldILocal measures a plain 3-clause for loop updating a LOCAL variable.
func BenchmarkCountOldILocal(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := 1; x <= 14; x++ {
			i += x
		}
		for x := 1; x <= 14; x++ {
			i += x
		}
	}
	sink += i
}

// BenchmarkCountOldIGlobal measures a plain 3-clause for loop updating a GLOBAL variable.
func BenchmarkCountOldIGlobal(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := 1; x <= 14; x++ {
			global += x
		}
		for x := 1; x <= 14; x++ {
			global += x
		}
	}
	sink += global
}

// BenchmarkCountOldIPGlobal measures a plain 3-clause for loop updating through a pointer
// (Was it the address formation that cost?)
func BenchmarkCountOldIPGlobal(b *testing.B) {
	b.ReportAllocs()
	pglobal := &global
	for range b.N {
		for x := 1; x <= 14; x++ {
			*pglobal += x
		}
		for x := 1; x <= 14; x++ {
			*pglobal += x
		}
	}
	sink += global
}

// BenchmarkLimitGenerate measures Limit(Generate ...) updating a local variable.
func BenchmarkLimitGenerate(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Limit(xiter.Generate(1, 1), t1Len) {
			i += x
		}
		for x := range xiter.Limit(xiter.Generate(1, 1), t2Len) {
			i += x
		}
	}
	sink += i
}

// BenchmarkOf measures xiter.Of updating a local variable.
func BenchmarkOf(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14) {
			i += x
		}
		for x := range xiter.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14) {
			i += x
		}
	}
	sink += i
}

// / BenchmarkSliceOldILocal measures an old range-of-slice loop updating a LOCAL variable.
func BenchmarkSliceOldILocal(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _, x := range slice {
			i += x
		}
		for _, x := range slice {
			i += x
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
	sink += i
}

// BenchmarkSliceOldIGlobal measures an old range-of-slice loop updating a GLOBAL variable.
func BenchmarkSliceOldIGlobal(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	for range b.N {
		for _, x := range slice {
			global += x
		}
		for _, x := range slice {
			global += x
		}
	}
	sink += global
}

// BenchmarkSlice measures xiter.OfSlice (generic), updating a local.
func BenchmarkSlice(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.OfSlice(slice) {
			i += x
		}
		for x := range xiter.OfSlice(slice) {
			i += x
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// BenchmarkSliceOpt measures a hand-optimized generic slice iteration, updating a local.
func BenchmarkSliceOpt(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range OfSliceOpt(slice) {
			i += x
		}
		for x := range OfSliceOpt(slice) {
			i += x
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// BenchmarkSliceInlineV2 measures a hand-inlined generic slice iteration, updating a local.
func BenchmarkSliceInlineV2(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range OfSliceInlineV2(slice) {
			i += x
		}
		for x := range OfSliceInlineV2(slice) {
			i += x
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// BenchmarkSliceSpecialized measures a non-generic slice iteration copied from xiter.OfSlice etc, updating a local.
// I.e., do we pay any cost for generics here?
func BenchmarkSliceSpecialized(b *testing.B) {
	slice := []Int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := Int32(0)
	for range b.N {
		for x := range OfSlice(slice) {
			i += x
		}
		for x := range OfSlice(slice) {
			i += x
		}
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// BenchmarkSliceSpecializedOpt measures a hand-optimized non-generic slice iteration, updating a local.
// I.e., do we pay any cost for generics here?
func BenchmarkSliceSpecializedOpt(b *testing.B) {
	slice := []Int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := Int32(0)
	for range b.N {
		for x := range OfSliceSpecializedOpt(slice) {
			i += x
		}
		for x := range OfSliceSpecializedOpt(slice) {
			i += x
		}
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

var gslice = []Int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

// BenchmarkSliceSpecializedOptX measures a hand-optimized non-generic slice iteration, updating a local.
// It calls the loop body explicitly instead of relying on the range-function transformation.
// Does this have any effect on what we see?
func BenchmarkSliceSpecializedOptX(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for range b.N {
		OfSliceSpecializedOpt(slice)(fi)
		OfSliceSpecializedOpt(slice)(fi)
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// BenchmarkSliceSpecializedX measures a non-generic slice iteration, updating a local.
// It calls the loop body explicitly instead of relying on the range-function transformation.
// Does this have any effect on what we see?
func BenchmarkSliceSpecializedX(b *testing.B) {
	b.ReportAllocs()
	slice := gslice
	i := Int32(0)
	fi := func(x Int32) bool {
		i += x
		return true
	}
	for range b.N {
		OfSlice(slice)(fi)
		OfSlice(slice)(fi)
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// BenchmarkSliceChecked wraps the xiter.OfSlice range function in a Check function.
func BenchmarkSliceChecked(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range Check(xiter.OfSlice(slice)) {
			i += x
		}
		for x := range Check(xiter.OfSlice(slice)) {
			i += x
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

// BenchmarkSliceSpecializedChecked wraps the (specialized) OfSlice range function in a Check function.
func BenchmarkSliceSpecializedChecked(b *testing.B) {
	slice := []Int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := Int32(0)
	for range b.N {
		for x := range CheckSeq(OfSlice(slice)) {
			i += x
		}
		for x := range CheckSeq(OfSlice(slice)) {
			i += x
		}
	}
	if i != Int32(b.N*15*14) {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

func BenchmarkSliceOldILocalBackward(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for j := len(slice) - 1; j >= 0; j-- {
			i += slice[j]
		}
		for j := len(slice) - 1; j >= 0; j-- {
			i += slice[j]
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
	sink += i
}

// BenchmarkSliceOldIGlobal measures an old range-of-slice loop updating a GLOBAL variable.
func BenchmarkSliceOldIGlobalBackward(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	for range b.N {
		for j := len(slice) - 1; j >= 0; j-- {
			global += slice[j]
		}
		for j := len(slice) - 1; j >= 0; j-- {
			global += slice[j]
		}
	}
	sink += global
}

// BenchmarkSlice measures xiter.OfSlice (generic), updating a local.
func BenchmarkSliceBackward(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _, x := range slices.Backward(slice) {
			i += x
		}
		for _, x := range slices.Backward(slice) {
			i += x
		}
	}
	if i != b.N*15*14 {
		panic(fmt.Errorf("Expected i = %d, got %d", b.N*15*14, i))
	}
}

func OfSliceOpt[T any, S ~[]T](s S) xiter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
		return
	}
}

// OfSliceInlineV2 returns a Seq over the elements of s. It is equivalent to
// range s with the index ignored.  Call to V2 has been hand-inlined.
func OfSliceInlineV2[T any, S ~[]T](s S) xiter.Seq[T] {
	seq := xiter.OfSliceIndex(s)
	return func(yield func(T) bool) {
		seq(func(v1 int, v2 T) bool {
			return yield(v2)
		})
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

// BenchmarkOldMap measures old-style map iteration, updating a local.
func BenchmarkOldMap(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for k, v := range m1 {
			i += k + len(v)
		}
		for k, v := range m2 {
			i += k + len(v)
		}
	}
	sink += i
}

// BenchmarkOldMap measures xiter.OfMap iteration, updating a local.
func BenchmarkOfMap(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for k, v := range xiter.OfMap(m1) {
			i += k + len(v)
		}
		for k, v := range xiter.OfMap(m2) {
			i += k + len(v)
		}
	}
	sink += i
}

// BenchmarkMapKeys measures xiter.MapKeys iteration, updating a local.
func BenchmarkMapKeys(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.MapKeys(m1) {
			i += x
		}
		for x := range xiter.MapKeys(m2) {
			i += x
		}
	}
	sink += i
}

// BenchmarkMapValues measures xiter.MapValues iteration, updating a local.
func BenchmarkMapValues(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.MapValues(m1) {
			i += len(x)
		}
		for x := range xiter.MapValues(m2) {
			i += len(x)
		}
	}
	sink += i
}

// BenchmarkToPairOfMap measures xiter.ToPair(xiter.OfMap(...)) iteration, updating a local
func BenchmarkToPairOfMap(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for p := range xiter.ToPair(xiter.OfMap(m1)) {
			k, v := p.V1, p.V2
			i += k + len(v)
		}
		for p := range xiter.ToPair(xiter.OfMap(m2)) {
			k, v := p.V1, p.V2
			i += k + len(v)
		}
	}
	sink += i
}

// BenchmarkOldBytes is a baseline, using plain old for loops.
func BenchmarkOldBytes(b *testing.B) {
	slice := []byte("abcdefghijklmn")
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _, x := range slice {
			i += int(x - 'a' + 1)
		}
		for _, x := range slice {
			i += int(x - 'a' + 1)
		}
	}
	sink += i
}

// BenchmarkBytes measures the performance of xiter.Bytes
func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Bytes("abcdefghijklmn") {
			i += int(x - 'a' + 1)
		}
		for x := range xiter.Bytes("abcdefghijklmn") {
			i += int(x - 'a' + 1)
		}
	}
	sink += i
}

// BenchmarkOldRunes is a baseline, using plain old for loops, applied to strings of emoji.
func BenchmarkOldRunes(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _, x := range "🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍" {
			i += int(x - '🚀' + 1)
		}
		for _, x := range "🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍" {
			i += int(x - '🚀' + 1)
		}
	}
	sink += i
}

// BenchmarkRunes measures the performance of xiter.Runes
func BenchmarkRunes(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Runes("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍") {
			i += int(x - '🚀' + 1)
		}
		for x := range xiter.Runes("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍") {
			i += int(x - '🚀' + 1)
		}
	}
	sink += i
}

// BenchmarkStringSplitEmpty measures the performance of xiter.StringSplit with
// an empty separator, applied to the same string of runes as BenchmarkRunes.
func BenchmarkStringSplitEmpty(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.StringSplit("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍", "") {
			i += int(x[0])
		}
		for x := range xiter.StringSplit("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍", "") {
			i += int(x[0])
		}
	}
	sink += i
}

// BenchmarkStringSplit measures the performance of xiter.StringSplit with
// a "." separator, applied to the same string of runes as BenchmarkRunes but
// separated by "." .
func BenchmarkStringSplit(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.StringSplit("🚀.🚁.🚂.🚃.🚄.🚅.🚆.🚇.🚈.🚉.🚊.🚋.🚌.🚍", ".") {
			i += int(x[0])
		}
		for x := range xiter.StringSplit("🚀 🚁 🚂 🚃 🚄 🚅 🚆 🚇 🚈 🚉 🚊 🚋 🚌 🚍", " ") {
			i += int(x[0])
		}
	}
	sink += i
}

// BenchmarkDoAllOld is a baseline, that calls a recursive tree walk visit on a supplied function.
func BenchmarkDoAllOld(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		t1.DoAll(func(x Int32) bool { i += int(x); return true })
		t2.DoAll(func(x Int32) bool { i += int(x); return true })
	}
	sink += i
}

// BenchmarkDoAll measures the cost of iterating a closure returned by a function.
func BenchmarkDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range t1.DoAllFunc() {
			i += int(x)
		}
		for x := range t2.DoAllFunc() {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkDoAll2 measures the cost of iterating a method value closure.
func BenchmarkDoAll2(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x, y := range t1.DoAll2 {
			i += int(x) + len(y.s)
		}
		for x, y := range t2.DoAll2 {
			i += int(x) + len(y.s)
		}
	}
	sink += i
}

// BenchmarkDoAll2 measures the cost of iterating a method value closure that does a non-recursive visit.
func BenchmarkDoAll2Flat(b *testing.B) {
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

// BenchmarkDoAllCheck measures the cost of the recursive iterator when wrapped in a "Check" function
// (not to be confused with the automatically inserted checking).
func BenchmarkDoAllCheck(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range Check(t1.DoAllFunc()) {
			i += int(x)
		}
		for x := range Check(t2.DoAllFunc()) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkDoAll2FlatEqualZip primarily checks that the flat iterator is correct
// (at least for these inputs), but also measures the cost of a complex composition
// of iterator transformers.
func BenchmarkDoAll2FlatEqualZip(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		if !xiter.Equal(xiter.V1(t1.DoAll2), xiter.V1(t1.DoAll2Flat)) {
			panic("Should have been equal")
		}
		if !xiter.EqualFunc(xiter.V2(t1.DoAll2), xiter.V2(t1.DoAll2Flat), func(x, y sstring) bool { return x == y }) {
			panic("Should have been equal")
		}
	}
}

// BenchmarkFromPair measures the cost of `for k, v := range xiter.FromPair(xiter.ToPair(t1.DoAll2Func()))`
func BenchmarkFromPair(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for k, v := range xiter.FromPair(xiter.ToPair(t1.DoAll2Func())) {
			i += int(k) + len(v.s)
		}
		for k, v := range xiter.FromPair(xiter.ToPair(t2.DoAll2Func())) {
			i += int(k) + len(v.s)
		}
	}
	sink += i
}

// BenchmarkMapDoAll measures the cost of wrapping a recursive walk in a xiter.Map.
// The body of this loop is unfortunately expensive.
func BenchmarkMapDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Map(t1.DoAllFunc(), func(j Int32) float64 { return float64(j) }) {
			i += int(math.Exp(x))
		}
		for x := range xiter.Map(t2.DoAllFunc(), func(j Int32) float64 { return float64(j) }) {
			i += int(math.Exp(x))
		}
	}
	sink += i
}

// BenchmarkFilterDoAll measures the cost of wrapping a recursive walk in xiter.Filter
func BenchmarkFilterDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Filter(t1.DoAllFunc(), func(j Int32) bool { return j&1 == 0 }) {
			i += int(x)
		}
		for x := range xiter.Filter(t2.DoAllFunc(), func(j Int32) bool { return j&1 == 0 }) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkLimitDoAll measures the cost of wrapping a recursive walk in xiter.Limit.
func BenchmarkLimitDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Limit(t1.DoAllFunc(), t1Len/2) {
			i += int(x)
		}
		for x := range xiter.Limit(t2.DoAllFunc(), t2Len/2) {
			i += int(x)
		}
		for x := range xiter.Limit(t1.DoAllFunc(), t1Len/2) {
			i += int(x)
		}
		for x := range xiter.Limit(t2.DoAllFunc(), t2Len/2) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkLimitCheckDoAll is the same as BenchmarkLimitDoAll, but wrapped in a call to Check.
func BenchmarkLimitCheckDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range Check(xiter.Limit(t1.DoAllFunc(), t1Len/2)) {
			i += int(x)
		}
		for x := range Check(xiter.Limit(t2.DoAllFunc(), t2Len/2)) {
			i += int(x)
		}
		for x := range Check(xiter.Limit(t1.DoAllFunc(), t1Len/2)) {
			i += int(x)
		}
		for x := range Check(xiter.Limit(t2.DoAllFunc(), t2Len/2)) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkReduceDoAll measures the cost of xiter.Reduce applied to DoAllFunc()
func BenchmarkReduceDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		i += xiter.Reduce(t1.DoAllFunc(), 0, func(total int, v Int32) int { return total + int(v) })
		i += xiter.Reduce(t2.DoAllFunc(), 0, func(total int, v Int32) int { return total + int(v) })
	}
	sink += i
}

// BenchmarkFoldDoAll measures the cost of xiter.Fold applied to DoAllFunc()
func BenchmarkFoldDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		i += int(xiter.Fold(t1.DoAllFunc(), func(total Int32, v Int32) Int32 { return total + v }))
		i += int(xiter.Fold(t2.DoAllFunc(), func(total Int32, v Int32) Int32 { return total + v }))
	}
	sink += i
}

// BenchmarkSkipDoAll measures the cost of xiter.Skip applied to DoAllFunc()
func BenchmarkSkipDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Skip(t1.DoAllFunc(), t1Len/2) {
			i += int(x)
		}
		for x := range xiter.Skip(t2.DoAllFunc(), t2Len/2) {
			i += int(x)
		}
		for x := range xiter.Skip(t1.DoAllFunc(), t1Len/2) {
			i += int(x)
		}
		for x := range xiter.Skip(t2.DoAllFunc(), t2Len/2) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkSkipSpecializedDoAll measures the cost of a type-specific Skip applied to DoAllFunc().
// Generics make little or no difference in the cost.
func BenchmarkSkipSpecializedDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range Skip(t1.DoAllFunc(), t1Len/2) {
			i += int(x)
		}
		for x := range Skip(t2.DoAllFunc(), t2Len/2) {
			i += int(x)
		}
		for x := range Skip(t1.DoAllFunc(), t1Len/2) {
			i += int(x)
		}
		for x := range Skip(t2.DoAllFunc(), t2Len/2) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkSkipCheckedDoAll wraps xiter.Skip applied to DoAllFunc() in a Check.
func BenchmarkSkipCheckedDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range Check(xiter.Skip(t1.DoAllFunc(), t1Len/2)) {
			i += int(x)
		}
		for x := range Check(xiter.Skip(t2.DoAllFunc(), t2Len/2)) {
			i += int(x)
		}
		for x := range Check(xiter.Skip(t1.DoAllFunc(), t1Len/2)) {
			i += int(x)
		}
		for x := range Check(xiter.Skip(t2.DoAllFunc(), t2Len/2)) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkSkipSpecializedCheckedDoAll is a non-generic version of BenchmarkSkipCheckedDoAll
func BenchmarkSkipSpecializedCheckedDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range CheckSeq(Skip(t1.DoAllFunc(), t1Len/2)) {
			i += int(x)
		}
		for x := range CheckSeq(Skip(t2.DoAllFunc(), t2Len/2)) {
			i += int(x)
		}
		for x := range CheckSeq(Skip(t1.DoAllFunc(), t1Len/2)) {
			i += int(x)
		}
		for x := range CheckSeq(Skip(t2.DoAllFunc(), t2Len/2)) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkSkipMethodDoAll measures the cost of xiter.Skip applied to t.DoAll (a method closure)
// That is, exactly like BenchmarkSkipDoAll but a method value closure instead.
func BenchmarkSkipMethodDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range Skip(t1.DoAll, t1Len/2) {
			i += int(x)
		}
		for x := range Skip(t2.DoAll, t2Len/2) {
			i += int(x)
		}
		for x := range Skip(t1.DoAll, t1Len/2) {
			i += int(x)
		}
		for x := range Skip(t2.DoAll, t2Len/2) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkConcatDoAll measures the cost of xiter.Concat
func BenchmarkConcatDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Concat(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkConcatRSCDoAll measures the cost of the version of Concat in the proposal
func BenchmarkConcatRSCDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range Concat(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkMergeFuncDoAll measures the cost of xiter.MergeFunc applied
// to a pair of DoAllFunc() iterators.
// Note that xiter.MergeFunc was rewritten into a 1-Pull coroutine form.
func BenchmarkMergeFuncDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.MergeFunc(t1.DoAllFunc(), t2.DoAllFunc(), compare) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkMergeFuncSpecializedDoAll measures the cost of a non-generic version of MergeFunc.
// This is a flawed comparison because xiter.MergeFunc was rewritten into the 1-Pull coroutine form.
func BenchmarkMergeFuncSpecializedDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range MergeFuncSpecialized(t1.DoAllFunc(), t2.DoAllFunc(), compare) {
			i += int(x)
		}
	}
	sink += i
}

// BenchmarkZipDoAll measures xiter.Zip of a pair of tree.DoAllFunc()
func BenchmarkZipDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range xiter.Zip(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkZipDoAll measures non-generic Zip of a pair of tree.DoAllFunc()
func BenchmarkZipSpecializedDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range ZipSpecialized(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkZipSpecializedSODoAll measures non-generic Zip implemented using a Pull
// with a fragile stop function (does not call sync.Once).
func BenchmarkZipSpecializedSODoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range ZipSpecializedSO(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkZipSpecialized1PullDoAll measures non-generic Zip rewritten to use only
// one Pull iterator.
func BenchmarkZipSpecialized1PullDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range ZipSpecialized1Pull(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkZipGo2PullDoAll measures the cost of 2-pull goroutine generic Zip.
func BenchmarkZipGo2PullDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range ZipGo(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkZipGo1PullDoAll measures the cost of 1-pull goroutine generic Zip.
func BenchmarkZipGo1PullDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range ZipGo1Pull(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkZipCoro2PullDoAll measures the cost of 2-pull COroutine generic Zip
func BenchmarkZipCoro2PullDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range ZipCoro(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkZipCoro2PullDoAll measures the cost of 1-pull COroutine generic Zip
func BenchmarkZipCoro1PullDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := range ZipCoro1Pull(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
	sink += i
}

// BenchmarkEqualDoAll measures the cost of generic Equal applied to a pair of DoAllFunc()
// It runs on 4 iterators (an equality test and an inequality test) but should terminate quickly
// on the inequalitytest.
func BenchmarkEqualDoAll(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		if !xiter.Equal(t1.DoAllFunc(), t1.DoAllFunc()) {
			panic("equal failure")
		}
		if xiter.Equal(t1.DoAllFunc(), t2.DoAllFunc()) {
			panic("unequal failure")
		}
	}
	sink += i
}

// BenchmarkEqualDoAllSmall measures the cost of generic Equal applied to a pair of small DoAllFunc()
// sequences.  This is to ensure that the allocations performed are not proportional to the length
// of the iteration.
func BenchmarkEqualDoAllSmall(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		if !xiter.Equal(t1Small.DoAllFunc(), t1Small.DoAllFunc()) {
			panic("equal failure")
		}
		if xiter.Equal(t1Small.DoAllFunc(), t2Small.DoAllFunc()) {
			panic("unequal failure")
		}
	}
	sink += i
}

// BenchmarkCountOldILocal measures a plain 3-clause for loop updating a LOCAL variable.
func BenchmarkNestedCountOldILocal(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for x := 1; x <= 14; x++ {
			for y := 1; y <= 14; y++ {
				i += y
			}
		}
	}
	sink += i
}

// BenchmarkLimitGenerate measures Limit(Generate ...) updating a local variable.
func BenchmarkNestedLimitGenerate(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range xiter.Limit(xiter.Generate(1, 1), t1Len) {
			for x := range xiter.Limit(xiter.Generate(1, 1), t2Len) {
				i += x
			}
		}
	}
	sink += i
}

// BenchmarkOf measures xiter.Of updating a local variable.
func BenchmarkNestedOf(b *testing.B) {
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range xiter.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14) {
			for x := range xiter.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14) {
				i += x
			}
		}
	}
	sink += i
}

// / BenchmarkSliceOldILocal measures an old range-of-slice loop updating a LOCAL variable.
func BenchmarkNestedSliceOldILocal(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _, _ = range slice {
			for _, x := range slice {
				i += x
			}
		}
	}
	sink += i
}

// BenchmarkSlice measures xiter.OfSlice (generic), updating a local.
func BenchmarkNestedSlice(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range xiter.OfSlice(slice) {
			for x := range xiter.OfSlice(slice) {
				i += x
			}
		}
	}
	sink += i
}

// BenchmarkSliceOpt measures a hand-optimized generic slice iteration, updating a local.
func BenchmarkNestedSliceOpt(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			for y := range OfSliceOpt(slice) {
				i += y
			}
		}
	}
	sink += i
}

// BenchmarkSliceInlineV2 measures a hand-inlined generic slice iteration, updating a local.
func BenchmarkNestedSliceInlineV2(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceInlineV2(slice) {
			for y := range OfSliceInlineV2(slice) {
				i += y
			}
		}
	}
	sink += i

}

func TestCheck(t *testing.T) {
	i := 0
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Saw panic %v", r)
		} else {
			t.Error("Wanted to see a failure")
		}
	}()
	for x := range Check(t1.DoAllTwice) {
		i += int(x)
		if i > 4*9 {
			break
		}
	}
}

func TestZips(t *testing.T) {
	if !xiter.EqualFunc(
		xiter.Zip(t1.DoAllFunc(), t2.DoAllFunc()),
		ZipSpecialized1Pull(t1.DoAllFunc(), t2.DoAllFunc()),
		func(z1 xiter.Zipped[Int32, Int32], z2 ZippedSpecialized) bool {
			return z1.V1 == z2.V1 && z1.V2 == z2.V2 && z1.OK1 == z2.OK1 && z1.OK2 == z2.OK2
		}) {
		t.Error("wanted equal 1")
	}
	if !xiter.EqualFunc(
		xiter.Zip(t1.DoAllFunc(), xiter.Skip(t2.DoAllFunc(), 1)),
		ZipSpecialized1Pull(t1.DoAllFunc(), Skip(t2.DoAllFunc(), 1)),
		func(z1 xiter.Zipped[Int32, Int32], z2 ZippedSpecialized) bool {
			return z1.V1 == z2.V1 && z1.V2 == z2.V2 && z1.OK1 == z2.OK1 && z1.OK2 == z2.OK2
		}) {
		t.Error("wanted equal 2")
	}
	if !xiter.EqualFunc(
		xiter.Zip(xiter.Skip(t1.DoAllFunc(), 1), t2.DoAllFunc()),
		ZipSpecialized1Pull(Skip(t1.DoAllFunc(), 1), t2.DoAllFunc()),
		func(z1 xiter.Zipped[Int32, Int32], z2 ZippedSpecialized) bool {
			return z1.V1 == z2.V1 && z1.V2 == z2.V2 && z1.OK1 == z2.OK1 && z1.OK2 == z2.OK2
		}) {
		t.Error("wanted equal 3")
	}
}

// local copy, specialized
type Seq func(yield func(Int32) bool)
type Seq2 func(yield func(int, Int32) bool)

// Concat from proposal itself, does this avoid allocations?
func Concat[V any](seqs ...xiter.Seq[V]) xiter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
		return
	}
}

// Check converts a seq to one that Checks it is not called after a false return.
func Check[V any](forall xiter.Seq[V]) xiter.Seq[V] {
	return func(body func(V) bool) {
		ret := true
		forall(func(v V) bool {
			if !ret {
				panic("Iterator access after exit")
			}
			ret = body(v)
			return ret
		})
	}
}

func MergeFuncSpecialized(seq1, seq2 Seq, compare func(Int32, Int32) int) Seq {
	return func(yield func(Int32) bool) {
		p1, stop := Pull(seq1)
		defer stop()
		p2, stop := Pull(seq2)
		defer stop()

		v1, ok1 := p1()
		v2, ok2 := p2()
		for ok1 || ok2 {
			var c int
			if ok1 && ok2 {
				c = compare(v1, v2)
			}

			switch {
			case !ok2 || c < 0:
				if !yield(v1) {
					return
				}
				v1, ok1 = p1()
			case !ok1 || c > 0:
				if !yield(v2) {
					return
				}
				v2, ok2 = p2()
			default:
				if !yield(v1) || !yield(v2) {
					return
				}
				v1, ok1 = p1()
				v2, ok2 = p2()
			}
		}

		return
	}
}

// Zipped holds values from an iteration of a Seq returned by [Zip].
type ZippedSpecialized struct {
	V1  Int32
	OK1 bool

	V2  Int32
	OK2 bool
}

// ZipSpecialized is a non-generic Zip that uses 2 Pulls and ordinary goroutines.
func ZipSpecialized(seq1 Seq, seq2 Seq) xiter.Seq[ZippedSpecialized] {
	return func(yield func(ZippedSpecialized) bool) {
		p1, stop := Pull(seq1)
		defer stop()
		p2, stop := Pull(seq2)
		defer stop()

		for {
			var val ZippedSpecialized
			val.V1, val.OK1 = p1()
			val.V2, val.OK2 = p2()
			if (!val.OK1 && !val.OK2) || !yield(val) {
				return
			}
		}
	}
}

// ZipSpecialized is a non-generic Zip that uses 2 PullSOs and ordinary goroutines.
// PullSO is like pull but does not use sync.Once in its stop function.
func ZipSpecializedSO(seq1 Seq, seq2 Seq) xiter.Seq[ZippedSpecialized] {
	return func(yield func(ZippedSpecialized) bool) {
		p1, stop1 := PullSO(seq1)
		p2, stop2 := PullSO(seq2)

		for {
			var val ZippedSpecialized
			val.V1, val.OK1 = p1()
			val.V2, val.OK2 = p2()
			if (!val.OK1 && !val.OK2) || !yield(val) {
				stop2()
				stop1()
				return
			}
		}
	}
}

// ZipSpecialized1Pull is a non-generic Zip that uses 1 Pull and ordinary goroutines.
func ZipSpecialized1Pull(seq1 Seq, seq2 Seq) xiter.Seq[ZippedSpecialized] {
	return func(body func(ZippedSpecialized) bool) {
		p2, stop2 := Pull(seq2)
		defer stop2()
		done := false
		for v := range seq1 {
			var val ZippedSpecialized
			val.V1, val.OK1 = v, true
			val.V2, val.OK2 = p2()
			if !body(val) {
				done = true
				break
			}
		}
		if done {
			return
		}
		// seq1 is exhausted
		for v2, ok2 := p2(); ok2; v2, ok2 = p2() {
			var v1 Int32
			var val ZippedSpecialized
			val.V1, val.OK1 = v1, false
			val.V2, val.OK2 = v2, true
			if !body(val) {
				return
			}
		}
		return
	}
}

// Zipped holds values from an iteration of a Seq returned by [Zip].
type Zipped[T1, T2 any] struct {
	V1  T1
	OK1 bool

	V2  T2
	OK2 bool
}

// ZipGo is a generic Zip using 2 xiter.GoPull (ordinary goroutines)
func ZipGo[T1, T2 any](seq1 iter.Seq[T1], seq2 iter.Seq[T2]) iter.Seq[Zipped[T1, T2]] {
	return func(yield func(Zipped[T1, T2]) bool) {
		p1, stop := xiter.GoPull(xiter.Seq[T1](seq1))
		defer stop()
		p2, stop := xiter.GoPull(xiter.Seq[T2](seq2))
		defer stop()

		for {
			var val Zipped[T1, T2]
			val.V1, val.OK1 = p1()
			val.V2, val.OK2 = p2()
			if (!val.OK1 && !val.OK2) || !yield(val) {
				return
			}
		}
	}
}

// ZipGo is a generic Zip using 2 iter.Pull (faster coroutines)
func ZipCoro[T1, T2 any](seq1 iter.Seq[T1], seq2 iter.Seq[T2]) iter.Seq[Zipped[T1, T2]] {
	return func(yield func(Zipped[T1, T2]) bool) {
		p1, stop := iter.Pull(seq1)
		defer stop()
		p2, stop := iter.Pull(seq2)
		defer stop()

		for {
			var val Zipped[T1, T2]
			val.V1, val.OK1 = p1()
			val.V2, val.OK2 = p2()
			if (!val.OK1 && !val.OK2) || !yield(val) {
				return
			}
		}
	}
}

// ZipGo1Pull is a generic Zip using 1 xiter.GoPull (ordinary goroutine)
func ZipGo1Pull[T1, T2 any](seq1 iter.Seq[T1], seq2 iter.Seq[T2]) iter.Seq[Zipped[T1, T2]] {
	return func(body func(Zipped[T1, T2]) bool) {
		p2, stop2 := xiter.GoPull(xiter.Seq[T2](seq2))
		defer stop2()
		done := false
		for v := range seq1 {
			var val Zipped[T1, T2]
			val.V1, val.OK1 = v, true
			val.V2, val.OK2 = p2()
			if !body(val) {
				done = true
				break
			}
		}
		if done {
			return
		}
		// seq1 is exhausted
		for v2, ok2 := p2(); ok2; v2, ok2 = p2() {
			var v1 T1
			var val Zipped[T1, T2]
			val.V1, val.OK1 = v1, false
			val.V2, val.OK2 = v2, true
			if !body(val) {
				return
			}
		}
		return
	}
}

// ZipCoro1Pull is a generic Zip using 1 iter.Pull (1 faster coroutine)
func ZipCoro1Pull[T1, T2 any](seq1 iter.Seq[T1], seq2 iter.Seq[T2]) iter.Seq[Zipped[T1, T2]] {
	return func(body func(Zipped[T1, T2]) bool) {
		p2, stop2 := iter.Pull(seq2)
		defer stop2()
		done := false
		for v := range seq1 {
			var val Zipped[T1, T2]
			val.V1, val.OK1 = v, true
			val.V2, val.OK2 = p2()
			if !body(val) {
				done = true
				break
			}
		}
		if done {
			return
		}
		// seq1 is exhausted
		for v2, ok2 := p2(); ok2; v2, ok2 = p2() {
			var v1 T1
			var val Zipped[T1, T2]
			val.V1, val.OK1 = v1, false
			val.V2, val.OK2 = v2, true
			if !body(val) {
				return
			}
		}
		return
	}
}

func Skip(forall Seq, n int) Seq {
	return func(body func(Int32) bool) {
		forall(func(v Int32) bool {
			if n > 0 {
				n--
				return true
			}
			return body(v)
		})
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

func Pull(seq Seq) (iter func() (Int32, bool), stop func()) {
	next := make(chan struct{})
	yield := make(chan Int32)

	go func() {
		defer close(yield)

		_, ok := <-next
		if !ok {
			return
		}

		seq(func(v Int32) bool {
			yield <- v
			_, ok := <-next
			return ok
		})
	}()

	return func() (v Int32, ok bool) {
			select {
			case <-yield:
				return v, false
			case next <- struct{}{}:
				v, ok := <-yield
				return v, ok
			}
		}, sync.OnceFunc(func() {
			close(next)
			<-yield
		})
}

func PullSO(seq Seq) (func() (Int32, bool), func()) {
	next := make(chan struct{})
	yield := make(chan Int32)
	stop := func() {
		close(next)
		<-yield
	}
	iter := func() (v Int32, ok bool) {
		select {
		case <-yield:
			return v, false
		case next <- struct{}{}:
			v, ok := <-yield
			return v, ok
		}
	}

	func() {
		go func() {
			defer close(yield)

			_, ok := <-next
			if !ok {
				return
			}

			seq(func(v Int32) bool {
				yield <- v
				_, ok := <-next
				return ok
			})
		}()
	}()

	return iter, stop
}

// todo flatten, handle, windows, chunks, split,

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
