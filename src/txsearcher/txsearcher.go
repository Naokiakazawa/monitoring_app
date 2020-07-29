package txsearcher

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gosuri/uiprogress"

	"app/tools"
)

func Searcher(i int64, addr string, wg *sync.WaitGroup, client *ethclient.Client) {
	defer wg.Done()
	Block_Number := big.NewInt(i)
	block, err := client.BlockByNumber(context.Background(), Block_Number)
	tools.FailOnError(err)
	if block != nil && block.Transactions() != nil {
		switch {
		case len(block.Transactions()) == 0:
			log.Printf("#%s is No Transactions", strconv.FormatInt(i, 10))
		default:
			txs := block.Transactions()
//			webhook("Block Number: " + block.Number().String() + "\n" + "Transaction counts:  " + strconv.Itoa(len(txs)))
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
				} else if tx.To().Hex() == addr && hexutil.Encode(tx.Data()[:4]) == tools.CalcMethodID("transfer(address,uint256)") {
					tx_hash := tx.Hash().Hex()

					chainID, err := client.NetworkID(context.Background())
					tools.FailOnError(err)
					msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
					tools.FailOnError(err)

					toAddress_byte := tx.Data()[16:36]
					Value_byte := common.TrimLeftZeroes(tx.Data()[36:])

					fromAddress := msg.From().Hex()
					toAddress := hexutil.Encode(toAddress_byte)
					Value := hexutil.Encode(Value_byte)

					items := []string{block.Number().String(), tx_hash, fromAddress, Value, toAddress}
					records = append(records, items)
				}
			}
			file_path := "data/record_" + strconv.FormatInt(i, 10) + ".csv"
			file, err := os.Create(file_path)
			tools.FailOnError(err)
	
			w := csv.NewWriter(file)
			err = w.WriteAll(records)
			tools.FailOnError(err)
			w.Flush()
			file.Close()
		}
	}
}
