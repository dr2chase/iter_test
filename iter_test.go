// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter_test

import (
	"deedles.dev/xiter"
	"github.com/dr2chase/iter"
	"math"
	"os"
	"sync"
	"testing"
)

var t1, t2 *iter.T[Int32, sstring]
var t1Len, t2Len int

var m1 map[int]string = make(map[int]string)
var m2 map[int]string = make(map[int]string)

func TestMain(m *testing.M) {
	t1 = &iter.T[Int32, sstring]{}
	t1.Insert(1, sstring{"ant"})
	t1.Insert(2, sstring{"bat"})
	t1.Insert(3, sstring{"cat"})
	t1.Insert(4, sstring{"dog"})
	t1.Insert(5, sstring{"emu"})
	t1.Insert(6, sstring{"fox"})
	t1.Insert(7, sstring{"gnu"})
	t1.Insert(8, sstring{"hen"})
	t1.Insert(9, sstring{"imp"})
	t1.Insert(10, sstring{"jay"})
	t1.Insert(11, sstring{"koi"})
	t1.Insert(12, sstring{"loi"})
	t1.Insert(13, sstring{"moi"})
	t1.Insert(14, sstring{"noi"})

	t2 = &iter.T[Int32, sstring]{}
	t2.Insert(21, sstring{"auntie"})
	t2.Insert(22, sstring{"batty"})
	t2.Insert(23, sstring{"catty"})
	t2.Insert(24, sstring{"doggie"})
	t2.Insert(25, sstring{"emu-like"})
	t2.Insert(26, sstring{"foxy"})
	t2.Insert(27, sstring{"gnu-like"})
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

	os.Exit(m.Run())
}

var i int

// for range 100_000 {
// 	i += doAll(t1, t2)
// 	i += doAll2(t1, t2)

// 	/// 61898 - xiter
// 	// map2, reduce2
// 	// mergefunc2, limit2, filter2
// 	// equal, equal2, equalfunc, equalfunc2
// 	// zip2

// 	// 61899 - slices
// 	// all, backward, values, append, collect
// 	// sorted, sortedfunc

// 	// 61900 - maps
// 	// all, keys, values, insert, collect

// 	// 61901 - bytes, strings
// 	// lines, bytes, runes
// 	// splitseq, runesplitseq
// 	// splitafterseq, fieldsseq, fieldsfuncseq

// 	// 61902 - regexp
// 	// (Find|All|FindAll)?(String)?(Submatch)?(Index)?

// }

func BenchmarkNothing(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
	}
}

func BenchmarkOldSlice(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} {
			i += x
		}
		for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} {
			i += x
		}
	}
}

func BenchmarkOldCount(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := 1; x <= 14; x++ {
			i += x
		}
		for x := 1; x <= 14; x++ {
			i += x
		}
	}
}

func BenchmarkMapKeys(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.MapKeys(m1) {
			i += x
		}
		for x := range xiter.MapKeys(m2) {
			i += x
		}
	}
}

func BenchmarkMapValues(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.MapValues(m1) {
			i += len(x)
		}
		for x := range xiter.MapValues(m2) {
			i += len(x)
		}
	}
}

func BenchmarkOfMap(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for k, v := range xiter.OfMap(m1) {
			i += k + len(v)
		}
		for k, v := range xiter.OfMap(m2) {
			i += k + len(v)
		}
	}
}

func BenchmarkToPairOfMap(b *testing.B) {
	b.ReportAllocs()
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
}

func BenchmarkFromPair(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for k, v := range xiter.FromPair(xiter.ToPair(t1.DoAll2Func())) {
			i += int(k) + len(v.s)
		}
		for k, v := range xiter.FromPair(xiter.ToPair(t2.DoAll2Func())) {
			i += int(k) + len(v.s)
		}
	}
}

func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Bytes("abcdefghijklmn") {
			i += int(x - 'a' + 1)
		}
		for x := range xiter.Bytes("abcdefghijklmn") {
			i += int(x - 'a' + 1)
		}
	}
}

func BenchmarkRunes(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Runes("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍") {
			i += int(x - '🚀' + 1)
		}
		for x := range xiter.Runes("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍") {
			i += int(x - '🚀' + 1)
		}
	}
}

func BenchmarkStringSplitEmpty(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.StringSplit("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍", "") {
			i += int(x[0])
		}
		for x := range xiter.StringSplit("🚀🚁🚂🚃🚄🚅🚆🚇🚈🚉🚊🚋🚌🚍", "") {
			i += int(x[0])
		}
	}
}

