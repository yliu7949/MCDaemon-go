package lib

import (
	"github.com/yliu7949/MCDaemon-go/command"
)

type Server interface {
	Say(...interface{})                      //在服务器中全局说话
	Tell(string, ...interface{})             //对某人私聊
	Execute(string)                          //执行mc原生命令
	Close()                                  //关闭服务器（适用于单服务器）
	CloseInContainer()                       //在容器中关闭服务器
	Back(string)				 //回档
	Restart()                                //重启
	Start(string, []string, string)          //开启一个服务器实例
	Clone(string) Server                     //克隆一个服务器实例
	GetPort() string                         //获取服务器端口
	ReloadConf()                             //重新获取配置
	RunPlugin(*command.Command)              //运行一个插件命令（结构体格式）
	RunPluginCommand(string, string)		 //运行一个插件命令（字符串格式）
	RunUniquePlugin(func())                  //在所有命令完成后执行
	WriteLog(level string, msg string)       //写入日志
	GetPluginList() map[string]Plugin        //获取可用插件列表
	GetDisablePluginList() map[string]Plugin //获取不可用插件列表
	GetParserList() []Parser                 //获取语法解释器列表
	RegParser(reg string) ([]string, bool)	 //从内存读取字符串并使用正则表达式提取信息
	GetName() string                         //获取服务器名称
	GetStartResult() int					 //获取服务器启动结果
}
