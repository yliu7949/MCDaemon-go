package lib

import "github.com/yliu7949/MCDaemon-go/command"

type Parser interface {
	Parsing(string) (*command.Command, bool)
}

