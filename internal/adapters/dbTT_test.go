package adapters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBTT(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		err := recover()
		assert.Nil(err)
	}()

	t.Run("dbTT", func(t *testing.T) {
		/* 		const (
		   			u32 uint32 = 1<<32 - 1
		   			i32 int32  = 1<<31 - 1
		   		)
		   		arr32 := [u32]bool{}
		   		// arr64 := [u64]int{}

		   		x := 2
		   		switch x {
		   		case 1:
		   			fmt.Println(arr32)
		   		case 2:
		   			//fmt.Println(arr64)
		   		}

		   		fmt.Printf("%b\n", u32)
		   		fmt.Printf("%b\n", i32) */
	})
}
