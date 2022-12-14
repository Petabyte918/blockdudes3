package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"log"
	"strconv"
	"sync"
)

type Client struct {
	room *Room
	ws *websocket.Conn
	wrtc *webrtc.PeerConnection
	dc *webrtc.DataChannel
	mu sync.Mutex

	id IdType
	name string
	voice bool
}

func NewClient(room* Room, ws *websocket.Conn, name string) *Client {
	client := &Client {
		room: room,
		ws: ws,
		wrtc: nil,
		dc: nil,

		id: room.nextClientId,
		name: name,
		voice: false,
	}
	go client.run()

	room.nextClientId += 1
	return client
}

func (c *Client) run() {
	defer func() {
		c.room.unregister <- c
	}()

	for {
		_, b, err := c.ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected socket error: %v", err)
			}
			return
		}

		imsg := IncomingMsg{
			b: b,
			client: c,
		}
		c.room.incoming <- imsg
	}
}

func (c Client) GetDisplayName() string {
	return c.name + " #" + strconv.Itoa(int(c.id))
}

func (c *Client) GetClientData() ClientData {
	return ClientData {
		Id: c.id,
		Name: c.name,
	}
}

func (c *Client) Send(msg interface{}) error {
	b := Pack(msg)
	return c.SendBytes(b)
}

func (c *Client) SendBytes(b []byte) error {
	// Lock required to synchronize writes from room and WebRTC callbacks
	c.mu.Lock()
	err := c.ws.WriteMessage(websocket.BinaryMessage, b)
	c.mu.Unlock()
	return err
}

func (c *Client) SendUDP(msg interface{}) error {
	b := Pack(msg)
	return c.SendBytesUDP(b)
}

func (c *Client) SendBytesUDP(b []byte) error {
	if c.dc == nil {
		return errors.New("Data channel not initialized")
	}
	return c.dc.Send(b)
}

func (c *Client) InitWebRTC(onSuccess func()) error {
	var err error
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun2.l.google.com:19302",
					"stun:openrelay.metered.ca:80",
				},
			},
		},
	}

	log.Printf("Starting new WebRTC connection")
	c.wrtc, err = webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}

	c.wrtc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("WebRTC connection state for %s has changed: %s", c.GetDisplayName(), s.String())
	})

	ordered := false
	maxRetransmits := uint16(0)
	dcInit := &webrtc.DataChannelInit {
		Ordered: &ordered,
		MaxRetransmits: &maxRetransmits,
	}
	c.dc, err = c.wrtc.CreateDataChannel("data", dcInit)
	if err != nil {
		return err
	}

	c.dc.OnOpen(func() {
		onSuccess()
		log.Printf("Opened data channel for %s: %s-%d", c.GetDisplayName(), c.dc.Label(), c.dc.ID())
	})

	c.dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		imsg := IncomingMsg{
			b: msg.Data,
			client: c,
		}
		c.room.incoming <- imsg
	})

	c.wrtc.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice == nil {
			return
		}

		candidateMsg := JSONMsg {
			T: candidateType,
			JSON: ice.ToJSON(),
		}

		c.Send(&candidateMsg)
	})

	return nil
}

func (c *Client) processWebRTCOffer(json interface{}) error {
	log.Printf("Received WebRTC offer for %s", c.GetDisplayName())

	offer, ok := json.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to parse offer: %+v", json)
	}
	var err error

	desc := webrtc.SessionDescription {
		Type: webrtc.SDPTypeOffer,
		SDP: offer["sdp"].(string),
	}
	err = c.wrtc.SetRemoteDescription(desc)
	if err != nil {
		return err
	}

	answer, err := c.wrtc.CreateAnswer(nil)
	if err != nil {
		return err
	}
	c.wrtc.SetLocalDescription(answer)

	answerMsg := JSONMsg {
		T: answerType,
		JSON: answer,
	}
	c.Send(&answerMsg)
	return nil
}

func (c *Client) processWebRTCCandidate(json interface{}) error {
	log.Printf("Received WebRTC ICE candidate for %s", c.GetDisplayName())

	candidate, ok := json.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to parse offer message: %+v", json)
	}
	var err error

	sdpMid := candidate["sdpMid"].(string)
	sdpMLineIndex := uint16(candidate["sdpMLineIndex"].(int8))
	candidateInit := webrtc.ICECandidateInit {
		Candidate: candidate["candidate"].(string),
		SDPMid: &sdpMid,
		SDPMLineIndex: &sdpMLineIndex,
	}
	err = c.wrtc.AddICECandidate(candidateInit)
	if err != nil {
		return err
	}
	return nil
}