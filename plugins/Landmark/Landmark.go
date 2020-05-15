/*
 * 坐标记录插件
 * 编写者：Underworld511
 */
package Landmark

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type Landmark struct{}

func (lm *Landmark) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	path := "plugins/Landmark/Landmark.ini"
	serverName := s.GetName()
	switch c.Argv[0] {
	case "add","a":
		if len(c.Argv) < 3 {
			s.Tell(c.Player, "参数不足!")
			return
		}
		cfg, _ := ini.Load(path)
		if section, err := cfg.NewSection(serverName); err == nil {
			_, _ = section.NewKey(c.Argv[1], fmt.Sprintf("%s", c.Argv[2:]))
			s.Tell(c.Player, "坐标点添加成功!")
		} else {
			s.Tell(c.Player, "创建Section时出现了错误。")
		}
		cfg.SaveTo(path)
	case "del","d":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "请加上要删除的坐标点的名字!")
			return
		}
		cfg, _ := ini.Load(path)
		cfg.Section(serverName).DeleteKey(c.Argv[1])
		cfg.SaveTo(path)
		s.Tell(c.Player, "坐标点删除成功!")
	case "list","l":
		cfg, _ := ini.Load(path)
		text := "Landmark记载的坐标点如下：\\n"
		landmarks := cfg.Section(serverName).KeyStrings()
		for _,landmark := range landmarks {
			//value := fmt.Sprintf("%s",cfg.Section("").Key(landmark).Value())
			text = text + landmark + "  " + cfg.Section(serverName).Key(landmark).Value() + "\\n"
		}
		s.Tell(c.Player, text)
	case "help","h":
		text := "使用规则：\\n!!lm add <name> <x> <y> <z> 添加坐标点\\n!!lm del <name> 删除坐标点\\n" +
			"!!lm list 列出坐标点\\n"
		s.Tell(c.Player, text)
	}

}

func (lm *Landmark) Init(s lib.Server) {
}

func (lm *Landmark) Close() {
}