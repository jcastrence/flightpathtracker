package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReduceFlightPath(t *testing.T) {
	testCases := map[string]struct {
		sampleFlights     []Flight
		expectedResult    Flight
		possibleErrorMsgs []string
	}{
		"SFO-EWR-1": {
			sampleFlights: []Flight{
				{"SFO", "EWR"},
			},
			expectedResult: Flight{"SFO", "EWR"},
		},
		"SFO-EWR-2": {
			sampleFlights: []Flight{
				{"ATL", "EWR"},
				{"SFO", "ATL"},
			},
			expectedResult: Flight{"SFO", "EWR"},
		},
		"SFO-EWR-3": {
			sampleFlights: []Flight{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
			expectedResult: Flight{"SFO", "EWR"},
		},
		"SFO-EWR-4": {
			sampleFlights: []Flight{
				{"DCA", "JFK"},
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IAD"},
				{"ATL", "GSO"},
				{"JFK", "IND"},
				{"IAD", "DCA"},
			},
			expectedResult: Flight{"SFO", "EWR"},
		},
		"SFO-EWR-5": {
			sampleFlights: []Flight{
				{"DCA", "JFK"},
				{"SFO", "DEN"},
				{"DEN", "CLT"},
				{"DXB", "IAD"},
				{"CLT", "ATL"},
				{"ATL", "LHR"},
				{"IND", "EWR"},
				{"LHR", "GSO"},
				{"JFK", "IND"},
				{"IAD", "DCA"},
				{"GSO", "ORD"},
				{"ORD", "DXB"},
			},
			expectedResult: Flight{"SFO", "EWR"},
		},
		"JFK-LAX-1": {
			sampleFlights: []Flight{
				{"JFK", "ATL"},
				{"ATL", "LAX"},
			},
			expectedResult: Flight{"JFK", "LAX"},
		},
		"JFK-LAX-2": {
			sampleFlights: []Flight{
				{"ATL", "LAX"},
				{"JFK", "ATL"},
			},
			expectedResult: Flight{"JFK", "LAX"},
		},
		// Should fail, multiple possible sources to LAX
		"JFK-LAX-3": {
			sampleFlights: []Flight{
				{"JFK", "LAX"},
				{"ATL", "LAX"},
			},
			// Any of these are possible, hashmap keys in Go are unordered
			possibleErrorMsgs: []string{
				"At least two possible sources: JFK, ATL",
				"At least two possible sources: ATL, JFK",
			},
		},
		// Should fail, multiple possible destinations from JFK
		"JFK-LAX-4": {
			sampleFlights: []Flight{
				{"ATL", "LAX"},
				{"JFK", "DEN"},
				{"JFK", "ATL"},
			},
			// Any of these are possible, hashmap keys in Go are unordered
			possibleErrorMsgs: []string{
				"Source JFK has at least two destinations: DEN, ATL",
				"Source JFK has at least two destinations: ATL, DEN",
			},
		},
		// Should fail, disconnected paths
		"JFK-LAX-5": {
			sampleFlights: []Flight{
				{"ATL", "DEN"},
				{"DCA", "LAX"},
				{"JFK", "ATL"},
			},
			// Any of these are possible, hashmap keys in Go are unordered
			possibleErrorMsgs: []string{
				"At least two possible sources: JFK, DCA",
				"At least two possible sources: DCA, JFK",
			},
		},
		// Should fail, cycle path
		"JFK-LAX-6": {
			sampleFlights: []Flight{
				{"JFK", "ATL"},
				{"ATL", "LAX"},
				{"LAX", "JFK"},
			},
			// Any of these are possible, hashmap keys in Go are unordered
			possibleErrorMsgs: []string{
				"Destination JFK has ambiguous source (multiple sources or path cycle)",
				"Destination ATL has ambiguous source (multiple sources or path cycle)",
				"Destination LAX has ambiguous source (multiple sources or path cycle)",
			},
		},
		// Should fail, repeated flight
		"JFK-LAX-7": {
			sampleFlights: []Flight{
				{"ATL", "LAX"},
				{"JFK", "ATL"},
				{"ATL", "LAX"},
			},
			// Any of these are possible, hashmap keys in Go are unordered
			possibleErrorMsgs: []string{
				"Repeated entry: [ATL, LAX]",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actualResult, err := ReduceFlightPath(testCase.sampleFlights)
			switch {
			// Happy path cases
			case testCase.possibleErrorMsgs == nil:
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, actualResult)
			// Error cases
			default:
				atLeastOneMatch := false
				for _, eMsg := range testCase.possibleErrorMsgs {
					if eMsg == err.Error() {
						atLeastOneMatch = true
					}
				}
				require.True(t, atLeastOneMatch)
				require.Equal(t, Flight{}, actualResult)
			}
		})
	}
}
