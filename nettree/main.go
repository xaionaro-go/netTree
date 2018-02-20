package main

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/vishvananda/netlink"
	"github.com/xaionaro-go/netTree"
	"strings"
)

func recursivePrint(node *netTree.Node, level int) {
	ifaceType := strings.ToLower(strings.Split(fmt.Sprintf("%T", node.Link), ".")[1])
	ifaceName := ""
	switch link := node.Link.(type) {
	case *netlink.Bond:
		ifaceName = link.LinkAttrs.Name
	case *netlink.Veth:
		ifaceName = link.LinkAttrs.Name
	case *netlink.Bridge:
		ifaceName = link.LinkAttrs.Name
	case *netlink.Vlan:
		ifaceName = link.LinkAttrs.Name
	case *netlink.Device:
		ifaceName = link.LinkAttrs.Name
	}
	fmt.Printf("%v- %v (%v)\n", strings.Repeat(" ", level*2), ifaceName, ifaceType)
	for _, child := range node.Children {
		recursivePrint(child, level+1)
	}
}

func main() {
	rootNode := netTree.GetTree()
	for _, child := range rootNode.Children {
		recursivePrint(child, 0)
	}
	//spew.Dump(rootNode)
}
