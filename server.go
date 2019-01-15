package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const prettyPrintValue = `false`

const oauthKey = `AIzaSyC2as3n3-616_w1GMwzhEL6JDa3HumoiUA`

//const oauthKey = `AIzaSyDcYUOG2yK1DX_CWPGe-WrlhKhQqzWQy3I`

const engine = `004560508768852172973:ymzun4smzzu`

//const engine = `017576662512468239146:omuauf_lfve`

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

	r.GET("/buscar/:search", func(c *gin.Context) {
		result := search(c.Param("search"))
		hola, _err := json.Marshal(result)
		if _err != nil {
			panic(_err)
		}
		c.JSON(200, string(hola))
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func search(searchable string) Page {
	googleSearch := makeRequest(replaceSpace(searchable))
	var response Page
	response.Context = googleSearch.Context
	response.Kind = googleSearch.Kind
	response.Queries = googleSearch.Queries
	response.SearchInformation = googleSearch.SearchInformation
	response.URL = googleSearch.URL

	for _, element := range googleSearch.Items {
		resource := makeFilterRequest(element.Snippet)
		if resource.Categories[0].Name == "technology" || resource.Categories[1].Name == "technology" || resource.Categories[2].Name == "technology" {
			response.Items = append(response.Items, element)
		}
	}
	response.SearchInformation.FormattedTotalResults = string(len(response.Items))
	response.SearchInformation.TotalResults = string(len(response.Items))
	return response
}

func makeRequest(query string) Page {
	resp, err := http.Get("https://www.googleapis.com/customsearch/v1?q=" + query + "&cx=" + engine + "&key=" + oauthKey + "&prettyPrint=" + prettyPrintValue)
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
	resp, err := http.Get("https://api.dandelion.eu/datatxt/cl/v1/?token=420239a12b934cc39193d9c1fa79958a&model=54cf2e1c-e48a-4c14-bb96-31dc11f84eac&text=" + query)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var structure TextClasification
	if err := json.Unmarshal(body, &structure); err != nil {
		panic(err)
	}
	return structure
}

func replaceSpace(s string) string {
	return strings.Replace(s, " ", "+", -1)
}
