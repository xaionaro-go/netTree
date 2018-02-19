package netTree

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

type Node struct {
	netlink.Link
	Parents Nodes
	Children Nodes
}
type Nodes []*Node

func (node *Node) ToSlice() (result Nodes) {
	result = append(result, node)
	result = append(result, node.Children.ToSlice()...)
	return
}
func (nodes Nodes) ToSlice() (result Nodes) {
	result = append(result, nodes...)
	for _, node := range nodes {
		result = append(result, node.Children.ToSlice()...)
	}
	return
}

func GetTree() *Node {
	rootNode := Node{}

	linkList, err := netlink.LinkList()
	if err != nil {
		panic(err)
	}

	indexMap := map[int]*Node{}
	parentsMap := map[int][]*Node{}
	childrenMap := map[int][]*Node{}

	for _, linkI := range linkList {
		nodeIndex := 0
		parentIndex := 0
		masterIndex := 0
		isPhysIface := false
		switch link := linkI.(type) {
		case *netlink.Veth:
			nodeIndex   = link.LinkAttrs.Index
			parentIndex = link.LinkAttrs.ParentIndex
			masterIndex = link.LinkAttrs.MasterIndex
		case *netlink.Bridge:
			nodeIndex   = link.LinkAttrs.Index
			parentIndex = link.LinkAttrs.ParentIndex
			masterIndex = link.LinkAttrs.MasterIndex
		case *netlink.Vlan:
			nodeIndex   = link.LinkAttrs.Index
			parentIndex = link.LinkAttrs.ParentIndex
			masterIndex = link.LinkAttrs.MasterIndex
		case *netlink.Device:
			nodeIndex   = link.LinkAttrs.Index
			parentIndex = link.LinkAttrs.ParentIndex
			masterIndex = link.LinkAttrs.MasterIndex
			isPhysIface = true
		default:
			fmt.Printf("Skipped type: %T\n", linkI)
			continue
		}

		node := Node{Link: linkI}
		indexMap[nodeIndex] = &node
		if parentIndex != 0 || masterIndex != 0 || isPhysIface {
			childrenMap[parentIndex] = append(childrenMap[parentIndex], &node)
		}

		if masterIndex != 0 {
			parentsMap[masterIndex] = append(parentsMap[masterIndex], &node)
		}
	}

	for index, node := range indexMap {
		newParents  := parentsMap[index]
		newChildren := childrenMap[index]
		for _, parent := range newParents {
			parent.Children = append(parent.Children, node)
			node.Parents    = append(node.Parents, parent)
		}
		for _, child := range newChildren {
			child.Parents = append(child.Parents, node)
			node.Children = append(node.Children, child)
		}
	}

	//fmt.Println(parentsMap)
	//fmt.Println(childrenMap)

	rootNode.Children = childrenMap[0]

	return &rootNode
}

