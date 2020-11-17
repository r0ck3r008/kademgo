package connector

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/r0ck3r008/kademgo/node"
	"github.com/r0ck3r008/kademgo/utils"
	"google.golang.org/protobuf/proto"
)

type Connector struct {
	conn   *net.UDPConn
	pinger *net.UDPConn
}

func ConnectorInit(addr *string) (*Connector, error) {
	conn_p, err := net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: utils.GENPORT, Zone: ""})
	if err != nil {
		return nil, fmt.Errorf("UDP Create: %s", err)
	}

	pinger_p, err := net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: utils.PINGPORT, Zone: ""})
	if err != nil {
		return nil, fmt.Errorf("UDP Create: %s", err)
	}

	var conn *Connector = &Connector{conn: conn_p, pinger: pinger_p}
	return conn, nil
}

// PingReq requests for a ping reply from the passed in UDP address.
func (conn_p *Connector) PingReq(hash *[utils.HASHSZ]byte, addr_p *net.IP) {
	var pkt_req *node.Pkt = &Pkt{Type: Pkt_PingReq, Hash: hex.EncodeToString((*hash)[:])}
	cmds, err := proto.Marshal(pkt_req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in marshalling: %s\n", err)
		os.Exit(1)
	}
	if _, err := conn_p.pinger.WriteToUDP(cmds, &net.UDPAddr{IP: *addr_p, Port: utils.PINGPORT, Zone: ""}); err != nil {
		fmt.Fprintf(os.Stderr, "Error in sending UDP msg: %s\n", err)
		os.Exit(1)
	}
}

// PingRep replies a to a ping request to the neighbour passed in as parameter.
func (conn_p *Connector) PingRep(hash *[utils.HASHSZ]byte, addr_p *net.IP) {
	var pkt_rep *node.Pkt = &Pkt{Type: Pkt_PingRep, Hash: hex.EncodeToString((*hash)[:])}
	cmds, err := proto.Marshal(pkt_rep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in marshalling: %s\n", err)
		os.Exit(1)
	}
	if _, err := conn_p.pinger.WriteToUDP(cmds, &net.UDPAddr{IP: *addr_p, Port: utils.PINGPORT, Zone: ""}); err != nil {
		fmt.Fprintf(os.Stderr, "Error in sending UDP msg: %s\n", err)
		os.Exit(1)
	}
}
