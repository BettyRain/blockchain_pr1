package p2

import (
	"p1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strconv"
	"time"
)

type Header struct {
	height     int32
	timestamp  int64 //The value must be in the UNIX timestamp format such as 1550013938
	hash       string
	parentHash string
	size       int32 //You have a mpt, you convert it to byte array, then size = len(byteArray)
}

type Block struct {
	//Block{Header{Height, Timestamp, Hash, ParentHash, Size}, Value}
	header Header
	value  p1.MerklePatriciaTrie
}

type jsonBlockSt struct{
	Hash string
	TimeStamp int64
	Height int32
	ParentHash string
	Size int32
	Mpt map[string]string
}


func Initial(new_height int32, new_parentHash string, new_value p1.MerklePatriciaTrie) Block {
	//This function takes arguments(such as height, parentHash, and value of MPT type) and forms a block.
	time_str := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	new_timestamp, err := strconv.ParseInt(time_str, 10, 64)
	new_size := int32(len([]byte(new_value.String())))
	new_hash := strconv.FormatInt(int64(new_height), 10) + time_str + new_parentHash + new_value.GetRoot() + strconv.FormatInt(int64(new_size), 10)
	sum := sha3.Sum256([]byte(new_hash))
	new_header := Header{height: new_height, timestamp: new_timestamp, hash: hex.EncodeToString(sum[:]), parentHash: new_parentHash, size: new_size}
	new_block := Block{value: new_value, header: new_header}
	if err != nil {
		return Block{}
	}
	return new_block
}

func (decodedBlock *Block) DecodeFromJson(jsonString string) {
	//This function takes a string that represents the JSON value of a block as an input,
	//and decodes the input string back to a block instance.
	byt := []byte(jsonString)
	blocks := jsonBlockSt{}
	err := json.Unmarshal(byt, &blocks)
	fmt.Println(err)	
	decodedBlock.DecodeFromStruct(blocks)
}

func (blc *Block) EncodeToJSON() string {
	//This function encodes a block instance into a JSON format string.
	//Note that the block's value is an MPT, and you have to record all of the (key, value) pairs
	//that have been inserted into the MPT in your JSON string.
	//timestampStr := strconv.FormatInt(blc.header.timestamp, 10)
	//heightStr := strconv.FormatInt(int64(blc.header.height), 10)
	//sizeStr := strconv.FormatInt(int64(blc.header.size), 10)	
	group := &jsonBlockSt{
		Hash: blc.header.hash,
		//TimeStamp: timestampStr,
		TimeStamp: blc.header.timestamp,
		//Height: heightStr,
		Height: blc.header.height,
		ParentHash: blc.header.parentHash,
		Size: blc.header.size,
		//Size: sizeStr,
		Mpt: blc.value.GetKeyValue(),
	}
	fmt.Println(group)
	b, err := json.Marshal(group)
	
	if err != nil {
		return ""
	}

	return string(b)
}

func (decodedBlock *Block) DecodeFromStruct (blocks jsonBlockSt) {
	//This function creates a Block from json Block structure
	//Create new mpt	
	mptDecoded := p1.MerklePatriciaTrie{}
	for key, value := range blocks.Mpt {
		mptDecoded.Insert(key, value)
	}
	decodedBlock.header.hash = blocks.Hash
	decodedBlock.header.parentHash = blocks.ParentHash
	decodedBlock.header.height = blocks.Height
	decodedBlock.header.size = blocks.Size
	decodedBlock.header.timestamp = blocks.TimeStamp
	decodedBlock.value = mptDecoded

}
