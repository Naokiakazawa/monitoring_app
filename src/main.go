package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"sync"
	
	"golang.org/x/crypto/sha3"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gosuri/uiprogress"
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
		Address:			c.Section("address").Key("tether").String(),
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

func remove(s [][]string, i int) [][]string {
    if i >= len(s) {
        return s
    }
    return append(s[:i], s[i+1:]...)
}

func calcMethodID(method string) string {
    Signature := []byte(method)
    hash := sha3.NewLegacyKeccak256()
    hash.Write(Signature)
    methodID := hash.Sum(nil)[:4]
    return hexutil.Encode(methodID)
}

func txsearcher(i int64, wg *sync.WaitGroup, client *ethclient.Client) {
	defer wg.Done()
	Block_Number := big.NewInt(i)
	block, err := client.BlockByNumber(context.Background(), Block_Number)
	FailOnError(err)
	if block != nil && block.Transactions() != nil {
		txs := block.Transactions()
		webhook("Block Number: " + block.Number().String() + "\n" + "Transaction counts:  " + strconv.Itoa(len(txs)))
		bar := uiprogress.AddBar(len(txs)).AppendCompleted().PrependElapsed()
		bar.Fill = '='
		bar.Head = '>'
		bar.Empty = ' '
		bar.Width = 50
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("#%s (%d/%d)", strconv.FormatInt(i, 10), b.Current(), len(txs))
		})

		records := [][]string{
			[]string{"Block_Number", "Tx_Hash", "fromAddress", "Value", "toAddress"},
		}
			for _, tx := range txs {
				bar.Incr()
				if tx.To() == nil {
					continue
				} else if tx.To().Hex() == Config.Address && hexutil.Encode(tx.Data()[:4]) == calcMethodID("transfer(address,uint256)") {
					tx_hash := tx.Hash().Hex()

					chainID, err := client.NetworkID(context.Background())
					FailOnError(err)
					msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
					FailOnError(err)

					toAddress_byte := tx.Data()[16:36]
					Value_byte := common.TrimLeftZeroes(tx.Data()[36:])

					fromAddress := msg.From().Hex()
					toAddress := hexutil.Encode(toAddress_byte)
					Value := hexutil.Encode(Value_byte)

					items := []string{block.Number().String(), tx_hash, fromAddress, Value, toAddress}
					records = append(records, items)
				}
			}
		file_path := "record_" + strconv.FormatInt(i, 10) + ".csv"
		file, err := os.Create(file_path)
		FailOnError(err)

		w := csv.NewWriter(file)
		err = w.WriteAll(records)
		FailOnError(err)
		w.Flush()
		file.Close()
	}
}

func joincsv(file_name string, start_block int64, end_block int64) {
	joined_file, err := os.Create(file_name)
	FailOnError(err)
	joined_records := [][]string{
		[]string{"Block_Number", "Tx_Hash", "fromAddress", "Value", "toAddress"},
	}
	for i := start_block; i < end_block; i++ {
		file_path := "record_" + strconv.FormatInt(i, 10) + ".csv"
		record_file, err := os.Open(file_path)
		FailOnError(err)
		read := csv.NewReader(record_file)
		d, err := read.ReadAll()
		FailOnError(err)
		record_file.Close()
		d = remove(d, 0)
		joined_records = append(joined_records, d...)
	}
	write := csv.NewWriter(joined_file)
	err = write.WriteAll(joined_records)
	FailOnError(err)
	write.Flush()
	joined_file.Close()
}

func main() {
	var start_block int64 = 10523035
	var end_block int64 = 10523041
	var RecordFilePath string = "record.csv"

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + Config.ProjectID)
	FailOnError(err)
	log.Println("connected!!")

	uiprogress.Start()
	var wg sync.WaitGroup
	for i := start_block; i < end_block; i++ {
		wg.Add(1)
		go txsearcher(i, &wg, client)
	}
	wg.Wait()
	uiprogress.Stop()
	log.Println("success!!")

	joincsv(RecordFilePath, start_block, end_block)
}