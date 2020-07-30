package csvutil

import (
	"encoding/csv"
	"os"
	"strconv"

	"app/tools"
)

func Joincsv(file_name string, start_block int64, end_block int64) {
	joined_file, err := os.Create(file_name)
	tools.FailOnError(err)
	joined_records := [][]string{
		[]string{"Block_Number", "Tx_Hash", "fromAddress", "Value", "toAddress"},
	}
	for i := start_block; i <= end_block; i++ {
		file_path := "data/record_" + strconv.FormatInt(i, 10) + ".csv"
		if tools.IsExist(file_path) == false {
			continue
		}
		record_file, err := os.Open(file_path)
		tools.FailOnError(err)
		read := csv.NewReader(record_file)
		d, err := read.ReadAll()
		tools.FailOnError(err)
		record_file.Close()
		d = tools.Remove(d, 0)
		joined_records = append(joined_records, d...)
	}
	write := csv.NewWriter(joined_file)
	err = write.WriteAll(joined_records)
	tools.FailOnError(err)
	write.Flush()
	joined_file.Close()
}
