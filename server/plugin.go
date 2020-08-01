package server

import (
	"fmt"
	"strings"

	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/config"
	"github.com/yliu7949/MCDaemon-go/lib"
)

//运行所有语法解析器
func (svr *Server) RunParsers(word string) {
	for _, val := range svr.parserList {
		cmd, ok := val.Parsing(word)
		if ok && svr.pluginList[cmd.Cmd] != nil {
			//异步运行插件
			svr.pulginPool <- 1
			if cmd.Player != "" {
				svr.WriteLog("info", fmt.Sprintf("玩家 %s 运行了 %s %s 命令", cmd.Player, cmd.Cmd, strings.Join(cmd.Argv, " ")))
			}
			go svr.RunPlugin(cmd)
		}
	}

}

//运行插件命令（结构体格式）
func (svr *Server) RunPlugin(cmd *command.Command) {
	svr.pluginList[cmd.Cmd].Handle(cmd, svr)
	<-svr.pulginPool
}

//运行插件命令（字符串格式）
func (svr *Server) RunPluginCommand(player string, cmd string) bool{
	cmdSlice := strings.Fields(cmd)
	if len(cmdSlice) == 0 {
		return false
	}
	svr.RunPlugin(&command.Command{
		Player: player,
		Argv:   cmdSlice[1:],
		Cmd:    cmdSlice[0],
	})
	return true
}

//等待现有插件的完成并停止后面插件的运行，在执行相关操作
func (svr *Server) RunUniquePlugin(handle func()) {
	svr.unqiueLock.Lock()
	defer svr.unqiueLock.Unlock()
	<-svr.pulginPool
	//根据插件最大并发数进行堵塞
	maxRunPlugins, _ := config.Cfg.Section("MCDaemon").Key("maxRunPlugins").Int()
	for i := 0; i < maxRunPlugins; i++ {
		svr.pulginPool <- 1
	}
	handle()
	for i := 0; i < maxRunPlugins; i++ {
		<-svr.pulginPool
	}
	svr.pulginPool <- 1
}

//获取当前实例的插件列表
func (svr *Server) GetPluginList() map[string]lib.Plugin {
	return svr.pluginList
}

//获取当前实例的禁用插件列表
func (svr *Server) GetDisablePluginList() map[string]lib.Plugin {
	return svr.disablePluginList
}

//获取语法解析器列表
func (svr *Server) GetParserList() []lib.Parser {
	return svr.parserList
}

//调用插件释放资源函数
func (svr *Server) RunPluginClose() {
	for _, v := range svr.pluginList {
		v.Close()
	}
	for _, v := range svr.disablePluginList {
		v.Close()
	}
}