func BenchmarkStringSplit(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.StringSplit("🚀.🚁.🚂.🚃.🚄.🚅.🚆.🚇.🚈.🚉.🚊.🚋.🚌.🚍", ".") {
			i += int(x[0])
		}
		for x := range xiter.StringSplit("🚀 🚁 🚂 🚃 🚄 🚅 🚆 🚇 🚈 🚉 🚊 🚋 🚌 🚍", " ") {
			i += int(x[0])
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Limit(xiter.Generate(1, 1), t1Len) {
			i += x
		}
		for x := range xiter.Limit(xiter.Generate(1, 1), t2Len) {
			i += x
		}
	}
}

func BenchmarkOf(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14) {
			i += x
		}
		for x := range xiter.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14) {
			i += x
		}
	}
}

func BenchmarkOfSlice(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.OfSlice(slice) {
			i += x
		}
		for x := range xiter.OfSlice(slice) {
			i += x
		}
	}
}

func BenchmarkDoAll(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range t1.DoAllFunc() {
			i += int(x)
		}
		for x := range t2.DoAllFunc() {
			i += int(x)
		}
	}
}
func BenchmarkDoAll2(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x, y := range t1.DoAll2 {
			i += int(x) + len(y.s)
		}
		for x, y := range t2.DoAll2 {
			i += int(x) + len(y.s)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Map(t1.DoAllFunc(), func(j Int32) float64 { return float64(j) }) {
			i += int(math.Exp(x))
		}
		for x := range xiter.Map(t2.DoAllFunc(), func(j Int32) float64 { return float64(j) }) {
			i += int(math.Exp(x))
		}
	}
}

func BenchmarkFilter(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Filter(t1.DoAllFunc(), func(j Int32) bool { return j&1 == 0 }) {
			i += int(x)
		}
		for x := range xiter.Filter(t2.DoAllFunc(), func(j Int32) bool { return j&1 == 0 }) {
			i += int(x)
		}
	}
}

func BenchmarkLimit(b *testing.B) {
	b.ReportAllocs()
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
}

func BenchmarkReduce(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		i += xiter.Reduce(t1.DoAllFunc(), 0, func(total int, v Int32) int { return total + int(v) })
		i += xiter.Reduce(t2.DoAllFunc(), 0, func(total int, v Int32) int { return total + int(v) })
	}
}

func BenchmarkFold(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		i += int(xiter.Fold(t1.DoAllFunc(), func(total Int32, v Int32) Int32 { return total + v }))
		i += int(xiter.Fold(t2.DoAllFunc(), func(total Int32, v Int32) Int32 { return total + v }))
	}
}

func BenchmarkSkip(b *testing.B) {
	b.ReportAllocs()
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
}

func BenchmarkSkipSpecialized(b *testing.B) {
	b.ReportAllocs()
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
}

func BenchmarkSkipMethod(b *testing.B) {
	b.ReportAllocs()
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
}

func BenchmarkConcat(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Concat(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x)
		}
	}
}

func BenchmarkConcatRSC(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range Concat(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x)
		}
	}
}

func BenchmarkMergeFunc(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.MergeFunc(t1.DoAllFunc(), t2.DoAllFunc(), compare) {
			i += int(x)
		}
	}
}

func BenchmarkMergeFuncSpecialized(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range MergeFuncSpecialized(t1.DoAllFunc(), t2.DoAllFunc(), compare) {
			i += int(x)
		}
	}
}

