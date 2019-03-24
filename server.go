package main

import (
  "fmt"
	"net/http"

	"github.com/labstack/echo"
)

func yallo(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}

func main() {
  fmt.Println("Welcome to the server")

	e := echo.New()

	e.GET("/", yallo)

	e.Logger.Fatal(e.Start("localhost:1323"))
}
