package parser

import (
	"regexp"
	"strings"

	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/config"
)

//解析玩家输入文字，判断是否是命令 ， 实现了Parser接口
type defaultParser struct{}

//默认语法解析器
func (c *defaultParser) Parsing(word string) (*command.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+<(?P<player>.+)>\s+(?P<commond>((!|!!)+.+))\s*`)
	match := re.FindStringSubmatch(word)
	groupNames := re.SubexpNames()

	result := make(map[string]string)

	//匹配到命令时
	if len(match) != 0 {
		// 转换为map
		for i, name := range groupNames {
			result[name] = match[i]
		}

		// 解析命令以及参数
		cmdArgv := strings.Fields(result["commond"])
		_commond := &command.Command{
			Player: result["player"],
			Cmd:    cmdArgv[0],
			Argv:   cmdArgv[1:],
		}
		//获取插件名称
		_commond.PluginName = config.GetPluginName(cmdArgv[0])
		return _commond, true
	}
	return nil, false
}
