# Beginning Go - A Chat Client

Welcome! This repository contains the template for
a [gossip](https://en.wikipedia.org/wiki/Gossip_protocol)-based chat client
which is designed to be an introduction project for programmers looking to learn
more about the [Go programming language](https://golang.org/). This file will
guide you from getting started all the way to connecting to other chat clients.

## What Makes this Chat Client Special

Other than you writing it?

By using the gossip protocol, this chat system is 100% decentralized. As long as
one client remains active, the chat history will be maintained. To join the chat
room a client needs only enter the connection information for any other
connected client.

Once connected, the gossip library we are using will tell other clients that you
have connected. In addition to now showing up in their list of connected users,
your client will now be able to periodically poll for new messages, receive
pushed messages, and push its own messages.

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
compatible chat clients, a template client has been created for you to complete.
You can obtain it either by cloning the git repository, or just by downloading
a zip of the files.

> [Download the
> files](https://github.com/tgrosinger/Beginning-Go-Project/archive/master.zip)

## Building and Running

Let's see where the code has gotten us started. After taking a quick look
through `main.go` (or before if you want), run the following set of commands to
build and execute the beginnings of your chat client implementation.

```cmd
Fill out instructions here
```

## Run the Unit Tests


## Fill in the Blanks


## Connect to Others

As your chat client progresses you will eventually be ready to test connecting
to an actual other client! If you are the first one done, come find me and
I will give you connection information for my running client. Otherwise, find
a group that has connected to the chat room and get their connection
information.

