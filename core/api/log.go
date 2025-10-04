package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func logFetchAndClear(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	info := c.JSON(http.StatusOK, myDice.LogWriter.Items)
	// myDice.LogWriter.Items = myDice.LogWriter.Items[:0]
	return info
}
