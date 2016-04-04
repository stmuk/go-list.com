package main

import (
	"bufio"
	"fmt"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func redraw(line int, files []os.FileInfo) int {

	if line == 1 {
		tbprint(0, 1, termbox.ColorBlack, termbox.ColorWhite, "\u2191..") // up arrow
	} else {
		tbprint(0, 1, termbox.ColorDefault, termbox.ColorDefault, "\u2191..") // up arrow
	}

	i := 2
	for _, f := range files {

		if line == i {
			tbprint(0, i, termbox.ColorBlack, termbox.ColorWhite, "\u2193") // down arrow
			tbprint(1, i, termbox.ColorBlack, termbox.ColorWhite, (f.Name()))
		} else {
			tbprint(0, i, termbox.ColorDefault, termbox.ColorDefault, "\u2193")
			tbprint(1, i, termbox.ColorDefault, termbox.ColorDefault, (f.Name()))
		}
		i++

	}
	return i
}

func fs(width int) string {
	return fmt.Sprintf("%%-%dd", width)
}

func main() {

	line := 1 // default reverse video line

	err := termbox.Init()
	width, height := termbox.Size()

	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	files, _ := ioutil.ReadDir("./")

	i := redraw(line, files)

	len := i - 2

	tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("LIST File Selection 1 of "+fs(width), len))

	tbprint(0, height-1, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("Files: "+fs(width/2)+"\u2666", len)) // diamond

	termbox.Flush()

	lcount := 1

	file_display := 0

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {

			case termbox.KeyEsc:
				if file_display == 0 {
					break mainloop
				}
				file_display = 0
				main()

			case termbox.KeyCtrlM:
				file_display++
				termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
				tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, files[line-2].Name())
				inFile, _ := os.Open(files[line-2].Name())
				defer inFile.Close()
				scanner := bufio.NewScanner(inFile)
				scanner.Split(bufio.ScanLines)

				for scanner.Scan() {
					tbprint(0, lcount, termbox.ColorWhite, termbox.ColorBlack, scanner.Text())
					lcount++
					if lcount == height {
						break
					}
				}

				termbox.Flush()
				continue mainloop

			case termbox.KeyArrowUp:
				if file_display == 0 {
					if line != 1 {
						line--
					}
					_ = redraw(line, files)
					termbox.Flush()
				}

			case termbox.KeyArrowDown:
				if file_display == 0 {
					if line != i-1 {
						line++
					}
					_ = redraw(line, files)
					termbox.Flush()
				}

			default:
				continue mainloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
