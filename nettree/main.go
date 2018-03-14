package main

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/vishvananda/netlink"
	"github.com/pborman/getopt/v2"
	"github.com/xaionaro-go/netTree"
	"os"
	"strings"
)

func recursivePrint(node *netTree.Node, level int, showAddrs bool) {
	ifaceType := ""
	typeWords := strings.Split(fmt.Sprintf("%T", node.Link), ".")
	if len(typeWords) == 2 {
		ifaceType = strings.ToLower(typeWords[1])
	}
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
	if showAddrs {
		addrs, err := netlink.AddrList(node.Link, netlink.FAMILY_V4)
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			fmt.Printf("%v    ip4 %v\n", strings.Repeat(" ", level*2), addr.Peer.IP.To4())
		}
	}
	for _, child := range node.Children {
		recursivePrint(child, level+1, showAddrs)
	}
}

func main() {
	showAddrsPtr := getopt.BoolLong("show-addresses", 'a', "", "print IP-addresses")
	err := getopt.Getopt(nil)
	if err != nil {
		// code to handle error
		fmt.Fprintln(os.Stderr, err)
	}
	getopt.Parse()

	showAddrs := false
	if showAddrsPtr != nil {
		showAddrs = *showAddrsPtr
	}

	rootNode := netTree.GetTree()
	for _, child := range rootNode.Children {
		recursivePrint(child, 0, showAddrs)
	}
	//spew.Dump(rootNode)
}
