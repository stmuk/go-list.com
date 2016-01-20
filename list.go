package main

import (
	"fmt"
	"os"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"bufio"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func redraw(line int, files []os.FileInfo) int {
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

func main() {

	line := 2 // default reverse video line

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	files, _ := ioutil.ReadDir("./")

    // line 0 below since len needed

	tbprint(0, 1, termbox.ColorDefault, termbox.ColorDefault, "\u2191..") // up arrow

    i := redraw(line, files)

	len := i - 2

	tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("LIST File Selection 1 of %d", len))

	tbprint(0, i, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("Files: %d \u2666", len)) // diamond

	termbox.Flush()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop

            case termbox.KeyCtrlM:
                termbox.Clear(termbox.ColorWhite,  termbox.ColorBlack)
                tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite,  files[line-2].Name())
				inFile, _ := os.Open(files[line-2].Name())
				defer inFile.Close()
				scanner := bufio.NewScanner(inFile)
				scanner.Split(bufio.ScanLines) 
				i:=1
				for scanner.Scan() {
					tbprint(0, i, termbox.ColorWhite, termbox.ColorBlack, scanner.Text())
					i++
				}
                
                termbox.Flush()
                continue mainloop

            case termbox.KeyArrowUp:
                if line !=0  {
                    line--
                    _ = redraw(line, files)
                    termbox.Flush()
                }

			case termbox.KeyArrowDown:
                line++
                _ = redraw(line, files)
	            termbox.Flush()

			default:
				continue mainloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
