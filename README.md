# go-ethclient-eoa

采用http api方式调用https://api.etherscan.io/api

``go run http-eth.go``

采用rpc 或 ethclient.Dial 连接 以太坊 https://mainnet.infura.io/v3/

单个 效率低 等待久了点

``go run eth.go``

多个go协程 采用channel通道 

``go run go-for-eth.go``


以我了解 大概是这样子
当然 可以塔私链 进行转账或者其他操作