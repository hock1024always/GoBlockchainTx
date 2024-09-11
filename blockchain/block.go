package blockchain

import (
	"GOBLOCKCHAIN/transaction"
	"GOBLOCKCHAIN/utils"
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Target       []byte
	Nonce        int64
	Transactions []*transaction.Transaction //用于存储一个或多个交易的信息
}

// 创建节点
func CreateBlock(prevhash []byte, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, []byte{}, 0, txs}
	//这个区块存储的一个要素解释一条交易信息
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// 创建创世节点
func GenesisBlock() *Block {
	tx := transaction.BaseTx([]byte("Always"))
	return CreateBlock([]byte{}, []*transaction.Transaction{tx})
}

func (b *Block) BackTrasactionSummary() []byte {
	txIDs := make([][]byte, 0)
	//二维切片 每一行存储一条交易信息
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	} //将每个交易的 ID 添加到 txIDs 这个二维切片中的新行中
	summary := bytes.Join(txIDs, []byte{})
	//相当于将这个二维切片连成一个一维切片，每行之间用[]byte{}隔开
	return summary
}

// 合信息 算哈希
func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.BackTrasactionSummary()}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}
