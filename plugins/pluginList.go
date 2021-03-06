package plugin

import (
	"github.com/yliu7949/MCDaemon-go/plugins/ChatPlugin"
	"github.com/yliu7949/MCDaemon-go/plugins/Landmark"
)

func CreatePluginsList(isload bool) (PluginMap, PluginMap) {
	//可使用插件列表
	PluginsList := make(PluginMap)
	//已被禁用插件列表
	DisablePluginsList := make(PluginMap)

	//加载热插件
	PluginsList.GetHotPlugins(isload)

	//注册冷插件
	PluginsList.RegisterPlugin("!!server", &BasePlugin{})          //基础插件
	PluginsList.RegisterPlugin("!!yinyinmaster", &Yinyin{})        //例子插件
	PluginsList.RegisterPlugin("!!SDChat", &SDChatPlugin{})        //沙雕聊天机器人插件
	PluginsList.RegisterPlugin("!!tps", &TpsPlugin{})              //tps插件
	PluginsList.RegisterPlugin("!!Chat", &ChatPlugin.ChatPlugin{}) //跨服聊天插件
	PluginsList.RegisterPlugin("!!qb", &QuickBackupY{})            //快速备份插件
	PluginsList.RegisterPlugin("!!qi", &QuickImageY{})             //快速镜像插件
	PluginsList.RegisterPlugin("!!here", &Here{})                  //广播坐标插件
	PluginsList.RegisterPlugin("!!calc", &Calculator{})            //计算器插件
	PluginsList.RegisterPlugin("!!lm", &Landmark.Landmark{})       //坐标记录插件
	return PluginsList, DisablePluginsList
}
