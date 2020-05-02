package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"github.com/alfredxing/calc/compute"
	"strconv"
	"strings"

)

type Calculator struct {}

func (ca *Calculator) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) != 0 {
		input := strings.Replace(strings.Join(c.Argv[0:], ""), " ", "", -1)
		res, err := compute.Evaluate(input)
		if err != nil {
			s.Tell(c.Player, "计算器出现了错误。")
			return
		}
		s.Tell(c.Player, strconv.FormatFloat(res, 'G', -1, 64))
		return
	}
	s.Tell(c.Player, "请输入要计算的表达式")
}

func (ca *Calculator) Init(s lib.Server) {
}

func (ca *Calculator) Close() {
}