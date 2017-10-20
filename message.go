package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clockworksoul/smudge"
)

// messageType is an alias for int8. Whenever you see "messageType" in the code,
// think int8. Using this alias however allows us to increase code clarity and
// safety beyond simply passing meaningless numbers around.
type messageType int8

const (
	// Go does not have any native support for enum types, however there is an
	// idiomatic way to accomplish the same goal.
	// https://golang.org/ref/spec#Iota

	messageTypeChat messageType = iota + 1
	messageTypeUsernames
	messageTypeUsernameReq
)

// message represents the structure of the contents in a smudge.Broadcast. We
// can use the Type to determine what the Body will contain.
type message struct {
	// The text in backticks after each field here is called a Struct Tag.
	// The standard JSON marshaller uses the "json" struct tag to name fields.
	//
	// https://golang.org/pkg/reflect/#StructTag
	// https://golang.org/pkg/encoding/json/#Marshal
	//
	// The JSON marshaller can only interact with "exported fields", therefore
	// these field names must start with an uppercase letter.

	// Type informs us what to expect in the body of this message, and what
	// action to take on it.
	Type messageType `json:"type"`

	// Body contains the bulk of the message
	Body string `json:"body"`

	// Usernames is filled only in a messageTypeUsernames. It contains a map
	// of the address->username pairings know by the sending client.
	Usernames map[NodeAddress]string `json:"usernames"`
}

// Encode converts the message into a form which can be sent to other clients
// through Smudge (a []byte, pronounced byte slice).
//
// Notice how this func has an extra set of parens before the name? This is how
// methods are defined in Go. Unlike a function, methods are not in the global
// namespace and must be called on an instance of an object. That object being
// called on is in the first set of parens, and is called the receiver. This of
// this as "self" in python, or "this" in many other languages.
// More info: https://tour.golang.org/methods/1
func (m *message) Encode() []byte {
	// There is a lot happening here in a pretty small space. We first create an
	// empty buffer in which we can temporarily store some bytes. This buffer
	// implements the io.Writer interface, but we want to write compressed
	// bytes, so we wrap that writer in the zlib writer which also implements
	// the io.Writer interface. Finally we create a json encoder which will
	// output the json format of m into the zlib writer.
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		printError("Failed to marshal a chat message to send: %s", err)
	}
	err = w.Close() // The bytes might not actually be written until closed (or flushed)
	if err != nil {
		printError("Failed to close the encoding writer: %s", err)
	}

	// read out the contents from our temporary buffer, and return them
	return b.Bytes()
}

// Messenger contains all the messages which we know have been sent in the past.
// Additionally, it provides the interface for sending and receiving new
// messages.
//
// To send and receive messages we must implement the Broadcast functionality in
// the smudge library.
//
// https://godoc.org/github.com/clockworksoul/smudge#BroadcastListener
// https://godoc.org/github.com/clockworksoul/smudge#BroadcastString
type Messenger struct {
	// clients is the list of all known and alive clients. Maintaining a
	// reference here will allow us to update status based on broadcasts.
	clients ClientList
}

// Decode converts the byte slice received from a broadcast into a usable
// message. This is the reverse of the Encode() operation.
//
// Just like how the encode method above uses the json and zlib packages to
// json marshal and then compress a message, here we are doing the reverse.
func (m *message) Decode(data []byte) error {
	bb := bytes.NewReader(data)
	r, err := zlib.NewReader(bb)
	if err != nil {
		return fmt.Errorf("Failed to decompress message: %s", err)
	}

	// msg is what the decompressed bytes will be un-json-marshalled into
	err = json.NewDecoder(r).Decode(m)
	if err != nil {
		return fmt.Errorf("Failed to decode message: %s", err)
	}

	return nil
}

// OnBroadcast is the only method defined on the smudge.BroadcastListener
// interface. By implementing this method on the Messenger struct, that struct
// will satisfy the interface and we can register it with smudge.
//
// When another node in the gossip cluster sends a broadcast message, this
// function will be called.
func (m *Messenger) OnBroadcast(b *smudge.Broadcast) {
	senderAddr := NodeAddress(b.Origin().Address())

	printDebug("Received %d bytes", len(b.Bytes()))
	var msg message
	err := msg.Decode(b.Bytes())
	if err != nil {
		printError("Failed to receive message from %s: %s", senderAddr, err)
		return
	}

	switch msg.Type {
	case messageTypeUsernames:
		printDebug("Received a broadcast containing usernames")

		if msg.Usernames == nil || len(msg.Usernames) == 0 {
			printError("Received an empty username list")
			return
		}

		err := m.clients.AddUsernames(msg.Usernames)
		if err != nil {
			printError("Failed to process received usernames: %s", err)
		}
	case messageTypeUsernameReq:
		printDebug("Received a broadcast requesting %s send usernames, my localAddress is %s", msg.Body, localAddress)

		if msg.Body == string(localAddress) {
			// The request targeted us...
			// Let's send all the usernames we know about to minimize requests
			// for a new client.
			err := m.clients.BroadcastUsernames()
			if err == nil {
				printInfo("Successfully broadcast usernames to the group")
			} else {
				printError("Tried to broadcast usernames but failed: %s", err)
			}
		}
	case messageTypeChat:
		// Received a chat message

		sender := m.clients[senderAddr]
		printChatMessage(msg.Body, sender.GetName())
	}
}

// SendMessage takes a chat message to be sent and broadcasts it to the cluster
// and posts to the local chat view.
func SendMessage(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	// First let's make the message show up in our own chat history
	printChatMessage(text, localUsername)

	// Now we can send it on to others
	msg := message{
		Type: messageTypeChat,
		Body: text,
	}

	return smudge.BroadcastBytes(msg.Encode())
}
