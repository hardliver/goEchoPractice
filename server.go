package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func yallo(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// http://localhost:1323/cats/string?name=arnold&type=fluffy
func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is: %s\nand his type is: %s\n", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to let us know if you want json or string data",
	})
}

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s", err)
		return c.JSON(http.StatusBadRequest, "")
	}

	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s", err)
		return c.JSON(http.StatusBadRequest, "")
	}

	log.Printf("this is your cat: %#v", cat)
	return c.String(http.StatusOK, "we got your cat!")
}

func addDog(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed processing addDogs request: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	log.Printf("this is your dog: %#v", dog)
	return c.String(http.StatusOK, "we got your dog!")
}

func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Failed processing addHamster request: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	log.Printf("this is your hamster: %#v", hamster)
	return c.String(http.StatusOK, "we got your hamster!")
}

func mainAdmin(c echo.Context) error {
  return c.String(http.StatusOK, "horay you are on the secret admin main page!")
}

// middleware
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    c.Response().Header().Set(echo.HeaderServer, "BlueBot/1.0")
    c.Response().Header().Set("notReallyHeader", "thisHaveNoMeaning")

    return next(c)
  }
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

  e.Use(ServerHeader)

  g := e.Group("/admin")

  // this logs the server interreaction
  g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
  }))

  g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
    // check in the DB
    if username == "jack" && password == "1234" {
      return true, nil
    }
    return false, nil
  }))

  g.GET("/main", mainAdmin)

	e.GET("/", yallo)
	e.GET("/cats/:data", getCats)

	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamsters", addHamster)

	e.Logger.Fatal(e.Start("localhost:1323"))
}
