package main

import (
	"fmt"
	"os"

	"github.com/yliu7949/MCDaemon-go/config"
	"github.com/yliu7949/MCDaemon-go/container"
	"github.com/yliu7949/MCDaemon-go/server"
)

func init() {
	//配置eula文件
	config.SetEula()
	//获取所有启动项配置
	_ = config.GetStartConfig()
	if _,err := os.Stat("./data"); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("./data", os.ModePerm)
			if err != nil {
				fmt.Println("创建data文件夹失败：", err)
			}
			return
		}
	}
}

func main() {
	c := container.GetInstance()
	defaultServer := &server.Server{}
	//加入到容器中并开启服务器
	c.Add("default", config.Cfg.Section("MCDaemon").Key("server_path").String(), defaultServer)
	c.Group.Wait()
}
