package api

import (
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"sealdice-core/dice"
)

func doSignInGetSalt(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"salt": myDice.Parent.UIPasswordSalt,
	})
}

func checkSecurity(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	isPublicService := strings.HasPrefix(myDice.Parent.ServeAddress, "0.0.0.0") || myDice.Parent.ServeAddress == ":3211"
	isEmptyPassword := myDice.Parent.UIPasswordHash == ""
	return c.JSON(200, map[string]bool{
		"isOk": !isEmptyPassword || !isPublicService,
	})
}

func doSignIn(c echo.Context) error {
	v := struct {
		Password string `json:"password"`
	}{}

	err := c.Bind(&v)
	if err != nil {
		return c.JSON(400, nil)
	}

	generateToken := func() error {
		now := time.Now().Unix()
		head := hex.EncodeToString(Int64ToBytes(now))
		token := dice.RandStringBytesMaskImprSrcSB2(64) + ":" + head

		myDice.Parent.AccessTokens.Store(token, true)
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Parent.Save()
		return c.JSON(http.StatusOK, map[string]string{
			"token": token,
		})
	}

	if myDice.Parent.UIPasswordHash == "" {
		return generateToken()
	}

	if myDice.Parent.UIPasswordHash == v.Password {
		return generateToken()
	}

	return c.JSON(400, nil)
}
