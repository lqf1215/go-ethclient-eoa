package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go-ethclient-eoa/library"
	"math/big"
	"os"
	"sync"
	"time"
)

func main() {
	// 设置连接参数
	rpcClient, err := rpc.Dial("https://mainnet.infura.io/v3/2252464df647405c84738d5264b57ce7")
	if err != nil {
		panic(err)
	}
	client := ethclient.NewClient(rpcClient)

	ctx := context.Background()
	fromAddr := common.HexToAddress("0xAADEd943BeF6115DaB272C90ef3a1E01a358336d")
	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		panic(err)
	}
	// 打开文件
	file, err := os.Create("go-for-eth.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	line := fmt.Sprintf("交易时间 发送方 接收方 余额(ETH)  \n")
	_, err = file.WriteString(line)
	if err != nil {
		panic(err)
	}
	// 设置并发数量
	concurrency := 80

	// 一般来说 没有固定的block区块需要从0开始遍历，但是我这里为了跳过 从网上浏览器 该地址 查询到记录最小的区块
	block := int64(11780605)
	// 创建 waitGroup 以等待所有协程执行完成
	var wg sync.WaitGroup
	wg.Add(concurrency)
	type transaction struct {
		ts   *types.Transaction
		time uint64
	}
	// 创建通道用于存储交易记录
	txCh := make(chan transaction)
	endBlockNumber := big.NewInt(int64(latestBlock))
	// 循环启动多个协程，每个协程处理一段区块的交易记录
	blockRange := endBlockNumber.Int64() / int64(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(start, end int64) {
			if start < block && end < block {
				return
			}
			fmt.Println(start, end)
			defer wg.Done()
			for i := start; i < end; i++ {
				block, err := client.BlockByNumber(ctx, big.NewInt(i))
				if err != nil {
					continue
				}
				for _, tx := range block.Transactions() {
					if tx.To() == nil {
						continue
					}
					if tx.To().Hex() != fromAddr.Hex() {
						continue
					}

					// 将交易记录发送到通道中
					txCh <- transaction{ts: tx, time: block.Time()}
				}
			}
		}(blockRange*int64(i), blockRange*(int64(i)+1))
	}

	// 等待所有协程执行完成后关闭通道
	go func() {
		wg.Wait()
		close(txCh)
	}()

	// 遍历通道读取交易记录并输出到文件
	for tx := range txCh {
		// 获取发送者地址
		from, err := types.Sender(types.NewEIP155Signer(tx.ts.ChainId()), tx.ts)
		if err != nil {
			fmt.Println("获取发送者地址失败：", err)
			return
		}
		// 获取交易时间
		tm := time.Unix(int64(tx.time), 0)
		datetime := tm.Format("2006-01-02 15:04:05")
		// 获取以太数量
		value := tx.ts.Value().String()
		to := *tx.ts.To()
		// 写入文件
		line := fmt.Sprintf("%s %s %s %v\n", datetime, from.Hex(), to.Hex(), library.EtherConvertAmount(value))
		fmt.Printf("line：%x\n", line)
		_, err = file.WriteString(line)
		if err != nil {
			panic(err)
		}
	}
}
