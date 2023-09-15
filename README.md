# iter_test

This is for benchmarking and performance research of an initial implementation of function iterators and [DeedleFake's xiter support package](https://github.com/DeedleFake/xiter).  Though that code is advertised as "not for actual usage", it works well and with a slightly modified compiler generally has high performance; I am very grateful that they wrote it.  The goal is to figure out
 - what changes does the compiler need to run this code well
 - what changes does the library need to run this code well
 - in some cases, tweaking the support function implementations for better performance.
 - are there any performance penalties associated with generics in this code?

To run this code, you will need [CL 510541](https://go.dev/cl/510541).  There are incompatible changes in the pipeline that will require updates to the support library and the benchmarks at some point, so be aware, at some point there will be some instability, but these should generally not matter to these benchmarks.

To obtain best performance with the current inliner, you'll need to increase the inlining threshold (in `src/cmd/compile/internal/inline/inl.go`) from 80 to 125.  One goal of this work is to create a test case for revisions to the inliner heuristics to improve performance of code containing functions that use (call) their function parameters.

Results so far:
 - so far there have been NO proportional-to-sequence-length allocations
 - very happily, the allocation count is often zero
 - the inliner will need adjusting
 - `sync.OnceFunc` performs "too many" allocations (5)
 - `Pull` is generally very expensive
 - the two-sequence support functions can be (tediously) rewritten to use only one call to `Pull`

I expect that this will ultimately be converted to a test for the compiler and library to ensure that all these functions do not accidentally increase their allocation count.

Results as of 2023-09-15:
```
go test -bench B -memprofile=m.prof .
goos: darwin
goarch: amd64
pkg: github.com/dr2chase/iter
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkNothing-8                	1000000000	         0.3001 ns/op	       0 B/op	       0 allocs/op
BenchmarkOldSlice-8               	19501254	        59.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkOldCount-8               	18576810	        62.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapKeys-8                	 2992586	       425.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapValues-8              	 2683900	       421.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkOfMap-8                  	 3583351	       346.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkToPairOfMap-8            	 3409785	       336.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkFromPair-8               	 4399246	       272.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkBytes-8                  	20255217	        53.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkRunes-8                  	 7823437	       144.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringSplitEmpty-8       	 1964611	       627.0 ns/op	     112 B/op	      28 allocs/op
BenchmarkStringSplit-8            	 4251924	       294.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkGenerate-8               	20647726	        58.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkOf-8                     	13209528	        89.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkOfSlice-8                	13939190	        82.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkDoAll-8                  	 8588011	       136.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDoAll2-8                 	 8801606	       134.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMap-8                    	 2790885	       458.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkFilter-8                 	 5687697	       204.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkLimit-8                  	 6557757	       184.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkReduce-8                 	 6160969	       191.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkFold-8                   	 6644575	       176.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkSkip-8                   	 3298513	       356.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkSkipSpecialized-8        	 3249984	       351.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkSkipMethod-8             	 2168030	       558.1 ns/op	     128 B/op	       8 allocs/op
BenchmarkConcat-8                 	 5392609	       227.0 ns/op	      48 B/op	       2 allocs/op
BenchmarkConcatRSC-8              	 3439323	       341.5 ns/op	     128 B/op	       6 allocs/op
BenchmarkMergeFunc-8              	   63649	     18738 ns/op	     944 B/op	      24 allocs/op
BenchmarkMergeFuncSpecialized-8   	   62031	     19146 ns/op	     864 B/op	      24 allocs/op
BenchmarkZip-8                    	   53372	     22953 ns/op	     944 B/op	      24 allocs/op
BenchmarkZipSpecialized-8         	   52080	     22778 ns/op	     864 B/op	      24 allocs/op
BenchmarkZipSpecializedND-8       	   53809	     22460 ns/op	     640 B/op	      14 allocs/op
BenchmarkZipSpecialized1Pull-8    	  113126	     10099 ns/op	     432 B/op	      12 allocs/op
BenchmarkEqual-8                  	   43077	     27639 ns/op	    1888 B/op	      48 allocs/op
PASS
```


