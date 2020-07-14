/**
 * 快速镜像插件
 * 前置插件： 快速备份插件QuickBackupY
 * 编写者：Underworld511
 */

package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-ini/ini"
	"github.com/otiai10/copy"
	"github.com/tidwall/gjson"
	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/container"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type QuickImageY struct{}

func (qi *QuickImageY) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	qbDataFile := "QuickBackup/qb_data.json"
	cor := container.GetInstance()
	switch c.Argv[0] {
	case "add":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "请为要添加的镜像取个名字吧!")
			return
		}
		if !checkSlot(qbDataFile, 1) {		//检查slot槽位
			s.Tell(c.Player,"参数不合法或无法找到有效存档！")
		} else {
			s.Tell(c.Player,"开始镜像...")
			text, _ := addMirror(qbDataFile, c.Argv[1])
			s.Tell(c.Player,text)
		}
	case "start":
		if len(c.Argv) < 3 {
			s.Tell(c.Player, "缺少参数，请同时加上镜像名和端口号！")
			return
		}
		if c.Argv[2] == s.GetPort() {
			s.Tell(c.Player, "端口"+ s.GetPort() +"已被占用，请指定其它端口号！")
		} else {
			if !checkFileIsExist("QuickImage/" +c.Argv[1]) {
				s.Tell(c.Player,"无法找到镜像文件！")
				return
			}
			if cor.IsRuntime(c.Argv[1]) {
				s.Tell(c.Player, "该镜像已经启动")
				return
			}
			port, err := strconv.Atoi(c.Argv[2])
			if err != nil || port <= 0 {
				s.Tell(c.Player, "参数不合法！")
				return
			}
			s.Tell(c.Player,"启动镜像...")
			path := "QuickImage/" + c.Argv[1] + "/server.properties"
			svr := s.Clone(c.Argv[2])
			sercfg, _ := ini.Load(path)
			rconPort, _ := strconv.Atoi(c.Argv[2])
			rconPort = rconPort + 7
			sercfg.Section("").Key("server-port").SetValue(c.Argv[2])
			sercfg.Section("").Key("rcon.port").SetValue(strconv.Itoa(rconPort))
			sercfg.Section("").Key("gamemode").SetValue("creative")
			sercfg.SaveTo(path)
			cor.Add(c.Argv[1], "QuickImage/" + c.Argv[1], svr)
			t := time.Now()
			s.Say("镜像" + c.Argv[1] + "正在启动中，请耐心等待地图加载完成...")
			for {
				if svr.GetStartResult() != 0 {
					if svr.GetStartResult() == 1 {
						s.Say("镜像" + c.Argv[1] + "启动成功，耗时" +
							fmt.Sprintf("%v", time.Since(t).Truncate(time.Second)) + "。")
						return
					}
					if svr.GetStartResult() == -1 {
						s.Say("镜像" + c.Argv[1] + "启动失败。")
						return
					}
				}
				if time.Since(t) > 70*time.Second {
					s.Say("镜像" + c.Argv[1] + "启动超时。")
					return
				}
				time.Sleep(1*time.Second)
			}
		}
	case "show":
		imageFiles, _ := filepath.Glob("QuickImage/*")
		text := "QuickImage镜像列表：\\n"
		for _, image := range imageFiles {
			if cor.IsRuntime(image[11:]) {
				path := image + "/server.properties"
				sercfg, _ := ini.Load(path)
				rconPort := fmt.Sprintf("%s",sercfg.Section("").Key("rcon.port"))
				serverPort := fmt.Sprintf("%s",sercfg.Section("").Key("server-port"))
				text += image[11:] + "  已启动  " + serverPort + "  [RCON] : " + rconPort + "\\n"
			} else {
				text += image[11:] + "  未启动  " + "\\n"
			}
		}
		s.Tell(c.Player, text)
	case "stop":
		if len(c.Argv) == 1 {
			s.Tell(c.Player, "缺少停止的镜像名称！")
		} else {
			if cor.IsRuntime(c.Argv[1]) {
				cor.Del(c.Argv[1])
			} else {
				s.Tell(c.Player, "镜像未启动.")
			}
		}
		s.Tell(c.Player,"镜像已停止。")
	case "del":
		if len(c.Argv) == 1 {
			s.Tell(c.Player, "缺少停止的镜像名称！")
		} else {
			if cor.IsRuntime(c.Argv[1]) {
				s.Tell(c.Player, "请先停止镜像。")
			} else {
				if !checkFileIsExist("QuickImage/" +c.Argv[1]) {
					s.Tell(c.Player,"镜像文件不存在或已经被删除！")
					return
				}
				err := os.RemoveAll("QuickImage/" +c.Argv[1])
				if err != nil {
					s.Tell(c.Player,"文件删除失败！")
					return
				}
				s.Tell(c.Player,"镜像删除成功。")
			}
		}
	case "update", "u":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "缺少要同步的镜像名称。")
			return
		}
		if !checkFileIsExist("QuickImage/" +c.Argv[1]) {
			s.Tell(c.Player,"镜像文件不存在！")
			return
		}
		path := "QuickImage/" + c.Argv[1] + "/server.properties"
		sercfg, _ := ini.Load(path)
		rconPort := fmt.Sprintf("%s",sercfg.Section("").Key("rcon.port"))
		serverPort := fmt.Sprintf("%s",sercfg.Section("").Key("server-port"))
		if cor.IsRuntime(c.Argv[1]) {
			cor.Del(c.Argv[1])		//停止运行镜像
			s.Tell(c.Player,"正在停止运行该镜像...")
		}
		if len(c.Argv) == 2 {
			c.Argv = append(c.Argv, "help")
		}
		switch c.Argv[2] {
		case "all", "a":
			time.Sleep(1 * time.Second)
			s.Tell(c.Player,"开始同步镜像...")
			s.Tell(c.Player,"正在删除旧的镜像数据...")
			err := os.RemoveAll("QuickImage/" + c.Argv[1])
			if err != nil {
				s.Tell(c.Player, "同步失败！在删除旧的镜像时出现了错误。")
				return
			}
			s.Tell(c.Player, "正在写入新的镜像数据...")
			if _, err := addMirror(qbDataFile, c.Argv[1]); err != nil {
				s.Tell(c.Player, "同步失败！在写入新的镜像数据时出现了错误。")
				return
			}
			sercfg, _ := ini.Load(path)
			sercfg.Section("").Key("server-port").SetValue(serverPort)
			sercfg.Section("").Key("rcon.port").SetValue(rconPort)
			sercfg.Section("").Key("gamemode").SetValue("creative")
			sercfg.SaveTo(path)
			//svr := s.Clone(serverPort)
			//cor.Add(c.Argv[1], "QuickImage/"+c.Argv[1], svr)
			s.Tell(c.Player, "镜像同步成功！使用!!qi start "+c.Argv[1]+" "+serverPort+"命令来启动该镜像。")
		case "world", "w", "nether", "n", "end", "e":
			time.Sleep(1 * time.Second)
			s.Tell(c.Player,"开始同步镜像...")
			s.Tell(c.Player,"正在删除旧的镜像数据...")
			folder := map[string]string{
				"w":"/world/",
				"world":"/world/",
				"n":"/world/DIM-1/",
				"nether":"/world/DIM-1/",
				"e":"/world/DIM1/",
				"end":"/world/DIM1/",
			}
			path = "QuickImage/" + c.Argv[1] + folder[c.Argv[2]]
			text, err := updateMirrorPartly(qbDataFile, folder[c.Argv[2]], path)
			if err != nil {
				s.Tell(c.Player, text)
			} else {
				s.Tell(c.Player, text+"使用!!qi start "+c.Argv[1]+" "+serverPort+"命令来启动该镜像。")
			}
		default:
			text := "update命令用法示例：\\n!!qi update <镜像名> [world/nether/end] 同步主世界/下界/末地\\n!!qi u <镜像名> [w/n/e] 同步主世界/下界/末地（简写命令）" +
				"\\n!!qi update <镜像名> all 同步所有维度\\n!!qi u <镜像名> a 同步所有维度（简写命令）\\n"
			s.Tell(c.Player, text)
		}
	case "op":
		if s.GetName() == "default" {
			s.Tell(c.Player, "该命令仅在镜像服可用。")
		} else {
			s.Execute("/op " + c.Player)
		}
	default:
	case "help":
		text := "使用规则：\\n!!qi add <镜像名> 添加镜像\\n!!qi start <镜像名> <port> 启动镜像 \\n!!qi show 查看镜像列表\\n" +
			"!!qi stop <镜像名> 停止镜像\\n!!qi del <镜像名> 删除镜像\\n!!qi update <镜像名> [world/nether/end/all] 同步镜像\\n" +
			"!!qi op 获取op身份（仅镜像服可用）\\n建议使用的镜像名及对应端口 MirrorY-25569 MirrorZ-25570\\n"
		s.Tell(c.Player, text)
	}
}

