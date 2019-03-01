package main

import (
	"p1"
	"p2"
	"fmt"
)

func main() {
	mpt := p1.MerklePatriciaTrie{}
	
	mpt.Insert("aaa", "apple")
	mpt.Insert("aap", "banana")
	mpt.Insert("bc", "new")

	mpt2 := p1.MerklePatriciaTrie{}
	
	mpt2.Insert("af", "ni")

	fmt.Println("Start")
	blc := p2.Initial (1, "genesis", mpt)
	blc2 := p2.Initial (2, "parent", mpt2)
//	blc3 := p1.Block {}
	blch := p2.Blockchain{}
	blch.Insert(blc)
	blch.Insert(blc2)

//	str1 := blc.EncodeToJSON()
//	fmt.Println(blc)
//	fmt.Println(mpt.Order_nodes())
//	blc3.DecodeFromJson(str1)
//	fmt.Println(blc3)

	
//	blch := p1.Blockchain{}
//	blch.Insert(blc)
//	blch.Insert(blc2)
	//fmt.Println(blch.Get(1))
	str:=blch.EncodeToJSON()
	
	blch2 := p2.Blockchain{}
	blch2 = p2.DecodeFromJSON(str)
	fmt.Println("===============")
	fmt.Println(blch)
	fmt.Println(blch2)
	fmt.Println("================")

}