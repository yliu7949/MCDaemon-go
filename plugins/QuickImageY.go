/**
 * 快速镜像插件
 * 前置插件： 快速备份插件QuickBackupY
 * 编写者：Underworld511
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/container"
	"MCDaemon-go/lib"
	"github.com/go-ini/ini"
	"github.com/otiai10/copy"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type QuickImageY struct{}

func (qi *QuickImageY) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	dir := "QuickBackup/"
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
			file, err := os.OpenFile(qbDataFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				s.Tell(c.Player,"数据文件打开失败！")
				return
			}
			b, err := ioutil.ReadAll(file)
			if err != nil {
				s.Tell(c.Player,"数据读取失败！")
				return
			}
			defer file.Close()
			if err := copy.Copy(dir + gjson.Get(string(b), "Slot1.Name").String(), "QuickImage/" + c.Argv[1]); err != nil {
				lib.WriteDevelopLog("error", err.Error())
				s.Tell(c.Player,"文件复制失败！")
				return
			}
			s.Tell(c.Player,"镜像添加成功。")
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
			sercfg.Section("").NewKey("server-port", svr.GetPort())
			sercfg.SaveTo(path)
			cor.Add(c.Argv[1], "QuickImage/" + c.Argv[1], svr)
		}
		s.Say("镜像启动成功。")
	case "show":
		imageFiles, _ := filepath.Glob("QuickImage/*")
		text := "QuickImage镜像列表：\\n"
		for k, _ := range imageFiles {
			if cor.IsRuntime(imageFiles[k]) {
				text += imageFiles[k] + "  已启动  " + cor.Servers[imageFiles[k]].GetPort() + "\\n"
			} else {
				text += imageFiles[k] + "  未启动  " + "\\n"
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
	case "restart":
		if len(c.Argv) < 3 {
			s.Tell(c.Player, "缺少参数，请同时加上镜像名和端口号！")
			return
		}
		port, err := strconv.Atoi(c.Argv[2])
		if err != nil || port <= 0 {
			s.Tell(c.Player, "参数不合法！")
			return
		}
		_commond1 := &command.Command{
			Player: c.Player,
			Cmd:    "qi",
			Argv:    []string{"stop", c.Argv[1]},
		}
		_commond2 := &command.Command{
			Player: c.Player,
			Cmd:    "qi",
			Argv:    []string{"start", c.Argv[1], c.Argv[2]},
		}
		s.RunPlugin(_commond1)
		s.RunPlugin(_commond2)
	default:
		text := "使用规则：\\n!!qi add <镜像名> 添加镜像\\n!!qi start <镜像名> <port> 启动镜像 \\n!!qi show 查看镜像列表\\n" +
			"!!qi stop <镜像名> 停止镜像\\n!!qi del <镜像名> 删除镜像\\n!!qi restart <镜像名> <port> 重启镜像\\n"
		s.Tell(c.Player, text)
	}
}

func (qi *QuickImageY) Init(s lib.Server) {
}

func (qi *QuickImageY) Close() {
}
