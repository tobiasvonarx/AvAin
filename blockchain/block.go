package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Data     []byte // data associated with the block
	Hash     []byte // hash of the block
	PrevHash []byte // the previous block's hash, in order to link the blocks together
}

// derive the hash of a block
func (b *Block) DeriveHash() {
	// to derive the hash of a block we need its data and the hash of the previous block
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})

	// hash the info with the sum256 hashing function
	hash := sha256.Sum256(info)

	// push the hash into the hash field of the block
	b.Hash = hash[:]
}

// creates a block from given data and the hash of the previous block
func CreateBlock(data string, prevHash []byte) *Block {
	// instantiate a block and store its reference
	block := &Block{
		Data:     []byte(data), // convert the data string into a slice of bytes
		Hash:     []byte{},     // an empty slice of bytes
		PrevHash: prevHash,     // the hash of the previous block
	}

	// calculate the block's hash
	block.DeriveHash()

	// return a pointer to the block
	return block
}

// Sentinel Block as the very first block, since otherwise we would have no "previous hash"
// thank you to Prof. Hoefler for the inspiration
func CreateSentinel() *Block {
	// create a block with an empty byte slice as its prevHash and return a reference to it
	return CreateBlock("Sentinel", []byte{})
}