func addMirror(qbDataFile string, mirrorName string) (string, error) {
	dir := "QuickBackup/"
	file, err := os.OpenFile(qbDataFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "数据文件打开失败！", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "数据读取失败！", err
	}
	defer file.Close()
	if err := copy.Copy(dir + gjson.Get(string(b), "Slot1.Name").String(), "QuickImage/" + mirrorName); err != nil {
		return "文件复制失败！", err
	}
	return "镜像添加成功。", nil
}

func updateMirrorPartly(qbDataFile string, qbPath string, qiPath string) (string, error) {
	file, err := os.OpenFile(qbDataFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "数据文件打开失败！", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "数据读取失败！", err
	}
	defer file.Close()
	qbPath ="QuickBackup/" + gjson.Get(string(b), "Slot1.Name").String() + qbPath
	dimData := [...]string{"data", "poi", "region"}
	for _, dir := range dimData {
		if err = os.RemoveAll(qiPath+dir); err != nil {
			return "同步失败！在删除旧的镜像数据时出现了错误。", err
		}
		if err = copy.Copy(qbPath+dir, qiPath+dir); err != nil {
			return "同步失败！在写入新的镜像数据时出现了错误。", err
		}
	}
	return "镜像同步成功！", nil
}

func (qi *QuickImageY) Init(s lib.Server) {
}

func (qi *QuickImageY) Close() {
}
