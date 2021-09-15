package blockchain

import (
	"fmt"
	"strconv"
)

type BlockChain struct {
	Blocks []*Block // the blockchain consists of an array of (pointers to) blocks
}

// adds a block with given data to the blockchain
func (bc *BlockChain) AddBlock(data string) {
	// the current last block will be the new block's previous block
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	// instantiate the new block
	newBlock := CreateBlock(data, prevBlock.Hash)

	// append the new block to our blockchain
	bc.Blocks = append(bc.Blocks, newBlock)
}

// creates a blockchain and returns a reference to it
func CreateBlockChain() *BlockChain {
	return &BlockChain{
		Blocks: []*Block{CreateSentinel()}, // instantiate the blocks array to create and contain a Sentinel block
	}
}

// display the blockchain as a string
func (bc *BlockChain) Print() {
	for i, block := range bc.Blocks {
		fmt.Printf("\nBlock %d\n", i)
		fmt.Printf("Data: 					%s\n", block.Data)
		fmt.Printf("Hash: 					%x\n", block.Hash)
		fmt.Printf("Previous Hash:			 	%x\n", block.PrevHash)

		pow := NewProof(block)
		fmt.Printf("Proof Of Work:				%s\n", strconv.FormatBool(pow.Validate()))
	}
}
