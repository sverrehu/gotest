package console

import (
	"log"
	"os"
	"strconv"
)

func output(s string) {
	_, err := os.Stdout.WriteString(s)
	if err != nil {
		log.Fatal("error writing to stdout")
	}
}

func MoveToXY(x, y int) {
	output("\u001b[" + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

func GetSize() (int, int) {
	output("\u001b[s" + // save cursor position
		"\u001b[5000;5000H" + // move to col 5000 row 5000
		"\u001b[6n" + // request cursor position
		"\u001b[u") // restore cursor position
	// on stdin: \u001b[25;80R"
	return -1, -1
	// check golang.org/x/term
}
