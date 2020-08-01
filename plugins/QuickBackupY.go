/**
 * 快速备份插件
 * 编写者：Underworld511
 */
package plugin

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/otiai10/copy"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/yliu7949/MCDaemon-go/command"
	"github.com/yliu7949/MCDaemon-go/config"
	"github.com/yliu7949/MCDaemon-go/lib"
)

type QuickBackupY struct {
	backupComment string
	name string
}

const Data = `{
  "Slot4": {
    "Flag": 0,
    "Name": "none",
    "Comment": "none"
  },
  "Slot3": {
    "Flag": 0,
    "Name": "none",
    "Comment": "none"
  },
  "Slot2": {
    "Flag": 0,
    "Name": "none",
    "Comment": "none"
  },
  "Slot1": {
    "Flag": 0,
    "Name": "none",
    "Comment": "none"
  }
}`

var qbIsMaking = false  //是否正在执行qb make命令

func (qb *QuickBackupY) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	dataFileName := "QuickBackup/qb_data.json"		//数据文件名字
	switch c.Argv[0] {
	case "make":
		qbIsMaking = true
		t := time.Now()
		nowDate := t.Format("2006-01-02")
		nowTime := t.Format("15:04:05")
		qb.name = nowDate + "@" + nowTime
		s.Say("开始快速备份...")
		s.Execute("/save-all flush")
		serverPath := config.Cfg.Section("MCDaemon").Key("server_path").String()
		if err := copy.Copy(serverPath, "QuickBackup/"+qb.name); err != nil {
			lib.WriteDevelopLog("error", err.Error())
		}
		if len(c.Argv) < 2 {
			qb.backupComment = "None"
		} else {
			qb.backupComment = c.Argv[1]
		}
		//生成或更新qb_data.json文件
		var result string
		if checkFileIsExist(dataFileName) {
			result = changeSlot(dataFileName, qb)      //如果数据文件已经存在，先将已有备份槽位后移，然后添加新的数据
		} else {
			result = newSlot(dataFileName, Data, qb)       //新的开始，新备份存放在槽位1
		}
		qbIsMaking = false
		s.Say(result+"耗时"+
			fmt.Sprintf("%v", time.Since(t).Truncate(time.Second)) + "。")
	case "back":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "缺少参数，请加上指定的槽位数字！")
		} else {
			slot, err := strconv.Atoi(c.Argv[1])
			if err != nil {
				s.Tell(c.Player, "参数不合法！")
				return
			}
			if !checkSlot(dataFileName,slot) {		//检查slot槽位
				s.Tell(c.Player,"参数不合法或无法找到有效存档！")
			} else {
				if !checkFileIsExist("./qbConfirm" + strconv.Itoa(slot)) {
					f, _ := os.Create("./qbConfirm" + strconv.Itoa(slot))					//用于回档确认
					f.Close()
				}
				s.Tell(c.Player, "已查询到指定存档，请于10秒钟内输入!!qb confirm确认回档[谨慎操作]。")
				timer := time.NewTimer(time.Second * 10)
				select {
				case <-timer.C:
					if checkFileIsExist("./qbConfirm" + strconv.Itoa(slot)) {
						_ = os.Remove("./qbConfirm" + strconv.Itoa(slot))
					}
				}
			}
		}
	case "confirm":
		result := false
		slot := 1
		for slot = 1; slot <= 4; slot++ {
			if checkFileIsExist("./qbConfirm" + strconv.Itoa(slot)) {
				result = true
				break
			}
		}
		if !result {
			s.Tell(c.Player, "确认无效！ 请先使用!!qb back命令。")
			return
		} else {
			_ = os.Remove("./qbConfirm" + strconv.Itoa(slot))
			s.Tell(c.Player, "已确认。")
			s.Say("12秒后将关闭服务器回档。")
			ticker1 := time.NewTicker(time.Millisecond * 100)	//创建打点器1，每100毫秒触发一次
			ticker2 := time.NewTicker(time.Second)				//创建打点器2，每秒触发一次
			stopper := time.NewTimer(time.Second * 12)			//创建一个计时器, 12秒后触发
			var i = 0
			for {
				// 多路复用通道
				select {
				case <-stopper.C:  // 计时器到时了
					goto BeginBack	// 跳出循环
				case <-ticker1.C:  // 打点器1触发了
					if checkFileIsExist("./qbAbout") {
						_ = os.Remove("./qbAbout")
						goto ReturnHere
					}
				case <-ticker2.C:  // 打点器2触发了
					i++
					s.Say("还有"+strconv.Itoa(12-i)+"秒，将回档为槽位" + strconv.Itoa(slot) +"。")
				}
			}
			ReturnHere:
				s.Say("已终止服务器回档。")
				return
			BeginBack:
				s.Say("开始回档...")		//开始回档
			s.Say("orz")
			lib.WriteDevelopLog("info", "服务器回档至槽位" + strconv.Itoa(slot))
			restoreBackUps(dataFileName, "QuickBackup/", slot, Data, s)	//回档函数
		}
	case "abort":
		f, _ := os.Create("./qbAbout")					//用于终止回档
		f.Close()
		timer := time.NewTimer(time.Second * 3)
		select {
		case <-timer.C:
			if checkFileIsExist("./qbAbout") {
				_ = os.Remove("./qbAbout")
			}
		}
	case "list":
		s.Tell(c.Player, listBackUps(dataFileName))
	case "tar":
		s.Tell(c.Player, "正在查询存档...")
		s.Tell(c.Player, tarBackUps(dataFileName))
	case "clean":
		s.Tell(c.Player, cleanBackUps(dataFileName, "QuickBackup/"))
	case "comment":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "缺少参数，请加上指定的槽位数字和新的存档注释！")
			return
		}
		if len(c.Argv) < 3{
			s.Tell(c.Player, "缺少参数，请同时加上指定的槽位数字和新的存档注释！")
		} else {
			slot, err := strconv.Atoi(c.Argv[1])
			if err != nil {
				s.Tell(c.Player, "参数不合法！")
				return
			}
			if !checkSlot(dataFileName,slot) {		//检查slot槽位
				s.Tell(c.Player,"参数不合法或无法找到有效存档！")
			} else {
				file, err := os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
				if err != nil {
					s.Tell(c.Player,"数据文件打开失败！")
					return
				}
				b, err := ioutil.ReadAll(file)
				if err != nil {
					s.Tell(c.Player,"数据读取失败！")
					return
				}
				defer file.Close()
				k := string(b)
				k, _ = sjson.Set(k, "Slot"+strconv.Itoa(slot)+".Comment", c.Argv[2])
				err = os.Remove(dataFileName)
				if err != nil {
					fmt.Printf("%s", err)
					s.Tell(c.Player,"文件删除失败！")
					return
				}
				file, err = os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
				if err != nil {
					s.Tell(c.Player,"数据文件打开失败！")
					return
				}
				defer file.Close()
				write := bufio.NewWriter(file)
				write.WriteString(k)
				err = write.Flush()
				if err != nil {
					s.Tell(c.Player,"注释修改失败！")
				}
				s.Tell(c.Player,"注释修改成功！")
			}
		}
	default:
		text := "使用规则：\\n!!qb make [<comment>] 制作位于槽位1的存档备份，并将已有备份槽位后移。<comment> 可选填存档注释。\\n" +
			"!!qb comment <slot> <comment> 修改指定槽位的存档注释。\\n!!qb tar 将槽位1的存档压缩并存储到 OBS 中。\\n!!qb list 显示各槽位的存档信息。\\n" +
			"!!qb clean 删除所有不再需要的存档。\\n!!qb back <slot> 回档为槽位 <slot> 的存档。\\n!!qb confirm 在执行 back 后使用，确认是否进行回档。\\n" +
			"!!qb abort 倒计时期间键入此指令可中断回档。\\n"
		s.Tell(c.Player, text)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func changeSlot(filename string, qb *QuickBackupY) string {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "数据文件打开失败"
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "数据读取失败"
	}
	defer file.Close()
	k := string(b)
	for i := 1; i <= 3; i++ {
		Data1, _ := sjson.Set(k, "Slot"+strconv.Itoa(5-i)+".Flag", gjson.Get(k, "Slot"+strconv.Itoa(4-i)+".Flag").String())
		Data2, _ := sjson.Set(Data1, "Slot"+strconv.Itoa(5-i)+".Name", gjson.Get(Data1, "Slot"+strconv.Itoa(4-i)+".Name").String())
		k, _ = sjson.Set(Data2, "Slot"+strconv.Itoa(5-i)+".Comment", gjson.Get(Data2, "Slot"+strconv.Itoa(4-i)+".Comment").String())
	}
	Data1, _ := sjson.Set(k, "Slot1.Flag", 1)
	Data2, _ := sjson.Set(Data1, "Slot1.Name", qb.name)
	k, _ = sjson.Set(Data2, "Slot1.Comment", qb.backupComment)
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("%s", err)
		return "文件删除失败!"
	}
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "数据文件打开失败"
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(k)
	err = write.Flush()
	if err != nil {
		return "写入备份失败"
	}
	return "备份制作完成。"
}

