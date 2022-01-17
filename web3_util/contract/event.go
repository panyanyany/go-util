package contract

import (
	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/abi"
)

// Event is a solidity event
type Event struct {
	event *abi.Event
}

// Encode encodes an event
func (e *Event) Encode() web3.Hash {
	return e.event.ID()
}

// ParseLog parses a log
func (e *Event) ParseLog(log *web3.Log) (map[string]interface{}, error) {
	return abi.ParseLog(e.event.Inputs, log)
}
