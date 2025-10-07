package api

import (
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/sealdice-ce/sealdice-ce/core/dice"
)

type fStopEcho struct {
	Key string `json:"key"`
}

// 这个函数是用来强制停止程序的，只有在Android上才能使用，其他平台会返回403
func forceStop(c echo.Context) error {
	if runtime.GOOS != "android" {
		return c.JSON(http.StatusForbidden, nil)
	}
	// this is a dangerous api, so we need to check the key
	haskey := false
	for _, s := range os.Environ() {
		if strings.HasPrefix(s, "FSTOP_KEY=") {
			key := strings.Split(s, "=")[1]
			v := fStopEcho{}
			err := c.Bind(&v)
			if err != nil {
				return c.JSON(http.StatusBadRequest, nil)
			}
			if v.Key == key {
				haskey = true
				break
			} else {
				return c.JSON(http.StatusForbidden, nil)
			}
		}
	}
	if !haskey {
		return c.JSON(http.StatusForbidden, nil)
	}
	defer func() {
		// Same with main.go `cleanUpCreate()` 由于无法导入 main.go 中的函数，所以这里直接复制过来了
		log := myDice.Logger
		log.Info("程序即将退出，进行清理……")
		err := recover()
		log.Errorf("异常: %v\n堆栈: %v", err, string(debug.Stack()))
		diceManager := dm

		for _, i := range diceManager.Dice {
			if i.IsAlreadyLoadConfig {
				i.Config.BanList.SaveChanged(i)
				err = i.AttrsManager.CheckForSave()
				if err != nil {
					log.Errorf("异常: %v", err)
				}
				i.Save(true)
				for _, j := range i.ExtList {
					if j.Storage != nil {
						err := j.StorageClose()
						if err != nil {
							log.Errorf("异常: %v\n堆栈: %v", err, string(debug.Stack()))
						}
					}
				}
				i.IsAlreadyLoadConfig = false
			}
		}

		for _, i := range diceManager.Dice {
			d := i
			d.DBOperator.Close()
		}

		// 清理gocqhttp
		for _, i := range diceManager.Dice {
			if i.ImSession != nil && i.ImSession.EndPoints != nil {
				for _, j := range i.ImSession.EndPoints {
					dice.BuiltinQQServeProcessKill(i, j)
				}
			}
		}

		if diceManager.Help != nil {
			diceManager.Help.Close()
		}
		if diceManager.IsReady {
			diceManager.Save()
		}
		if diceManager.Cron != nil {
			diceManager.Cron.Stop()
		}
		log.Info("清理完成，程序即将退出")
		os.Exit(0)
	}()
	return c.JSON(http.StatusOK, nil)
}
