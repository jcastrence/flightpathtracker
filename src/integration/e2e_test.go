package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	// Must match running port for e2e tests
	calculateURL = "http://localhost:8080/calculate"
)

type EndToEndSuite struct {
	suite.Suite
}

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

func (s *EndToEndSuite) TestCalculateFlightPathEndpoint() {
	testCases := map[string]struct {
		jsonBody         [][]string
		expectedResult   string
		expectedErrorMsg string
	}{
		"happy-path-1": {
			jsonBody: [][]string{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
			expectedResult: `["SFO", "EWR"]`,
		},
		"happy-path-2": {
			jsonBody: [][]string{
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
			expectedResult: `["SFO", "EWR"]`,
		},
		"too-many-codes": {
			jsonBody: [][]string{
				{"IND", "EWR", "DEN"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
			expectedErrorMsg: `Bad flight input [IND EWR DEN]: Flights must be represented as a JSON string of size 2`,
		},
		"too-little-codes": {
			jsonBody: [][]string{
				{"IND"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
			expectedErrorMsg: `Bad flight input [IND]: Flights must be represented as a JSON string of size 2`,
		},
		"bad-codes-1": {
			jsonBody: [][]string{
				{"IND", "EWR"},
				{"SFO", "aTL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
			expectedErrorMsg: `Bad flight input [SFO aTL]: Flight elements must consist of 3 uppercase [A-Z] characters`,
		},
		"bad-codes-2": {
			jsonBody: [][]string{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"G3O", "IND"},
				{"ATL", "GSO"},
			},
			expectedErrorMsg: `Bad flight input [G3O IND]: Flight elements must consist of 3 uppercase [A-Z] characters`,
		},
		"bad-codes-3": {
			jsonBody: [][]string{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSOX"},
			},
			expectedErrorMsg: `Bad flight input [ATL GSOX]: Flight elements must consist of 3 uppercase [A-Z] characters`,
		},
		"same-codes": {
			jsonBody: [][]string{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"IND", "IND"},
				{"ATL", "GSO"},
			},
			expectedErrorMsg: `Bad flight input [IND IND]: Flight elements cannot be the same`,
		},
	}

	for testName, testCase := range testCases {
		log.Printf("Running %s", testName)
		// Prepare request
		b, err := json.Marshal(testCase.jsonBody)
		s.NoError(err)

		req, err := http.NewRequest("GET", calculateURL, bytes.NewBuffer(b))
		s.NoError(err)

		req.Header.Set("Content-Type", "application/json")

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		s.NoError(err)
		defer resp.Body.Close()

		b, err = io.ReadAll(resp.Body)
		s.NoError(err)

		// Check response
		switch {
		// Happy path cases
		case testCase.expectedErrorMsg == "":
			s.Equal(http.StatusOK, resp.StatusCode)
			s.JSONEq(testCase.expectedResult, string(b))
		// Error cases
		default:
			s.Equal(http.StatusBadRequest, resp.StatusCode)
			s.Contains(string(b), testCase.expectedErrorMsg)
		}
	}

}
