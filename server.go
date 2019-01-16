package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const prettyPrintValue = `false`

const oauthKey = ``


const engine = ``


func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func main() {

	r := gin.Default()

	r.Use(CORSMiddleware())

	r.GET("/buscar/:search/:tipo", func(c *gin.Context) {
		result := search(c.Param("search"), c.Param("tipo"))
		hola, _err := json.Marshal(result)
		if _err != nil {
			panic(_err)
		}
		c.JSON(200, string(hola))
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func search(searchable, tipo string) Page {
	googleSearch := makeRequest(replaceSpace(searchable), tipo)
	var response Page
	response.Context = googleSearch.Context
	response.Kind = googleSearch.Kind
	response.Queries = googleSearch.Queries
	response.SearchInformation = googleSearch.SearchInformation
	response.URL = googleSearch.URL
	for _, element := range googleSearch.Items {
		resource := makeFilterRequest(replaceSpace(element.Snippet))
		if len(resource.Categories) > 0 {
			if resource.Categories[0].Name == "technology" {
				response.Items = append(response.Items, element)
			}
		}
	}
	return response
}

func makeRequest(query, tipo string) Page {
	link := "https://www.googleapis.com/customsearch/v1?q=" + query + "&cx=" + engine + "&key=" + oauthKey + "&prettyPrint=" + prettyPrintValue

	if tipo == "img" {
		link = link + "&searchType=image"
	}

	resp, err := http.Get(link)

	if err != nil {
		log.Fatalln(err)
	}
	var structure Page
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &structure); err != nil {
		panic(err)
	}
	return structure
}

func makeFilterRequest(query string) TextClasification {
	fmt.Println(url.QueryEscape(query))
	req, _ := http.NewRequest("GET", "https://api.dandelion.eu/datatxt/cl/v1/", nil)
	q := req.URL.Query()
	q.Add("token", "e2870fbc4dfa49e48a889c9681526aa4")
	q.Add("model", "54cf2e1c-e48a-4c14-bb96-31dc11f84eac")
	q.Add("text", replaceSpace(query))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	structure := new(TextClasification)
	err = json.NewDecoder(resp.Body).Decode(structure)
	if err != nil {
		panic(err)
	}
	return *structure
}

func replaceSpace(s string) string {
	s = strings.Replace(s, " ", "+", -1)
	return strings.Replace(s, "/", "", -1)

}
