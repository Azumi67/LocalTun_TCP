package network

import (
	"crypto/rsa"
	"log"
	"time"

	"github.com/songgao/water"
	"github.com/xtaci/smux"
)

func HandleSmux(session *smux.Session, tun *water.Interface, publicKey *rsa.PublicKey, verbose bool) {
	defer session.Close()

	for {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Printf("Couldn't accept smux stream: %v", err)

			for retry := 0; retry < 3; retry++ {
				time.Sleep(5 * time.Second)
				stream, err = session.AcceptStream()
				if err == nil {
					log.Printf("smux stream re-established after %d attempt(s)", retry+1)
					break
				}
				log.Printf("Retrying to accept smux stream: %v", err)
			}

			if err != nil {
				continue
			}
		}

		log.Printf("smux stream accepted successfully")
		Handle(stream, tun, publicKey, verbose)
	}
}
