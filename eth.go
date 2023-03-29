package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-ethclient-eoa/library"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	// infura url
	ethNodeUrl = "https://mainnet.infura.io/v3/2252464df647405c84738d5264b57ce7"
	// 文件名
	outputFileName = "eth.txt"
)

func main() {
	client, err := ethclient.Dial(ethNodeUrl)
	if err != nil {
		fmt.Println("连接以太坊节点失败：", err)
		return
	}
	// 起始区块号
	blockNumber := big.NewInt(11780605)

	// 获取区块链上最新的区块号
	lastBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		fmt.Println("获取最新区块号失败：", err)
		return
	}
	endBlockNumber := big.NewInt(int64(lastBlock))

	// EOA 地址
	address := common.HexToAddress("0xAADEd943BeF6115DaB272C90ef3a1E01a358336d")

	// 打开文件
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	line := fmt.Sprintf("交易时间 发送方 接收方 余额(ETH)  \n")
	_, err = file.WriteString(line)
	if err != nil {
		panic(err)
	}
	// 遍历区块
	for i := blockNumber.Int64(); i <= endBlockNumber.Int64(); i++ {
		// 获取最新区块的所有交易
		block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			fmt.Println("获取区块信息失败：", err)
			return
		}

		for _, tx := range block.Transactions() {
			if tx.To() != nil && *tx.To() == address {
				// 获取发送者地址
				from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
				if err != nil {
					fmt.Println("获取发送者地址失败：", err)
					return
				}
				// 获取交易时间
				tm := time.Unix(int64(block.Time()), 0)
				datetime := tm.Format("2006-01-02 15:04:05")
				// 获取以太数量
				value := tx.Value().String()
				to := *tx.To()
				// 写入文件
				line := fmt.Sprintf("%s %s %s %v\n", datetime, from.Hex(), to.Hex(), library.EtherConvertAmount(value))
				fmt.Printf("line：%x\n", line)
				_, err = file.WriteString(line)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	fmt.Printf("Done! Results written to %s\n", outputFileName)

}
