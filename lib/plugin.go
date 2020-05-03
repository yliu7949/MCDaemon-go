package lib

import "github.com/yliu7949/MCDaemon-go/command"

type Plugin interface {
	Handle(*command.Command, Server)
	Init(Server)
	Close()
}
