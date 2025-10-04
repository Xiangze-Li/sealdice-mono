package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monaco-io/request"
)

func onebotTool(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(200, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Port int64 `form:"port" json:"port"`
	}{}
	_ = c.Bind(&v)

	port := int64(13325)
	if v.Port != 0 {
		port = v.Port
	}

	errText := ""
	ip, err := socksOpen(myDice, port)
	if err != nil {
		errText = err.Error()
	}

	resp := c.JSON(200, map[string]interface{}{
		"ok":      err == nil,
		"ip":      ip,
		"errText": errText,
	})
	return resp
}

func getToken(c echo.Context) error {
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	key := c.QueryParam("key")
	req := request.Client{
		URL:    fmt.Sprintf("https://gocqhelper.sealdice.com/code_query?key=%s", key),
		Method: "GET",
	}

	resp := req.Send()
	if resp.OK() {
		text := resp.String()
		return c.String(http.StatusOK, text)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}
