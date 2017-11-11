# Beginning Go - A Chat Client

Welcome! This repository contains the framework for an implementation of
a [gossip](https://en.wikipedia.org/wiki/Gossip_protocol)-based chat client,
designed to be an introductory project for programmers looking to learn
more about the [Go programming language](https://golang.org/). This file will
guide you through the process: from getting started to connecting your
implementation to others.

## What Makes this Chat Client Special

Other than you writing it?

By using the gossip protocol, this chat system is 100% decentralized. As long as
one client remains active, the chat cluster is maintained. To join the chat room
a client needs only enter the connection address of any other connected client.

Once connected, the gossip library we are using will tell other clients that you
have connected. In addition to now showing up in their list of connected users,
your client will receive pushed messages, and push its messages to everyone else
in the cluster.

Currently the [gossip library](https://github.com/clockworksoul/smudge) we are
using is limited to LAN connections only.

## Demo

[![asciicast](https://asciinema.org/a/G4YYRdotQIDQtb66n2lU1aQle.png)](https://asciinema.org/a/G4YYRdotQIDQtb66n2lU1aQle)

# Creating the Chat Client

## Assumptions

If you are following along with this exercise during the class at UW, you are
likely working on a Windows computer with Go already installed. Go is
a statically typed and compiled language that can create executables for many
different operating systems and architectures. This means code written for one
operating system will work just as well on another. You can even cross-compile
(create Windows executables from Linux, for example).

If you do not already have Go installed, please follow the [installation
instructions](https://golang.org/doc/install). 

## Obtaining the Source

To fit within the time limits of the class, and to ensure everyone writes
compatible chat clients, a template client has been created which implements the
UI updates and the data structures. Your job will be to fill in the missing
pieces, make sure the unit tests pass, and test with others. You can obtain it
either by cloning the git repository, or just by downloading a zip of the files.

> [Download the
> files](https://github.com/tgrosinger/Beginning-Go-Project/archive/master.zip)

## Building and Running

### A Windows Caveat

Go is fairly opinionated about how the source code is structured. The computers
should already have Go installed, but when Go looks for code to build it will
check in the `%USERPROFILE%\go` directory. To get everything working nicely,
open a new File Explorer and in the address bar type `%USERPROFILE%` and press
enter. Next, create a series of nested directories so you end up with this exact
path (don't replace the name):

```
%USERPROFILE%\go\src\tgrosinger\beginning-go
```

Finally, put the source code you downloaded in that final folder after
unzipping. You should end up with a file at the path:
`%USERPROFILE%\go\src\tgrosinger\beginning-go\main.go`

### Another Windows Caveat

For the following commands to run we need to have our command prompt in the
right directory. The easiest way to do that is to hold shift+ctrl and right
click in the file explorer we navigated to in the previous step. Click "Open
command window here"

### Back to Building and Running

Let's see where the code has gotten us started. After taking a quick look
through `main.go` (or before if you want), run the following set of commands to
build and execute the beginnings of your chat client implementation.

```cmd
go build
```

That's it. Go will drop a binary in the root of the project called
`beginning-go`. If you run it with only the flag `-h` it will provide you with
info about the other runtime flags.

## Run the Unit Tests

Running the unit tests is the fastest method for checking most of the
functionality you are responsible for implementing. When running the tests, the
output will show both the result from running, and the expected result.
I promise that the expected result is correct.

```
go test -v
```

Removing the `-v` will cause only test failures to be displayed.

If your program does not compile it will show those errors and not run the
tests.

## Editor Support

I highly recommend taking a minute to install the Go plugin for your favorite
editor. Both Sublime Text
([gosublime]([200~https://github.com/DisposaBoy/GoSublime])) and VS Code (Just
search for the Go extension) have great plugins for editing Go code.

## Fill in the Blanks

Read through the code, run the unit tests, implement the missing parts. You
should skim through `main.go` and `message.go` but you only need to write code
in `client.go`. The related tests are located in `client_test.go`.

Remember to work with your partner, and ask another group if you get stuck.

## Connect to Others

As your chat client progresses you will eventually be ready to test connecting
to an actual other client! If you are the first one done, come find me and
I will give you connection information for my running client. Otherwise, find
a group that has connected to the chat room and get their connection
information.

