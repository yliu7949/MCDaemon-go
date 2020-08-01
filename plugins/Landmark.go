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
			case "world","w","0":
				dim = "0"
			case "end","e","1":
				dim = "1"
			case "nether","n","-1":
				dim = "-1"
			default:
				s.Tell(c.Player, "参数中的世界维度有误!")
				return
			}
			_, _ = section.NewKey(c.Argv[5], dim + " " + fmt.Sprintf("%s", c.Argv[2:5]))
			cfg.SaveTo(path)
			s.Tell(c.Player, "坐标点添加成功!")
		}
	case "del","d":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "请加上要删除的坐标点的名字!")
			return
		}
		cfg, _ := ini.Load(path)
		cfg.Section(serverName).DeleteKey(c.Argv[1])
		cfg.SaveTo(path)
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
		if _, err := os.Stat(path); os.IsNotExist(err) {
			f, _ := os.Create(path)
			defer f.Close()
		}
		cfg, _ := ini.Load(path)
		if section, err := cfg.NewSection(serverName); err != nil {
			s.Tell(c.Player, "Landmark数据文件错误：创建Section时出现了错误！")
		} else {
			_, _ = section.NewKey(c.Argv[1], fmt.Sprintf("%s %s %s %s", p.Dim, p.PosX, p.PosY, p.PosZ))
			cfg.SaveTo(path)
			s.Tell(c.Player, "坐标点添加成功!")
		}
	case "list","l":
		cfg, _ := ini.Load(path)
		landmarks := cfg.Section(serverName).KeyStrings()
		for n,landmark := range landmarks {
			value := cfg.Section(serverName).Key(landmark).Value()
			dim := strings.Split(value," ")[0]
			position := strings.Split(value," ")[1][0:]
			var hoverText string
			switch dim {
			case "0":
				hoverText = "主世界 " + position
				dim = "overworld"
			case "-1":
				hoverText = "下界 " + position
				dim = "the_nether"
			case "1":
				hoverText = "末地 " + position
				dim = "the_end"
			}
			spawnCmd := "/player Bot" + strconv.Itoa(n) + " spawn at " + position + " facing 0 0 in " + dim
			killCmd := "/player Bot" + strconv.Itoa(n) + " kill"
			s.Tell(c.Player,"§7<"+strconv.Itoa(n)+">",
				MinecraftText(landmark).SetHoverText("§6"+hoverText).SetClickEvent("copy_to_clipboard",landmark+" "+hoverText),
				MinecraftText("☻").SetClickEvent("run_command",spawnCmd),
				MinecraftText("☠").SetClickEvent("run_command",killCmd))
		}
	case "help":
		text := "使用规则：\\n!!lm add [end/nether/world] <x> <y> <z> <name> 添加坐标点\\n" +
			"!!lm here <name> 添加此时的位置为坐标点\\n!!lm del <name> 删除坐标点\\n" +
			"!!lm list 列出坐标点\\n"
		s.Tell(c.Player, text)
	}
}

func (lm *Landmark) Init(s lib.Server) {
	if _,err := os.Stat("./data/Landmark"); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("./data/Landmark", os.ModePerm)
			if err != nil {
				fmt.Println("创建Landmark数据文件夹失败：", err)
			}
			return
		}
	}
}

func (lm *Landmark) Close() {
}