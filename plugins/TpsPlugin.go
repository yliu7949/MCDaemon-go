package plugin

import (
	"fmt"
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
		s.Tell(c.Player,"开始调试分析...")
		s.Execute("debug start")
		time.Sleep(time.Second * time.Duration(second))
		go func() {
			s.Execute("debug stop")
		}()
		reg := `Stopped debug profiling after (\d+\.\d+) (.*) (\d+) ticks \((\d+\.\d+)`
		match, flag := s.RegParser(reg)
		if !flag {
			s.Tell(c.Player,"出现错误，请重新尝试。")
			return
		}
		s.Tell(c.Player,fmt.Sprintf("调试分析结束，用时%s秒和%s刻（每秒%s刻）。",match[1],match[3],match[4]))
	} else {
		text := "!!tps [秒数] 启动调试并计算指定秒数的tps（自动限幅30s）。"
		s.Tell(c.Player, text)
	}
}

func (hp *TpsPlugin) Init(s lib.Server) {
}

func (hp *TpsPlugin) Close() {
}
