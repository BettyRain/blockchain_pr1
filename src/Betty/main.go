package main

import (
	"fmt"

)

func main() {
	MerklePatriciaTrie := &MerklePatriciaTrie{}
	fmt.Print("After insert: ")
	MerklePatriciaTrie.Insert("ab","ancdgg")
//	MerklePatriciaTrie.Insert("abc","ancghjd")
	fmt.Print("Get value: ")
	fmt.Println(MerklePatriciaTrie.Get("ab"))
//	fmt.Println(keyToHex("verb"))
//	fmt.Println(compact_encode([]uint8{6,1,6,2}))
//	test_compact_encode()
//	fmt.Println(compact_encode(keyToHex("ab")))
	
}

