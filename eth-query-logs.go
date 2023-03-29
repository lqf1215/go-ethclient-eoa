package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/2252464df647405c84738d5264b57ce7")
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress("0xAADEd943BeF6115DaB272C90ef3a1E01a358336d")
	fromBlock := uint64(0)
	toBlock := uint64(99999999)

	query := buildQuery(account, fromBlock, toBlock)
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("logs=", logs)
	for _, log := range logs {
		fmt.Println(log.BlockNumber, log.TxHash.Hex(), log.Address.Hex(), log.Topics[0].Hex(), log.Data)
	}
}

func buildQuery(address common.Address, fromBlock, toBlock uint64) ethereum.FilterQuery {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{address},
	}
	return query
}
