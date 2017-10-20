package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

//
// This file does not need to be edited!
//
// Please feel free to dig through this file if you are curious, however the
// contents are fully implemented, so no edits are required to arrive at a
// functional chat client.
//

var (
	gui *gocui.Gui

	logsVisible = false
)

func runGUI(cl ClientList) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println("Fatal GUI error: ", err)
		os.Exit(1)
	}
	defer gui.Close()

	gui = g

	// Set GUI managers and key bindings

	gui.Cursor = true
	gui.SetManagerFunc(layout)

	err = gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		fmt.Println("Fatal GUI error: ", err)
		os.Exit(1)
	}
	err = gui.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, toggleLogs)
	if err != nil {
		fmt.Println("Fatal GUI error: ", err)
		os.Exit(1)
	}
	err = gui.SetKeybinding("enter-text", gocui.KeyEnter, gocui.ModNone, readGuiMsg)
	if err != nil {
		fmt.Println("Fatal GUI error: ", err)
		os.Exit(1)
	}

	// We will update the client list after the GUI is initialized because we
	// need to print the name of the initial client we connected to when
	// creating Smudge.
	// If this is skipped, we will not see the initial node connected until
	// another node is added or removed.
	printClientList(cl)

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println("Fatal GUI error: ", err)
		os.Exit(1)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	helpY := maxY - 2
	chatX := 25
	if maxX < 25 {
		// Support for small terminals
		chatX = 15
	}

	chatMaxY := maxY - 6

	if v, err := g.SetView("logs", 3, 2, maxX-3, chatMaxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Logs"
		v.Autoscroll = true
		v.Wrap = true
		g.SetViewOnBottom("logs")
	}

	if v, err := g.SetView("help", 0, helpY, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false

		fmt.Fprintf(v, "%s %s    %s %s    %s %s",
			frameText("Ctrl-L"), "Toggle Logs",
			frameText("Ctrl-C"), "Quit",
			frameText("Enter"), "Send Message")
	}

	if v, err := g.SetView("clients", 0, 0, chatX-1, helpY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Clients"
	}

	if v, err := g.SetView("messages", chatX, 0, maxX-1, chatMaxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Autoscroll = true
		v.Wrap = true
		v.Title = "Message-History"
	}

	if v, err := g.SetView("enter-text", chatX, chatMaxY+1, maxX-1, helpY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("enter-text"); err != nil {
			return err
		}

		v.Title = "Send:"
		v.Editable = true
		v.Wrap = true
	}
	return nil
}

func readGuiMsg(g *gocui.Gui, v *gocui.View) error {
	msgText := v.Buffer()
	v.Clear()

	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	SendMessage(msgText)
	return nil
}

func printChatMessage(msg, sender string) {
	gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("messages")
		if err != nil {
			return err
		}

		fmt.Fprintf(v, "%s: %s\n", sender, msg)
		return nil
	})
}

// NOTE TO SELF: CHANGE THIS FROM CLIENT LIST TO SOMETHING THAT IS AN INTERFACE
// SO I AM NOT DICTATING THE STRUCTURE OF THEIR PROGRAM
func printClientList(cl ClientList) {
	if gui == nil {
		return
	}

	gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("clients")
		if err != nil {
			return err
		}

		v.Clear()
		v.SetCursor(0, 0)

		for _, client := range cl {
			fmt.Fprintln(v, client.GetName())
		}
		return nil
	})
}

func printLogs(msg string) {
	if gui == nil {
		fmt.Println(msg)
	} else {
		gui.Update(func(g *gocui.Gui) error {
			v, err := g.View("logs")
			if err != nil {
				return err
			}

			fmt.Fprintln(v, msg)
			return nil
		})
	}
}

func toggleLogs(g *gocui.Gui, v *gocui.View) error {
	if logsVisible {
		_, err := g.SetViewOnBottom("logs")
		if err != nil {
			return err
		}
	} else {
		_, err := g.SetViewOnTop("logs")
		if err != nil {
			return err
		}
	}

	logsVisible = !logsVisible
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func stringFormatBoth(fg, bg int, str string, args []string) string {
	return fmt.Sprintf("\x1b[48;5;%dm\x1b[38;5;%d;%sm%s\x1b[0m", bg, fg, strings.Join(args, ";"), str)
}

// Frame text with colors
func frameText(text string) string {
	return stringFormatBoth(15, 0, text, []string{"1"})
}
