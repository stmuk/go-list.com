package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nsf/termbox-go"
)

func printRange(inFile *os.File, start int, finish int, width int) {
	count := 1
	_, err := inFile.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if count >= start {
			tbprint(0, count-start, termbox.ColorWhite, termbox.ColorBlack, scanner.Text())
			log.Printf("pr: %v %v ", count, count-start)

		}
		count++
		if count == finish+start {
			tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf(fs(width)+"", inFile.Name()))
			break
		}
	}
}

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
	return fmt.Sprintf("%%-%vv", width)
}

func main() {
	var err error
	f, _ := os.OpenFile("list.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.SetOutput(f)

	line := 1 // default reverse video line

	err = termbox.Init()
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

	err = termbox.Flush()

	if err != nil {
		log.Fatal(err)
	}
	currLine := 1
	fileDisplay := 0
	var inFile *os.File
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {

			case termbox.KeyEsc:
				if fileDisplay == 0 {
					break mainloop
				}
				fileDisplay = 0
				main()

			case termbox.KeyCtrlM:
				var err error
				fileDisplay++
				termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
				err = tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, files[line-2].Name())
				if err != nil {
					log.Fatal(err)
				}
				inFile, _ = os.Open(files[line-2].Name())

				defer func() {
					err = f.Close()
					if err != nil {
						log.Fatal(err)
					}
				}()
				printRange(inFile, 0, height, width)

				termbox.Flush()
				continue mainloop

			case termbox.KeyArrowUp:
				if fileDisplay == 0 {
					if line != 1 {
						line--
					}
				} else if currLine != 0 {
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					currLine--
					printRange(inFile, currLine, height, width)
				}
				termbox.Flush()
				continue mainloop

			case termbox.KeyArrowDown:
				if fileDisplay == 0 {
					if line != i-1 {
						line++
					}
					_ = redraw(line, files)
				} else if currLine != height {
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					currLine++
					log.Printf("down: %v", currLine)
					printRange(inFile, currLine, height, width)
				}
				termbox.Flush()
				continue mainloop

			default:
				continue mainloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
