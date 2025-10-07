package api

import (
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sealdice-ce/sealdice-ce/core/dice"
)

var startTime = time.Now().Unix()

type VersionDetail struct {
	Major         uint64 `json:"major"`
	Minor         uint64 `json:"minor"`
	Patch         uint64 `json:"patch"`
	Prerelease    string `json:"prerelease"`
	BuildMetaData string `json:"buildMetaData"`
}

func preInfo(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"testMode": dm.JustForTest,
	})
}

func baseInfo(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var versionNew string
	var versionNewNote string
	var versionNewCode int64
	if dm.AppVersionOnline != nil {
		versionNew = dm.AppVersionOnline.VersionLatestDetail
		versionNewNote = dm.AppVersionOnline.VersionLatestNote
		versionNewCode = dm.AppVersionOnline.VersionLatestCode
	}

	getName := func() string {
		defer func() {
			// 防止报错
			_ = recover()
		}()

		ctx := &dice.MsgContext{Dice: myDice, EndPoint: nil, Session: myDice.ImSession}
		return dice.DiceFormatTmpl(ctx, "核心:骰子名字")
	}
	extraTitle := getName()

	versionDetail := VersionDetail{
		Major:         dice.VERSION.Major(),
		Minor:         dice.VERSION.Minor(),
		Patch:         dice.VERSION.Patch(),
		Prerelease:    dice.VERSION.Prerelease(),
		BuildMetaData: dice.VERSION.Metadata(),
	}

	return c.JSON(http.StatusOK, struct {
		AppName        string        `json:"appName"`
		AppChannel     string        `json:"appChannel"`
		Version        string        `json:"version"`
		VersionSimple  string        `json:"versionSimple"`
		VersionDetail  VersionDetail `json:"versionDetail"`
		VersionNew     string        `json:"versionNew"`
		VersionNewNote string        `json:"versionNewNote"`
		VersionCode    int64         `json:"versionCode"`
		VersionNewCode int64         `json:"versionNewCode"`
		MemoryAlloc    uint64        `json:"memoryAlloc"`
		Uptime         int64         `json:"uptime"`
		MemoryUsedSys  uint64        `json:"memoryUsedSys"`
		ExtraTitle     string        `json:"extraTitle"`
		OS             string        `json:"OS"`
		Arch           string        `json:"arch"`
		JustForTest    bool          `json:"justForTest"`
		ContainerMode  bool          `json:"containerMode"`
	}{
		AppName:        dice.APPNAME,
		AppChannel:     dice.APP_CHANNEL,
		Version:        dice.VERSION.String(),
		VersionSimple:  dice.VERSION_MAIN + dice.VERSION_PRERELEASE,
		VersionDetail:  versionDetail,
		VersionNew:     versionNew,
		VersionNewNote: versionNewNote,
		VersionCode:    dice.VERSION_CODE,
		VersionNewCode: versionNewCode,
		MemoryAlloc:    m.Alloc,
		MemoryUsedSys:  m.Sys - m.HeapReleased,
		Uptime:         time.Now().Unix() - startTime,
		ExtraTitle:     extraTitle,
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		JustForTest:    dm.JustForTest,
		ContainerMode:  dm.ContainerMode,
	})
}

func hello2(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
