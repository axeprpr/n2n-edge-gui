//go:build windows

package main
 
import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
	// "bytes"
	"github.com/lxn/win"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"gopkg.in/ini.v1"
)

const (
    SIZE_W = 300
    SIZE_H = 200
)

type StWindow struct {
	*walk.MainWindow
	ni *walk.NotifyIcon
}

func SystemTray() *StWindow {
	st := new(StWindow)
	var err error
	st.MainWindow, err = walk.NewMainWindow()
	checkError(err)
	return st
}

func ConfigWindowMain() {
	var cf *walk.MainWindow
    var communityName, Address, Supernode, Mtu *walk.LineEdit
    communityname,address,supernode,mtu := GetConfig()
    MainWindow {
    	AssignTo: &cf,
    	Icon: "main.ico",
  		Title:   "AN Config",
        MaxSize: Size{SIZE_W, SIZE_H},
        Layout:  VBox{},
        Children: []Widget{
            Composite{
                Layout: Grid{Columns: 2, Spacing: 10},
                Children: []Widget{
                    VSplitter{
                        Children: []Widget{
                            Label{
                            	MinSize: Size{80, 20},
                            	MaxSize: Size{80, 20},
                                Text: "communityName",
                            },
                        },
                    },
                    VSplitter{
                        Children: []Widget{
                            LineEdit{
                            	MaxLength: 30,
                               	MinSize: Size{180, 20},
                               	MaxSize: Size{180, 20},                               	
                                AssignTo: &communityName,
                                Text: communityname,
                            },
                        },
                    },
                    VSplitter{
                        Children: []Widget{
                            Label{
                            	MinSize: Size{80, 20},
                            	MaxSize: Size{80, 20},
                                Text: "Address",
                            },
                        },
                    },
                    VSplitter{
                        Children: []Widget{
                            LineEdit{
                               	MinSize: Size{180, 20},
                               	MaxSize: Size{180, 20},
                                AssignTo: &Address,
                                Text: address,
                            },
                        },
                    },
                    VSplitter{
                        Children: []Widget{
							Label{
                            	MinSize: Size{80, 20},
                            	MaxSize: Size{80, 20},
                                Text: "Supernode",
                            },
                        },
                    },
                    VSplitter{
                        Children: []Widget{
                            LineEdit{
                               	MinSize: Size{180, 20},
                               	MaxSize: Size{180, 20},
                                AssignTo: &Supernode,
                                Text: supernode,
                            },
                        },
                    },
                    VSplitter{
                        Children: []Widget{
                        	Label{
                            	MinSize: Size{80, 20},
                            	MaxSize: Size{80, 20},
                                Text: "Mtu",
                            },
                        },
                    },
                    VSplitter{
                        Children: []Widget{
                            LineEdit{
                               	MinSize: Size{180, 20},
                               	MaxSize: Size{180, 20},
                                AssignTo: &Mtu,
                                Text: mtu,
                            },
                        },
                    },
                },
            },
			Composite{
				MaxSize: Size{0, 30},
				Layout:  HBox{},
				Children: []Widget{
            		PushButton{
                		Text:    "Confirm",
                		MaxSize: Size{280, 20},
                		MinSize: Size{280, 20},
                		OnClicked: func() {
                    		if communityName.Text() == "" {
                        		var tmp walk.Form
                        		walk.MsgBox(tmp, "Error", "communityName can not be empty.", walk.MsgBoxIconInformation)
                        		return
                    		}
                    		if Supernode.Text() == "" {
                        		var tmp walk.Form
                        		walk.MsgBox(tmp, "Error", "Supernode can not be empty.", walk.MsgBoxIconInformation)
                        		return
                    		}
							SetConfig(communityName.Text(),Address.Text(),Supernode.Text(),Mtu.Text())
                		},
                	},
            	},
			},
		},
    }.Create()

    xScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
    yScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
    win.SetWindowPos(
        cf.Handle(),
        0,
        (xScreen - SIZE_W)/2,
        (yScreen - SIZE_H)/2,
        SIZE_W,
        SIZE_H,
        win.SWP_FRAMECHANGED,
    )
    win.ShowWindow(cf.Handle(), win.SW_SHOW)
    cf.Run()
}

