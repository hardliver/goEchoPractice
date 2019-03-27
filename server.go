package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
  "time"
  "strings"

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

func mainCookie(c echo.Context) error {
  return c.String(http.StatusOK, "you are on the secret cookie page!")
}

func login(c echo.Context) error {
  username := c.QueryParam("username")
  password := c.QueryParam("password")

  // check username and password against DB after hashing the password
  if username == "jack" && password == "1234" {
    cookie := &http.Cookie{}
    // this is the same
    // cookie := new(http.Cookie{)}

    cookie.Name = "sessionID"
    cookie.Value = "some_string"  // cookie content
    cookie.Expires = time.Now().Add(48 * time.Hour)

    c.SetCookie(cookie)

    return c.String(http.StatusOK, "You were logged in!")
  }

  return c.String(http.StatusUnauthorized, "Your username or password were wrong")
}

// middleware
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    c.Response().Header().Set(echo.HeaderServer, "BlueBot/1.0")
    c.Response().Header().Set("notReallyHeader", "thisHaveNoMeaning")

    return next(c)
  }
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    cookie, err := c.Cookie("sessionID")
    if err != nil {
      if strings.Contains(err.Error(), "named cookie not present") {
        return c.String(http.StatusUnauthorized, "You dont have any cookie")
      }
      log.Println(err)
      return err
    }
    if cookie.Value == "some_string" {
      return next(c)
    }

    return c.String(http.StatusUnauthorized, "You dont have the right cookie, cookie")
  }
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

  e.Use(ServerHeader)

  adminGroup := e.Group("/admin")
  cookieGroup := e.Group("/cookie")

  // this logs the server interreaction
  adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
  }))

  adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
    // check in the DB
    if username == "jack" && password == "1234" {
      return true, nil
    }
    return false, nil
  }))

  cookieGroup.Use(checkCookie)

  cookieGroup.GET("/main", mainCookie)

  adminGroup.GET("/main", mainAdmin)

  e.GET("/login", login)
	e.GET("/", yallo)
	e.GET("/cats/:data", getCats)

	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamsters", addHamster)

	e.Logger.Fatal(e.Start("localhost:1323"))
}
