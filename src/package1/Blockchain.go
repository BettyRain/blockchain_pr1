package package1

import (

)


type Blockchain struct {
	Chain map[int32][]Block //This is a map which maps a block height to a list of blocks
	//MPT value is an instance of MerklePatriciaTrie. You would create a mpt like P1 does, then pass it to Initial(). 
	Length int32 //Length equals to the highest block height
}

func (blch *Blockchain) Get (height int32) []Block {
	//This function takes a height as the argument,
	//returns the list of blocks stored in that height or None if the height doesn't exist.
	getList := []Block 
	//find value by key in map
	for key, value := range blch.Chain {
		if reflect.DeepEqual(height, key) {
			getList = value
		}
	}
	if getList == null {
		//what does return none mean? Like empty list?
	}
}

func (blch *Blockchain) Insert (blc *Block) string {
	//This function takes a block as the argument, insert that block to the BlockChain.Chain map. 
	//Argument: block
	getList := blch.Get
	//TODO: hpw we can understand whic block is the last? By length?
	//TODO: What if there some of them? We start with the first from list?
	
	parentHash = getList[1].Hash //is is parent hash?
	return ""
}

func (blch *Blockchain) EncodeToJSON string {
	//Description: This function iterates over all the blocks,
	//generate blocks' JsonString by the function you implemented previously,
	//and return the list of those JsonStritgns. 
	//Return type: string	
}

func (blch *Blockchain) DecodeFromJSON (base64 string) {
	//Description: This function is called upon a blockchain instance.
	//It takes a blockchain JSON string as input,
	//decodes the JSON string back to a list of block JSON strings,
	//decodes each block JSON string back to a block instance,
	//and inserts every block into the blockchain.  
	//TODO: Copy means overwrite?
	//Argument: self, string
}
