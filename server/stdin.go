package server

import (
	"fmt"
	"io"

	c "github.com/yliu7949/MCDaemon-go/command"
)

func (svr *Server) Say(argv ...interface{}) {
	TellGhost(argv...)
	svr.Tell("@a", argv...)
}

func (svr *Server) Tell(player string, argv ...interface{}) {
	if player == "Ghost" {		//从后台运行插件命令时Tell函数仅会将消息输出至后台
		TellGhost(argv...)
		return
	}
	var (
		jsonLists []c.MText
		command string
	)
	for _, v := range argv {
		switch t := v.(type) {
		case string:
			jsonLists = append(jsonLists,*c.MinecraftText(t))
		case *c.MText:
			jsonLists = append(jsonLists,*t)
		default:
			fmt.Println("Tell函数出错：不支持的消息类型!")
		}
	}
	if len(jsonLists) == 1 {
		command = fmt.Sprintf("%s",jsonLists[0])
		command = "tellraw " + player + " " + command[1:len(command)-1]
		svr.Execute(command)
	} else {
		command = "tellraw " + player + ` ["",`
		for _,json := range jsonLists {
			_command := fmt.Sprintf("%s",json)
			command +=  _command[1:len(_command)-1]+ ","
		}
		command =command[:len(command)-1] + "]"
		svr.Execute(command)
	}
}

func (svr *Server) Execute(_command string) {
	//输入的命令要换行！否则无法执行
	_command = _command + "\n"
	//同步写入
	svr.lock.Lock()
	defer svr.lock.Unlock()
	_, err := io.WriteString(svr.stdin, _command)
	if err != nil {
		 fmt.Println("Execute函数出错：", err)
	}
}

//将参数传入的消息去掉交互和颜色样式，并在后台(Ghost)中输出。
func TellGhost(argv ...interface{}) {
	fmt.Println("[SYSTEM]")
	for _, v := range argv {
		switch t := v.(type) {
		case string:
			for i := 1;i < len(t);i++ {
				if t[i] == '\\' && t[i+1] == 'n' {
					fmt.Println(t[:i])
					t = t[i+2:]
					i = 1
				}
			}
			if t != "" {
				fmt.Println(t)
			}
		default:
			fmt.Println("不支持的消息类型。")
		}
	}
}