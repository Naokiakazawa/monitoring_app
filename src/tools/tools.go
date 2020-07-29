package tools

import (
    "log"
	"os"
	
	"golang.org/x/crypto/sha3"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CalcMethodID(method string) string {
    Signature := []byte(method)
    hash := sha3.NewLegacyKeccak256()
    hash.Write(Signature)
    methodID := hash.Sum(nil)[:4]
    return hexutil.Encode(methodID)
}

func Remove(s [][]string, i int) [][]string {
    if i >= len(s) {
        return s
    }
    return append(s[:i], s[i+1:]...)
}

func FailOnError(e error) {
	if e != nil{
		log.Fatal(e)
	}
}