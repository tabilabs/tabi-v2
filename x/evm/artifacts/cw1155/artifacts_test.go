package cw1155_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tabilabs/tabi-v2/x/evm/artifacts/cw1155"
)

// run with `-race`
func TestGetBinConcurrent(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			require.NotEmpty(t, cw1155.GetBin())
		}(i)
	}

	wg.Wait()
}
