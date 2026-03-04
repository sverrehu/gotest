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
