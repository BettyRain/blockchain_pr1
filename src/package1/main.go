package main

import (
//	"package1"
	"fmt"
)

func main() {
	mpt := &MerklePatriciaTrie{}
	
	mpt.Insert("a", "apple")
	mpt.Insert("b", "banana")
	mpt.Insert("ab", "new")
	inserted_trie1 := mpt.Order_nodes()

	mpt.Delete("c")
	inserted_trie2 := mpt.Order_nodes()

	mpt.Delete("ab")
	inserted_trie3 := mpt.Order_nodes()

	//inserted_trie1 := mpt.Order_nodes()
	
	
	fmt.Println(inserted_trie1)
	fmt.Println(inserted_trie2)
	fmt.Println(inserted_trie3)
	
	fmt.Println(mpt)

}