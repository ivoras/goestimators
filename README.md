# HyperLogLog (and LogLog and SuperLogLog) for Go

This family of algorithms implement fast approximate counting of arbitrary data (here implemented as binary keys), with a trivially small memory footprint. Their workings are similar enough that the only point of difference is the summarisation step, where the approximation is calculated from internal data. The same internal data can be used for both LogLog, SuperLogLog and HyperLogLog approximation.

The number of LogLog buckets (which affects accuracy, but soon saturates to diminishing returns) is adjustable, with 1024 being the recommended default, and 32 being the smallest practical size. This number must be an integral power of 2.

This library was implemented from scratch, by interpreting the following papers and tutorials:

* http://www.ens-lyon.fr/LIP/Arenaire/SYMB/teams/algo/algo4.pdf
* http://moderndescartes.com/essays/hyperloglog
* https://research.neustar.biz/2012/10/25/sketch-of-the-day-hyperloglog-cornerstone-of-a-big-data-infrastructure/
* http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf

Any implementation errors are probably my own.

## Usage

	ll, _ := NewLogLog(1024)
	var buf [8]byte
	const NumEntries = 100000
	for i := 0; i < NumEntries; i++ {
		n, err := rand.Read(buf[:])
		ll.Observe(buf[:])  // Observes the buffer, i.e. updates internal representation from it
	}
	est := ll.Estimate()    // Returns the estimated number of observed buffers
	fmt.Println("Estimate of a set of 100k random entries: ", est)
