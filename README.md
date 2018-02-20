```sh
go get github.com/xaionaro-go/netTree/nettree
go install github.com/xaionaro-go/netTree/nettree
"$GOPATH"/bin/nettree
```

```
- lo (device)
- eth0 (device)
  - bond0 (bond)
    - bond0.2090 (vlan)
      - vlan2090 (bridge)
    - bond0.2091 (vlan)
      - vlan2091 (bridge)
```

---

```sh
"$GOPATH"/bin/nettree -a
```

```
- lo (device)
- eth0 (device)
  - bond0 (bond)
    - bond0.2090 (vlan)
      - vlan2090 (bridge)
    - bond0.2091 (vlan)
      - vlan2091 (bridge)
- eth1 (device)
   10.5.6.131
```

---

A bash script to print the network interfaces tree: [https://github.com/zabojcampula/show-net-devices-tree](https://github.com/zabojcampula/show-net-devices-tree)
