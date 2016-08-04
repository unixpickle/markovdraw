package main

import "github.com/unixpickle/markovchain"

func BuildChain(paths []Path) *markovchain.Chain {
	segChan := make(chan Segment)
	sampleChan := make(chan markovchain.State)
	go func() {
		for _, path := range paths {
			segments := SegmentPath(path)
			for _, seg := range segments {
				segChan <- seg
			}
		}
		close(segChan)
	}()
	go func() {
		var tuple SegmentTuple
		for seg := range segChan {
			tuple = append(tuple, seg)
			if len(tuple) == 2 {
				sampleChan <- tuple
				tuple = nil
			}
		}
		close(sampleChan)
	}()
	return markovchain.NewChainChan(sampleChan)
}

func SampleChain(chain *markovchain.Chain, count int) []Segment {
	var state markovchain.State
	chain.Iterate(func(s *markovchain.StateTransitions) bool {
		state = s.State
		return false
	})

	// Monto carlo way of going to a random state.
	for i := 0; i < 100; i++ {
		newStart := chain.Lookup(state).Sample()
		if newStart == nil {
			break
		}
		state = newStart
	}

	var res []Segment
	for i := 0; i < count; i++ {
		state = chain.Lookup(state).Sample()
		if state == nil {
			break
		}
		res = append(res, state.(SegmentTuple)...)
	}
	return res
}
