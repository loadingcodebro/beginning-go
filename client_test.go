package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/clockworksoul/smudge"
)

func TestMain(m *testing.M) {
	unittestMode = true
	os.Exit(m.Run())
}

func CheckNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected nil error, received: %s", err)
	}
}

func TestGetName(t *testing.T) {
	testNode, err := smudge.CreateNodeByIP(net.ParseIP("127.0.0.1"), 9999)
	CheckNoError(t, err)

	var cases = []struct {
		client         ChatClient
		expectedResult string
	}{
		{ // When the username is set, it is returned
			client: ChatClient{
				username: "testing",
				node:     testNode,
			},
			expectedResult: "testing",
		},
		{ // When the username is not set, the node address is returned
			client: ChatClient{
				node: testNode,
			},
			expectedResult: "127.0.0.1:9999",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			result := c.client.GetName()
			if result != c.expectedResult {
				t.Fatalf("Expected %q but got %q", c.expectedResult, result)
			}
		})
	}
}

func TestRemoveClient(t *testing.T) {
	testNode, err := smudge.CreateNodeByIP(net.ParseIP("127.0.0.1"), 9999)
	CheckNoError(t, err)

	var cases = []struct {
		clientList     *ClientList
		expectedResult *ClientList
	}{
		{ // Test that a client is removed
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "testing",
				},
				NodeAddress("127.0.0.2:9998"): ChatClient{
					username: "testing2",
				},
			},
			expectedResult: &ClientList{
				NodeAddress("127.0.0.2:9998"): ChatClient{
					username: "testing2",
				},
			},
		},
		{ // Test that clients that doesn't exist is not removed
			clientList: &ClientList{
				NodeAddress("127.0.0.2:9998"): ChatClient{
					username: "testing2",
				},
				NodeAddress("127.0.0.1:9997"): ChatClient{
					username: "testing",
				},
			},
			expectedResult: &ClientList{
				NodeAddress("127.0.0.2:9998"): ChatClient{
					username: "testing2",
				},
				NodeAddress("127.0.0.1:9997"): ChatClient{
					username: "testing",
				},
			},
		},
		{ // Test that it still works if there are no clients connected
			clientList:     &ClientList{},
			expectedResult: &ClientList{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			c.clientList.RemoveClient(testNode)
			if !reflect.DeepEqual(*c.clientList, *c.expectedResult) {
				t.Fatalf("Expected %v but got %v", *c.expectedResult, *c.clientList)
			}
		})
	}
}

func TestAddClient(t *testing.T) {
	localAddress = "192.168.0.101:8888"
	*username = "unittest"

	testNode, err := smudge.CreateNodeByIP(net.ParseIP("192.168.0.10"), 9999)
	CheckNoError(t, err)
	testNodeLocal, err := smudge.CreateNodeByIP(net.ParseIP("192.168.0.101"), 8888)
	CheckNoError(t, err)

	var cases = []struct {
		clientList     *ClientList
		expectedResult *ClientList
		nodeToAdd      *smudge.Node
	}{
		{ // Test that client is added
			clientList: &ClientList{
				NodeAddress("192.168.0.5:9998"): ChatClient{
					username: "testing2",
				},
			},
			expectedResult: &ClientList{
				NodeAddress("192.168.0.10:9999"): ChatClient{
					username: "",
					node:     testNode,
				},
				NodeAddress("192.168.0.5:9998"): ChatClient{
					username: "testing2",
				},
			},
			nodeToAdd: testNode,
		},
		{ // Test that if the client is us, the username is added
			clientList: &ClientList{
				NodeAddress("192.168.0.10:9997"): ChatClient{
					username: "testing",
				},
			},
			expectedResult: &ClientList{
				NodeAddress(localAddress): ChatClient{
					username: *username,
					node:     testNodeLocal,
				},
				NodeAddress("192.168.0.10:9997"): ChatClient{
					username: "testing",
				},
			},
			nodeToAdd: testNodeLocal,
		},
		{ // Test that it still works if the client list is empty
			clientList: &ClientList{},
			expectedResult: &ClientList{
				NodeAddress("192.168.0.10:9999"): ChatClient{
					node: testNode,
				},
			},
			nodeToAdd: testNode,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			c.clientList.AddClient(c.nodeToAdd)
			if !reflect.DeepEqual(*c.clientList, *c.expectedResult) {
				t.Fatalf("Expected %v but got %v", *c.expectedResult, *c.clientList)
			}
		})
	}
}

