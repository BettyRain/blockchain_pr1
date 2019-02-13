package main

import (
//	"package1"
	"fmt"
)

func main() {
	mpt := &MerklePatriciaTrie{}
	
	mpt.Insert("aaa", "apple")
	mpt.Insert("aap", "banana")
	inserted_trie1 := mpt.Order_nodes()
	mpt.Insert("bc", "new")


	inserted_trie2 := mpt.Order_nodes()
	mpt.Delete("bc")
	inserted_trie4 := mpt.Order_nodes()
    fmt.Println(mpt.Get("aap"))
	
	fmt.Println(inserted_trie1)
	fmt.Println(inserted_trie2)
//	fmt.Println(inserted_trie3)
	fmt.Println(inserted_trie4)


//	fmt.Println(mpt)

}