package merkleTree

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode
	if len(nodes)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for i, _ := range data {
		node := NewMerkleNode(nil, nil, data[i])
		nodes = append(nodes, *node)
	}

	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	return &MerkleTree{&nodes[0]}
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := new(MerkleNode)

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHash)
		node.Data = hash[:]
	}
	node.Left = left
	node.Right = right
	return node
}
