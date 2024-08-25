package network

import (
	"log"

	"github.com/xtaci/smux"
	"github.com/songgao/water"
)

func HandleSmux(session *smux.Session, tun *water.Interface, verbose bool) {
	defer session.Close()

	for {
		stream, err := session.OpenStream()
		if err != nil {
			log.Printf("Couldn't open smux stream: %v", err)
			continue
		}

		log.Printf("smux stream opened successfully")
		go HandleConnection(stream, tun, verbose)
	}
}
