package plugin

import (
	"fmt"
	"strings"

	. "github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type Here struct {
	Dim 				string			//玩家所在的维度，"0":"主世界","-1":"地狱","1":"末地"
	PosX				string			//玩家的x坐标
	PosY				string			//玩家的y坐标
	PosZ				string			//玩家的z坐标
}


func (h *Here) Handle(c *Command, s lib.Server) {
	p := GetDim(c.Player, s, h)
	switch p.Dim {
	case "0", `"minecraft:overworld"`:
		p.Dim = "主世界"
	case "-1", `"minecraft:the_nether"`:
		p.Dim = "下界"
	case "1", `"minecraft:the_end"`:
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
	positionShow := p.Dim + " XYZ: " + p.PosX + " / " + p.PosY + " / " + p.PosZ
	tp := "/tp " + p.PosX + " " + p.PosY + " " + p.PosZ
	s.Say("§b"+c.Player,"在",MinecraftText("§6"+positionShow).SetClickEvent("run_command",tp),"向大家打招呼！")
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
	h.PosX = strings.Split(h.PosX, ".")[0] + "." + strings.Split(h.PosX, ".")[1][:1]
	h.PosY = strings.Split(h.PosY, ".")[0] + "." + strings.Split(h.PosY, ".")[1][:1]
	h.PosZ = strings.Split(h.PosZ, ".")[0] + "." + strings.Split(h.PosZ, ".")[1][:1]
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