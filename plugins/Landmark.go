/*
 * 坐标记录插件
 * 编写者：Underworld511
 */
package plugin

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
	. "github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type Landmark struct{}

func (lm *Landmark) Handle(c *Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	path := "data/Landmark/Landmark.ini"
	serverName := s.GetName()
	switch c.Argv[0] {
	case "add","a":
		if len(c.Argv) < 6 {
			s.Tell(c.Player, "参数不足!")
			return
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			f, _ := os.Create(path)
			defer f.Close()
		}
		cfg, _ := ini.Load(path)
		if section, err := cfg.NewSection(serverName); err != nil {
			s.Tell(c.Player, "Landmark数据文件错误：创建Section时出现了错误！")
		} else {
			dim := ""
			switch c.Argv[1] {
			case "world","w","0","overworld":
				dim = "minecraft:overworld"
			case "end","e","1","the_end":
				dim = "minecraft:the_end"
			case "nether","n","-1","the_nether":
				dim = "minecraft:the_nether"
			default:
				s.Tell(c.Player, "参数中的世界维度有误!")
				return
			}
			position := fmt.Sprintf("%s", c.Argv[2:5])
			_, _ = section.NewKey(c.Argv[5], dim + " " + position[1:len(position)-1])
			_ = cfg.SaveTo(path)
			s.Tell(c.Player, "坐标点添加成功!")
		}
	case "del","d":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "请加上要删除的坐标点的名字!")
			return
		}
		cfg, _ := ini.Load(path)
		cfg.Section(serverName).DeleteKey(c.Argv[1])
		_ = cfg.SaveTo(path)
		s.Tell(c.Player, "坐标点删除成功!")
	case "here","h":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "请加上要添加的坐标点的名字!")
			return
		}
		h := Here{}
		p := GetDim(c.Player, s, &h)
		p = GetPosition(c.Player, s,p)
		if p.PosZ == "" {
			s.Tell(c.Player, "添加坐标点失败！")
			return
		}
		cfg, _ := ini.Load(path)
		if section, err := cfg.NewSection(serverName); err != nil {
			s.Tell(c.Player, "Landmark数据文件错误：创建Section时出现了错误！")
		} else {
			_, _ = section.NewKey(c.Argv[1], fmt.Sprintf("%s %s %s %s", p.Dim, p.PosX, p.PosY, p.PosZ))
			_ = cfg.SaveTo(path)
			s.Tell(c.Player, "坐标点添加成功!")
		}
	case "rename","r":
		if len(c.Argv) < 3 {
			s.Tell(c.Player, "参数不足，请同时加上坐标点的旧名字和新名字!")
			return
		}
		cfg, _ := ini.Load(path)
		var value string
		if cfg.Section(serverName).HasKey(c.Argv[1]) {
			value = cfg.Section(serverName).Key(c.Argv[1]).Value()
			cfg.Section(serverName).DeleteKey(c.Argv[1])
			_, _ = cfg.Section(serverName).NewKey(c.Argv[2], value)
			_ = cfg.SaveTo(path)
			s.Tell(c.Player,"重命名成功。")
		} else {
			s.Tell(c.Player,"不存在这个坐标点！")
		}
	case "list","l":
		cfg, err := ini.Load(path)
		if err != nil {
			s.Tell(c.Player,"打开数据文件时出错！")
			return
		}
		landmarks := cfg.Section(serverName).KeyStrings()
		for n,landmark := range landmarks {
			value := cfg.Section(serverName).Key(landmark).Value()
			dim := strings.Split(value," ")[0]
			position := fmt.Sprintf("%s",strings.Split(value," ")[1:])
			position = position[1:len(position)-1]
			var hoverText string
			switch dim {
			case "minecraft:overworld":
				hoverText = "主世界：" + position
			case "minecraft:the_end":
				hoverText = "末地：" + position
			case "minecraft:the_nether":
				hoverText = "下界：" + position
			}
			spawnCmd := "/player Bot" + strconv.Itoa(n+1) + " spawn at " + position + " facing 0 0 in " + dim
			killCmd := "/player Bot" + strconv.Itoa(n+1) + " kill"
			renameCmd := "!!lm rename " + landmark + " "
			delCmd := "!!lm del " + landmark
			s.Tell(c.Player,"§e<"+strconv.Itoa(n+1)+">",
				MinecraftText("§7"+landmark).SetHoverText("§6"+hoverText).SetClickEvent("run_command","[分享] "+landmark+"： "+hoverText),
				MinecraftText("§6§l [☻]").SetClickEvent("run_command",spawnCmd).SetHoverText("§6单击召唤工具人"),
				MinecraftText("§c§l[☠]").SetClickEvent("run_command",killCmd).SetHoverText("§6单击退出工具人"),
				MinecraftText("§3§l[✎]").SetClickEvent("suggest_command",renameCmd).SetHoverText("§6单击重命名坐标点"),
				MinecraftText("§9§l[⤫]").SetClickEvent("run_command",delCmd).SetHoverText("§6单击删除坐标点"))
		}
		if len(landmarks) < 1 {
			s.Tell(c.Player,"无记录的坐标点！")
		}
	case "help":
		s.Tell(c.Player,"§e§l⚑命令指南")
		s.Tell(c.Player,MinecraftText("!!lm add [end/nether/world] <x> <y> <z> <name> §6添加坐标点\\n").
			SetClickEvent("suggest_command","!!lm add "),
			MinecraftText("!!lm here <name> §6添加当前位置为坐标点\\n").SetClickEvent("suggest_command","!!lm here "),
			MinecraftText("!!lm rename <oldName> <newName> §6重命名坐标点\\n").SetClickEvent("suggest_command","!!lm rename "),
			MinecraftText("!!lm del <name> §6删除坐标点\\n").SetClickEvent("suggest_command","!!lm del "),
			MinecraftText("!!lm list §6列出所有坐标点\\n").SetClickEvent("run_command","!!lm list"))
		s.Tell(c.Player,"§e§l⚑当前在线工具人")
		s.Tell(c.Player,"§7§o该功能正在开发中...\\n")
		s.Tell(c.Player,"§e§l⚑快捷操作")
		s.Tell(c.Player,"§7§o该功能正在开发中...\\n")
	}
}

func (lm *Landmark) Init(s lib.Server) {
	if _,err := os.Stat("./data/Landmark"); os.IsNotExist(err) {
		if err := os.Mkdir("./data/Landmark", os.ModePerm); err != nil {
			fmt.Println("创建Landmark数据文件夹失败：", err)
			return
		}
		if _, err := os.Stat("data/Landmark/Landmark.ini"); os.IsNotExist(err) {
			f, _ := os.Create("data/Landmark/Landmark.ini")
			defer f.Close()
		}
	}
}

func (lm *Landmark) Close() {
}