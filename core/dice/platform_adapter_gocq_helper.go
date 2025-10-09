package dice

import (
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/acarl005/stripansi"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/sealdice-ce/sealdice-ce/core/logger"
	"github.com/sealdice-ce/sealdice-ce/core/utils"
	"github.com/sealdice-ce/sealdice-ce/core/utils/procs"
)

type deviceFile struct {
	Display      string         `json:"display"`
	Product      string         `json:"product"`
	Device       string         `json:"device"`
	Board        string         `json:"board"`
	Model        string         `json:"model"`
	FingerPrint  string         `json:"finger_print"`
	BootID       string         `json:"boot_id"`
	ProcVersion  string         `json:"proc_version"`
	Protocol     int            `json:"protocol"` // 0: iPad 1: Android 2: AndroidWatch  // 3 macOS 4 企点
	IMEI         string         `json:"imei"`
	Brand        string         `json:"brand"`
	Bootloader   string         `json:"bootloader"`
	BaseBand     string         `json:"base_band"`
	SimInfo      string         `json:"sim_info"`
	OSType       string         `json:"os_type"`
	MacAddress   string         `json:"mac_address"`
	IPAddress    []int32        `json:"ip_address"`
	WifiBSSID    string         `json:"wifi_bssid"`
	WifiSSID     string         `json:"wifi_ssid"`
	ImsiMd5      string         `json:"imsi_md5"`
	AndroidID    string         `json:"android_id"`
	APN          string         `json:"apn"`
	VendorName   string         `json:"vendor_name"`
	VendorOSName string         `json:"vendor_os_name"`
	Version      *osVersionFile `json:"version"`
}

type osVersionFile struct {
	Incremental string `json:"incremental"`
	Release     string `json:"release"`
	Codename    string `json:"codename"`
	Sdk         uint32 `json:"sdk"`
}

