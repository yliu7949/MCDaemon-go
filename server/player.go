package server

import (
	"strconv"
	"time"
)


type Player struct {
	Dim 				string			//玩家所在的维度，"0":"主世界","-1":"地狱","1":"末地"
	PosX				string			//玩家的x坐标
	PosY            	string			//玩家的y坐标
	PosZ				string			//玩家的z坐标
	onlineTime			int 			//玩家本次登录后的在线时长
	totalOnlineTime		int				//玩家总在线时长
}

var pl Player

func (svr *Server ) GetPosition(playerName string) (string, string, string) {
	reg := `\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+(?P<player>.+)\s+has the following entity data:\s+\[(?P<x>.+)d,\s+(?P<y>.+)d,\s+(?P<z>.+)d\]`
	go func() {
		time.Sleep(2e8)
		svr.Execute("/data get entity " + playerName + " Pos")
	}()
	match, flag := svr.RegParser(reg)
	if !flag {
		return "", "", ""
	}
	pl.PosX = match[2]
	pl.PosY = match[3]
	pl.PosZ = match[4]
	return pl.PosX, pl.PosY, pl.PosZ
}

func (svr *Server ) GetDim(playerName string) string {
	reg := `\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+(?P<player>.+)\s+has the following entity data:\s+(?P<dim>.+)`
	go func() {
		time.Sleep(2e8)
		svr.Execute("/data get entity " + playerName + " Dimension")
	}()
	match, flag := svr.RegParser(reg)
	if !flag {
		return ""
	}
	pl.Dim = match[2]
	return pl.Dim
}

