package config

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

//配置变量
var (
	Cfg       *ini.File
	err       error
	plugins   map[string]string
	PluginCfg *ini.File
)

func init() {
	//加载配置文件
	Cfg, err = ini.Load("MCD_conf.ini")
	if err != nil {
		fmt.Printf("读取配置文件失败: %v", err)
		os.Exit(1)
	}
}

//获取服务器启动配置
func GetStartConfig() []string {
	//读取配置
	Section := Cfg.Section("MCDaemon")
	serverName := Section.Key("server_name").String()
	// serverPath := Section.Key("server_path").String()
	//设置默认值
	xms := Section.Key("Xms").Validate(func(in string) string {
		if len(in) == 0 {
			return "-Xms1024M"
		}
		return fmt.Sprint("-Xms", in)
	})
	xmx := Section.Key("Xmx").Validate(func(in string) string {
		if len(in) == 0 {
			return "-Xmx1024M"
		}
		return fmt.Sprint("-Xmx", in)
	})
	gui := Section.Key("gui").Validate(func(in string) string {
		if len(in) == 0 {
			return "false"
		}
		return in
	})
	agent := Section.Key("agent").Validate(func(in string) string {
		url := Section.Key("yggdrasil-url").Value()
		if len(in) == 0 {
			return ""
		}
		if len(url) == 0 {
			fmt.Println("未指定 yggdrasil-url！不会启用 -javaagent 参数")
			return ""
		}
		return fmt.Sprint("-javaagent:", in, "=", url)
	})
	var result []string
	if len(agent) == 0 {
		result = []string{
			xmx,
			xms,
			"-jar",
			serverName,
		}
	} else {
		result = []string{
			xmx,
			xms,
			agent,
			"-jar",
			serverName,
		}
	}
	if gui != "true" {
		result = append(result, "nogui")
	}
	return result
}

//获取插件配置
func GetPlugins(is_rebuild bool) map[string]string {
	if is_rebuild {
		Cfg, err = ini.Load("MCD_conf.ini")
		if err != nil {
			fmt.Printf("读取配置文件失败: %v", err)
			os.Exit(1)
		}
		//重置配置文件
		plugins = nil
	}
	if plugins == nil {
		plugins = make(map[string]string)
		keys := Cfg.Section("plugins").KeyStrings()
		for _, val := range keys {
			plugins[val] = Cfg.Section("plugins").Key(val).String()
		}
	}
	return plugins
}

//根据命令获取插件
func GetPluginName(cmd string) string {
	pluins := GetPlugins(false)
	return pluins[cmd]
}

//获取插件配置文件对象
func GetPluginCfg(is_rebuild bool) *ini.File {
	//加载配置文件
	if PluginCfg == nil || is_rebuild {
		PluginCfg, err = ini.ShadowLoad("Plugin_conf.ini")
		if err != nil {
			fmt.Printf("读取插件配置文件失败: %v", err)
		}
	}
	return PluginCfg
}

func SetEula() {
	path := fmt.Sprintf("%s/eula.txt", Cfg.Section("MCDaemon").Key("server_path").String())
	eulaCfg, eulaerr := ini.Load(path)
	//不存在eula.txt
	if eulaerr != nil {
		eulaCfg = ini.Empty()
		eulaCfg.Section("").NewKey("eula", "true")
		_ = eulaCfg.SaveTo(path)
	}
	//如果为false
	if eulaCfg.Section("").Key("eula").String() == "false" {
		eulaCfg.Section("").NewKey("eula", "true")
		_ = eulaCfg.SaveTo(path)
	}
}
