// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package align_test

import (
	"fmt"
	"testing"

	"github.com/dr2chase/xiter"
)

var sink int

//go:noinline
func pad() {}

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

// BenchmarkSliceOpt measures a hand-optimized generic slice iteration, updating a local.
func bNSO0(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

// BenchmarkSliceOpt measures a hand-optimized generic slice iteration, updating a local.
func bNSO1(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

// BenchmarkSliceOpt measures a hand-optimized generic slice iteration, updating a local.
func bNSO2(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

// BenchmarkSliceOpt measures a hand-optimized generic slice iteration, updating a local.
func bNSO3(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSO4(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSO5(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSO6(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSO7(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSO8(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSO9(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOA(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOB(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOC(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOD(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOE(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOF(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOG(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOH(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOI(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
			pad()
		}
	}
	sink += i
}

func bNSOJ(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
			pad()
		}
	}
	sink += i
}

func bNSOK(b *testing.B, cs []int) {
	slice := cs
	b.ReportAllocs()
	i := 0
	for range b.N {
		for _ = range OfSliceOpt(slice) {
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			pad()
			for y := range OfSliceOpt(slice) {
				i += y
			}
		}
	}
	sink += i
}

var cs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
}

func BenchmarkNestedSliceOptX(b *testing.B) {
	benches := []func(*testing.B, []int){bNSO0, bNSO1, bNSO2, bNSO3, bNSO4, bNSO5, bNSO6, bNSO7, bNSO8, bNSO9, bNSOA, bNSOB, bNSOC, bNSOD, bNSOE, bNSOF, bNSOG, bNSOH, bNSOI, bNSOJ, bNSOK}

	for i, bm := range benches {
		b.Run(fmt.Sprintf("padding=%d", i), func(b *testing.B) {
			bm(b, cs)
		})
	}
}
