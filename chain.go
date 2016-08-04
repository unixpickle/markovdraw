package main

import "github.com/unixpickle/markovchain"

const memorySize = 2

func BuildChain(paths []Path) *markovchain.Chain {
	stateChan := make(chan markovchain.State)
	go func() {
		var state SegmentTuple
		for _, path := range paths {
			segments := SegmentPath(path)
			segments = append([]Segment{Segment{}}, segments...)
			for _, seg := range segments {
				state = append(state, seg)
				if len(state) > memorySize {
					state = state[1:]
				}
				stateChan <- state
			}
		}
		stateChan <- SegmentTuple{Segment{}}
		close(stateChan)
	}()
	return markovchain.NewChainChan(stateChan)
}

func SampleChain(chain *markovchain.Chain) []Segment {
	var res []Segment
	var state markovchain.State = SegmentTuple{Segment{}}
	for {
		state = chain.Lookup(state).Sample()
		tuple := state.(SegmentTuple)
		last := tuple[len(tuple)-1]
		if (last == Segment{}) {
			break
		}
		res = append(res, last)
	}
	return res
}
