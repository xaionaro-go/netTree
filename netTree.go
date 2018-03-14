package netTree

import (
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
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
	slavesMap := map[int][]*Node{}
	childrenMap := map[int][]*Node{}

	for _, linkI := range linkList {
		node := Node{Link: linkI}
		indexMap[linkI.Attrs().Index] = &node
	}

	for _, node := range indexMap {
		attrs := node.Link.Attrs()
		//spew.Dump(node.Link)
		parentIndex := attrs.ParentIndex
		masterIndex := attrs.MasterIndex

		_, nodeIsDevice := node.Link.(*netlink.Device)

		if parentIndex != 0 {
			parent := indexMap[parentIndex]
			if parent == nil {
				childrenMap[0] = append(childrenMap[0], node)
			} else {
				if _, ok := parent.Link.(*netlink.Bridge); ok {
					childrenMap[0] = append(childrenMap[0], node)
				} else {
					_, parentIsVeth   := parent.Link.(*netlink.Veth)
					_, parentIsDevice := parent.Link.(*netlink.Device)
					if parentIsVeth || parentIsDevice {
						childrenMap[parentIndex] = append(childrenMap[parentIndex], node)
					} else {
						if masterIndex == 0 && !nodeIsDevice {
							childrenMap[0] = append(childrenMap[0], node)
						}
					}
				}
			}
		}

		if masterIndex != 0 {
			master := indexMap[masterIndex]
			if _, ok := master.Link.(*netlink.Bridge); ok {
				slavesMap[masterIndex] = append(slavesMap[masterIndex], node)
			} else {
				//childrenMap[0] = append(childrenMap[0], node)
			}
		} else {
			if nodeIsDevice {
				childrenMap[0] = append(childrenMap[0], node)
			}
		}
	}

	for index, node := range indexMap {
		newSlaves   := slavesMap[index]
		newChildren := childrenMap[index]
		for _, slave := range newSlaves {
			slave.Children = append(slave.Children, node)
			node.Parents   = append(node.Parents, slave)
		}
		for _, child := range newChildren {
			child.Parents = append(child.Parents, node)
			node.Children = append(node.Children, child)
		}
	}

	//fmt.Println(slavesMap)
	//fmt.Println(childrenMap)

	rootNode.Children = childrenMap[0]

	/*for index, children := range childrenMap {
		if indexMap[index] != nil {
			continue
		}
		aggNode := Node{}
		aggNode.Children = children
		rootNode.Children = append(rootNode.Children, &aggNode)
	}*/

	return &rootNode
}