func TestGetUsernameMap(t *testing.T) {
	localAddress = "192.168.0.101:8888"
	*username = "unittest"

	var cases = []struct {
		clientList     *ClientList
		expectedResult map[NodeAddress]string
	}{
		{ // Test that clients with usernames are included
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "testing",
				},
				NodeAddress("127.0.0.2:9998"): ChatClient{
					username: "testing2",
				},
			},
			expectedResult: map[NodeAddress]string{
				NodeAddress("127.0.0.1:9999"): "testing",
				NodeAddress("127.0.0.2:9998"): "testing2",
				NodeAddress(localAddress):     *username,
			},
		},
		{ // Test that clients with no usernames are not included
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "",
				},
			},
			expectedResult: map[NodeAddress]string{
				NodeAddress(localAddress): *username,
			},
		},
		{ // Test that it still works if there are no clients connected
			clientList: &ClientList{},
			expectedResult: map[NodeAddress]string{
				NodeAddress(localAddress): *username,
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			usernameMap := c.clientList.getUsernameMap()
			if !reflect.DeepEqual(usernameMap, c.expectedResult) {
				t.Fatalf("Expected %v but got %v", c.expectedResult, usernameMap)
			}
		})
	}
}
func TestGetMissingUsername(t *testing.T) {
	var cases = []struct {
		clientList         *ClientList
		expectedResultBool bool
		expectedResultAddr NodeAddress
	}{
		{ // Test that false is returned when all usernames are present
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "testing",
				},
				NodeAddress("127.0.0.2:9998"): ChatClient{
					username: "testing2",
				},
			},
			expectedResultBool: false,
			expectedResultAddr: NodeAddress(""),
		},
		{ // Test that a missing username is returned
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{},
			},
			expectedResultBool: true,
			expectedResultAddr: NodeAddress("127.0.0.1:9999"),
		},
		{ // Test that it still works if there are no clients connected
			clientList:         &ClientList{},
			expectedResultBool: false,
			expectedResultAddr: NodeAddress(""),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			resultAddr, resultBool := c.clientList.GetMissingUsername()
			if resultAddr != c.expectedResultAddr {
				t.Fatalf("Expected %v but got %v", c.expectedResultAddr, resultAddr)
			}
			if resultBool != c.expectedResultBool {
				t.Fatalf("Expected %v but got %v", c.expectedResultBool, resultBool)
			}
		})
	}
}

func TestAddUsernames(t *testing.T) {
	cases := []struct {
		clientList     *ClientList
		usernames      map[NodeAddress]string
		expectedResult *ClientList
	}{
		{ // Simple case where the username has an entry we need
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{},
			},
			usernames: map[NodeAddress]string{
				NodeAddress("127.0.0.1:9999"): "tester",
			},
			expectedResult: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "tester",
				},
			},
		},
		{ // Usernames have a new name for a client we know
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "tester",
				},
			},
			usernames: map[NodeAddress]string{
				NodeAddress("127.0.0.1:9999"): "new-tester",
			},
			expectedResult: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "new-tester",
				},
			},
		},
		{ // Usernames have an entry we don't need
			clientList: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "tester",
				},
			},
			usernames: map[NodeAddress]string{
				NodeAddress("127.0.0.1:8888"): "tester",
			},
			expectedResult: &ClientList{
				NodeAddress("127.0.0.1:9999"): ChatClient{
					username: "tester",
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			err := c.clientList.AddUsernames(c.usernames)
			CheckNoError(t, err)

			if !reflect.DeepEqual(*c.clientList, *c.expectedResult) {
				t.Fatalf("Expected %v but got %v", *c.expectedResult, *c.clientList)
			}
		})
	}
}
