package engine

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/yqchilde/pkgs/log"

	"github.com/lxlxw/go-wxbot/engine/robot"
)

const layout = "2006-01-02 15:04:05"

var duration = time.Minute * 15

func InitRobot() {
	start := time.Now()
	// 使用桌面方式登录
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 关闭心跳回调
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		if resp.RetCode == "1100" {
			end := time.Now()
			s := end.Sub(start)
			log.Println(fmt.Sprintf("微信已退出%s至%s，持续：%s", start, end, s))
		}
		switch resp.Selector {
		case "0":
			log.Println("正常")
		case "2", "6":
			log.Println("有新消息")
		case "7":
			log.Println("进入/离开聊天界面")
			err := bot.WebInit()
			if err != nil {
				log.Println(fmt.Sprintf("重新初始化失败：%v", err))
			}
		default:
			log.Println(fmt.Sprintf("RetCode：%s Selector: %s", resp.RetCode, resp.Selector))
		}
	}

	// 登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 开启热登录
	reloadStorage := &robot.JsonLocalStorage{FileName: "storage.json"}
	if err := bot.HotLogin(reloadStorage, true); err != nil {
		panic(err)
	}

	// 处理消息回调
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsSendBySelf() {
			return
		}

		reply := "Bobo Bot Beta🤖 欢迎您使用\n\n"
		for _, plugin := range Plugins {
			if plugin.RawConfig["enable"] != false {
				plugin.Config.OnEvent(&robot.Message{Message: msg})
			}
			if !plugin.HiddenMenu {
				reply += plugin.Desc + "\n"
			}
		}

		if msg.IsText() && msg.Content == "/menu" {
			msg.ReplyText(reply)
		}
		if msg.IsSendByFriend() {
			sender, err := msg.Sender()
			if err != nil {
				log.Printf("get friend chat sender error: %v", err)
				return
			}

			if msg.IsText() {
				log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", sender.NickName, msg.Content))
			} else {
				log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", sender.NickName, msg.String()))
			}
		} else {
			sender, err := msg.SenderInGroup()
			if err != nil {
				log.Printf("get group chat sender error: %v", err)
				return
			}
			if msg.IsText() {
				log.Println(fmt.Sprintf("收到群(%s[%s])消息 ==> %v", getGroupNicknameByGroupUsername(msg.FromUserName, sender.NickName, sender.RemarkName), sender.NickName, msg.Content))
			} else {
				log.Println(fmt.Sprintf("收到群(%s[%s])消息 ==> %v", getGroupNicknameByGroupUsername(msg.FromUserName, sender.NickName, sender.RemarkName), sender.NickName, msg.String()))
			}
		}
	}

	var count int32
	bot.MessageErrorHandler = func(err error) bool {
		atomic.AddInt32(&count, 1)
		if count == 5 {
			bot.Logout()

		}
		return true
	}

	// 获取登陆的用户
	if self, err := bot.GetCurrentUser(); err == nil {
		robot.Self = self
	} else {
		panic(err)
	}

	// 获取所有的好友
	if friends, err := robot.Self.Friends(true); err != nil {
		panic(err)
	} else {
		robot.Friends = friends
	}

	// 获取所有的群组
	if groups, err := robot.Self.Groups(true); err != nil {
		panic(err)
	} else {
		robot.Groups = groups

		log.Println(robot.Groups)
	}

	bot.Block()
}

func getGroupNicknameByGroupUsername(username string, nickname string, remarkname string) string {
	groups := robot.Groups.SearchByUserName(1, username)
	if groups == nil {
		groups = robot.Groups.SearchByNickName(1, nickname)
		if groups != nil {
			return groups[0].NickName
		} else {
			groups = robot.Groups.SearchByRemarkName(1, remarkname)
			return ""
		}
	} else {
		return groups[0].NickName
	}
}
