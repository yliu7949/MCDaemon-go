package plugin

import (
	"fmt"
	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type Here struct {
	Dim 				string			//玩家所在的维度，"0":"主世界","-1":"地狱","1":"末地"
	PosX				string			//玩家的x坐标
	PosY				string			//玩家的y坐标
	PosZ				string			//玩家的z坐标
}


func (h *Here) Handle(c *command.Command, s lib.Server) {
	p := GetDim(c.Player, s, h)
	switch p.Dim {
	case "0":
		p.Dim = "主世界"
	case "-1":
		p.Dim = "地狱"
	case "1":
		p.Dim = "末地"
	default:
		s.Tell(c.Player, "获取维度失败。")
		return
	}
	p = GetPosition(c.Player, s,p)
	if p.PosZ == "" {
		s.Tell(c.Player, "获取坐标失败。")
		return
	}
	p.PosX = p.PosX[0:len(p.PosX)-11]
	p.PosZ = p.PosZ[0:len(p.PosZ)-11]
	positionShow := p.Dim + " XYZ: " + p.PosX + " / " + p.PosY + " / " + p.PosZ
	tp := "/tp " + p.PosX + " " + p.PosY + " " + p.PosZ
	text := `/tellraw @a [{"text":"[Here]"},{"text":"`+c.Player+`","color":"gold"},{"text":"在"},{"text":"`+ positionShow +`","color":"gold",` +
	`"clickEvent":{"action":"run_command","value":"`+tp+`"}}` +
	`,{"text":"向大家打招呼！"}]`
	s.Execute(text)
	s.Execute("/effect give " + c.Player + " minecraft:glowing 30 1 true")

}

func GetPosition(playerName string, svr lib.Server, h *Here) *Here {
	reg := `(?P<player>.+)\s+has the following entity data:\s+\[(?P<x>.+)d,\s+(?P<y>.+)d,\s+(?P<z>.+)d\]`
	go func() {
		svr.Execute("/data get entity " + playerName + " Pos")
	}()
	match, flag := svr.RegParser(reg)
	if !flag {
		fmt.Println("获取坐标失败，请重新尝试。")
		return h
	}
	h.PosX = match[2]
	h.PosY = match[3]
	h.PosZ = match[4]
	return h
}

func GetDim(playerName string, svr lib.Server, h *Here) *Here {
	reg := `(?P<player>.+)\s+has the following entity data:\s+(?P<dim>.+)`
	go func() {
		svr.Execute("/data get entity " + playerName + " Dimension")
	}()
	match, flag := svr.RegParser(reg)
	if !flag {
		fmt.Println("获取维度失败，请重新尝试。")
		return h
	}
	h.Dim = match[2]
	return h
}

func (h *Here) Init(s lib.Server) {
}

func (h *Here) Close() {
}
