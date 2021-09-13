package main

import "github.com/tobiasvonarx/AvAin/blockchain"

// entrypoint
func main() {
	chain := blockchain.CreateBlockChain()
	chain.AddBlock("Lorem")
	chain.AddBlock("Ipsum")
	chain.AddBlock("dolor")
	chain.AddBlock("sit")
	chain.AddBlock("amet")
	chain.Print()
}
