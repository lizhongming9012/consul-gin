package broadcast

import (
	"github.com/pion/webrtc/v2"
	"sync"
)

// BroadcastHub keeps a list of all channels
type BroadcastHub struct {
	BroadcastChannel chan []byte
	listenChannels   map[*uint16]*webrtc.DataChannel
	dataMutex        *sync.RWMutex
}

func NewHub() *BroadcastHub {
	hub := &BroadcastHub{
		BroadcastChannel: make(chan []byte),
		listenChannels:   make(map[*uint16]*webrtc.DataChannel),
		dataMutex:        new(sync.RWMutex),
	}
	go hub.run()
	return hub
}

func (h *BroadcastHub) AddListener(d *webrtc.DataChannel) {
	h.dataMutex.Lock()
	h.listenChannels[d.ID()] = d
	h.dataMutex.Unlock()
}

func (h *BroadcastHub) run() {
	for message := <-h.BroadcastChannel; ; message = <-h.BroadcastChannel {
		h.dataMutex.RLock()
		channels := h.listenChannels
		h.dataMutex.RUnlock()
		for _, client := range channels {
			if err := client.SendText(string(message)); err != nil {
				h.dataMutex.Lock()
				delete(h.listenChannels, client.ID())
				h.dataMutex.Unlock()
			}
		}
	}
}
