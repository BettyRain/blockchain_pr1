package package1

import (
	
)

type Block struct {
	Height int32
	Timestamp int64 //The value must be in the UNIX timestamp format such as 1550013938
	Hash string
	ParentHash string
	Size int32 //You have a mpt, you convert it to byte array, then size = len(byteArray)
}

func (mpt *MerklePatriciaTrie) Initial (height int32, parentHash string, value string) blc *Block {
	//Description: This function takes arguments(such as height, parentHash, and value of MPT type) and forms a block.
	//VALUE = MPT value is an instance of MerklePatriciaTrie. You would create a mpt like P1 does, then pass it to Initial().
	//TODO: what is height here?
	//TODO: how I can find parent?
}


func DecodeFromJson(jsonString string) blc *Block {
	//Description: This function takes a string that represents the JSON value of a block as an input
	//and decodes the input string back to a block instance.
	//Note that you have to reconstruct an MPT from the JSON string, and use that MPT as the block's value. 
	//Argument: a string of JSON format
	//Return type: block instance
	//TODO: how can I put block into input function? Is it right?
	return blc
}

func (blc *Block) EncodeToJSON string {
	//Description: This function encodes a block instance into a JSON format string.
	//Note that the block's value is an MPT, and you have to record all of the (key, value) pairs
	//that have been inserted into the MPT in your JSON string.
	//Argument: a block, or you may define this as a method of the block struct
	//Return type: a string of JSON format
	

	return ""
}




