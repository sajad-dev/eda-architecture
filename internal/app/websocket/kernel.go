package websocket

var Addrs []Addr

func AddAddr(addr Addr) {
	Addrs = append(Addrs, addr)
}
