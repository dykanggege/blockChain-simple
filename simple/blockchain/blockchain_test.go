package blockchain

import (
	"strconv"
	"testing"
)

func TestBlockChain_Iterator(t *testing.T) {
	bc := New()
	for i := 0; i < 10; i++ {
		bc.AddBlock("test" + strconv.Itoa(i))
	}
	iterater := bc.Iterator()
	for iterater.Has() {
		b := iterater.Next()
		if b == nil {
			break
		}
		b.Print()
	}
}
