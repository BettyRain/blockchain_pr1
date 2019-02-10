package main

import (
	"fmt"
)

func main() {
	MerklePatriciaTrie := &MerklePatriciaTrie{}
	fmt.Println("After insert: ")
//	MerklePatriciaTrie.Insert("ab","ancdgg")
//	fmt.Println("After insert: ")
//	MerklePatriciaTrie.Insert("aff","123")
//	fmt.Println("After insert: ")
	//"ab" some same one left
	// "f" "a12" "1ab" no same one
	
//	MerklePatriciaTrie.Insert("1ab","raggrhrehtrh")
//	MerklePatriciaTrie.Insert("pab","kitten")
//	fmt.Println("After insert: ")
	MerklePatriciaTrie.Insert("p","apple")
	MerklePatriciaTrie.Insert("aaaaa","banana")
	MerklePatriciaTrie.Insert("aaaap","orange")
	fmt.Println("After insert: ")
	MerklePatriciaTrie.Insert("aa","new")
	//2 ext 3 leaf 2 branch (1 with value)
	
	
	fmt.Println("Get value: ")
	fmt.Println(MerklePatriciaTrie.Get("p"))


	
}

