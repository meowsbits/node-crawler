package crawler

import (
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
)

func (c Crawler) makeDiscoveryConfig() (*enode.LocalNode, discover.Config) {
	var cfg discover.Config
	var err error

	if c.NodeKey != "" {
		key, err := crypto.HexToECDSA(c.NodeKey)
		if err != nil {
			panic(fmt.Errorf("-%s: %v", c.NodeKey, err))
		}
		cfg.PrivateKey = key
	} else {
		cfg.PrivateKey, _ = crypto.GenerateKey()
	}

	cfg.Bootnodes, err = c.parseBootnodes()
	if err != nil {
		panic(err)
	}

	return enode.NewLocalNode(c.NodeDB, cfg.PrivateKey), cfg
}

func listen(ln *enode.LocalNode, addr string) *net.UDPConn {
	socket, err := net.ListenPacket("udp4", addr)
	if err != nil {
		panic(err)
	}
	usocket := socket.(*net.UDPConn)
	uaddr := socket.LocalAddr().(*net.UDPAddr)
	if uaddr.IP.IsUnspecified() {
		ln.SetFallbackIP(net.IP{127, 0, 0, 1})
	} else {
		ln.SetFallbackIP(uaddr.IP)
	}
	ln.SetFallbackUDP(uaddr.Port)
	return usocket
}

func (c Crawler) parseBootnodes() ([]*enode.Node, error) {
	bootnodes := params.MainnetBootnodes
	if c.Sepolia {
		bootnodes = params.SepoliaBootnodes
	}
	if c.Goerli {
		bootnodes = params.GoerliBootnodes
	}
	if c.Classic {
		bootnodes = params.ClassicBootnodes
		bootnodes = []string{
			"enode://ce7cbef463c5bac7e310a1ad2279975d8a5b38627f462cd6157a8bd310641ae9e33aae7634afc8b102be864491cba8dd33c9343b94674ba86be8afc74f4721a9@159.203.56.33:30369",
			// "enode://5e85df7bc6d529647cf9a417162784a89b7ccf2b8e1570fadb6fdf9fa025c8ec2257825d1ec5d7357a6f49898fdfbd9c4c56d22645dbe8b8a6aa67dacbcf3ecc@157.230.152.87:30303",
		}
	}
	if c.Mordor {
		bootnodes = params.MordorBootnodes
	}
	if len(c.Bootnodes) != 0 {
		bootnodes = c.Bootnodes
	}

	nodes := make([]*enode.Node, len(bootnodes))
	var err error
	for i, record := range bootnodes {
		nodes[i], err = parseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %v", err)
		}
	}
	return nodes, nil
}
