package blockchain

import (
	"bytes"
	"encoding/gob"
)

type Block struct {
	Data     []byte // data associated with the block
	Hash     []byte // hash of the block
	PrevHash []byte // the previous block's hash, in order to link the blocks together
	Nonce    int    // important for the validation of the proof of work algorithm
}

/* derive the hash of a block
func (b *Block) DeriveHash() {
	// to derive the hash of a block we need its data and the hash of the previous block
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})

	// hash the info with the sum256 hashing function
	hash := sha256.Sum256(info)

	// push the hash into the hash field of the block
	b.Hash = hash[:]
}*/

// creates a block from given data and the hash of the previous block
func CreateBlock(data string, prevHash []byte) *Block {
	// instantiate a block and store its reference
	block := &Block{
		Data:     []byte(data), // convert the data string into a slice of bytes
		Hash:     []byte{},     // an empty slice of bytes
		PrevHash: prevHash,     // the hash of the previous block
		Nonce:    0,            // initial nonce
	}

	// runs the proof of work algorithm on each created block, to sign it
	pow := NewProof(block)
	nonce, hash := pow.Run()

	// set attributes
	block.Hash = hash[:]
	block.Nonce = nonce

	// return a pointer to the block
	return block
}

// Sentinel Block as the very first block, since otherwise we would have no "previous hash"
// thank you to Prof. Hoefler for the inspiration
func CreateSentinel() *Block {
	// create a block with an empty byte slice as its prevHash and return a reference to it
	return CreateBlock("Sentinel", []byte{})
}

// BadgerDB only accepts arrays of bytes, so we need to serialize the blocks into bytes
func (deserialized *Block) Serialize() []byte {
	// create a buf to store the serialized block
	var serialized bytes.Buffer

	// create an encoder for the buffer
	encoder := gob.NewEncoder(&serialized)

	// encode our block and panic on error
	if err := encoder.Encode(deserialized); err != nil {
		panic(err)
	}

	// returns a representation of the block in bytes
	return serialized.Bytes()
}

// deserialize an slice of bytes to recover its block representation
func Deserialize(serialized []byte) *Block {
	// declare the block that we want to our deserialized data to represent
	var deserialized Block

	// create the decoder using the bytes we want to deserialize
	decoder := gob.NewDecoder(bytes.NewReader(serialized))

	// decode into the block and panic on error
	if err := decoder.Decode(&deserialized); err != nil {
		panic(err)
	}

	// return a reference to the deserielized block in memory
	return &deserialized
}
