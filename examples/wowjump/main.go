package main

import (
	"github.com/lxn/win"
	"github.com/whtiehack/wingui"
	"github.com/whtiehack/wingui/kernel32"
	"log"
	"os"
	"syscall"
)

var dlg *wingui.Dialog
var out *wingui.Edit

// 防止进程多开，返回 true 表示进程已经开启
func ProcessMutex(name string) bool {
	r, err := kernel32.OpenMutex(2031617, 1, name)
	if err == nil || r != 0 {
		return true
	}
	r, err = kernel32.CreateMutex(name)
	if err != nil || r == 0 {
		return true
	}
	return false
}

func init() {
	log.SetOutput(os.Stdout)
	if ProcessMutex("@@wowjump@@") {
		win.MessageBox(0, &syscall.StringToUTF16("进程已经开启了，不可以多开")[0], nil, 0)
		os.Exit(-1)
	}
	// control
	go process()
}

var btn *wingui.Button

// optional  genereate resource IDs
//go:generate go run github.com/whtiehack/wingui/tools/genids -filename ui/resource.h -packagename main
func main() {
	var err error
	dlg, err = wingui.NewDialog(IDD_DIALOG_MAIN, 0, nil)
	if err != nil {
		log.Panic("main dialog create error", err)
	}
	dlg.SetIcon(IDI_ICON1)
	editLog := wingui.NewEdit(IDE_LOG)
	out = editLog
	btn = wingui.NewButton(IDB_RUN)
	btn.OnClicked = btnClick
	SetLogOutput(editLog)
	config.editNormaltime = wingui.NewEdit(IDE_NORMAL_TIME)
	config.editEnterTime = wingui.NewEdit(IDE_ENTER_TIME)
	config.editInputTime = wingui.NewEdit(IDE_INPUT_TIME)
	config.editCharWaitTime = wingui.NewEdit(IDE_CHAR_WAIT_TIME)
	dlg.AddWidgets([]wingui.Widget{editLog, btn, config.editNormaltime, config.editEnterTime, config.editInputTime, config.editCharWaitTime})

	config.InitVal()
	dlg.Show()
	// Make sure Tabstop can work.
	wingui.SetCurrentDialog(dlg.Handle())
	wingui.MessageLoop()
	log.Println("stoped")
}

type MyLogoutput struct {
	edit  *wingui.Edit
	count int
}

func (m *MyLogoutput) Write(p []byte) (n int, err error) {
	n, err = os.Stdout.Write(p)
	m.count++
	if m.count >= 1000 {
		m.edit.SetText("")
	}
	m.edit.AppendText(string(p) + "\r\n")
	return
}
func SetLogOutput(edit *wingui.Edit) {
	m := &MyLogoutput{edit: edit}
	log.SetOutput(m)
}

func btnClick() {
	if clicking {
		return
	}
	clicking = true
	defer func() { clicking = false }()
	mux.Lock()
	defer mux.Unlock()
	config.Save()
	running = !running
	//move.Retset()
	//jump.Reset()
	var text = "开启"
	if running {
		logouts = GetMultiLogout()
		if len(logouts) == 0 {
			running = !running
			log.Println("没有找到WOW窗口")
			return
		}
		out.SetText("")
		log.Printf("找到 %d 个 WOW窗口\n", len(logouts))
		text = "停止"
		//str := skillKey.Text()
		//config.SkillKey = str
		//randomSkill.ParseSkillKey(str)
		log.Println("开始运行")
	} else {
		log.Println("已经停止运行")
	}
	config.EditEnable(!running)
	btn.SetText(text)
}