func (st *StWindow) AddNotifyIcon() {
	var err error
	st.ni, err = walk.NewNotifyIcon()
	checkError(err)
	st.ni.SetVisible(true)
 
	icon, err := walk.NewIconFromResourceId(3)
	checkError(err)
	st.SetIcon(icon)
	st.ni.SetIcon(icon)
 
 	configAction := st.addAction(nil, "config")
	startAction := st.addAction(nil, "start")
	stopAction := st.addAction(nil, "stop")
	stopAction.SetChecked(true)
	configAction.Triggered().Attach(func() {
		ConfigWindowMain()

	})
	startAction.Triggered().Attach(func() {
		start := n2nStart()
		if start {
			startAction.SetChecked(true)
			stopAction.SetChecked(false)
			// startAction.SetEnabled(false)
			// stopAction.SetEnabled(true)
		}else{
			st.msgbox("Error", "Start n2n edge failed. Please make sure the configuration is correct.", walk.MsgBoxIconError)
		}

	})
 
	stopAction.Triggered().Attach(func() {
		stop := n2nStop()
		if stop {
			startAction.SetChecked(false)
			stopAction.SetChecked(true)
		}else{
			st.msgbox("Error", "Stop n2n failed. Process is not exist.", walk.MsgBoxIconError)
		}
	})
 
	helpMenu := st.addMenu("help")
	st.addAction(helpMenu, "help").Triggered().Attach(func() {
		walk.MsgBox(st, "help", "https://github.com/ntop/n2n", walk.MsgBoxIconInformation)
	})
 
	st.addAction(helpMenu, "about").Triggered().Attach(func() {
		walk.MsgBox(st, "about", "xv.yayun@astute-tec.com", walk.MsgBoxIconInformation)
	})
 
	st.addAction(nil, "exit").Triggered().Attach(func() {
		st.ni.Dispose()
		st.Dispose()
		walk.App().Exit(0)
	})
 
}
 
func (st *StWindow) addMenu(name string) *walk.Menu {
	helpMenu, err := walk.NewMenu()
	checkError(err)
	help, err := st.ni.ContextMenu().Actions().AddMenu(helpMenu)
	checkError(err)
	help.SetText(name)
 
	return helpMenu
}
 
func (st *StWindow) addAction(menu *walk.Menu, name string) *walk.Action {
	action := walk.NewAction()
	action.SetText(name)
	if menu != nil {
		menu.Actions().Add(action)
	} else {
		st.ni.ContextMenu().Actions().Add(action)
	}
 
	return action
}
 
func (st *StWindow) msgbox(title, message string, style walk.MsgBoxStyle) {
	st.ni.ShowInfo(title, message)
	walk.MsgBox(st, title, message, style)
}
 
func n2nStart()(bool) {
	communityname,address,supernode,mtu := GetConfig()
	//edge -c astute -a 17.17.17.3 -l www.astute-tec.com:88 -M 1300
	if address != ""{
		cmd := exec.Command("n2n/edge.exe","-c",communityname,"-a",address,"-l",supernode,"-M",mtu)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Start()
		fmt.Println(err)
	}else{
		fmt.Println("DHCP enabled.")
		cmd := exec.Command("n2n/edge.exe","-c",communityname,"-r","-a","dhcp:0.0.0.0","-l",supernode,"-M",mtu)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Start()
		fmt.Println(err)
	}
	cmd_line := "tasklist -v | findstr /i edge.exe"
    err := exec.Command("cmd", "/c",cmd_line).Run()
	if err != nil {
		fmt.Println("Edge start failed.")
		return false
	}else{
		fmt.Println("Edge start success.")
		return true
	}
}

func n2nStop()(bool) {
	err := exec.Command("cmd", "/c","taskkill","/im","edge.exe","/f").Run()
	if err != nil {
		fmt.Println("Edge stop failed.")
		return false
	}else{
		fmt.Println("Edge stop success.")
		return true
	}
}

func GetConfig()(communityname,address,supernode,mtu string) {
	cfg, err := ini.Load("conf.ini")
    if err != nil {
        fmt.Println("Fail to read file: %v", err)
       	return
    }
	communityname = cfg.Section("DEFAULT").Key("communityname").String()
	address = cfg.Section("DEFAULT").Key("address").String()
	supernode = cfg.Section("DEFAULT").Key("supernode").String()
	mtu = cfg.Section("DEFAULT").Key("mtu").String()
	return
}

func SetConfig(communityname,address,supernode,mtu string) {
	cfg, err := ini.Load("conf.ini")
    if err != nil {
        fmt.Println("Fail to read file: %v", err)
       	return
    }
	cfg.Section("DEFAULT").Key("communityname").SetValue(communityname)
	cfg.Section("DEFAULT").Key("address").SetValue(address)
	cfg.Section("DEFAULT").Key("supernode").SetValue(supernode)
	cfg.Section("DEFAULT").Key("mtu").SetValue(mtu)

	err = cfg.SaveTo("conf.ini")
}

func main() {
	st := SystemTray()
	st.AddNotifyIcon()
	st.Run()
}
 
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