func newSlot(filename string, Data string, qb *QuickBackupY) string {
	Data1, _ := sjson.Set(Data, "Slot1.Flag", 1)
	Data2, _ := sjson.Set(Data1, "Slot1.Name", qb.name)
	Data, _ = sjson.Set(Data2, "Slot1.Comment", qb.backupComment)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "数据文件打开失败"
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(Data)
	err = write.Flush()
	if err != nil {
		return "写入备份失败"
	}
	return "备份制作完成。"
}

func checkSlot(filename string, slot int) bool {
	if !checkFileIsExist(filename) {
		return false
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("数据文件打开失败")
		return false
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("数据读取失败")
		return false
	}
	defer file.Close()
	k := string(b)
	if slot >= 1 && slot <= 4 {
		if gjson.Get(k, "Slot"+strconv.Itoa(slot)+".Flag").Bool() {
			return true
		}
	}
	return false
}

func listBackUps(filename string) string {
	if !checkFileIsExist(filename) {
		return "未查询到快速备份存档。"
	}
	file, err := os.Open(filename)
	if err != nil {
		return "数据文件打开失败。"
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "数据文件读取失败。"
	}
	defer file.Close()
	k := string(b)
	text := "QuickBackup存档列表：\\n"
	for i := 1; i <= 4; i++ {
		text += "<" + strconv.Itoa(i) + ">  "
		if gjson.Get(k, "Slot"+strconv.Itoa(i)+".Flag").Bool() {
			text += gjson.Get(k, "Slot"+strconv.Itoa(i)+".Name").String() + "  "
			text += gjson.Get(k, "Slot"+strconv.Itoa(i)+".Comment").String() + "\\n"
		} else {
			text += "<空>  <空>\\n"
		}
	}
	return text
}

