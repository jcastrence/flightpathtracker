package router

import (
	"github.com/jcastrence/flightpathtracker/src/handlers"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	// Single GET route for calculating flight path
	e.GET("/calculate", handlers.CalculateFlightPath)

	return e
}
