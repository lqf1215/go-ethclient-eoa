package main

import (
	"encoding/json"
	"fmt"
	"go-ethclient-eoa/library"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	ApiUrl   = "https://api.etherscan.io/api"
	StatusOk = "1"
)

type EtherscanResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Result  []EtherscanResult `json:"result"`
}

type EtherscanResult struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
	MethodId          string `json:"methodId"`
	FunctionName      string `json:"functionName"`
}

func main() {

	apikey := "GP395NYAIB17SG9DTHJY9HFQJDQN9E5CAH"
	address := "0xAADEd943BeF6115DaB272C90ef3a1E01a358336d"
	module := "account"
	startBlock := big.NewInt(0)
	endBlock := big.NewInt(99999999)
	page := 1
	offset := 10
	sort := "desc" // asc  升序 desc 降序

	var eths []EtherscanResult
	var resp *http.Response
	var err error

	for {
		url := fmt.Sprintf("%s?module=%s&action=txlist&address=%s&startblock=%s&endblock=%s&page=%d&offset=%d&sort=%s&apikey=%s", ApiUrl, module, address, startBlock, endBlock, page, offset, sort, apikey)
		fmt.Println("url = ", url)
		resp, err = http.Get(url)
		if err != nil {
			fmt.Println("Something went wrong err = ", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ioutil.ReadAll err = ", err)
		}

		var ethResp EtherscanResponse
		err = json.Unmarshal(body, &ethResp)
		if err != nil {
			fmt.Println("Json反序列化错误 err=", err)
		}

		if ethResp.Status == StatusOk {
			eths = append(eths, ethResp.Result...)

			if len(ethResp.Result) >= offset {
				page++
			} else {
				break
			}
		} else {
			break
		}
	}
	defer resp.Body.Close()
	// 文件名
	filename := "http-eth.txt"

	// 打开文件
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	line := fmt.Sprintf("交易时间 发送方 接收方 余额(ETH)  \n")
	_, err = file.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, eth := range eths {
		timeStamp, err := strconv.Atoi(eth.TimeStamp)
		if err != nil {
			log.Fatal("字符转换整数失败 err = ", err)
		}
		datetime := time.Unix(int64(timeStamp), 0).Format("2006-01-02 15:04:05")

		line := fmt.Sprintf("%s %s %s %v ETH\n", datetime, eth.From, eth.To, library.EtherConvertAmount(eth.Value))
		fmt.Println(line)
		_, err = file.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}

	}
}
