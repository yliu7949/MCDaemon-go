/*
 * 插件编写例子
 * 当玩家输入命令：!!yinyinmaster
 * 他会对所有人说：嘤嘤嘤
 */

package plugin

import (
	"fmt"

	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type Yinyin struct {
}

func (hp *Yinyin) Handle(c *command.Command, s lib.Server) {
	s.Say(fmt.Sprintf("%s对所有人说：嘤嘤嘤！", c.Player))
}

func (hp *Yinyin) Init(s lib.Server) {
}

func (hp *Yinyin) Close() {
}
