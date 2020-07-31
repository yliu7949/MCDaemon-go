package command

import "fmt"

type MText []string

/*
设置点击事件函数，共有5个action。
1.run_command
2.suggest_command
3.insertion
4.open_url
5.copy_to_clipboard(1.15+)
*/
func (mt *MText) SetClickEvent(action string, value string) *MText {
	var aimJsonList []string
	for _,Json := range *mt {
		Json = Json[:len(Json)-1] + `,"clickEvent":{"action":"` + action + `","value":"` + value + `"}}`
		aimJsonList = append(aimJsonList,Json)
	}
	*mt = aimJsonList
	return mt
}

func (mt *MText) SetHoverText(value string) *MText {
	var aimJsonList []string
	for _,Json := range *mt {
		Json = Json[:len(Json)-1] + `,"hoverEvent":{"action":"show_text","value":"` + value + `"}}`
		aimJsonList = append(aimJsonList,Json)
	}
	*mt = aimJsonList
	return mt
}

//对文字的颜色和样式进行json编码
func MinecraftText(text string) *MText {
	var (
		aimJsonList []string//存放所有颜色作用域对应的json文本的列表
		aimJson string		//存放一个颜色作用域内对应的json文本
		colorField string	//颜色作用域，即颜色控制符和样式控制符起作用的字符串
		color string
		styles []string
	)
	text += "§"	//在字符串末尾添加一个"§"
	colorDict := map[string]string{
		"§0":"black",
		"§1":"dark_blue",
		"§2":"dark_green",
		"§3":"dark_aqua",
		"§4":"dark_red",
		"§5":"dark_purple",
		"§6":"gold",
		"§7":"gray",
		"§8":"dark_gray",
		"§9":"blue",
		"§a":"green",
		"§b":"aqua",
		"§c":"red",
		"§d":"light_purple",
		"§e":"yellow",
		"§f":"white",
		"§r":"reset",
	}
	styleDict := map[string]string{
		"§l":"bold",
		"§o":"italic",
		"§n":"underlined",
		"§m":"strikethrough",
		"§k":"obfuscated",
	}
	frontCharSpecial := false //若前一个字符是'§'则为true
	for _,letter := range text {
		if letter == '§' {
			frontCharSpecial = true
			if colorField != "" {
				aimJson = fmt.Sprintf("{\"text\":\"%s\"",colorField)
				if color != "" {
					aimJson += fmt.Sprintf(",\"color\":\"%s\"",color)
				}
				if len(styles) != 0 {
					for _,style := range styles {
						aimJson += fmt.Sprintf(",\"%s\":\"true\"",style)
					}
				}
				aimJson += "}"
				colorField = ""
				color = ""
				styles = nil
				aimJsonList = append(aimJsonList,aimJson)
			}
			continue
		}
		if frontCharSpecial == false {
			colorField += string(letter)
		} else {
			frontCharSpecial = false
			if colorDict["§"+string(letter)] != "" {
				color = colorDict["§"+string(letter)]
			}
			if styleDict["§"+string(letter)] != "" {
				styles = append(styles,styleDict["§"+string(letter)])
			}
		}
	}
	return (*MText)(&aimJsonList)
}