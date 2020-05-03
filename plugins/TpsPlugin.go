package plugin

import (
	"strconv"
	"time"

	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type TpsPlugin struct{}

func (hp *TpsPlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	if second, ok := strconv.Atoi(c.Argv[0]); ok == nil {
		if second > 30 {
			second = 30
		}
		s.Execute("debug start")
		time.Sleep(time.Second * time.Duration(second))
		s.Execute("debug stop")
	} else if c.Argv[0] == "res" {
		s.Say(c.Argv[1][:len(c.Argv[1])-1])
	} else {
		text := "使用 !!tps [秒数] 指定获取多少秒内的tps（自动限幅30s内）"
		s.Tell(c.Player, text)
	}
}

func (hp *TpsPlugin) Init(s lib.Server) {
}

func (hp *TpsPlugin) Close() {
}