func tarBackUps(filename string) string {
	if !checkFileIsExist(filename) {
		return "数据文件不存在！"
	}
	file, err := os.Open(filename)
	if err != nil {
		return "数据文件打开失败。"
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "数据文件读取失败。"
	}
	defer file.Close()
	k := string(b)
	if !gjson.Get(k, "Slot1.Flag").Bool() {
		return "槽位1无存档！"
	}
	tarFilesName := gjson.Get(k, "Slot1.Name").String()
	if checkFileIsExist("/OBS/" + tarFilesName + ".tar.gz") {
		return "OBS已存储最新存档，无需重复存储。"
	}
	cmd := exec.Command("tar", "zcvf", "/OBS/" + tarFilesName + ".tar.gz", "QuickBackup/" + tarFilesName)
	if err := cmd.Run(); err != nil {
		return "压缩姬出问题了，压缩失败！ "
	} else {
		return "压缩备份完成，存储成功！"
	}
}

func cleanBackUps(filename string, dir string) string {
	if qbIsMaking {
		return "备份期间无法执行清理操作！"
	}
	num := 0	//当前存档数
	num2 := 0	//删除的存档数
	names,err := filepath.Glob(filepath.Join(dir,"*"))		//获取指定目录下的文件名或目录名(包含路径)
	if err != nil {
		return "获取信息失败！"
	}
	if !checkFileIsExist(filename) {
		return "未查询到快速备份存档。"
	}
	file, err := os.Open(filename)
	if err != nil {
		return "数据文件打开失败。"
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "数据文件读取失败。"
	}
	defer file.Close()
	k := string(b)
	for i := range names {
		flag := false	//标志
		num++
		for j := 1; j <= 4; j++ {
			if names[i] == dir + gjson.Get(k, "Slot"+strconv.Itoa(j)+".Name").String() {
				flag = true
			}
		}
		if !flag && names[i] != filename {
			os.RemoveAll(names[i])
			num2++
		}
	}
	if num2 != 0 {
		return "已成功清理" + strconv.Itoa(num2) + "个不需要的存档。"
	} else {
		return "当前共有" + strconv.Itoa(num-1) + "个存档，均无需清理。"
	}
}

func restoreBackUps(filename string, dir string, slot int, Data string,s lib.Server) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("数据文件打开失败")
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("数据读取失败")
	}
	defer file.Close()
	k := string(b)
	restorePath := dir + gjson.Get(k, "Slot"+strconv.Itoa(slot)+".Name").String()
	s.Back(restorePath)		//调用了server/server.go 中的back方法
	if slot == 1 {
		return
	}
	for i := 1; i <= 5-slot; i++ {
		Data1, _ := sjson.Set(Data, "Slot"+strconv.Itoa(i)+".Flag", gjson.Get(k, "Slot"+strconv.Itoa(slot+i-1)+".Flag").String())
		Data2, _ := sjson.Set(Data1, "Slot"+strconv.Itoa(i)+".Name", gjson.Get(k, "Slot"+strconv.Itoa(slot+i-1)+".Name").String())
		Data, _ = sjson.Set(Data2, "Slot"+strconv.Itoa(i)+".Comment", gjson.Get(k, "Slot"+strconv.Itoa(slot+i-1)+".Comment").String())
	}
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("%s", err)
		fmt.Println("文件删除失败！")
	}
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("数据文件打开失败！")
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(Data)
	err = write.Flush()
	if err != nil {
		fmt.Println("数据文件更新失败！")
	}
	fmt.Println("数据文件更新成功。")
}

func (qb *QuickBackupY) Init(s lib.Server) {
}

func (qb *QuickBackupY) Close() {
}
