package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/clockworksoul/smudge"
)

const (
	// variables declared within "const" are constants in Go. The type is
	// determined by the compiler.
	// More info: https://blog.golang.org/constants
	//            https://gobyexample.com/constants

	// heartbeatMillis is used to configure how frequently the gossip protocol
	// announces that it is still connected. No need to change this value.
	heartbeatMillis = 500
)

// TODO: Add unit and benchmark tests which fail in the not-complete
// implementation. Then students can still be given something which compiles, it
// just doesn't work.

var (
	// variables declared within "var" are mutable in Go. They can be explicitly
	// initialized to a value, or if not set explicitly, default to the "empty
	// value" for their type.
	// More info: https://golang.org/doc/effective_go.html#variables

	// otherClient specifies the address of one running instance of the client.
	// If omitted, this client will not initiate a connection to any existing
	// client (i.e. this is the first client in a cluster)
	otherClient = flag.String("client", "",
		"Address of an existing client, if empty do not attempt to connect")

	// listenPort is where this client will listen for other clients connecting.
	// Must not be left empty.
	listenPort = flag.Int("listenport", 0,
		"Port on which client listens for connections to other clients")

	// username is the friendly name we will present to other clients instead of
	// our address. Must not be left empty.
	username = flag.String("username", "", "Friendly name for this client")

	// localAddress is the NodeAddress which other Clients will use to reach us.
	localAddress NodeAddress
)

// printDebug outputs a log message with the "DEBUG:" prefix. This function can
// be edited to easily enable and disable debugging logs without removing all
// the log lines in the codebase.
func printDebug(msg string, args ...interface{}) {
	printLogs(fmt.Sprintf("DEBUG: "+msg, args...))
}

// printInfo outputs a log message with the "INFO:" prefix. This function can
// be edited to easily enable and disable debugging logs without removing all
// the log lines in the codebase.
func printInfo(msg string, args ...interface{}) {
	printLogs(fmt.Sprintf("INFO: "+msg, args...))
}

// printError outputs a log message with the "ERROR:" prefix. This function can
// be edited to easily enable and disable error logs without removing all
// the log lines in the codebase.
func printError(msg string, args ...interface{}) {
	printLogs(fmt.Sprintf("ERROR: "+msg, args...))
}

// cacheLocalIP populates the value of the localAddress global variable.
// localAddress is used to determine if a broadcast was directed to us
// specifically, as it is the address which other clients use to communicate
// with us.
func cacheLocalIP() {
	// this pattern of returning a result and an error is extremely prevalent in
	// Go. Unlike many languages, exceptions (or in Go, Panics) are very rarely
	// used. When a function returns an error, it Must be handled and the result
	// disregarded.
	// More info: https://blog.golang.org/error-handling-and-go
	ip, err := smudge.GetLocalIP()
	if err != nil {
		fmt.Println("Unable to retrieve local IP", err)
		os.Exit(1)
	}

	localIP := ip.String()

	// listenPort, defined above, is a pointer to a number. Take a look at the
	// return type of https://golang.org/pkg/flag/#Int
	// Prepending our use of listenPort with a * will dereference the pointer,
	// giving us a normal int.
	localAddress = NodeAddress(fmt.Sprintf("%s:%d", localIP, *listenPort))
}

// main is the entry point to the application.
func main() {
	// Populate the flag variables at the top of this file with input from the
	// user. Afterwards, determine if any required values were omitted.
	flag.Parse()

	if *listenPort == 0 {
		printError("Listen port is required")
		flag.Usage()
		os.Exit(1)
	} else if *username == "" {
		printError("Username is required")
		flag.Usage()
		os.Exit(1)
	}

	// Now the user input is parsed, lets start configuring the gossip
	// communication with other clients. These options were all grabbed from the
	// example on the project homepage: https://github.com/clockworksoul/smudge#everything-in-one-place

	// Set configuration options
	smudge.SetListenPort(*listenPort)
	smudge.SetHeartbeatMillis(heartbeatMillis)

	// Add the status listener
	clientList := ClientList(make(map[NodeAddress]ChatClient))
	smudge.AddStatusListener(clientList)

	// Add the broadcast listener
	messenger := Messenger{clients: clientList}
	smudge.AddBroadcastListener(&messenger)

	// Only attempt to connect to another client if the address for one was
	// provided. If not, the client will sit and wait until a client connects.
	if *otherClient != "" {
		// Add a new remote node. To join an existing cluster you must
		// add at least one of its healthy member nodes.
		if node, err := smudge.CreateNodeByAddress(*otherClient); err != nil {
			printError("Failed to create a new node from addr: ", err)
			os.Exit(1)
		} else {
			smudge.AddNode(node)
		}
	}

	// The default logs from smudge just print to stdout and look messy in our
	// fancy UI.
	smudge.SetLogThreshold(smudge.LogOff)

	// Start the server!
	// We will run the smudge server in a background go routine. This is similar
	// to a new thread, however it is scheduled on real OS threads by the Go
	// runtime.
	//
	// For the scope of this class, you can assume this function is
	// running in the background. I encourage reading more about these later
	// from a resource such as this: https://gobyexample.com/goroutines
	printDebug("Starting Smudge...\n")
	go smudge.Begin()

	cacheLocalIP()

	// Start the username watcher!
	// Another go routine, both will be scheduled by the runtime and run as
	// frequently as possible, depending on the number of threads given to the
	// process.
	go clientList.FillMissingInfo()

	// Start the gui!
	// Notice that here we are not starting in a go routine. If we did then this
	// thread (the main one) would reach the end of the main function, exit, and
	// kill all the other go routines. We will hand-off control of the program
	// to the UI which will listen for input from the user from here out.
	printDebug("Starting the GUI...\n")
	runGUI(clientList)
}
