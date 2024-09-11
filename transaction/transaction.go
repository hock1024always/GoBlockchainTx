package transaction

import (
	"GOBLOCKCHAIN/constcoe"
	"GOBLOCKCHAIN/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// 下列结构用于记录一笔交易
type Transaction struct {
	ID      []byte     //自身的ID值(哈希值) 由一组TxInput与一组TxOutput构成
	Inputs  []TxInput  //TxInput用于标记支持我们本次转账的前置的交易信息的TxOutput
	Outputs []TxOutput //而TxOutput记录我们本次转账的amount和Reciever
}

//下面两个函数 实现计算每个transaction的哈希值的功能

// TxHash返回交易信息的哈希值
func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	//gob，其功能主要是序列化结构体，与json有些像但是更方便
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

// SetID设置每个交易信息的ID值，也就是哈希值
func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash() //调用 TxHash() 方法计算的哈希值设置为交易 (tx) 的 ID 字段
}

// 创建一个基本的交易（transaction），并返回该交易的指针
func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}}
	//txIn：一个空的交易输入（TxInput），表示没有引用任何之前的交易输出（UTXO）作为输入
	txOut := TxOutput{constcoe.InitCoin, toaddress}
	//txOut：一个交易输出（TxOutput），包含一个常量 InitCoin 作为金额，以及传入函数的 toaddress 作为收款地址。
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	//tx：一个交易（Transaction）结构，包含了一些元数据（如交易的描述），以及上述创建的交易输入和输出。
	return &tx
}

// 用是判断一个交易是否是基础交易
func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
	//检查交易的输入是否只有一个，并且该输入的输出索引（OutIdx）是否为 -1
	//如果是，那么这个交易被认为是基础交易
}
