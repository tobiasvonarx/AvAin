package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/tobiasvonarx/AvAin/globals"
	"github.com/tobiasvonarx/AvAin/util"
)

// TODO perform useful work by processing MCMC chains http://web.mit.edu/alex_c/www/nooshare.pdf
// Proof of Work Algorithm
// secure the blockchain by forcing the network to do work in order to add a block to the chain
// if you do the proof of work algorithm you can sign the respective block and power the network that way
// this makes the blocks and associated data more secure
// when the work to sign a block is done, they need to provide proof of this work
// => need something that is hard to do but easy to prove

// what makes it secure?
// if you wanted to manipulate a block in the blockchain
// you'd need recalculate the hash (expensive) and recalculate every consecutive block in the blockchain to validate the data

// inspired by http://www.cypherspace.org/hashcash/, the proof-of-work spec that bitcoin used originally
// the first few bits must contain 0s. the more bits must be 0, the harder
// globals.Difficulty == x => x/4 preceding hex-digits are 0 / x preceding 0-bits

type ProofOfWork struct {
	Block  *Block   // the block we want to sign
	Target *big.Int // represents the requirement of containing #Difficulty-many preceding 0s
}

// instantiates and returns a reference to a proof for the given block
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	// upper bound for the number we can get s.t. the preceding #Difficulty bits are 0s
	// target = 1 << (256-Difficulty) = 2^(256-Difficulty), since we have 256 bits in our hash
	target.Lsh(target, uint(256-globals.Difficulty))

	return &ProofOfWork{
		Block:  b,
		Target: target,
	}
}

// returns a slice of bytes, containing the relevant Data for the proof of work
func (pow *ProofOfWork) InitiateData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,           // take the data from the block
			util.ToHex(int64(nonce)), // nonce, a number only used once (unique to the block)
			util.ToHex(int64(globals.Difficulty)),
		},
		[]byte{},
	)
}

// creates the hash from the nonce and the block data and checks if the hash meets the hashcash-like requirement
// it does so until it finds a nonce s.t. the hash fits the requirement, and the block is signed (or the loop runs out)
func (pow *ProofOfWork) Run() (int, []byte) {
	var attempt big.Int
	var hash [32]byte //32*8=256 bit hash

	nonce := 0

	// a lot of iterations
	for nonce < math.MaxInt64 {
		// prepare data
		data := pow.InitiateData(nonce)

		// hash the data and nonce with sha256
		hash = sha256.Sum256(data)

		if globals.Debug {
			fmt.Printf("\r%x\n", hash)
		}

		// convert the hash into a bigint
		attempt.SetBytes(hash[:])

		// compare that integer with out target upper bound bigint (Difficulty)
		// this is to see whether the requirement of the #Difficulty many preceding bits of the hash being 0s is met
		if attempt.Cmp(pow.Target) == -1 {
			// attempt < target => preceding #Difficulty bits are 0s and requirement fulfilled
			// we are successful and have signed the block
			break
		} else {
			// attempt >= target => at least one of the preceding #Difficulty bits is a 1, the attempt fails
			nonce++
		}
	}

	// returns the nonce with which the Block was successfully signed, together with the requirement-fulfilling hash
	return nonce, hash[:]
}

// this method validates the hash derived from the Nonce attribute of the block which the pow operates on
func (pow *ProofOfWork) Validate() bool {
	// we run the Run() loop one more time
	var attempt big.Int

	data := pow.InitiateData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	attempt.SetBytes(hash[:])

	return attempt.Cmp(pow.Target) == -1
}
