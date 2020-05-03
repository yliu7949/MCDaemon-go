package parser

import (
	"regexp"

	"github.com/yliu7949/MCDaemon-go/command"
)

type AutoBackupParser struct{}

func (p *AutoBackupParser) Parsing(word string) (*command.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s\S*\sleft the game`)
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
