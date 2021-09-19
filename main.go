package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/tobiasvonarx/AvAin/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

// prints the intended usage of the cli
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("	add -block BLOCK_DATA 'adds a block to the chain'")
	fmt.Println("	print 'prints the blocks in the chain'")
}

// validate the command line args and print the usage if neccessary
func (cli *CommandLine) validateArgs() {
	// check if the user has entered a invalid command
	if len(os.Args) < 2 {
		cli.printUsage()
		// exits without corrupting the DB
		runtime.Goexit()
	}
}

// adds a block through the command line
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("added block with data", data)
}

// prints the chain to the command line by iterating over it
func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	i := 0
	for {
		block := iter.Next()

		fmt.Printf("\nBlock %d\n", i)
		fmt.Printf("Data: 					%s\n", block.Data)
		fmt.Printf("Hash: 					%x\n", block.Hash)
		fmt.Printf("Previous Hash:			 	%x\n", block.PrevHash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Proof Of Work:				%s\n", strconv.FormatBool(pow.Validate()))

		if len(block.PrevHash) == 0 {
			// stop iterating
			break
		}

		i++
	}
}

// runs the commandhandling
func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockFlag := flag.NewFlagSet("add", flag.ExitOnError)
	printChainFlag := flag.NewFlagSet("print", flag.ExitOnError)

	// subset to the add flag
	addBlockData := addBlockFlag.String("block", "", "block data")

	switch os.Args[1] {
	case "add":
		if err := addBlockFlag.Parse(os.Args[2:]); err != nil {
			panic(err)
		}

	case "print":
		if err := printChainFlag.Parse(os.Args[2:]); err != nil {
			panic(err)
		}

	default:
		// catch nubs
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockFlag.Parsed() {
		// check if the arg for block data is empty
		if *addBlockData == "" {
			addBlockFlag.Usage()
			runtime.Goexit()
		}

		// we have block data given
		// make a new block with the data and add it to the blockchain
		cli.addBlock(*addBlockData)
	} else if printChainFlag.Parsed() {
		cli.printChain()
	}
}

// entrypoint
func main() {
	// failsafe for exiting
	defer os.Exit(0)

	bc := blockchain.CreateBlockChain()

	// properly close the DB before the main function ends
	defer bc.Close()

	cli := CommandLine{
		blockchain: bc,
	}

	cli.run()
}
