package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"testing"
)

func TestLearn(t *testing.T) {
	const data = "000000000000000000000000f87c133c1642f6c5aefc8889aefbed6d7d144257" +
		"0000000000000000000000000000000000000000000000000000000000000005" +
		"0000000000000000000000000000000000000000000000000000000000000065" +
		"000000000000000000000000ee4a73cf0cbe6e850e7be821aeb3a7382d2c02c5"

	m := decode(data)
	spew.Dump(m)

}

func decode1(data string) (m map[string]interface{}) {
	// load contract ABI
	m = make(map[string]interface{})
	locked, err := abi.JSON(strings.NewReader(LockedABI))
	if err != nil {
		log.Fatal(err)
	}

	if err := locked.UnpackIntoMap(m, "Locked", Hex2Bytes(data)); err != nil {
		return nil
	}

	return m
}
