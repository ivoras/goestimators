# Estimator algorithms for Go: (Hyper|Super)?LogLog, Bloom filter

## HyperLogLog (and LogLog and SuperLogLog) for Go

This algorithms in this family implement fast approximate counting of arbitrary data (here implemented as binary keys), with a trivially small memory footprint. Their workings are similar enough that the only point of difference is the summarisation step, where the approximation is calculated from internal data. The same internal data can be used for both LogLog, SuperLogLog and HyperLogLog approximation.

The number of LogLog buckets (which affects accuracy, but soon saturates to diminishing returns) is adjustable, with 1024 being the recommended default, and 32 being the smallest practical size. This number must be an integral power of 2.

This library was implemented from scratch, by interpreting the following papers and tutorials:

* http://www.ens-lyon.fr/LIP/Arenaire/SYMB/teams/algo/algo4.pdf
* http://moderndescartes.com/essays/hyperloglog
* https://research.neustar.biz/2012/10/25/sketch-of-the-day-hyperloglog-cornerstone-of-a-big-data-infrastructure/
* http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf

Any implementation errors are probably my own.

This library uses the [Murmur3](https://en.wikipedia.org/wiki/MurmurHash) hash function implementation from [github.com/spaolacci/murmur3](https://github.com/spaolacci/murmur3).

### Example

	ll, _ := NewLogLog(1024)
	var buf := make([]byte, 8)
	for i := 0; i < 100000; i++ {
		rand.Read(buf)
		ll.Observe(buf)  // Observes the buffer, i.e. updates internal representation from it
	}
	est := ll.Estimate() // Returns the estimated number of unique observed buffers
	fmt.Println("Estimate of a set of 100k random entries: ", est)

This example uses the plain old LogLog estimation algorithm, implemented in the `Estimate()` function. All three algorithms can be used from the same observation results:

* `Estimate()` - LogLog (fastest)
* `SuperEstimate()` - SuperLogLog
* `HyperEstimate()` - HyperLogLog (slowest)

### Performance

Measured on an i5-5200U, the benchmark results (using 1024 buckets) are:

    BenchmarkLogLog-4                        2000000               712 ns/op
    BenchmarkSuperLogLog-4                   2000000               763 ns/op
    BenchmarkHyperLogLog-4                    500000              2367 ns/op
    BenchmarkObservationLogLog-4            20000000                61.6 ns/op

Since HyperLogLog improves accuracy only slightly compared to LogLog and SuperLogLog, users should decide for themselves if they can accept the hit in performance.

### Accuracy 

A typical run of the test functions (using 1024 buckets) looks like this:

    LogLogRandom_1M 918170
    SuperLogLogRandom_1M 932190
    HyperLogLogRandom_1M 1059335

These are the estimates for 1M of sequential 8-byte buffers by the respective algorithms. YMMV.

# Bloom filter

"Bloom filters" test for the *absence* of an element in a large set. If a Bloom filter returns negative, the element certainly isn't in the set. If it returns positive, the element *might be* in the set. Memory usage for this is trivially small.

This algorithm was implemented from scratch, mostly by reading the [Wikipedia article](https://en.wikipedia.org/wiki/Bloom_filter). It uses the Murmur3 function for the 64-bit hash implementation.

For performance reasons, there are only two supported Bloom filter sizes: 256-bit (which is 32 bytes), and 65536-bit (8 KiB). The 256-bit Bloom filter implements two hash functions (k=2), and the 65536-bit one implements four hash functions (k=4).
