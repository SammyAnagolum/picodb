package proxy

import (
	"errors"
	"net"
	"slices"

	"github.com/segmentio/fasthash/fnv1a"
)

type ConsistentHash struct {
	NumSlots      uint32
	ServersOnRing []uint32
}

func NewConsistentHash(numSlots uint32) ConsistentHash {
	return ConsistentHash{
		NumSlots:      numSlots,
		ServersOnRing: make([]uint32, 0, 1024)}
}

func (ch *ConsistentHash) hashServerIP(serverConn *net.Conn) uint32 {
	return ch.hashString((*serverConn).RemoteAddr().String())
}

func (ch *ConsistentHash) hashString(key string) uint32 {
	return fnv1a.HashString32(key) % ch.NumSlots
}

func (ch *ConsistentHash) AddToHashRing(id uint32) {
	ch.ServersOnRing = append(ch.ServersOnRing, id)
	slices.Sort(ch.ServersOnRing)
}

func (ch *ConsistentHash) GetNextServerIDOnRing(keyID uint32) (uint32, error) {
	if len(ch.ServersOnRing) == 0 {
		return 0, errors.New("No servers available!")
	}

	for i := 0; i < len(ch.ServersOnRing); i++ {
		if ch.ServersOnRing[i] >= keyID {
			return ch.ServersOnRing[i], nil
		}
	}

	return ch.ServersOnRing[0], nil
}
