package connector

import (
	"fmt"
	"net"
	"os"
	sync "sync"

	"github.com/r0ck3r008/kademgo/utils"
	"google.golang.org/protobuf/proto"
)

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

	return conn, nil
}

	}
	}
}

	}
}
