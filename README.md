```sh
go install github.com/xaionaro-go/netTree/nettree
"$GOPATH"/bin/nettree
```

```
- lo (device)
- eth0 (device)
- eth1 (device)
  - eth1br (bridge)
  - eth1.6 (vlan)
    - vlan6 (bridge)
  - eth1.8 (vlan)
    - vlan8 (bridge)
```
