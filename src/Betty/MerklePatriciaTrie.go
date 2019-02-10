package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
)

type Flag_value struct {
	encoded_prefix []uint8
	value          string	//here can be hashed_value of node (if branch or ext)
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
	key_hex := keyToHex(key)
	root_node := mpt.db[mpt.root]
	n := Node{}
	switch root_node.node_type {
		case 1:
			n = mpt.GetByNode(root_node,[]byte (key_hex))
		case 2:
		root_key := root_node.flag_value.encoded_prefix
		root_key_decoded := compact_decode(root_key)
		if (intersectionCount(string (root_key_decoded), string (key_hex)) > 0) {
			n = mpt.GetByNode(root_node,[]byte (key_hex))
		}  else {
			fmt.Println("case 342 out")
		return "no such key"
		}
		
	}
		fmt.Println("out")
		fmt.Println(n)
		switch n.node_type {
			case 1:
			fmt.Println("case 1 out")
			if (n.branch_value[16] != "") {
				return n.branch_value[16]
			} else {
				return "no such key"
			}

			case 2:
			fmt.Println("case 2 out")
			if (n.flag_value.value != "") {
				return n.flag_value.value
			} else {
				return "no such key"
			}
		}
		
	return "no such key"
}


func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	n := Node{}
	encoded_key := keyToHex(key)
	fmt.Println(" ----INSERTION----- ")
	fmt.Println(encoded_key)
	fmt.Println(" ----INSERTION----- ")
	
	//if mpt is empt
	if mpt.db == nil {
		encoded_key = append(encoded_key, 16)
		encoded_key = compact_encode(encoded_key)
		n.node_type = 2
		n.branch_value = [17]string{}
		n.flag_value.encoded_prefix = encoded_key
		n.flag_value.value = new_value
		root := n.hash_node()
		mpt.root = root
		mpt.db = map[string]Node { 
			root : n,
		}
	} else {
		encoded_key = append(encoded_key)
		encoded_key = compact_encode(encoded_key)
		mpt.InsertByNode(mpt.db[mpt.root],[]byte (encoded_key),  new_value, true)
	}
	
	fmt.Println(mpt)	

}


