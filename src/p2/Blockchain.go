package p2

import (
	"encoding/json"
	"fmt"
)


type Blockchain struct {
	chain map[int32][]Block //This is a map which maps a block height to a list of blocks
	//MPT value is an instance of MerklePatriciaTrie. You would create a mpt like P1 does, then pass it to Initial(). 
	length int32 //Length equals to the highest block height
}

func (blch *Blockchain) Get (height int32) []Block {
	//This function takes a height as the argument,
	//returns the list of blocks stored in that height or None if the height doesn't exist.
	getList, ok := blch.chain[height]	
	if ok == false {
		return nil
	} else {
		return getList
	}
}

func (blch *Blockchain) Insert (blc Block) {
	//This function takes a block as the argument, insert that block to the BlockChain.Chain map. 
	heightNewBlock := blc.header.height
	
	if (len(blch.chain) == 0) {
		//chain is null 
		blockList := []Block {blc}
		blch.chain = make (map[int32][]Block)
		blch.chain[1] = blockList
		blch.length = 1
	} else {
		blockList := blch.Get(heightNewBlock)
		if (len(blockList) == 0){
			//no such blocks with that height
			blockList = []Block {blc}
			if (heightNewBlock > blch.length) {
				blch.length = heightNewBlock
			}
		} else {
			//if block already exists in blockchain
			isEqual := false
			for _, block:= range blockList {
				if (block.header.hash == blc.header.hash){
					isEqual = true
					fmt.Println("equal")
				}
			}
			if (!isEqual) {
				//add to the end of the list
				blockList = append(blockList, blc)
			}
		}
		blch.chain[heightNewBlock] = blockList
		if(blch.length < blc.header.size) {
			blch.length = blc.header.size
		}	
	}
}

func (blch *Blockchain) EncodeToJSON () string {
	//This function iterates over all the blocks,
	//generate blocks' JsonString by the function you implemented previously,
	//and return the list of those JsonStritgns. 
	jsonString := "[\n"
	for _, list_block := range blch.chain {
		for _, block := range list_block {
			blockJson := block.EncodeToJSON()
			jsonString += blockJson
			jsonString += ","	
		}
	}
	leng := len(jsonString)
	jsonString = jsonString[:leng-1]
	jsonString += "]"
	return jsonString
}

func DecodeFromJSON (jsonString string) Blockchain {
	//Description: This function is called upon a blockchain instance.
	//It takes a blockchain JSON string as input,
	//decodes the JSON string back to a list of block JSON strings,
	//decodes each block JSON string back to a block instance,
	//and inserts every block into the blockchain.  
	blockch := Blockchain{}
	decodedBlock := Block{}
	byt := []byte(jsonString)
	blocks := make([]jsonBlockSt,0)
	err := json.Unmarshal(byt, &blocks)
	for _, block := range blocks {
		decodedBlock = Block{}
		decodedBlock.DecodeFromStruct(block)
		blockch.Insert(decodedBlock)
	}
	if (err != nil) {
		return Blockchain{}
	}
	return blockch
}