func BenchmarkZip(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range xiter.Zip(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
}

func BenchmarkZipSpecialized(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range ZipSpecialized(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
}

func BenchmarkZipSpecializedND(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range ZipSpecializedND(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
}

func BenchmarkZipSpecialized1Pull(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		for x := range ZipSpecialized1Pull(t1.DoAllFunc(), t2.DoAllFunc()) {
			i += int(x.V1)
			i += int(x.V2)
		}
	}
}

func BenchmarkEqual(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		if !xiter.Equal(t1.DoAllFunc(), t1.DoAllFunc()) {
			panic("equal failure")
		}
		if xiter.Equal(t1.DoAllFunc(), t2.DoAllFunc()) {
			panic("unequal failure")
		}
	}
}

func TestZips(t *testing.T) {
	if !xiter.EqualFunc(
		xiter.Zip(t1.DoAllFunc(), t2.DoAllFunc()),
		ZipSpecialized1Pull(t1.DoAllFunc(), t2.DoAllFunc()),
		func(z1 xiter.Zipped[Int32, Int32], z2 Zipped) bool {
			return z1.V1 == z2.V1 && z1.V2 == z2.V2 && z1.OK1 == z2.OK1 && z1.OK2 == z2.OK2
		}) {
		t.Error("wanted equal 1")
	}
	if !xiter.EqualFunc(
		xiter.Zip(t1.DoAllFunc(), xiter.Skip(t2.DoAllFunc(), 1)),
		ZipSpecialized1Pull(t1.DoAllFunc(), Skip(t2.DoAllFunc(), 1)),
		func(z1 xiter.Zipped[Int32, Int32], z2 Zipped) bool {
			return z1.V1 == z2.V1 && z1.V2 == z2.V2 && z1.OK1 == z2.OK1 && z1.OK2 == z2.OK2
		}) {
		t.Error("wanted equal 2")
	}
	if !xiter.EqualFunc(
		xiter.Zip(xiter.Skip(t1.DoAllFunc(), 1), t2.DoAllFunc()),
		ZipSpecialized1Pull(Skip(t1.DoAllFunc(), 1), t2.DoAllFunc()),
		func(z1 xiter.Zipped[Int32, Int32], z2 Zipped) bool {
			return z1.V1 == z2.V1 && z1.V2 == z2.V2 && z1.OK1 == z2.OK1 && z1.OK2 == z2.OK2
		}) {
		t.Error("wanted equal 3")
	}
}

// local copy, specialized
type Seq func(yield func(Int32) bool) bool

// from proposal itself, does this avoid allocations?
func Concat[V any](seqs ...xiter.Seq[V]) xiter.Seq[V] {
	return func(yield func(V) bool) bool {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return false
				}
			}
		}
		return true
	}
}

func MergeFuncSpecialized(seq1, seq2 Seq, compare func(Int32, Int32) int) Seq {
	return func(yield func(Int32) bool) bool {
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
					return false
				}
				v1, ok1 = p1()
			case !ok1 || c > 0:
				if !yield(v2) {
					return false
				}
				v2, ok2 = p2()
			default:
				if !yield(v1) || !yield(v2) {
					return false
				}
				v1, ok1 = p1()
				v2, ok2 = p2()
			}
		}

		return false
	}
}

// Zipped holds values from an iteration of a Seq returned by [Zip].
type Zipped struct {
	V1  Int32
	OK1 bool

	V2  Int32
	OK2 bool
}

// Zip returns a new Seq that yields the values of seq1 and seq2
// simultaneously.
func ZipSpecialized(seq1 Seq, seq2 Seq) xiter.Seq[Zipped] {
	return func(yield func(Zipped) bool) bool {
		p1, stop := Pull(seq1)
		defer stop()
		p2, stop := Pull(seq2)
		defer stop()

		for {
			var val Zipped
			val.V1, val.OK1 = p1()
			val.V2, val.OK2 = p2()
			if (!val.OK1 && !val.OK2) || !yield(val) {
				return false
			}
		}
	}
}

// Zip returns a new Seq that yields the values of seq1 and seq2
// simultaneously.
func ZipSpecializedND(seq1 Seq, seq2 Seq) xiter.Seq[Zipped] {
	return func(yield func(Zipped) bool) bool {
		p1, stop1 := PullND(seq1)
		p2, stop2 := PullND(seq2)

		for {
			var val Zipped
			val.V1, val.OK1 = p1()
			val.V2, val.OK2 = p2()
			if (!val.OK1 && !val.OK2) || !yield(val) {
				stop2()
				stop1()
				return false
			}
		}
	}
}

// Zip returns a new Seq that yields the values of seq1 and seq2
// simultaneously.
func ZipSpecialized1Pull(seq1 Seq, seq2 Seq) xiter.Seq[Zipped] {
	return func(body func(Zipped) bool) bool {
		p2, stop2 := Pull(seq2)
		defer stop2()
		done := false
		for v := range seq1 {
			var val Zipped
			val.V1, val.OK1 = v, true
			val.V2, val.OK2 = p2()
			if !body(val) {
				done = true
				break
			}
		}
		if done {
			return false
		}
		// seq1 is exhausted
		for v2, ok2 := p2(); ok2; v2, ok2 = p2() {
			var v1 Int32
			var val Zipped
			val.V1, val.OK1 = v1, false
			val.V2, val.OK2 = v2, true
			if !body(val) {
				return false
			}
		}
		return true
	}
}

func Skip(forall Seq, n int) Seq {
	return func(body func(Int32) bool) bool {
		return forall(func(v Int32) bool {
			if n > 0 {
				n--
				return true
			}
			return body(v)
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

func PullND(seq Seq) (func() (Int32, bool), func()) {
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
