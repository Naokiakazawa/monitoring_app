package main

import (
	"log"
	"os"
	"sync"
	
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gosuri/uiprogress"
	"gopkg.in/ini.v1"

	"app/csvutil"
	"app/slack"
	"app/tools"
	"app/txsearcher"
)

type ConfigList struct {
	Slack_url			string
	Slack_channel		string
	Slack_name			string
	Slack_icon_emoji	string
	ProjectID			string
	Address				string
}

var Config ConfigList

func init() {
	c, err := ini.Load("config.ini")
	tools.FailOnError(err)
	Config = ConfigList{
		Slack_url:			c.Section("slack").Key("url").String(),
		Slack_channel:		c.Section("slack").Key("channel").MustString("general"),
		Slack_name:			c.Section("slack").Key("name").MustString("slack bot"),
		Slack_icon_emoji:	c.Section("slack").Key("icon_emoji").String(),
		ProjectID:			c.Section("infura").Key("projectid").String(),
		Address:			c.Section("address").Key("tether").String(),
	}
}

func webhook(text string) {
	channel := Config.Slack_channel
	name := Config.Slack_name
	icon_emoji := Config.Slack_icon_emoji
	webhook := "https://hooks.slack.com/services/" + Config.Slack_url

	err := slack.SlackWebhook(channel, name, text, icon_emoji, webhook)
	tools.FailOnError(err)
}

func main() {
	var WORKER_COUNT int = 5
	var START_BLOCK int64 = 10523000
	var END_BLOCK int64 = 10523030
	var RECORD_FILE_PATH string = "record.csv"
	var ADDR string = Config.Address

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + Config.ProjectID)
	tools.FailOnError(err)
	log.Println("connected!!")

	uiprogress.Start()
	var wg sync.WaitGroup

	err = os.Mkdir("data", 0777)
	tools.FailOnError(err)

	txsearcher.Dispatch(WORKER_COUNT, START_BLOCK, END_BLOCK, ADDR, &wg, client)
	wg.Wait()
	uiprogress.Stop()

	csvutil.Joincsv(RECORD_FILE_PATH, START_BLOCK, END_BLOCK)
	log.Println("success!!")
}