func randomMacAddress() string {
	buf := make([]byte, 6)
	_, err := crand.Read(buf)
	if err != nil {
		return "00:16:ea:ae:3c:40"
	}
	// Set the local bit
	buf[0] |= 2
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

func RandString(n int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, n)
	for i := range n {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// model	设备
// "iPhone11,2"	iPhone XS
// "iPhone11,8"	iPhone XR
// "iPhone12,1"	iPhone 11
// "iPhone13,2"	iPhone 12
// "iPad8,1"	iPad Pro
// "iPad11,2"	iPad mini
// "iPad13,2"	iPad Air 4
// "Apple Watch"	Apple Watch

func GenerateDeviceJSONIos(protocol int) (string, []byte, error) {
	bootID := uuid.New()
	imei := goluhn.Generate(15) // 注意，这个imei是完全胡乱创建的，并不符合imei规则
	androidID := fmt.Sprintf("%X", rand.Uint64())

	deviceJSON := deviceFile{
		Display:      "iPhone",      // Rom的名字 比如 Flyme 1.1.2（魅族rom）  JWR66V（Android nexus系列原生4.3rom）
		Product:      RandString(6), // 产品名，比如这是小米6的代号
		Device:       RandString(6),
		Board:        RandString(6),  // 主板:骁龙835                                                                    //
		Brand:        "Apple",        // 品牌
		Model:        "iPhone13,2",   // 型号
		Bootloader:   "unknown",      // unknown不需要改
		FingerPrint:  RandString(24), // 指纹
		BootID:       bootID.String(),
		ProcVersion:  "1.0", // 很长，后面 builder省略了
		BaseBand:     "",    // 基带版本 4.3CPL2-... 一大堆，直接不写
		SimInfo:      "",
		OSType:       "iOS",
		MacAddress:   randomMacAddress(),
		IPAddress:    []int32{192, 168, rand.Int31() % 255, rand.Int31()%253 + 2}, // 192.168.x.x
		WifiBSSID:    randomMacAddress(),
		WifiSSID:     "<unknown ssid>",
		IMEI:         imei,
		AndroidID:    androidID, // 原版的 androidId和Display内容一样，我没看协议，但是按android文档上说应该是64-bit number的hex，姑且这么做
		APN:          "wifi",
		VendorName:   "Apple", // 这个和下面一个选项(VendorOSName)都属于意义不明，找不到相似对应，不知道是啥
		VendorOSName: "Apple",
		Protocol:     protocol,
		Version: &osVersionFile{
			Incremental: "OCACNFA", // Build.Version.INCREMENTAL, MIUI12: V12.5.3.0.RJBCNXM
			Release:     "11",
			Codename:    "REL",
			Sdk:         29,
		},
	}

	if protocol == 2 {
		deviceJSON.Model = "Apple Watch"
	}

	if protocol == 3 {
		deviceJSON.Model = "mac OS X"
	}

	a, b := json.Marshal(deviceJSON)
	return deviceJSON.Model, a, b
}

func GenerateDeviceJSONAndroidWatch(protocol int) (string, []byte, error) {
	bootID := uuid.New()
	imei := goluhn.Generate(15) // 注意，这个imei是完全胡乱创建的，并不符合imei规则
	androidID := fmt.Sprintf("%X", rand.Uint64())

	deviceJSON := deviceFile{
		Display:      "MIRAI.142521.001", // Rom的名字 比如 Flyme 1.1.2（魅族rom）  JWR66V（Android nexus系列原生4.3rom）
		Product:      "mirai",            // 产品名，比如这是小米6的代号
		Device:       "mirai",
		Board:        "mirai",                                                           // 主板:骁龙835                                                                    //
		Brand:        "Apple",                                                           // 品牌
		Model:        "mirai",                                                           // 型号
		Bootloader:   "unknown",                                                         // unknown不需要改
		FingerPrint:  "mamoe/mirai/mirai:10/MIRAI.200122.001/9108230:user/release-keys", // 指纹
		BootID:       bootID.String(),
		ProcVersion:  "Linux version 3.0.31-zli0DMkg (android-build@xxx.xxx.xxx.xxx.com)", // 很长，后面 builder省略了
		BaseBand:     "",                                                                  // 基带版本 4.3CPL2-... 一大堆，直接不写
		SimInfo:      "T-Mobile",
		OSType:       "android",
		MacAddress:   randomMacAddress(),
		IPAddress:    []int32{192, 168, rand.Int31() % 255, rand.Int31()%253 + 2}, // 192.168.x.x
		WifiBSSID:    randomMacAddress(),
		WifiSSID:     "<unknown ssid>",
		IMEI:         imei,
		AndroidID:    androidID, // 原版的 androidId和Display内容一样，我没看协议，但是按android文档上说应该是64-bit number的hex，姑且这么做
		APN:          "wifi",
		VendorName:   "MIUI", // 这个和下面一个选项(VendorOSName)都属于意义不明，找不到相似对应，不知道是啥
		VendorOSName: "mirai",
		Protocol:     protocol,
		Version: &osVersionFile{
			Incremental: "5891938", // Build.Version.INCREMENTAL, MIUI12: V12.5.3.0.RJBCNXM
			Release:     "10",
			Codename:    "REL",
			Sdk:         29,
		},
	}

	a, b := json.Marshal(deviceJSON)
	return deviceJSON.Model, a, b
}

func GenerateDeviceJSONAllRandom(protocol int) (string, []byte, error) {
	bootID := uuid.New()
	imei := goluhn.Generate(15) // 注意，这个imei是完全胡乱创建的，并不符合imei规则
	androidID := fmt.Sprintf("%X", rand.Uint64())

	deviceJSON := deviceFile{
		Display:      RandString(6), // Rom的名字 比如 Flyme 1.1.2（魅族rom）  JWR66V（Android nexus系列原生4.3rom）
		Product:      RandString(6), // 产品名，比如这是小米6的代号
		Device:       RandString(6),
		Board:        RandString(6),  // 主板:骁龙835                                                                    //
		Brand:        RandString(12), // 品牌
		Model:        RandString(24), // 型号
		Bootloader:   "unknown",      // unknown不需要改
		FingerPrint:  RandString(24), // 指纹
		BootID:       bootID.String(),
		ProcVersion:  "1.0", // 很长，后面 builder省略了
		BaseBand:     "",    // 基带版本 4.3CPL2-... 一大堆，直接不写
		SimInfo:      "",
		OSType:       "android",
		MacAddress:   randomMacAddress(),
		IPAddress:    []int32{192, 168, rand.Int31() % 255, rand.Int31()%253 + 2}, // 192.168.x.x
		WifiBSSID:    randomMacAddress(),
		WifiSSID:     "<unknown ssid>",
		IMEI:         imei,
		AndroidID:    androidID, // 原版的 androidId和Display内容一样，我没看协议，但是按android文档上说应该是64-bit number的hex，姑且这么做
		APN:          "wifi",
		VendorName:   RandString(12), // 这个和下面一个选项(VendorOSName)都属于意义不明，找不到相似对应，不知道是啥
		VendorOSName: RandString(12),
		Protocol:     protocol,
		Version: &osVersionFile{
			Incremental: "OCACNFA", // Build.Version.INCREMENTAL, MIUI12: V12.5.3.0.RJBCNXM
			Release:     "11",
			Codename:    "REL",
			Sdk:         29,
		},
	}

	a, b := json.Marshal(deviceJSON)
	return deviceJSON.Model, a, b
}

func GenerateDeviceJSON(dice *Dice, protocol int) (string, []byte, error) {
	switch protocol {
	case 0, 3:
		return GenerateDeviceJSONIos(protocol)
	case 2:
		return GenerateDeviceJSONAndroidWatch(protocol)
	case 1:
		return GenerateDeviceJSONAndroid(dice, protocol)
	default:
		return GenerateDeviceJSONAllRandom(protocol)
	}
}

func GenerateDeviceJSONAndroid(dice *Dice, protocol int) (string, []byte, error) {
	// check if ./my_device.json exists
	if _, err := os.Stat("./my_device.json"); err == nil {
		dice.Logger.Info("检测到my_device.json，将使用该文件中的设备信息")
		// file exists
		data, err := os.ReadFile("./my_device.json")
		if err == nil {
			deviceJSON := deviceFile{}
			err = json.Unmarshal(data, &deviceJSON)
			if err == nil {
				deviceJSON.Protocol = protocol
				a, b := json.Marshal(deviceJSON)
				return deviceJSON.Model, a, b
			}
			dice.Logger.Warn("读取./my_device.json失败，将使用随机设备信息。原因为JSON解析错误: " + err.Error())
		}
		dice.Logger.Warn("读取./my_device.json失败，将使用随机设备信息")
	}

	pool := androidDevicePool
	imei := goluhn.Generate(15) // 注意，这个imei是完全胡乱创建的，并不符合imei规则
	androidID := fmt.Sprintf("%X", rand.Uint64())

	m := pool[rand.Int()%len(pool)]
	deviceJSON := m.data

	deviceJSON.MacAddress = randomMacAddress()
	deviceJSON.IPAddress = []int32{192, 168, rand.Int31() % 255, rand.Int31()%253 + 2} // 192.168.x.x
	deviceJSON.IMEI = imei
	deviceJSON.AndroidID = androidID
	deviceJSON.Protocol = protocol

	a, b := json.Marshal(deviceJSON)
	return deviceJSON.Model, a, b
}

func NewGoCqhttpConnectInfoItem(account string) *EndPointInfo {
	conn := new(EndPointInfo)
	conn.ID = uuid.New().String()
	conn.Platform = "QQ"
	conn.ProtocolType = "onebot"
	conn.Enable = false
	conn.RelWorkDir = "extra/go-cqhttp-qq" + account

	conn.Adapter = &PlatformAdapterGocq{
		EndPoint:        conn,
		UseInPackClient: true,
		BuiltinMode:     "gocq",
	}
	return conn
}

func BuiltinQQServeProcessKillBase(dice *Dice, conn *EndPointInfo, isSync bool) {
	f := func() {
		defer func() {
			if r := recover(); r != nil {
				dice.Logger.Error("内置 QQ 客户端清理报错: ", r)
				// go-cqhttp/lagrange 进程退出: exit status 1
			}
		}()

		pa, ok := conn.Adapter.(*PlatformAdapterGocq)
		if !ok {
			return
		}
		if !pa.UseInPackClient {
			return
		}

		// 重置状态
		conn.State = 0
		pa.GoCqhttpState = 0
		pa.GoCqhttpQrcodeData = nil

		if pa.BuiltinMode == "lagrange" {
			workDir := lagrangeGetWorkDir(dice, conn)
			qrcodeFile := filepath.Join(workDir, fmt.Sprintf("qr-%s.png", conn.UserID[3:]))
			if _, err := os.Stat(qrcodeFile); err == nil {
				// 如果已经存在二维码文件，将其删除
				_ = os.Remove(qrcodeFile)
				dice.Logger.Info("onebot: 删除已存在的二维码文件")
			}

			// 注意这个会panic，因此recover捕获了
			if pa.GoCqhttpProcess != nil {
				p := pa.GoCqhttpProcess
				pa.GoCqhttpProcess = nil
				// sigintwindows.SendCtrlBreak(p.Cmds[0].Process.Pid)
				_ = p.Stop()
				_ = p.Wait() // 等待进程退出，因为Stop内部是Kill，这是不等待的
			}
		} else {
			workDir := gocqGetWorkDir(dice, conn)
			qrcodeFile := filepath.Join(workDir, "qrcode.png")
			if _, err := os.Stat(qrcodeFile); err == nil {
				// 如果已经存在二维码文件，将其删除
				_ = os.Remove(qrcodeFile)
				dice.Logger.Info("onebot: 删除已存在的二维码文件")
			}

			// 注意这个会panic，因此recover捕获了
			if pa.GoCqhttpProcess != nil {
				p := pa.GoCqhttpProcess
				pa.GoCqhttpProcess = nil
				// sigintwindows.SendCtrlBreak(p.Cmds[0].Process.Pid)
				_ = p.Stop()
				_ = p.Wait() // 等待进程退出，因为Stop内部是Kill，这是不等待的
			}
		}
	}
	if isSync {
		f()
	} else {
		go f()
	}
}

func BuiltinQQServeProcessKill(dice *Dice, conn *EndPointInfo) {
	BuiltinQQServeProcessKillBase(dice, conn, false)
}

func gocqGetWorkDir(dice *Dice, conn *EndPointInfo) string {
	workDir := filepath.Join(dice.BaseConfig.DataDir, conn.RelWorkDir)
	return workDir
}

func GoCqhttpServeRemoveSessionToken(dice *Dice, conn *EndPointInfo) {
	workDir := gocqGetWorkDir(dice, conn)
	if _, err := os.Stat(filepath.Join(workDir, "session.token")); err == nil {
		_ = os.Remove(filepath.Join(workDir, "session.token"))
	}
}

type GoCqhttpLoginInfo struct {
	UIN           int64
	Password      string
	Protocol      int
	AppVersion    string
	IsAsyncRun    bool
	UseSignServer bool
}

func GoCqhttpServe(dice *Dice, conn *EndPointInfo, loginInfo GoCqhttpLoginInfo) {
	pa := conn.Adapter.(*PlatformAdapterGocq)
	// if pa.GoCqHttpState != StateCodeInit {
	//	return
	//}

	if pa.UseInPackClient {
		log := dice.Logger

		if dice.ContainerMode {
			log.Warn("onebot: 尝试启动内置客户端，但内置客户端在容器模式下被禁用")
			conn.State = 3
			pa.GoCqhttpState = StateCodeLoginFailed
			dice.Save(false)
			return
		}

		if pa.BuiltinMode == "gocq" || pa.BuiltinMode == "" {
			pa.CurLoginIndex++
			pa.GoCqhttpState = StateCodeInLogin
			builtinGoCqhttpServe(dice, conn, loginInfo)
		}
	} else {
		pa.GoCqhttpState = StateCodeLoginSuccessed
		pa.GoCqhttpLoginSucceeded = true
		dice.Save(false)
		go ServeQQ(dice, conn)
	}
}

func builtinGoCqhttpServe(dice *Dice, conn *EndPointInfo, loginInfo GoCqhttpLoginInfo) {
	log := zap.S().Named(logger.LogKeyAdapter)
	pa := conn.Adapter.(*PlatformAdapterGocq)
	loginIndex := pa.CurLoginIndex

	// 保留此if语句块，使历史可追溯，后续commit可移除if
	if pa.UseInPackClient { //nolint:nestif
		workDir := gocqGetWorkDir(dice, conn)
		_ = os.MkdirAll(workDir, 0o755)
		downloadGoCqhttp(dice.Logger)

		qrcodeFile := filepath.Join(workDir, "qrcode.png")
		deviceFilePath := filepath.Join(workDir, "device.json")
		configFilePath := filepath.Join(workDir, "config.yml")
		if _, err := os.Stat(qrcodeFile); err == nil {
			// 如果已经存在二维码文件，将其删除
			_ = os.Remove(qrcodeFile)
			dice.Logger.Info("onebot: 删除已存在的二维码文件")
		}

		if !pa.GoCqhttpLoginSucceeded {
			// 并未登录成功，删除记录文件
			dice.Logger.Info("onebot: 之前并未登录成功，删除设备文件和配置文件")
			_ = os.Remove(configFilePath)
			_ = os.Remove(deviceFilePath)
		}

		// 创建设备配置文件
		if _, err := os.Stat(deviceFilePath); errors.Is(err, os.ErrNotExist) {
			_, deviceInfo, err := GenerateDeviceJSON(dice, loginInfo.Protocol)
			if err == nil {
				_ = os.WriteFile(deviceFilePath, deviceInfo, 0o644)
				dice.Logger.Info("onebot: 成功创建设备文件")
			}
		}

		// 启动客户端
		wd, _ := os.Getwd()
		gocqhttpExePath, _ := filepath.Abs(filepath.Join(wd, "go-cqhttp/go-cqhttp"))
		gocqhttpExePath = filepath.ToSlash(gocqhttpExePath) // windows平台需要这个替换

		// 随手执行一下
		_ = os.Chmod(gocqhttpExePath, 0o755)

		dice.Logger.Info("onebot: 正在启动onebot客户端…… ", gocqhttpExePath)
		p := procs.NewProcess(fmt.Sprintf(`"%s" -faststart`, gocqhttpExePath))
		p.Dir = workDir

		if runtime.GOOS == "android" {
			p.Env = os.Environ()
		}
		chQrCode := make(chan int, 1)
		riskCount := 0
		isSelfKilling := false

		p.OutputHandler = func(line string, _type string) string {
			if loginIndex != pa.CurLoginIndex {
				// 当前连接已经无用，进程自杀
				if !isSelfKilling {
					dice.Logger.Infof("检测到新的连接序号 %d，当前连接 %d 将自动退出", pa.CurLoginIndex, loginIndex)
					// 注: 这里不要调用kill
					isSelfKilling = true
					_ = p.Stop()
				}
				return ""
			}

			// 登录中
			if pa.IsInLogin() {
				// 请使用手机QQ扫描二维码 (qrcode.png) :
				if strings.Contains(line, "qrcode.png") {
					chQrCode <- 1
				}

				// 获取二维码失败，登录失败
				if strings.Contains(line, "fetch qrcode error: Packet timed out ") {
					dice.Logger.Infof("从QQ服务器获取二维码错误（超时），帐号: <%s>(%s)", conn.Nickname, conn.UserID)
					pa.GoCqhttpState = StateCodeLoginFailed
				}

				// 未知错误，gocqhttp崩溃
				if strings.Contains(line, "Packet failed to sendPacket: connection closed") {
					dice.Logger.Infof("登录异常，gocqhttp崩溃")
					pa.GoCqhttpState = StateCodeLoginFailed
				}

				if strings.Contains(line, "按 Enter 继续....") {
					// 直接输入继续，基本都是登录失败
					return "\n"
				}

				// 登录成功
				if strings.Contains(line, "CQ WebSocket 服务器已启动") {
					// CQ WebSocket 服务器已启动
					// 登录成功 欢迎使用
					pa.GoCqhttpState = StateCodeLoginSuccessed
					pa.GoCqhttpLoginSucceeded = true
					dice.Logger.Infof("gocqhttp登录成功，帐号: <%s>(%s)", conn.Nickname, conn.UserID)
					dice.LastUpdatedTime = time.Now().Unix()
					dice.Save(false)

					go ServeQQ(dice, conn)
				}
			}

			if strings.Contains(line, "请使用手机QQ扫描二维码以继续登录") {
				// TODO
				log.Info("请使用手机QQ扫描二维码以继续登录")
			}

			if (pa.IsLoginSuccessed() && strings.Contains(line, "[ERROR]:") && strings.Contains(line, "Protocol -> sendPacket msg error: 120")) || strings.Contains(line, "账号可能被风控####2测试触发语句") {
				// 这种情况应该是被禁言，提前减去以免出事
				riskCount--
				dice.Logger.Infof("因禁言无法发言: 下方可能会提示遭遇风控")
			}

			if pa.IsLoginSuccessed() {
				if strings.Contains(line, "[ERROR]:") && strings.Contains(line, "错误: 请先添加") && strings.Contains(line, "为好友") {
					// [ERROR]: 错误: 请先添加 xxx 为好友
					pa.riskAlertShieldCount++
				}

				if strings.Contains(line, "[WARNING]:") && strings.Contains(line, "图片上传失败") {
					// [2023-09-27 11:19:18] [WARNING]: 警告: 群 xxx 图片上传失败
					pa.riskAlertShieldCount++
				}
			}

			if (pa.IsLoginSuccessed() && strings.Contains(line, "WARNING") && strings.Contains(line, "账号可能被风控")) || strings.Contains(line, "账号可能被风控####测试触发语句") {
				// 群消息发送失败: 账号可能被风控
				pa.GoCqhttpLastRestrictedTime = time.Now().Unix()
			}

			if pa.IsInLogin() || strings.Contains(line, "[WARNING]") || strings.Contains(line, "[ERROR]") || strings.Contains(line, "[FATAL]") {
				//  [WARNING]: 登录需要滑条验证码, 请使用手机QQ扫描二维码以继续登录
				if pa.IsLoginSuccessed() {
					skip := false

					if strings.Contains(line, "WARNING") {
						if strings.Contains(line, "检查更新失败") || strings.Contains(line, "Protocol -> device lock is disable.") {
							skip = true
						}
						if strings.Contains(line, "语音文件") && strings.Contains(line, "下载失败") {
							skip = true
						}
					}

					if strings.Contains(line, "ERROR") {
						if strings.Contains(line, "panic on decoder MsgPush.PushGroupProMsg") {
							skip = true
						}
					}

					if !skip {
						dice.Logger.Infof("onebot | %s", stripansi.Strip(line))
					} else if strings.HasSuffix(line, "\n") {
						dice.Logger.Infof("onebot | %s", line[:len(line)-1])
					}
				} else {
					if strings.HasSuffix(line, "\n") {
						dice.Logger.Infof("onebot | %s", line[:len(line)-1])
					}

					skip := false
					if strings.Contains(line, "WARNING") && strings.Contains(line, "使用了过时的配置格式，请更新配置文件") {
						skip = true
					}

					if strings.Contains(line, "WARNING") {
						if strings.Contains(line, "检查更新失败") || strings.Contains(line, "Protocol -> device lock is disable.") {
							skip = true
						}
					}

					// error 之类错误无条件警告
					if !skip {
						if strings.Contains(line, "WARNING") || strings.Contains(line, "ERROR") || strings.Contains(line, "FATAL") {
							dice.Logger.Infof("onebot | %s", stripansi.Strip(line))
						}
					}
				}
			}
			return ""
		}

		go func() {
			<-chQrCode
			if _, err := os.Stat(qrcodeFile); err == nil {
				dice.Logger.Info("onebot: 二维码已经就绪")
				fmt.Fprintln(os.Stdout, "如控制台二维码不好扫描，可以手动打开 ./data/default/extra/go-cqhttp-qqXXXXX 目录下qrcode.png")
				qrdata, err := os.ReadFile(qrcodeFile)
				if err == nil {
					pa.GoCqhttpState = StateCodeInLoginQrCode
					pa.GoCqhttpQrcodeData = qrdata
					dice.Logger.Info("获取二维码成功")
					dice.LastUpdatedTime = time.Now().Unix()
					dice.Save(false)
					_ = os.Rename(qrcodeFile, qrcodeFile+".bak.png")
				} else {
					pa.GoCqhttpQrcodeData = nil
					pa.GoCqhttpState = StateCodeLoginFailed
					pa.GocqhttpLoginFailedReason = "获取二维码失败"
					dice.LastUpdatedTime = time.Now().Unix()
					dice.Save(false)
					dice.Logger.Info("获取二维码失败，错误为: ", err.Error())
				}
			}
		}()

		run := func() {
			defer func() {
				if r := recover(); r != nil {
					dice.Logger.Errorf("onebot: 异常: %v 堆栈: %v", r, string(debug.Stack()))
				}
			}()

			// 启动gocqhttp，开始登录
			pa.GoCqhttpProcess = p
			err := p.Start()

			if err == nil {
				if dice.Parent.progressExitGroupWin != 0 && p.Cmd != nil {
					errAdd := dice.Parent.progressExitGroupWin.AddProcess(p.Cmd.Process)
					if errAdd != nil {
						dice.Logger.Warn("添加到进程组失败，若主进程崩溃，gocqhttp进程可能需要手动结束")
					}
				}
				_ = p.Wait()
			}

			isInLogin := pa.IsInLogin()
			isDeviceLockLogin := pa.GoCqhttpState == StateCodeInLoginDeviceLock
			if !isDeviceLockLogin {
				// 如果在设备锁流程中，不清空数据
				BuiltinQQServeProcessKill(dice, conn)

				if isInLogin {
					conn.State = 3
					pa.GoCqhttpState = StateCodeLoginFailed
				} else {
					conn.State = 0
					pa.GoCqhttpState = GoCqhttpStateCodeClosed
				}
			}

			if err != nil {
				dice.Logger.Info("go-cqhttp 进程退出: ", err)
			} else {
				dice.Logger.Info("go-cqhttp 进程退出")
			}
		}

		if loginInfo.IsAsyncRun {
			go run()
		} else {
			run()
		}
	}
}

var isGocqDownloading = false

func downloadGoCqhttp(logger *zap.SugaredLogger) {
	fn := "go-cqhttp/go-cqhttp"
	if runtime.GOOS == "windows" {
		fn += ".exe"
	}

	url := fmt.Sprintf("https://d1.sealdice.com/go-cqhttp/go-cqhttp_%s_%s", runtime.GOOS, runtime.GOARCH)

	isExists := false
	if _, err := os.Stat(fn); err == nil {
		// 存在，并可执行
		_ = os.Chmod(fn, 0o755)
		cmd := exec.Command(fn, "-h")
		out, err := cmd.Output()
		if err == nil {
			if strings.Contains(string(out), "go-cqhttp") {
				isExists = true
			}
		}
	}

	if !isExists {
		if logger != nil {
			logger.Info("go-cqhttp不存在，进行下载")
		}
		if !isGocqDownloading {
			isGocqDownloading = true
			_ = os.MkdirAll("./go-cqhttp", 0755)
			if err := utils.DownloadFile(fn, url); err != nil && logger != nil {
				logger.Info("go-cqhttp下载失败，请进行禁用/启用操作以再次下载")
			}
			isGocqDownloading = false
		} else {
			for {
				// 等待到下载完成
				time.Sleep(2 * time.Second)
				if !isGocqDownloading {
					break
				}
			}
		}
	}
}
