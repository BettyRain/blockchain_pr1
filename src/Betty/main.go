package main

import (
	"fmt"

)

func main() {
	MerklePatriciaTrie := &MerklePatriciaTrie{}
	MerklePatriciaTrie.Insert("ab","ancd")
	MerklePatriciaTrie.Insert("abc","ancghjd")
//	fmt.Println(keyToHex("verb"))
//	fmt.Println(compact_encode([]uint8{6,1,6,2}))
//	test_compact_encode()
	fmt.Println(compact_encode(keyToHex("ab")))
	
}

