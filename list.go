package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nsf/termbox-go"
)

var f *os.File

func init() {
	var err error
	if os.Getenv("LOG") == "1" {

		f, err = os.OpenFile("list.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(f)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	err = termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	defer termbox.Close()

	currLine := 1
	var fileName string

	for {
		fileName, currLine = fileSelect(currLine)

	dir:
		for {
			fi, err := os.Lstat(fileName)
			if err != nil {
				log.Fatal(err)
			}

			if fi.Mode().IsDir() {

				os.Chdir(fileName)
				fileName, currLine = fileSelect(1)

			} else {
				break dir
			}
		}

		displayFile(fileName)
	}
}

func fileSelect(currLine int) (string, int) {

	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	width, height := termbox.Size()

	files, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	i := listFiles(currLine, files)

	len := i - 2

	// note poor support for symlinks XXX
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	s0 := fmt.Sprintf("LIST File Selection 1 of %d ", len)
	s1 := fmt.Sprintf(" Path "+padSpace(width), cwd)

	// first line
	tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, s0+s1)

	// last line
	tbprint(0, height-1, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf("Files: "+padSpace(width)+"\u2666", len)) // diamond

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
				termbox.Close() // XXX
				os.Exit(1)

				// enter file
			case termbox.KeyCtrlM:
				break fileselect

			case termbox.KeyArrowUp:
				if currLine != 1 {
					currLine--
				}
				_ = listFiles(currLine, files)
				termbox.Flush()
				continue fileselect

			case termbox.KeyArrowDown:
				if currLine != (i - 1) { // XXX
					currLine++
				}
				_ = listFiles(currLine, files)
				termbox.Flush()
				continue fileselect

			default:
				termbox.Close()
				if ev.Ch != 0 && (string(ev.Ch) == "q" || string(ev.Ch) == "x") {
					os.Exit(1)
				}

				continue fileselect
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
	var fileName string
	if currLine > 1 {
		fileName = files[currLine-2].Name()
	} else {
		fileName = ".."
	}
	return fileName, currLine
}

func displayFile(fileName string) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

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
				if (currLine - (height - 1)) > 1 { // XXX
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					currLine -= (height - 1)
				} else {
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					currLine = 0
				}

				printRange(inFile, currLine, height, width)
				termbox.Flush()
				continue filedisplay

			case termbox.KeyArrowDown:
				if currLine < height { //XXX len of file
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
					break filedisplay
				}

				continue filedisplay
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

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
			tbprint(0, 0, termbox.ColorBlack, termbox.ColorWhite, fmt.Sprintf(padSpace(width)+"", inFile.Name()))
			break
		}
	}
}

func listFiles(currLine int, files []os.FileInfo) int {

	if currLine == 1 {
		tbprint(0, 1, termbox.ColorBlack, termbox.ColorWhite, "\u2191..") // up arrow
	} else {
		tbprint(0, 1, termbox.ColorDefault, termbox.ColorDefault, "\u2191..") // up arrow
	}

	i := 2
	for _, f := range files {

		if currLine == i {
			if f.IsDir() {
				tbprint(0, i, termbox.ColorBlack, termbox.ColorWhite, "\u2193") // down arrow
			} else {
				tbprint(0, i, termbox.ColorBlack, termbox.ColorWhite, " ")
			}
			tbprint(1, i, termbox.ColorBlack, termbox.ColorWhite, (f.Name()))
		} else {
			if f.IsDir() {
				tbprint(0, i, termbox.ColorDefault, termbox.ColorDefault, "\u2193")
			} else {
				tbprint(0, i, termbox.ColorDefault, termbox.ColorDefault, " ")
			}

			tbprint(1, i, termbox.ColorDefault, termbox.ColorDefault, (f.Name()))
		}
		i++

	}
	return i
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func padSpace(width int) string {
	return fmt.Sprintf("%%-%vv", width)
}
