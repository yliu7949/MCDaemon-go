/**
websocket客户端
*/
package ChatPlugin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/yliu7949/MCDaemon-go/lib"
	"golang.org/x/net/websocket"
)

var clientID = 10000

type WSServerClient struct {
	ServerId int //服务器id
	Conn     *websocket.Conn
}

type WSServer struct {
	ServerId        int    //服务器id
	ServerName      string //服务器名称
	Port            int
	Suburl          string                     //子路由
	ReceiveMessage  chan *msgPackage           //接受到的消息
	SendMessage     chan *Message              //要发送的消息
	origin          string                     //源地址
	minecraftServer lib.Server                 //服务器实例接口
	ConnPool        map[string]*WSServerClient //连接池，键为服务器名称，值为封装的websocket连接对象
	RWPool          *sync.RWMutex              //连接池读写锁
	WhiteList       map[string]interface{}     //白名单
	Ctx             context.Context            //上下文
	Cancel          context.CancelFunc
	Alive           bool //是否是活跃连接
}

func (WSS *WSServer) handler(conn *websocket.Conn) {
	fmt.Println("聊天服务器连接成功")
	defer conn.Close()
	var err error
	clientID = clientID + 1
	Cilent := &WSServerClient{clientID, conn}
	for {
		var reply []byte
		err = websocket.Message.Receive(conn, &reply)
		if err != nil {
			lib.WriteDevelopLog("error", fmt.Sprint("聊天服务器出错：", err))
			break
		}
		//将proto消息解码
		newMessage := &Message{}
		err = proto.Unmarshal(reply, newMessage)
		if err != nil {
			lib.WriteDevelopLog("warn", fmt.Sprint("非法连接：", conn.RemoteAddr().String()))
			break
		}
		serverName := newMessage.GetServerName()
		//加入到连接池中,若不在聊天白名单中，则关闭连接
		if ok := WSS.appendToConnPool(serverName, Cilent); !ok {
			data, _ := proto.Marshal(&Message{
				ServerName: &LocalServerName,
				State:      &NotInWhitelist,
			})
			websocket.Message.Send(conn, data)
			lib.WriteDevelopLog("warn", fmt.Sprint("不在白名单中：", serverName))
			break
		}
		fmt.Println("1---", newMessage)
		if newMessage.GetState() == FirstTouch {
			fmt.Println("1---", newMessage)
			continue
		}
		//将消息发送给其他连接
		clientMsg := &msgPackage{Cilent.ServerId, newMessage}
		WSS.SendtoClient(clientMsg)
		//将消息加入到接收管道中
		ServerMsg := &msgPackage{LocalServerId, newMessage}
		WSS.ReceiveMessage <- ServerMsg
	}
}

//向连接池里的所有连接发送消息
func (WSS *WSServer) SendtoClient(msg *msgPackage) {
	//编码
	data, _ := proto.Marshal(msg.Msg)
	for serverName, client := range WSS.ConnPool {
		//如果serverid相同，则不发送
		if client.ServerId == msg.From {
			continue
		}
		//若出现错误，则从连接池中删除并关闭这条连接
		fmt.Println("向子连接发送", client.ServerId, "------", msg.From)
		if err := websocket.Message.Send(client.Conn, data); err != nil {
			WSS.deletePool(serverName)
			client.Conn.Close()
			break
		}
	}
}

//接收要发送的消息
func (WSS *WSServer) Send(msg *Message) {
	WSS.SendtoClient(&msgPackage{LocalServerId, msg})
}

func (WSS *WSServer) Read() {
	for {
		select {
		case <-WSS.Ctx.Done():
			fmt.Println("server context over!")
			return
		case msg := <-WSS.ReceiveMessage:
			packageChan <- msg
		}
	}
}

//将websocket连接加入到连接池中
func (WSS *WSServer) appendToConnPool(serverName string, Client *WSServerClient) bool {
	if _, ok := WSS.WhiteList[serverName]; ok {
		//如果没有进入连接池，则加入到连接池中
		if _, ok := WSS.readPool(serverName); !ok {
			WSS.writePool(serverName, Client)
		}
		return true
	}
	return false
}

func (WSS *WSServer) readPool(serverName string) (*WSServerClient, bool) {
	WSS.RWPool.RLock()
	defer WSS.RWPool.RUnlock()
	if val, ok := WSS.ConnPool[serverName]; ok {
		return val, true
	}
	return nil, false
}

func (WSS *WSServer) writePool(serverName string, Client *WSServerClient) {
	WSS.RWPool.Lock()
	defer WSS.RWPool.Unlock()
	WSS.ConnPool[serverName] = Client
}

func (WSS *WSServer) deletePool(serverName string) {
	WSS.RWPool.Lock()
	defer WSS.RWPool.Unlock()
	delete(WSS.ConnPool, serverName)
}

func (WSS *WSServer) Start() error {
	url := "localhost:" + strconv.Itoa(WSS.Port)
	http.Handle("/"+WSS.Suburl, websocket.Handler(WSS.handler))
	WSS.ConnPool = make(map[string]*WSServerClient)
	lib.WriteDevelopLog("info", fmt.Sprint("websocket启动，连接地址：", url, "/", WSS.Suburl))
	go http.ListenAndServe(url, nil)
	WSS.Alive = true
	return nil
}

func (WSS *WSServer) GetId() int {
	return WSS.ServerId
}

func (WSS *WSServer) GetName() string {
	return WSS.ServerName
}

func (WSS *WSServer) IsAlive() bool {
	if WSS.Alive {
		return true
	}
	WSS.Cancel()
	return false
}
