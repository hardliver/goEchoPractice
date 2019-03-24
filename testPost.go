package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	// "reflect"
)

func post(url, data string) *http.Response {
	jsonStr := []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// fmt.Println(reflect.TypeOf(resp).String())
	return resp
}

func main() {
	resp := post("http://localhost:1323/cats", `{"name":"fishmaster", "type":"cat-fish"}`)
	defer resp.Body.Close()
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	resp = post("http://localhost:1323/dogs", `{"name":"doggymaster", "type":"dog-fish"}`)
	defer resp.Body.Close()
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	resp = post("http://localhost:1323/hamsters", `{"name":"da hamster", "type":"master"}`)
	defer resp.Body.Close()
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
