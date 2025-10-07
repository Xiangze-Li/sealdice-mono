package api

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sealdice-ce/sealdice-ce/core/dice"
)

var (
	lastPrivateExecTime int64
	lastGroupExecTime   int64
)

func DiceExec(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID          string `form:"id"          json:"id"`
		Message     string `form:"message"`
		MessageType string `form:"messageType"`
	}{}
	err := c.Bind(&v)
	if err != nil {
		return c.JSON(400, "格式错误")
	}
	if v.Message == "" {
		return c.JSON(400, "格式错误")
	}
	timeNeed := int64(500)
	if dm.JustForTest {
		timeNeed = 80
	}

	now := time.Now().UnixMilli()
	userID := "UI:1001"
	messageType := "private"
	groupID := ""
	groupName := ""
	groupRole := ""
	if v.MessageType == "group" {
		userID = "UI:1002"
		messageType = "group"
		groupID = "UI-Group:2001"
		groupName = "UI-Group 2001"
		groupRole = "owner"
		if now-lastGroupExecTime < timeNeed {
			return c.JSON(400, "过于频繁")
		}
		lastGroupExecTime = now
	} else {
		if now-lastPrivateExecTime < timeNeed {
			return c.JSON(400, "过于频繁")
		}
		lastPrivateExecTime = now
	}

	msg := &dice.Message{
		MessageType: messageType,
		Message:     v.Message,
		Platform:    "UI",
		Sender: dice.SenderBase{
			Nickname:  "User",
			UserID:    userID,
			GroupRole: groupRole,
		},
		GroupID:   groupID,
		GroupName: groupName,
	}
	myDice.ImSession.Execute(myDice.UIEndpoint, msg, false)
	return c.JSON(200, "ok")
}

func DiceRecentMessage(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	pa := myDice.UIEndpoint.Adapter.(*dice.PlatformAdapterHTTP)
	defer func() {
		pa.RecentMessage = []dice.HTTPSimpleMessage{}
	}()
	return c.JSON(200, pa.RecentMessage)
}

func DiceAllCommand(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	var cmdLst []string
	for k := range myDice.CmdMap {
		cmdLst = append(cmdLst, k)
	}

	for _, i := range myDice.ExtList {
		for k := range i.CmdMap {
			cmdLst = append(cmdLst, k)
		}
	}
	sort.Sort(dice.ByLength(cmdLst))
	return c.JSON(200, cmdLst)
}

func DiceNewVersionUpload(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return Success(&c, Response{
			"testMode": true,
		})
	}

	if dm.UpdateSealdiceByFile == nil {
		return Error(&c, "骰子没有正确初始化，无法使用此功能", Response{})
	}

	if dm.ContainerMode {
		return Error(&c, "容器模式下禁止更新，请手动拉取最新镜像", Response{})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return Error(&c, err.Error(), Response{})
	}
	files := form.File["files"]
	if len(files) == 0 {
		return Error(&c, "参数错误", Response{})
	}

	myDice.Logger.Info("收到新版本骰子上传请求")

	file := files[0]
	src, err := file.Open()
	if err != nil {
		return Error(&c, err.Error(), Response{})
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	// TODO: 临时将逻辑写在这里，后续根据情况再调整
	fn := "./new_package"
	if runtime.GOOS == "windows" {
		fn += ".zip"
	} else {
		fn += ".tar.gz"
	}

	f2, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return Error(&c, err.Error(), Response{})
	}
	_, err = io.Copy(f2, src)
	if err != nil {
		return Error(&c, err.Error(), Response{})
	}

	f2.Close()

	myDice.Logger.Info("新版本骰子上传成功")

	_ = Success(&c, Response{"result": true})

	if !dm.UpdateSealdiceByFile(fn) {
		myDice.Logger.Error("更新骰子失败")
	}

	return nil
}
