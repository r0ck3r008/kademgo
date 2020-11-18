package connector

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/r0ck3r008/kademgo/utils"
	"google.golang.org/protobuf/proto"
)

// Ping requests for a ping reply from the passed in UDP address and returns a bool return value if
// a reply shows up. It waits for utils.PINGWAIT amount of time before it fails.
func (conn_p *Connector) Ping(srchash *[utils.HASHSZ]byte, addr_p *net.UDPAddr) bool {
	var rand_num int64 = int64(rand.Int())
	addr := *addr_p
	var pkt Pkt = Pkt{Type: Pkt_PingReq, Hash: hex.EncodeToString((*srchash)[:]), RandNum: rand_num, Hops: utils.MAXHOPS}
	cmds, err := proto.Marshal(&pkt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in marshaling: %s\n", err)
		os.Exit(1)
	}
	var env Envelope = Envelope{rand_num, cmds, addr}
	conn_p.sch <- env
	time.Sleep(utils.PINGWAIT)
	// Fetch result from map
	var ret bool = false
	conn_p.mut.Lock()
	if _, ok := conn_p.pcache[rand_num]; ok {
		ret = true
		delete(conn_p.pcache, rand_num)
	}
	conn_p.mut.Unlock()

	return ret
}