func (mpt *MerklePatriciaTrie) InsertByNode(node Node, key []byte, new_value string, isRoot bool) string {
	//check 0 values

	
	key_node := node.flag_value.encoded_prefix
	hash_node := node.flag_value.value //hash_value of Node
	hashed_old_leaf := mpt.KeyByValue(node)

	switch node.node_type {
		
		//branch node
		case 1:
		decoded_key := compact_decode(key)
		branch_value := decoded_key[0]
		fmt.Println(" ----DECODED-BR----- ")
		fmt.Println(decoded_key)
		fmt.Println(" ----DECODED-BR----- ")
		
		rest_path := []uint8{}
		rest_path = append(rest_path, decoded_key[1:]...)
		rest_path = append(rest_path, 16)
		rest_path = compact_encode(rest_path)
		
		if (node.branch_value[branch_value] == "") {
			//cell is empty
			
			//create new leaf
			f_leaf := Flag_value {encoded_prefix: rest_path, value: new_value}
			n_leaf := Node {node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
			n_leaf_hash := n_leaf.hash_node()
			
			// add leaf to cell
			node.branch_value[branch_value] = n_leaf.hash_node()
			mpt.db[hashed_old_leaf] = node
			mpt.db[n_leaf_hash] = n_leaf
			
		} else {
			//cell is not empty
			//move to next node
			fmt.Println("move to next")
			hash_next_node := node.branch_value[branch_value]
			next_node := mpt.db[hash_next_node]
			fmt.Println(next_node)
			hash_next_node = mpt.InsertByNode(next_node, rest_path, new_value, false)
			node.branch_value[branch_value] = hash_next_node
			mpt.db[hashed_old_leaf] = node
			
			//TODO: add it		
		}
		
		//ext or leaf
		case 2:
		node_prefix := key_node[0]/16

		//check for leaf and ext
		switch node_prefix {
			//extension node
			case 0,1:
			decoded_key_node := compact_decode(key_node) //nibbles of ext decoded
			decoded_key := compact_decode(key)
			fmt.Println(" ----DECODED-EXT----- ")
			fmt.Println(decoded_key_node)
			fmt.Println(decoded_key)
			fmt.Println(" ----DECODED-EXT----- ")
			
			//new leaf with empty prefix
			f_leaf := Flag_value {encoded_prefix: []uint8 {}, value: new_value}
			n_leaf := Node {node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
			n_leaf_hash := n_leaf.hash_node()			
			
			index_ext := intersectionCount (string(decoded_key_node), string(decoded_key))
			if (index_ext == len(decoded_key_node)){
				//all same nibbles
				fmt.Println("all same")
				branch_node := mpt.db[hash_node]
				//TODO: here we go to insert branch node
				mpt.InsertByNode(branch_node, key, new_value, false)
				
				
			} else if (index_ext <= (len(decoded_key_node)-1)) && (index_ext > 0) {
				fmt.Println("some same")
				//some same nibbles - more than one left (<)
				//take same for new ext + take one for branch + leave other for old ext
				//some same nibbles - one nibble left (=)
				//take same for new ext + take one for branch + delete old ext + add old branch to index[] of new branch
				
				//have to change the root
				common_path := decoded_key_node[:index_ext]
				rest := decoded_key[index_ext:]
				branch_nib := decoded_key_node[index_ext:]
									
					fmt.Println(rest)
					fmt.Println(branch_nib)
				//create branch
				branch := [17] string {}

				if (len(rest)) > 0 && (len(branch_nib)) > 1 {
					branch_path := rest[0]
					branch_nibble := branch_nib[0]
					rest_path := decoded_key[(index_ext+1):]
					rest_nibble := decoded_key_node[(index_ext+1):]
					
					branch[branch_path] = n_leaf_hash
					branch[branch_nibble] = hashed_old_leaf
					
					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)
					rest_nibble = compact_encode(rest_nibble)
					
					//change prefixes
					node.flag_value.encoded_prefix = rest_nibble
					n_leaf.flag_value.encoded_prefix = rest_path
					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf
					
				} else if len(rest) == 0 && len(branch_nib) > 1{
					branch_nibble := branch_nib[0]
					rest_nibble := decoded_key_node[(index_ext+1):]
					branch[branch_nibble] = hashed_old_leaf
					branch[16] = new_value
					
					rest_nibble = compact_encode(rest_nibble)
					node.flag_value.encoded_prefix = rest_nibble
					mpt.db[hashed_old_leaf] = node
											
				} else if len(rest) == 0 && len(branch_nib) == 1 {
					//delete ext
					//add old branch to new
					branch_nibble := branch_nib[0]
					branch[branch_nibble] = hash_node
					delete(mpt.db, hashed_old_leaf)
					branch[16] = new_value
										
				} else if len(branch_nib) == 1  {
					//delete ext
					//add old branch to new
					branch_nibble := branch_nib[0]
					branch[branch_nibble] = hash_node
					delete(mpt.db, hashed_old_leaf)
					
					branch_path := rest[0]
					rest_path := decoded_key[(index_ext+1):]
					branch[branch_path] = n_leaf_hash
					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)
					
					n_leaf.flag_value.encoded_prefix = rest_path
					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf
				}
				f_br := Flag_value {encoded_prefix: []uint8 {}, value:""}
				n_br := Node {node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()
				mpt.db[n_br_hash] = n_br
						
				//create ext						
				common_path = compact_encode(common_path)
				f_ext := Flag_value {encoded_prefix: common_path, value: n_br_hash }
				n_ext := Node {node_type: 2, branch_value: [17]string{}, flag_value: f_ext}
				n_ext_hash := n_ext.hash_node()
				mpt.db[n_ext_hash] = n_ext
				if (isRoot == true) {
					mpt.root = n_ext.hash_node()
				}
							
			} else {
				fmt.Println("no same")
				//no same nibbles - more than one in ext
				//create a branch + take one for branch + add ext to branch []
				//no same - one in ext
				//create a branch + take nibble for branch + delete ext + add old branch to [] in new branch
				//have to change the root -> to branch

				//create branch
				branch := [17] string {}
				first_key := decoded_key[0]
				first_node_key := decoded_key_node[0]

				if (len(decoded_key_node) > 1) {
					fmt.Println("more than one")
					branch[first_key] = n_leaf_hash
					branch[first_node_key] = hashed_old_leaf 
					
					rest_nibble := []uint8{}
					rest_nibble = append(rest_nibble, decoded_key_node[1:]...)
					rest_nibble = compact_encode(rest_nibble)
					
					rest_path := []uint8{}
					rest_path = append(rest_path, decoded_key[1:]...)
					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)
					
					//change prefixes
					node.flag_value.encoded_prefix = rest_nibble
					n_leaf.flag_value.encoded_prefix = rest_path
					
					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf
					
				} else {
					fmt.Println("one")
					branch[first_key] = n_leaf_hash
					branch[first_node_key] = hash_node 
					delete(mpt.db, hashed_old_leaf)
					
					rest_path := []uint8{}
					rest_path = append(rest_path, decoded_key[1:]...)
					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)
					
					//change prefixes
					n_leaf.flag_value.encoded_prefix = rest_path				
					mpt.db[n_leaf_hash] = n_leaf
				}
				f_br := Flag_value {encoded_prefix: []uint8 {}, value:""}
				n_br := Node {node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()
				mpt.db[n_br_hash] = n_br
				if (isRoot == true) {
					mpt.root = n_br.hash_node()
				}
			}
			
			//leaf node
			case 2,3:
				//check the key
				fmt.Println(" ----ENCODED-LEAF----- ")
				fmt.Println(key_node)
				fmt.Println(key)
				fmt.Println(" ----ENCODED-LEAF----- ")
				//keys are equal
				decoded_key_node := compact_decode(key_node)
				decoded_key := compact_decode(key)
				
				fmt.Println(" ----DECODED------ ")
				fmt.Println(decoded_key_node)
				fmt.Println(decoded_key)
				fmt.Println(" ----DECODED------ ")


				if decoded_key_node[0] == decoded_key[0] {
					

						
						hashed_old_leaf := mpt.KeyByValue(node)
	
						delete(mpt.db, hashed_old_leaf)
						
						//get old leaf  -> n
						index_ext := intersectionCount (string(decoded_key_node), string(decoded_key))
						
						common_path := decoded_key_node[:index_ext]
						rest := decoded_key[index_ext:]
						branch_nib := decoded_key_node[index_ext:]
						
						//create new leaf
						f_leaf := Flag_value {encoded_prefix: []uint8{}, value: new_value}
						n_leaf := Node {node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
						n_leaf_hash := n_leaf.hash_node()
						
						//create branch
						//if common = key -> add value
						//decrease on one value
						branch := [17] string {}

						if (len(rest)) > 0 && (len(branch_nib)) > 0 {
							branch_path := rest[0]
							branch_nibble := branch_nib[0]
							rest_path := decoded_key[(index_ext+1):]
							rest_nibble := decoded_key_node[(index_ext+1):]
					
							branch[branch_path] = n_leaf_hash
							branch[branch_nibble] = hashed_old_leaf
							
							rest_path = append(rest_path, 16)
							rest_path = compact_encode(rest_path)
							rest_nibble = append(rest_nibble, 16)
							rest_nibble = compact_encode(rest_nibble)
							
							//change paths in leafs
							node.flag_value.encoded_prefix = rest_nibble
							n_leaf.flag_value.encoded_prefix = rest_path
							
							mpt.db[hashed_old_leaf] = node
							mpt.db[n_leaf_hash] = n_leaf
							
						} else if len(rest) == 0 {
							branch_nibble := branch_nib[0]
							rest_nibble := decoded_key_node[(index_ext+1):]
							branch[branch_nibble] = hashed_old_leaf
							branch[16] = new_value
							
							rest_nibble = append(rest_nibble, 16)
							rest_nibble = compact_encode(rest_nibble)
							
							node.flag_value.encoded_prefix = rest_nibble
							mpt.db[hashed_old_leaf] = node						
						
						} else if len(branch_nib) == 0 {
							branch_path := rest[0]
							rest_path := decoded_key[(index_ext+1):]
							branch[branch_path] = n_leaf_hash
							
							rest_path = append(rest_path, 16)
							rest_path = compact_encode(rest_path)
							//change ald leaf

							branch[16] = hash_node
							n_leaf.flag_value.encoded_prefix = rest_path
							mpt.db[n_leaf_hash] = n_leaf
						}
											
						f_br := Flag_value {encoded_prefix: []uint8 {}, value:""}
						n_br := Node {node_type: 1, branch_value: branch, flag_value: f_br}
						n_br_hash := n_br.hash_node()
						mpt.db[n_br_hash] = n_br
						
						//create ext
						common_path = compact_encode(common_path)
						f_ext := Flag_value {encoded_prefix: common_path, value: n_br_hash }
						n_ext := Node {node_type: 2, branch_value: [17]string{}, flag_value: f_ext}
						n_ext_hash := n_ext.hash_node()
						mpt.db[n_ext_hash] = n_ext
						if (isRoot == true) {
							mpt.root = n_ext_hash
						}
						return n_ext_hash
					} else {
						//not equal keys
						hashed_old_leaf := mpt.KeyByValue(node)
						delete(mpt.db, hashed_old_leaf)
						
						//create new leaf
						f_leaf := Flag_value {encoded_prefix: decoded_key, value: new_value}
						n_leaf := Node {node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
						n_leaf_hash := n_leaf.hash_node()
						
						//create branch
						//decrease on one value
						branch := [17] string {}

						if (len(decoded_key)) > 0 && (len(decoded_key_node)) > 0 {
							br_path_value := decoded_key[0]
							path := []uint8{}
							decoded_key = append(path, decoded_key[1:]...)
									
							br_nibble_value := decoded_key_node[0]
							nibble := []uint8{}
							decoded_key_node = append(nibble, decoded_key_node[1:]...)

							branch[br_path_value] = n_leaf_hash
							branch[br_nibble_value] = hashed_old_leaf
							
							decoded_key = append(decoded_key, 16)
							decoded_key = compact_encode(decoded_key)
							decoded_key_node = append(decoded_key_node, 16)
							decoded_key_node = compact_encode(decoded_key_node)
							
							//change paths in leafs
							node.flag_value.encoded_prefix = decoded_key_node
							n_leaf.flag_value.encoded_prefix = decoded_key
							
							mpt.db[hashed_old_leaf] = node
							mpt.db[n_leaf_hash] = n_leaf
							
						} else if len(decoded_key) == 0 {
							br_nibble_value := decoded_key_node[0]
							nibble := []uint8{}
							decoded_key_node = append(nibble, decoded_key_node[1:]...)
							
							branch[br_nibble_value] = hashed_old_leaf
							branch[16] = new_value
							
							decoded_key_node = append(decoded_key_node, 16)
							decoded_key_node = compact_encode(decoded_key_node)
							
							node.flag_value.encoded_prefix = decoded_key_node
							mpt.db[hashed_old_leaf] = node
							
						
						} else if len(decoded_key_node) == 0 {
							br_path_value := decoded_key[0]
							path := []uint8{}
							decoded_key = append(path, decoded_key[1:]...)
							
							branch[br_path_value] = n_leaf_hash
							
							decoded_key = append(decoded_key, 16)
							decoded_key = compact_encode(decoded_key)
							//change old leaf

							branch[16] = hash_node
							n_leaf.flag_value.encoded_prefix = decoded_key
							mpt.db[n_leaf_hash] = n_leaf

						}
						
						
						f_br := Flag_value {encoded_prefix: []uint8 {}, value:""}
						n_br := Node {node_type: 1, branch_value: branch, flag_value: f_br}
						n_br_hash := n_br.hash_node()
						mpt.db[n_br_hash] = n_br
						if (isRoot == true) {
							mpt.root = n_br_hash
						}
						return n_br_hash
					}
		}
}
	
	return ""
}


func (mpt *MerklePatriciaTrie) GetByNode(node Node, key []byte) Node {
	n := Node{}
	fmt.Println("start")
	fmt.Println(node)

	if len(key) == 0 {
		return node
	}

	if (node.node_type == 0) {
		return n
	}
	
	switch node.node_type {	
	case 2:
	fmt.Println("case2")
		key_node := node.flag_value.encoded_prefix
		key_decoded := compact_decode(key_node)
	fmt.Println(key)
	fmt.Println(key_decoded)
		hash_node := node.flag_value.value //hash_value of child Node
		n = mpt.db[hash_node] //child node (next node)

		node_prefix := key_node[0]/16

		//check for leaf and ext
		switch node_prefix {
			//extension node
			case 0,1:
			fmt.Println("case 01")
				//decerase nibbles
				index := intersectionCount (string(key_decoded), string(key))
				if len(key) > len(key_decoded) && index > 0 {
				//decrease key from beginning
				fmt.Println("1if")
				return mpt.GetByNode(n, key[index:])
				//same key all - go to branch node - return value
				} else if len(key) == len(key_decoded) && string(key_decoded) == string(key) {
					fmt.Println("2if")
					return mpt.db[node.flag_value.value]
				} else if (len(key) <= len(key_decoded)) && string(key_decoded) != string(key) {
					fmt.Println("3if")
					return Node{}			
				} else {
					fmt.Println("4if")
					return n
				}
				
			//leaf node
			case 2,3:
			fmt.Println("case23")
				//same key - return node
				//not same - return nil
				if string(key_decoded) == string(key) {
					return node
				} else {
					return n
				}

			default:
				return node
		}
	//branch node
	case 1:
	fmt.Println("case1")
		fmt.Println(key)
		//no hash in branch

		if (len(key) == 0){
			return node
		} else {
			br_path_value := key[0]

			path := []uint8{}
			path = append(path, key[1:]...)
			//rest_path := key[:len(key)-1]
//			fmt.Println("br_path_value")
			//fmt.Println(key[:len(key)-1])
//			fmt.Println(br_path_value)

			
			
			new_hash_node := node.branch_value[br_path_value]
//			fmt.Println(new_hash_node)

			
			n = mpt.db[new_hash_node]
			return mpt.GetByNode(n, path)
		}
				
		//case 0
	default:
//	fmt.Println("default")
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

func intersectionCount(key string, key_node string) int {
	count := 0
	min := 0
	if len(key) > len (key_node) {
		min = len(key_node)
	} else {
		min = len(key)
	}
	for i:=0; i < min; i += 1 {
		if (key[i] == key_node[i]){
				count += 1
			} else {
				break
			}
	}
		return count
}

func (mpt *MerklePatriciaTrie) KeyByValue (node Node) string {
	for key, value := range mpt.db {
		if reflect.DeepEqual(node, value) {
			return key	
		}
	}
	return ""
}


