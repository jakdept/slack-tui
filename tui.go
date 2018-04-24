package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// All available layout types
const (
	LayoutChannelsLeft = 1 << iota
	LayoutChannelsRight
	LayoutChannelsHidden
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		// probably need something cleaner
		log.Fatalln(err)
	}
	defer g.Close()

	maxX, maxY := g.Size()

	mainWindow := mainWindowType{
		layout: LayoutChannelsLeft,
		pos: viewSize{
			l: 30,
			r: maxX,
			t: -1,
			b: maxY,
		},
	}

	sideWindow := channelListType{
		layout: LayoutChannelsLeft,
		pos: viewSize{
			l: -1,
			r: 30,
			t: -1,
			b: maxY,
		},
	}

	g.SetManager(mainWindow, sideWindow)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// viewSize represents the dimensions of a view - for easier handling.
// l 0> left, r -> right, t -> top, b -> bottom
// top left point is (l,t), bottom right corner is (r, b)
type viewSize struct {
	l int
	r int
	t int
	b int
}

type channel struct {
	name   string
	unread bool
	order  int
}

type channelListType struct {
	activeChannel channel
	allChannels   channel
	layout        int
	pos           viewSize
}

func (cl channelListType) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("list", cl.pos.l, cl.pos.t, cl.pos.r, cl.pos.b); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprintln(v, "Item 3")
		fmt.Fprint(v, "\rWill be")
		fmt.Fprint(v, "deleted\rItem 4\nItem 5")
	}
	return nil
}

type mainWindowType struct {
	pos    viewSize
	layout int
}

func (mw mainWindowType) Layout(g *gocui.Gui) error {
	if mainBox, err := g.SetView("main", mw.pos.l, mw.pos.t, mw.pos.r, mw.pos.b-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(mainBox, "jakdept: here is a message")
		fmt.Fprintln(mainBox, "jakdept: how about another message")
		fmt.Fprintln(mainBox, "jakdept: what about a \nmultiline\nmessage")
		mainBox.Editable = false
		mainBox.Wrap = true
	}
	if msgBox, err := g.SetView("msgBox", mw.pos.l, mw.pos.b-3, mw.pos.r, mw.pos.b); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		msgBox.Editable = false
		msgBox.Wrap = true
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}
