package parser

import (
	"regexp"

	"github.com/yliu7949/MCDaemon-go/command"
)

type SeenParser struct{}

func (Sp *SeenParser) Parsing(word string) (*command.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+(?P<player>.+)\s+`)
	match := re.FindStringSubmatch(word)
	if len(match) != 0 {
		_commond := &command.Command{
			Cmd:  "!!autobk",
			Argv: []string{"save"},
		}
		return _commond, true
	}
	return nil, false
}
