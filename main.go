package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
	"time"
)


func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		// c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	engine := gee.New()
	engine.Use(gee.Logger())

	engine.GET("/", func(c *gee.Context) {
		fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
	})

	v2 := engine.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	hello := engine.Group("/hello")
	{

		hello.GET("/", func(c *gee.Context) {
			for k, v := range c.Req.Header {
				fmt.Fprintf(c.Writer, "Header[%q] = %q\n", k, v)
			}
		})

		hello.GET("/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
		})


	}
	engine.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	engine.Run(":9999")
}