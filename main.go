package main

import (
	"GOBLOCKCHAIN/blockchain"
	"GOBLOCKCHAIN/transaction"
	"fmt"
)

func main() {
	// 创建一个空的交易池
	txPool := make([]*transaction.Transaction, 0)
	// 声明临时交易和其他变量
	var tempTx *transaction.Transaction
	var ok bool
	var property int

	// 创建一个区块链
	chain := blockchain.CreateBlockChain()
	// 查询 "Always" 的余额
	property, _ = chain.FindUTXOs([]byte("Always"))
	fmt.Println("Balance of Always: ", property)

	// 创建一笔交易，从 "Always" 转账给 "KFC"
	tempTx, ok = chain.CreateTransaction([]byte("Always"), []byte("ManBa"), 50)
	if ok {
		// 如果交易创建成功，将其添加到交易池中
		txPool = append(txPool, tempTx)
	}
	// 将交易池中的交易打包到新的区块中，并添加到区块链中
	chain.Mine(txPool)
	// 清空交易池
	txPool = make([]*transaction.Transaction, 0)
	// 再次查询 "Leo Cao" 的余额
	property, _ = chain.FindUTXOs([]byte("Always"))
	fmt.Println("Balance of Always: ", property)

	// 创建一笔无效的交易，从 "KFC" 转账给 "MDL"
	tempTx, ok = chain.CreateTransaction([]byte("ManBa"), []byte("KFC"), 100) // 交易金额不足 交易失败
	if ok {
		// 如果交易创建成功，将其添加到交易池中
		txPool = append(txPool, tempTx)
	}

	// 创建两笔有效的交易
	tempTx, ok = chain.CreateTransaction([]byte("Always"), []byte("ManBa"), 200)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("ManBa"), []byte("KFC"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}
	// 将交易池中的交易打包到新的区块中，并添加到区块链中
	chain.Mine(txPool)
	// 清空交易池
	txPool = make([]*transaction.Transaction, 0)
	// 再次查询 "Always" 的余额
	property, _ = chain.FindUTXOs([]byte("Always"))
	fmt.Println("Balance of Always: ", property)
	// 查询其他用户的余额
	property, _ = chain.FindUTXOs([]byte("ManBa"))
	fmt.Println("Balance of ManBa: ", property)
	property, _ = chain.FindUTXOs([]byte("KFC"))
	fmt.Println("Balance of KFC: ", property)

	// 遍历区块链中的每个区块，并打印信息
	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}

	// 尝试创建两笔额度不足的交易
	tempTx, ok = chain.CreateTransaction([]byte("ManBa"), []byte("Always"), 300)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("KFC"), []byte("Always"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}

	// 将交易池中的交易打包到新的区块中，并添加到区块链中
	chain.Mine(txPool)
	// 清空交易池
	txPool = make([]*transaction.Transaction, 0)

	// 再次遍历区块链中的每个区块，并打印信息
	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}

	// 再次查询用户的余额
	property, _ = chain.FindUTXOs([]byte("Always"))
	fmt.Println("Balance of Always: ", property)
	// 查询其他用户的余额
	property, _ = chain.FindUTXOs([]byte("ManBa"))
	fmt.Println("Balance of ManBa: ", property)
	property, _ = chain.FindUTXOs([]byte("KFC"))
	fmt.Println("Balance of KFC: ", property)
}
————————————————

版权声明：本文为博主原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接和本声明。

原文链接：https://blog.csdn.net/Hock2023/article/details/137382553
