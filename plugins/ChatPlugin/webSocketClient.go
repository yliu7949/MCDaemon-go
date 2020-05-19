package ChatPlugin

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/yliu7949/MCDaemon-go/lib"
	"golang.org/x/net/websocket"
)

type WSClient struct {
	ServerId       int    //服务器id
	ServerName     string //服务器名称
	Addr           string
	Origin         string
	ws             *websocket.Conn //websocket连接
	ReceiveMessage chan *Message   //接受到的消息
	Ctx            context.Context //上下文
	Cancel         context.CancelFunc
}

func (WSC *WSClient) Start() error {
	var err error
	WSC.ws, err = websocket.Dial("ws://"+WSC.Addr, "", WSC.Origin)
	if err != nil {
		lib.WriteDevelopLog("error", err.Error())
		return err
	}
	defer WSC.ws.Close()
	defer lib.WriteDevelopLog("error", "连接")
	WSC.Send(&Message{
		ServerName: &LocalServerName,
		State:      &FirstTouch,
	})
	for {
		msg := make([]byte, 5096)
		slen, err := WSC.ws.Read(msg) //此处阻塞，等待有数据可读
		msg = msg[:slen]
		if err != nil {
			lib.WriteDevelopLog("error", fmt.Sprint("读取错误：", err))
			//如果连接出错，则释放连接
			break
		}
		newMessage := &Message{}
		err = proto.Unmarshal(msg, newMessage)
		if err != nil {
			lib.WriteDevelopLog("error", fmt.Sprint("解码错误:", err, "内容：", msg))
			break
		}
		if newMessage.GetState() != 0 {
			lib.WriteDevelopLog("error", "聊天服务器连接失败：不再白名单内！")
			break
		}
		WSC.ReceiveMessage <- newMessage
	}
	WSC.Cancel()
	return err
}

func (WSC *WSClient) Send(msg *Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		lib.WriteDevelopLog("error", fmt.Sprint("加密错误：", err))
		return
	}
	err = websocket.Message.Send(WSC.ws, data)
	if err != nil {
		lib.WriteDevelopLog("error", fmt.Sprint(WSC.GetName, "发送信息错误：", err))
	}
}

func (WSC *WSClient) Read() {
	for {
		select {
		case <-WSC.Ctx.Done():
			return
		case msg := <-WSC.ReceiveMessage:
			packageChan <- &msgPackage{
				From: WSC.ServerId,
				Msg:  msg,
			}
		}
	}
}

func (WSC *WSClient) GetId() int {
	return WSC.ServerId
}

func (WSC *WSClient) GetName() string {
	return WSC.ServerName
}

func (WSC *WSClient) IsAlive() bool {
	if WSC.ws != nil {
		return true
	}
	WSC.Cancel()
	return false
}
