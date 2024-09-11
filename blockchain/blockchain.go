package blockchain

import (
	"GOBLOCKCHAIN/transaction"
	"GOBLOCKCHAIN/utils"
	"encoding/hex"
	"fmt"
)

type BlockChain struct {
	Blocks []*Block //blocks是该结构体中的一个 区块类型 的成员变量
}

// 下面两个函数 就是添加进链的信息换成了交易内容 其余不变
func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	return &blockchain
}

// 用于查找指定地址（address）的未花费交易（unspent transactions）
// 该函数会遍历区块链中的每一个区块，然后对每个区块中的交易进行处理，找出未被花费的交易输出 UTXO
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	//初始化一个空的未花费交易数组 unSpentTxs，用于存储找到的未花费交易
	spentTxs := make(map[string][]int)
	//初始化一个 spentTxs 映射，用于跟踪已经被花费的交易输出
	//该映射的键是交易 ID 的十六进制表示，值是一个整数切片，表示对应交易中已经被花费的输出索引

	//外部循环
	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		//从最新的区块开始向前遍历区块链中的每一个区块
		block := bc.Blocks[idx] //现在，我们进入了一个区块中
		//对于每个区块，遍历其中的每一笔交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //将交易ID转换为十六进制字符串格式

		IterOutputs: //判断是不是第一个交易
			for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					//检查交易ID没有在已花费列表中 下列循环
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs //如果发现已经存在，就不执行ToAddress判断
						}
					}
				}

				if out.ToAddressRight(address) { //判断是不是第一个交易
					unSpentTxs = append(unSpentTxs, *tx)
					//检查输出是否属于指定地址，如果是，则将该交易添加到未花费交易列表中。
				}
			} //简单来说 就是找到属于某个用户的交易

			//如果不是第一个交易 进行判断
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
						//检查输入是否来自指定地址，如果是，则将该输入标记为已花费
					}
				}
			}
		}
	}
	//返回找到的未花费交易数组
	return unSpentTxs
}

// 用于查找指定地址（address）的未花费交易输出（Unspent Transaction Outputs，UTXOs）的总金额和对应的输出索引
func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	// 创建一个 map 用于存储未花费的输出
	unspentOuts := make(map[string]int)
	// 查找包含未花费输出的交易
	unspentTxs := bc.FindUnspentTransactions(address)
	// 初始化累计变量
	accumulated := 0

Work:
	// 遍历未花费交易
	for _, tx := range unspentTxs {
		// 将交易ID转换为字符串形式
		txID := hex.EncodeToString(tx.ID)
		// 遍历交易的输出
		for outIdx, out := range tx.Outputs {
			// 检查输出是否属于指定地址
			if out.ToAddressRight(address) {
				// 如果是，则累加输出的价值到累计变量中
				accumulated += out.Value
				// 将交易ID和输出索引添加到未花费输出的 map 中
				unspentOuts[txID] = outIdx
				// 继续到标签为 Work 的循环，因为一笔交易只能有一个输出与地址相关联
				continue Work
			}
		}
	}
	// 返回累计的价值和未花费输出的 map
	return accumulated, unspentOuts
}

// 用于查找指定地址（address）可以用来支付指定金额（amount）的未花费交易输出（Unspent Transaction Outputs，UTXOs）的总金额和对应的输出索引。
func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	// 创建一个 map 用于存储未花费的输出
	unspentOuts := make(map[string]int)
	// 查找包含未花费输出的交易
	unspentTxs := bc.FindUnspentTransactions(address)
	// 初始化累计变量，用于跟踪累计的金额
	accumulated := 0

Work:
	// 遍历未花费交易
	for _, tx := range unspentTxs {
		// 将交易ID转换为字符串形式
		txID := hex.EncodeToString(tx.ID)
		// 遍历交易的输出
		for outIdx, out := range tx.Outputs {
			// 检查输出是否属于指定地址且累计的金额小于支付金额
			if out.ToAddressRight(address) && accumulated < amount {
				// 累加输出的价值到累计变量中
				accumulated += out.Value
				// 将交易ID和输出索引添加到未花费输出的 map 中
				unspentOuts[txID] = outIdx
				// 如果累计的金额达到或超过支付金额，则跳出循环
				if accumulated >= amount {
					break Work
				}
				// 继续到标签为 Work 的循环，因为一笔交易只能有一个输出与地址相关联
				continue Work
			}
		}
	}
	// 返回累计的价值和未花费输出的 map
	return accumulated, unspentOuts
}

// 该函数接收发送者地址 from、接收者地址 to 和交易金额 amount 作为参数，然后根据这些参数创建一个新的交易
func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput
	//定义了两个空的切片，inputs 用于存储输入，outputs 用于存储输出
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	//使用 FindSpendableOutputs 函数查找发送者的地址的可花费输出
	if acc < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}
	//如果找到的 UTXO 数量不够amount，则打印“不 enough coins！”并返回错误
	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{txID, outidx, from}
		inputs = append(inputs, input)
	}
	//将找到的 UTXO 转换为输入，并将其添加到 inputs 切片中
	outputs = append(outputs, transaction.TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{acc - amount, from})
	}
	//创建两个输出，一个是给接收者的amount，另一个是转出者的funds，并将它们添加到 outputs 切片中
	tx := transaction.Transaction{nil, inputs, outputs}
	tx.SetID()
	//创建了一个新的交易对象，并使用 SetID 函数设置交易 ID
	return &tx, true
}

func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}
