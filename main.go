package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zodimo/go-gio-ark-todo/systems"

	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	go func() {
		w := &app.Window{}
		todoApp := systems.NewTodoApp(w)
		settings := todoApp.GetSettings()
		w.Option(app.Title(settings.Title))
		w.Option(app.Size(unit.Dp(settings.ScreenWidth), unit.Dp(settings.ScreenHeight)))

		if err := todoApp.Run(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Exiting application")
		os.Exit(0)

	}()
	app.Main()
}
