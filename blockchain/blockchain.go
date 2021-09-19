package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/tobiasvonarx/AvAin/globals"
)

// a struct representing a blockchain
type BlockChain struct {
	//Blocks   []*Block   // the blockchain consists of an array of (pointers to) blocks
	LastHash []byte     // a blockchain is identified by the hash of its last block
	Database *badger.DB // the blockchain (blocks+transaction data+metadata) is stored in a file using BadgerDB's KeyValue DB
}

// used to iterate through a blockchain which is stored in a DB
type BlockChainIterator struct {
	CurrentHash []byte     // the hash of the current block that is being iterated over
	Database    *badger.DB // a reference to the db of the blockchain
}

func (bc *BlockChain) Close() error {
	return bc.Database.Close()
}

// adds a block with given data to the blockchain
func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// read only transaction to get the current lastHash out of the DB
	err := bc.Database.View(func(txn *badger.Txn) error {
		// get the lashHash from the DB
		item, err := txn.Get([]byte("lh"))

		// need to handle on err immediately here
		if err != nil {
			panic(err)
		}

		// convert the last hash into a byte slice
		lastHash, err = item.Value()

		return err
	})

	if err != nil {
		panic(err)
	}

	// create a block with the data and the current last hash, aka the previous hash of the new block
	block := CreateBlock(data, lastHash)

	// put the new block into the db and assign the new block's hash as the last hash
	err = bc.Database.Update(func(txn *badger.Txn) error {
		// store the new block in the db
		err := txn.Set(block.Hash, block.Serialize())

		if err != nil {
			panic(err)
		}

		// store the new block's hash as the last hash in the db
		err = txn.Set([]byte("lh"), block.Hash)

		// update the LastHash field of the blockchain struct to the hash of the new block
		bc.LastHash = block.Hash

		return err
	})

	if err != nil {
		panic(err)
	}
}

// creates a blockchain and returns a reference to it
func CreateBlockChain() *BlockChain {
	var lastHash []byte

	cfg := badger.DefaultOptions

	// path for storing keys & metadata
	cfg.Dir = globals.DBPath

	// path for storing values
	cfg.ValueDir = globals.DBPath

	// open the db and panic on err
	db, err := badger.Open(cfg)
	if err != nil {
		panic(err)
	}

	// r&w transactions (txn) to the db by using update()
	err = db.Update(func(txn *badger.Txn) error {
		// do we already have a blockchain in our db? check for the key lh (lasthash)
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found in the DB")

			// create a sentinel block
			sentinel := CreateSentinel()
			fmt.Println("Sentinel created")

			// store the one and only block in the database, with the hash as key and its byte representation as value
			err = txn.Set(sentinel.Hash, sentinel.Serialize())

			// as the sentinel is the first and last block, its hash is the last hash (lh), store that in the db
			err = txn.Set([]byte("lh"), sentinel.Hash)

			// save the last hash as a variable as well, to initialize the blockchain struct
			lastHash = sentinel.Hash

			return err
		} else {
			// we already have a DB with a blockchain
			// "import" that existing DB

			// get the last hash of the blockchain from the DB
			item, err := txn.Get([]byte("lh"))

			// convert the last hash into a byte slice
			lastHash, err = item.Value()

			return err
		}
	})

	if err != nil {
		panic(err)
	}

	// create a blockchain instance and return a reference to it
	return &BlockChain{
		LastHash: lastHash,
		Database: db,
	}
}

// create an iterator for a blockchain
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		CurrentHash: bc.LastHash, // we start iterating from the end and go backwards
		Database:    bc.Database,
	}
}

// we start iterating from the last block (with the last hash), and go "backwards in time"
// Next() performs one iteration step, and goes to the immediately preceding block of the current block
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// read only transaction to get the block
	err := iter.Database.View(func(txn *badger.Txn) error {
		// get the byte slice of the block
		item, err := txn.Get(iter.CurrentHash)
		serializedBlock, err := item.Value()

		// desererialize to retrieve the block in block form
		block = Deserialize(serializedBlock)

		return err
	})

	if err != nil {
		panic(err)
	}

	// next time we call Next(), we fetch the previous block
	iter.CurrentHash = block.PrevHash

	return block
}
