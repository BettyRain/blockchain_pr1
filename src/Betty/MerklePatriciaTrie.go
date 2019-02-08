package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
//	"bytes"
)

type Flag_value struct {
	encoded_prefix []uint8
	//here can be hashed_value of node (if branch or ext)
	value          string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db map[string]Node
	//string is a hashed_value of Node
	root string
}

func (mpt *MerklePatriciaTrie) Get(key string) string {
	encodedKey := keyToHex(key)
	encodedKey = append(encodedKey, 16)
	encodedKey = compact_encode(encodedKey)
//	fmt.Println(mpt.root)
	n := mpt.GetByNode(mpt.db[mpt.root], encodedKey)
	if (n.flag_value.value != "") {
		return n.flag_value.value
	} else {
	return "no such key"
	}
}


func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	
	//if mpt is empty
	if mpt.db == nil {
		encodedKey := keyToHex(key)
		encodedKey = append(encodedKey, 16)
		encodedKey = compact_encode(encodedKey)
		n := Node{}
		n.node_type = 2
		n.branch_value = [17]string{}
		n.flag_value.encoded_prefix = encodedKey
		n.flag_value.value = new_value
		root := n.hash_node()
		mpt.root = root
		mpt.db = map[string]Node { 
			root : n,
		}
	
		fmt.Println(mpt)
	} else {
		
	}
	
	 
	
}

func (mpt *MerklePatriciaTrie) GetByNode(node Node, key []byte) Node {
	n := Node{}
	
	if len(key) == 0 {
		return node
	}

	if (node.node_type == 0) {
		return n
	}

	switch node.node_type {
	case 2:
		key_node := node.flag_value.encoded_prefix
		hash_node := node.flag_value.value //hash_value of Node
		n = mpt.db[hash_node]
		
		//fmt.Println(key[2:])

		if len(key) > len(key_node) && strings.HasPrefix(string(key_node), string(key)) {
			return mpt.GetByNode(n, key[len(key_node):])
		} else if len(key) == len(key_node) && strings.HasPrefix(string(key_node), string(key)) {
			return node
		}
		return n
	case 1:
		hash_node := node.flag_value.value
		n = mpt.db[hash_node]
		return mpt.GetByNode(n, key[1:])
		//case 0
	default:
		return n
	}
	
	return n
}



func (mpt *MerklePatriciaTrie) Delete(key string) {
	// TODO
}

func compact_decode(hex_array []uint8) []uint8 {
	decoded_arr := []uint8{}
	for i := 0; i < len(hex_array); i += 1 {
		firstPart := hex_array[i] / 16
		secondPart := hex_array[i] % 16
		decoded_arr = append(decoded_arr, firstPart)
		decoded_arr = append(decoded_arr, secondPart)
	}

	if decoded_arr[0] == 0 || decoded_arr[0] == 2 {
		decoded_arr = append(decoded_arr[:0], decoded_arr[1:]...)
		decoded_arr = append(decoded_arr[:0], decoded_arr[1:]...)
	} else if decoded_arr[0] == 1 || decoded_arr[0] == 3 {
		decoded_arr = append(decoded_arr[:0], decoded_arr[1:]...)
	}
	return decoded_arr
}

// TODO: how we send prefix to function?

// If Leaf, ignore 16 at the end
func compact_encode(encoded_arr []uint8) []uint8 {
	//encoded_arr = [] {1, 6, 1}
	term := 0
	if encoded_arr[len(encoded_arr)-1] == 16 {
		term = 1
		encoded_arr = encoded_arr[:len(encoded_arr)-1]
	} else {
		term = 0
	}
	
	oddlen := len(encoded_arr) % 2
	flags := 2*term + oddlen

	if oddlen == 1 {
		encoded_arr = append([]uint8{uint8(flags)}, encoded_arr...)
	} else {
		encoded_arr = append([]uint8{0}, encoded_arr...)
		encoded_arr = append([]uint8{uint8(flags)}, encoded_arr...)
	}

	result := []uint8{}
	for i := 0; i < len(encoded_arr); i += 2 {
		result = append(result, (16*encoded_arr[i] + 1*encoded_arr[i+1]))
	}
	
	return result
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

const hextable = "0123456789abcdef"

//convert key to hex array
func keyToHex(key string) []uint8 {
	result := []uint8{}
	encodedByteString := hex.EncodeToString([]byte(key))
	for _, encodedByte := range encodedByteString {
		result = append(result, uint8(strings.IndexByte(hextable, uint8(encodedByte))))
	}
	return result
}

