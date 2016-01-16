package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io/ioutil"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func main() {

	line := 2

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	files, _ := ioutil.ReadDir("./")

	tbprint(0, 1, termbox.ColorDefault, termbox.ColorDefault, "\u2191..")

	y := 2
	for _, f := range files {

		if line == y {
			tbprint(0, y, termbox.ColorBlack, termbox.ColorWhite, "\u2193")
			tbprint(1, y, termbox.ColorBlack, termbox.ColorWhite, (f.Name()))
		} else {
			tbprint(0, y, termbox.ColorDefault, termbox.ColorDefault, "\u2193")
			tbprint(1, y, termbox.ColorDefault, termbox.ColorDefault, (f.Name()))
		}
		y++

	}

	len := y - 2

	tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("LIST File Selection 1 of %d", len))

	tbprint(0, y, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("Files: %d \u2666", len))

	termbox.Flush()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop

			case termbox.KeyArrowUp:
				break mainloop

			case termbox.KeyArrowDown:
				break mainloop

			default:
				break mainloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
