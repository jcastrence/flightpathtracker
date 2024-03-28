package services

import (
	"fmt"
)

type Flight [2]string
type flightMap map[string]string

// Reduce a set of Flights to a single source and destination
func ReduceFlightPath(Flights []Flight) (Flight, error) {

	// Map sources to destinations in hashmap
	flMap := make(flightMap)
	for _, fl := range Flights {
		src, dst := fl[0], fl[1]
		// Handle cases where destination is ambiguous, i.e. [["ATL", "JFK"], ["ATL", "LAX"]]
		if flMap[src] != "" {
			// Check for repeated entries
			if flMap[src] == dst {
				return Flight{}, newRepeatedEntryError(src, dst)
			}
			return Flight{}, newAmbiguousDestinationError(src, flMap[src], dst)
		}

		flMap[src] = dst
	}

	// Determine Flight route using DFS-like approach, iterating each source of hashmap
	for src := range flMap {
		err := subReduce(src, flMap)
		if err != nil {
			return Flight{}, err
		}
	}

	// Fully reduced Flight path is marked as head, find and return head
	reducedFlight := Flight{}
	headCount := 0 // If more than one head found throw error
	for src := range flMap {
		if flMap.isHead(src) {
			if headCount > 0 {
				// Multiple heads, source is ambiguous
				return Flight{}, newMultipleSourcesError(reducedFlight[0], src)
			}
			reducedFlight = Flight{src, flMap[src][1:]}
			headCount++
		}
	}
	return reducedFlight, nil
}

// Helpers

// Reduce single source to final destination
func subReduce(src string, flMap flightMap) error {
	// Skip visited sources
	if flMap.isVisited(src) {
		return nil
	}

	// Search for final destination or current head
	curr := flMap[src]
	flMap.setVisited(src)
	for !flMap.isHead(src) {
		switch {
		// Final desination, path fully reduced
		case flMap[curr] == "":
			flMap.setHead(src)
		// Visited destination, source is ambiguous
		case flMap.isVisited(curr):
			return newAmbiguousSourceError(flMap[curr][1:])
		// Previous head found, dominate it
		case flMap.isHead(curr):
			flMap[src] = flMap[curr]
			flMap.setVisited(curr)
		// new destination, reduce path
		default:
			next := flMap[curr]
			flMap.setVisited(curr)
			flMap[src] = flMap[curr]
			curr = next
		}
	}

	return nil
}

// Visited and head state control
func (f flightMap) setVisited(src string) {
	if f.isHead(src) {
		f[src] = f[src][1:]
	}
	f[src] = fmt.Sprintf("-%s", f[src])
}

func (f flightMap) isVisited(src string) bool {
	return string(f[src][0]) == "-"
}

func (f flightMap) setHead(src string) {
	f[src] = fmt.Sprintf("*%s", f[src][1:])
}

func (f flightMap) isHead(src string) bool {
	return string(f[src][0]) == "*"
}

// Error Types
type AmbiguousDestinationError struct {
	src, dst1, dst2 string
}

func newAmbiguousDestinationError(src, dst1, dst2 string) *AmbiguousDestinationError {
	return &AmbiguousDestinationError{
		src, dst1, dst2,
	}
}
func (e *AmbiguousDestinationError) Error() string {
	return fmt.Sprintf("Source %s has at least two destinations: %s, %s", e.src, e.dst1, e.dst2)
}

type AmbiguousSourceError struct {
	dst string
}

func newAmbiguousSourceError(dst string) *AmbiguousSourceError {
	return &AmbiguousSourceError{
		dst,
	}
}

func (e *AmbiguousSourceError) Error() string {
	return fmt.Sprintf("Destination %s has ambiguous source (multiple sources or path cycle)", e.dst)
}

type MultipleSourcesError struct {
	src1, src2 string
}

func newMultipleSourcesError(src1, src2 string) *MultipleSourcesError {
	return &MultipleSourcesError{
		src1, src2,
	}
}

func (e *MultipleSourcesError) Error() string {
	return fmt.Sprintf("At least two possible sources: %s, %s", e.src1, e.src2)
}

type RepeatedEntryError struct {
	src, dst string
}

func newRepeatedEntryError(src, dst string) *RepeatedEntryError {
	return &RepeatedEntryError{
		src, dst,
	}
}
func (e *RepeatedEntryError) Error() string {
	return fmt.Sprintf("Repeated entry: [%s, %s]", e.src, e.dst)
}
