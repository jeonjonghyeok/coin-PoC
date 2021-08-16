package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/jeonjonghyeok/coin/db"
	"github.com/jeonjonghyeok/coin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlock(b *Block) {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func createBlock(newestHash string, height, diff int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   newestHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}
	block.mine()
	block.Transactions = Mempool.TxToConfirm()
	persistBlock(block)
	return block
}
