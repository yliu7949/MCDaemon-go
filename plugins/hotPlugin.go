package plugin

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/lib"
)

//热加载插件类型
type HotPlugin struct {
	string
}

func (hp *HotPlugin) Handle(c *command.Command, s lib.Server) {
	commandName := "./hotPlugins/" + c.PluginName
	pluginProcess := exec.Command(commandName, c.Argv...)
	buffer, err := pluginProcess.Output()
	if err != nil {
		s.Tell(c.Player, fmt.Sprint("插件出现错误：", err))
	}
	retStr := string(buffer)
	/**
	插件返回数据以空格区分参数
	第一个为调用方法名
	第二个为方法参数
	第三个如果有则代表玩家名
	*/
	argv := strings.Fields(retStr)
	if len(argv) >= 2 {
		switch argv[0] {
		case "say":
			s.Say(strings.Join(argv[1:], " "))
		case "tell":
			if len(argv) >= 3 {
				s.Tell(argv[1], strings.Join(argv[2:], " "))
			}
		case "Execute":
			s.Execute(strings.Join(argv[1:], " "))
		}
	}
}

func (hp *HotPlugin) Init(s lib.Server) {
}

func (hp *HotPlugin) Close() {
}
