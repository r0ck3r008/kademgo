package connector

import (
	"fmt"
	"net"
	"os"
	sync "sync"

	"github.com/r0ck3r008/kademgo/utils"
	"google.golang.org/protobuf/proto"
)

type Envelope struct {
	id   int64
	cmds []byte
	addr net.UDPAddr
}

type Connector struct {
	conn   *net.UDPConn
	pcache map[int64]Envelope
	mut    *sync.Mutex
	sch    chan Envelope
	rch    chan Envelope
}

func ConnectorInit(addr *string) (*Connector, error) {
	conn_p, err := net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: utils.GENPORT, Zone: ""})
	if err != nil {
		return nil, fmt.Errorf("UDP Create: %s", err)
	}
	sch := make(chan Envelope, 100)
	rch := make(chan Envelope, 100)
	var mut *sync.Mutex = &sync.Mutex{}

	var conn *Connector = &Connector{conn: conn_p, sch: sch, rch: rch, mut: mut}
	return conn, nil
}

func (conn_p *Connector) Collector() {
	for env := range conn_p.sch {
		// Acquire write lock and write to cache
		conn_p.mut.Lock()
		conn_p.pcache[env.id] = env
		conn_p.mut.Unlock()
	}
}
func (conn_p *Connector) ReadLoop() {
	for {
		var cmdr []byte
		_, addr_p, err := conn_p.conn.ReadFromUDP(cmdr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in reading: %s\n", err)
			close(conn_p.rch)
			break
		}
		// Extra UnMarshal due to pesky Mutex in Google Protobuf which stops from being sent on a channel
		var pkt *Pkt = &Pkt{}
		err = proto.Unmarshal(cmdr, pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in unmarshalling: %s\n", err)
			os.Exit(1)
		}
		hops := pkt.GetHops()
		if hops != 0 {
			pkt.Hops = hops - 1
			var id int64 = pkt.GetRandNum()
			cmdr, err = proto.Marshal(pkt)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error in marshaling: %s\n", err)
				os.Exit(1)
			}
			var env Envelope = Envelope{id, cmdr, *addr_p}
			conn_p.sch <- env
		}
	}
}

func (conn_p *Connector) WriteLoop() {
	for env := range conn_p.sch {
		if _, err := conn_p.conn.WriteToUDP(env.cmds, &env.addr); err != nil {
			fmt.Fprintf(os.Stderr, "Error in writing: %s\n", err)
			break
		}
	}
}
