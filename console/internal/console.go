package console

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/term"
)

const defaultWidth = 80
const defaultHeight = 25

func output(s string) {
	_, err := os.Stdout.WriteString(s)
	if err != nil {
		log.Fatal("error writing to stdout")
	}
}

func input() byte {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		log.Fatal("error reading from stdin")
	}
	return buf[0]
}

func MoveToXY(x, y int) {
	output("\u001b[" + strconv.Itoa(y+1) + ";" + strconv.Itoa(x+1) + "H")
}

func oldGetSize() (int, int) {
	output("\u001b[s" + // save cursor position
		"\u001b[5000;5000H" + // move to col 5000 row 5000
		"\u001b[6n" + // request cursor position
		"\u001b[u") // restore cursor position
	// on stdin: \u001b[25;80R"
	b := input() // hangs waiting for newline. Need raw mode.
	fmt.Printf("got %s", string(b))
	return -1, -1
	// check golang.org/x/term
}

func GetSizeWH() (int, int) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Printf("error getting size of terminal. will default to %dx%d: %v\n", defaultWidth, defaultHeight, err)
		return defaultWidth, defaultHeight
	}
	return width, height
}

func SetCursorVisible(visible bool) {
	if visible {
		output("\u001b[?25h")
	} else {
		output("\u001b[?25l")
	}
}
