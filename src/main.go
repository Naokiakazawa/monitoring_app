package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/ini.v1"
	
	"app/slack"
)

type ConfigList struct {
	Slack_url			string
	Slack_channel		string
	Slack_name			string
	Slack_icon_emoji	string
	ProjectID			string
	Address				string
}

func FailOnError(e error) {
	if e != nil{
		log.Fatal(e)
	}
}

var Config ConfigList

func init() {
	c, err := ini.Load("config.ini")
	FailOnError(err)
	Config = ConfigList{
		Slack_url:			c.Section("slack").Key("url").String(),
		Slack_channel:		c.Section("slack").Key("channel").MustString("general"),
		Slack_name:			c.Section("slack").Key("name").MustString("slack bot"),
		Slack_icon_emoji:	c.Section("slack").Key("icon_emoji").String(),
		ProjectID:			c.Section("infura").Key("projectid").String(),
		Address:			c.Section("address").Key("0x").String(),
	}
}

func webhook(text string) {
	channel := Config.Slack_channel
	name := Config.Slack_name
	icon_emoji := Config.Slack_icon_emoji
	webhook := "https://hooks.slack.com/services/" + Config.Slack_url

	err := slack.SlackWebhook(channel, name, text, icon_emoji, webhook)
	FailOnError(err)
}

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + Config.ProjectID)
	FailOnError(err)
	log.Println("connected!!")

	contractAddress := common.HexToAddress(Config.Address)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6383482),
		ToBlock: big.NewInt(6384488),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	FailOnError(err)
	log.Printf("Log counts: %d", len(logs))
	fmt.Println("\n---\n")
	webhook(strconv.Itoa(len(logs)))

	logFillEvent_hash := crypto.Keccak256Hash([]byte("LogFill(address,address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes32)")).Hex()
	logFillEvent := common.HexToHash(logFillEvent_hash)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Block Hash: %s\n", vLog.BlockHash.Hex())
		fmt.Printf("Log Transaction Hash: %s\n", vLog.TxHash.Hex())
		fmt.Printf("Log Index: %d\n", vLog.Index)

        if vLog.Topics[0].Hex() == logFillEvent.Hex(){
			transaction, status, err := client.TransactionByHash(context.Background(), vLog.TxHash)
			FailOnError(err)
			fmt.Printf("Status: %t\n", status)
			fmt.Println(hexutil.Encode(transaction.Data()))
		} else {
			fmt.Println("skip!")
		}
        fmt.Printf("\n")
    }
}