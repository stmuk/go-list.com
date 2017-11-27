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

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	i := redraw(line, files)

	len := i - 2

	// first line
	tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("LIST File Selection 1 of "+fs(width), len))

	// last line
	tbprint(0, height-1, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("Files: "+fs(width/2)+"\u2666", len)) // diamond

	err = termbox.Flush()

	if err != nil {
		log.Fatal(err)
	}

fileselect:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {

			case termbox.KeyEsc:
				os.Exit(1)

				// enter file
			case termbox.KeyCtrlM:
				break fileselect

			case termbox.KeyArrowUp:
				if line != 1 {
					line--
				}
				_ = redraw(line, files)
				termbox.Flush()
				continue fileselect

			case termbox.KeyArrowDown:
				if line != i-1 {
					line++
				}
				_ = redraw(line, files)
				termbox.Flush()
				continue fileselect

			default:
				if ev.Ch != 0 && (string(ev.Ch) == "q" || string(ev.Ch) == "x") {
					os.Exit(1)
				}

				continue fileselect
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
	fileName := files[line-2].Name()
	displayFile(fileName)
}

func displayFile(fileName string) {
	width, height := termbox.Size() // XXX
	currLine := 1
	inFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = inFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// enter file
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, fileName)
	printRange(inFile, 0, height, width)

	termbox.Flush()

filedisplay:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {

			case termbox.KeyEsc:
				break filedisplay

			case termbox.KeyArrowUp:
				if currLine != 0 {
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					currLine--
					printRange(inFile, currLine, height, width)
				}
				termbox.Flush()
				continue filedisplay

			case termbox.KeyPgup:
				//if (currLine - height) > 1 { // XXX
				termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
				currLine -= height - 1
				printRange(inFile, currLine, height, width)
				//}
				termbox.Flush()
				continue filedisplay

			case termbox.KeyArrowDown:
				if currLine != height {
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					currLine++
					log.Printf("down: %v", currLine)
					printRange(inFile, currLine, height, width)
				}
				termbox.Flush()
				continue filedisplay

			case termbox.KeyPgdn:
				if currLine != height {
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					currLine += (height - 2)
					log.Printf("down: %v", currLine)
					printRange(inFile, currLine, height, width)
				}
				termbox.Flush()
				continue filedisplay

			default:
				if ev.Ch != 0 && (string(ev.Ch) == "q" || string(ev.Ch) == "x") {
					os.Exit(1)
				}

				continue filedisplay
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
