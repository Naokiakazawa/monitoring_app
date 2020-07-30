package txsearcher

import (
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Dispatch(WORKER_COUNT int, START_BLOCK int64, END_BLOCK int64, address string, wg *sync.WaitGroup, client *ethclient.Client) {
	block_count := END_BLOCK - START_BLOCK + 1
	jobs := make(chan int64, block_count)

	for w := 1; w <= WORKER_COUNT; w++ {
		go Worker(w, jobs, address, wg, client)
	}

	for j := START_BLOCK; j <= END_BLOCK; j++ {
		wg.Add(1)
		jobs <- j
	}
	defer close(jobs)
}