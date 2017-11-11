package main

import (
	"time"

	"github.com/clockworksoul/smudge"
)

// NodeAddress is just an alias for strings, but increases clarity in the map
// keys of the ClientList.
type NodeAddress string

// ClientList contains all clients which are currently connected to the cluster.
//
// Additionally, because this struct has methods defined on it which fulfill the
// requirements to be a smudge.StatusListener, it is used to handle
// notifications about added or removed clients.
//
// https://godoc.org/github.com/clockworksoul/smudge#StatusListener
type ClientList map[NodeAddress]ChatClient

// OnChange is the only method defined on the smudge.StatusListener. By
// implementing this method on the ClientStatusListener struct, that struct will
// satisfy the interface and we can register it with smudge.
//
// When a client is added or removed from the gossip cluster, update our
// internal list of the membership. We can use this internally maintained
// membership list to display a friends list.
func (cl ClientList) OnChange(node *smudge.Node, status smudge.NodeStatus) {
	if status == smudge.StatusAlive {
		printDebug("Adding a new node: %s", node.Address())
		cl.AddClient(node)
	} else {
		printDebug("Removing a node: " + node.Address())
		cl.RemoveClient(node)
	}

	printClientList(cl)
}

// AddClient creates a ChatClient for the provided node and inserts it into the
// ClientList. If the node is ourselves, sets our username on the created
// ChatClient.
func (cl ClientList) AddClient(node *smudge.Node) {
	// We need to create a new ChatClient to store in our ClientList. Learn
	// about creating structs here: https://gobyexample.com/structs
	newClient := ChatClient{
		node: node,
	}

	// If the node being added is us (the address matches localAddress) then we
	// should add our username to the ChatClient object (localUsername).

	if NodeAddress(node.Address()) == localAddress {
		newClient.username = localUsername
	}

	cl[NodeAddress(node.Address())] = newClient
}

// RemoveClient deletes a ChatClient from the ClientList if it exists, based on
// the information from the provided node.
func (cl ClientList) RemoveClient(node *smudge.Node) {
	// ClientList is a map from NodeAddress to a ChatClient. We can get the node
	// address from the provided node by calling the "Address()" function on the
	// Node, but that gives us a string. Casting to a NodeAddress is required
	// before looking up in the map.

	if client, exists := cl[NodeAddress(node.Address())]; exists {
		if client.username != "" {
			printDebug("Removing a node: %s", node.Address())
		}
		delete(cl, NodeAddress(node.Address()))
	}
}

// AddUsernames takes a map of NodeAddress->Username pairings and fills the
// ClientList with the usernames provided. It is possible that a node may change
// username, in which case the map should be updated.
func (cl ClientList) AddUsernames(usernames map[NodeAddress]string) error {
	printDebug("Received username list containing: %+v", usernames)

	// loop over the provided map of usernames, updating our client list with
	// the username as we go.
	//
	// range is used to iterate over maps, slices, and arrays.
	// More info: https://tour.golang.org/moretypes/16
	for addr, username := range usernames {
		if client, exists := cl[addr]; exists {
			printDebug("Updating username of %s to %s", addr, username)
			client.username = username

			// When reading from the map, we created a copy of the struct. We
			// now need to put the modified copy back into the map.
			cl[addr] = client
		}
	}

	// Tell the UI the client list has changed and should be redrawn
	printClientList(cl)
	return nil
}

// getUsernameMap returns a map from node addresses to username,
// including only clients for which we know the username. Also include ourselves
// with the localAddress and localUsername.
func (cl ClientList) getUsernameMap() map[NodeAddress]string {
	// Learn more about creating an empty map: https://gobyexample.com/maps
	// Learn more about iterating maps: https://gobyexample.com/range

	usernames := make(map[NodeAddress]string)
	for addr, client := range cl {
		// No sense telling about the usernames that we don't know yet!
		if client.username != "" {
			usernames[addr] = client.username
		}
	}

	// Don't forget to add ourselves!
	usernames[localAddress] = localUsername

	return usernames
}

// BroadcastUsernames builds a map of the known usernames and broadcasts them
// to the chat cluster.
func (cl ClientList) BroadcastUsernames() error {
	printDebug("Processing request to broadcast our known usernames...")

	usernames := cl.getUsernameMap()
	msg := message{
		Type:      messageTypeUsernames,
		Usernames: usernames,
	}
	return smudge.BroadcastBytes(msg.Encode())
}

// FillMissingInfo looks for any connected clients for which we do not already
// know the username. If any missing usernames are found, request a username
// list from the first client found which does not have a username.
func (cl ClientList) FillMissingInfo() {
	c := time.Tick(15 * time.Second)
	for _ = range c {
		printDebug("Checking for clients with a missing username...")

		if addrMissing, ok := cl.GetMissingUsername(); ok {
			if err := cl.RequestUsernameList(addrMissing); err != nil {
				printError("Error requesting missing usernames: %s", err)
			}
		}
	}
}

// GetMissingUsername iterates through the client list, looking for an connected
// clients for which we do not yet have the username. Returns the address of the
// first client encountered which is missing the username.
// If all usernames are known, an empty address and false are returned.
func (cl ClientList) GetMissingUsername() (NodeAddress, bool) {
	for addr, client := range cl {
		if client.username == "" {
			printDebug("Found client with missing username: %s", addr)
			return addr, true
		}
	}
	return NodeAddress(""), false
}

// RequestUsernameList sends a broadcast to all nodes, requesting that the
// specified node respond with a list of all the usernames it is aware of.
//
// A broadcast is used because we have no way of directly connecting to this
// node. Other nodes will just have to ignore this message.
func (cl ClientList) RequestUsernameList(addrMissing NodeAddress) error {
	printDebug("Sending username request to %s", addrMissing)
	msg := message{
		Type: messageTypeUsernameReq,
		Body: string(addrMissing),
	}

	return smudge.BroadcastBytes(msg.Encode())
}

// ChatClient is a structure containing a reference to the smudge.Node
// represented and any additional information we know about this client, such as
// their username.
type ChatClient struct {
	node *smudge.Node

	// username is a value we will query the client for when first discovered
	username string
}

// GetName returns the username of the connected client if the username is
// known, otherwise returns the address used by smudge to connect.
func (c *ChatClient) GetName() string {
	if c.username != "" {
		return c.username
	}
	return c.node.Address()
}
