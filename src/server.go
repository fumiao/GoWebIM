package main

import (
	"fmt"
	_ "fmt"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/yiqguo/GoWebIM/src/chat"
	"text/template"
)

func main() {

	t := &chat.Template{
		Templates: template.Must(template.ParseGlob("template/*.html")),
	}

	e := echo.New()
	e.SetRenderer(t)

	// Middleware
	e.Use(mw.Logger(), mw.Recover())

	// Routes
	e.Get("/", chat.Index)
	e.Get("/index", chat.Index)
	e.WebSocket("/ws", chat.Ws)
	e.Static("/scripts", "scripts")

	// Start server
	fmt.Println("listening ... 1323")
	e.Run(":1323")
}
