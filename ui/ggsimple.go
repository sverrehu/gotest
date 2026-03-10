package main

import (
	"log"

	"github.com/gogpu/gogpu"
	"github.com/gogpu/gogpu/gmath"
)

func main() {
	app := gogpu.NewApp(gogpu.DefaultConfig().WithTitle("hello world").WithSize(800, 600))

	app.OnDraw(func(ctx *gogpu.Context) {
		if err := ctx.DrawTriangleColor(gmath.DarkGray); err != nil {
			log.Fatal(err.Error())
		}
	})

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
