package main
import ( 
    "fmt"
    "io/ioutil"
    "github.com/nsf/termbox-go"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func main() {

    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    files, _ := ioutil.ReadDir("./")

    tbprint(0,1, termbox.ColorDefault, termbox.ColorDefault,  "\u2191..")

    y :=2;
    for _, f := range files {
        tbprint(0,y, termbox.ColorDefault, termbox.ColorDefault, "\u2193")
        tbprint(1,y, termbox.ColorDefault, termbox.ColorDefault, (f.Name()))
        y++
    }

    len := y-2

    tbprint(0,0, termbox.ColorBlack, termbox.ColorWhite,  fmt.Sprintf("LIST File Selection 1 of %d",len));

    tbprint(0,y, termbox.ColorBlack, termbox.ColorWhite,  fmt.Sprintf("Files: %d \u2666",len));

    termbox.Flush()

    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            default:
                panic("exit")
            }
        case termbox.EventError:
            panic(ev.Err)
        }
    }
}
