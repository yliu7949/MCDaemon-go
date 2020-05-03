package plugin

import (
	"encoding/json"

	"github.com/shiguanghuxian/txai" // 引入sdk
	"github.com/tidwall/gjson"
	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/config"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type SDChatPlugin struct{}

func (hp *SDChatPlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) < 1 {
		c.Argv = append(c.Argv, "help")
	}
	switch c.Argv[0] {
	case "all":
		command.Group.AddPlayer("SDChat-all", c.Player)
		s.Tell(c.Player, "开启全局聊天模式成功")
	case "start":
		command.Group.AddPlayer("SDChat", c.Player)
		s.Tell(c.Player, "开启聊天模式成功")
	case "stop":
		command.Group.DelPlayer("SDChat", c.Player)
		command.Group.DelPlayer("SDChat-all", c.Player)
		s.Tell(c.Player, "退出聊天模式成功")
	case "say":
		s.Tell(c.Player, "沙雕："+chat(c.Argv[1]))
	case "say-all":
		s.Say("沙雕对" + c.Player + "说：" + chat(c.Argv[1]))
	case "reload":
		_ = config.GetPluginCfg(true)
		s.Tell(c.Player, "已重新读取配置文件")
	default:
		text := "!!SDChat all start 开启全局聊天模式\\n!!SDChat start 开启私聊模式（别的玩家看不见沙雕机器人给你发的信息）\\n!!SDChat stop 关闭聊天模式\\n!!SDChat reload 重新加载配置文件"
		s.Tell(c.Player, text)
	}
}

func (hp *SDChatPlugin) Init(s lib.Server) {
}

func (hp *SDChatPlugin) Close() {
}


//发出请求获取聊天回复
func chat(question string) string {
	txAi := txai.New("2********9", "S*************", false)	// 创建sdk操作对象
	val, _ := txAi.NlpTextchatForText("10000",question)	// 调用对应腾讯ai接口的对应函数
	js, _ := json.Marshal(val)
	return gjson.Get(string(js), "data.answer").String()
}
