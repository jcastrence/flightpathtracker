package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/jcastrence/flightpathtracker/src/services"
	"github.com/labstack/echo/v4"
)

// This type is used instead of services.Flight for more strict input checking
type flights [][]string

func CalculateFlightPath(c echo.Context) error {
	defer c.Request().Body.Close()

	// Decoding json response
	fls := flights{}
	err := json.NewDecoder(c.Request().Body).Decode(&fls)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to decode JSON from request body: %v\n", err)
		log.Print(errMsg)
		return c.String(http.StatusInternalServerError, errMsg)
	}

	// Input checking
	for _, fl := range fls {
		// Size check (this check would miss lists of size > 2 if using services.Flight type since json Decoder truncates input)
		if len(fl) != 2 {
			errMsg := fmt.Sprintf("Bad flight input %v: Flights must be represented as a JSON string of size 2\n", fl)
			log.Print(errMsg)
			return c.String(http.StatusBadRequest, errMsg)
		}
		// Airport code string check
		re := regexp.MustCompile(`^[A-Z]{3}$`)
		if !re.MatchString(fl[0]) || !re.MatchString(fl[1]) {
			errMsg := fmt.Sprintf("Bad flight input %v: Flight elements must consist of 3 uppercase [A-Z] characters\n", fl)
			log.Print(errMsg)
			return c.String(http.StatusBadRequest, errMsg)
		}
		// Same code string check
		if fl[0] == fl[1] {
			errMsg := fmt.Sprintf("Bad flight input %v: Flight elements cannot be the same\n", fl)
			log.Print(errMsg)
			return c.String(http.StatusBadRequest, errMsg)
		}
	}

	// Calculation
	flightPath, err := services.ReduceFlightPath(inputTypeConversion(fls))
	if err != nil {
		log.Print(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, flightPath)
}

// Helper
// Convert local flights type to services.[]Flight type
func inputTypeConversion(input flights) []services.Flight {
	output := make([]services.Flight, 0, len(input)) // Known capacity allows for faster appends
	for _, f := range input {
		output = append(output, services.Flight{f[0], f[1]})
	}
	return output
}
