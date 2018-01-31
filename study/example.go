package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Passwd   string `form:"passwd" json:"passwd" binding:"required"`
	Age      int    `form:"age" json:"age"`
}

func main() {
	r := gin.Default()

	//curl -X POST http://127.0.0.1:8000/login -H "Content-Type:application/json" \
	//		-d '{"username": "rsj217", "passwd": "123", "age": 21}' | python -m json.tool
	r.POST("/login", func(c *gin.Context) {
		var user User
		contentType := c.Request.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err := c.BindJSON(&user)
			if err != nil {
				log.Fatal(err)
			}
		case "application/x-www-form-urlencoded":
			err := c.BindWith(&user, binding.Form)
			if err != nil {
				log.Fatal(err)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"username": user.Username,
			"passwd":   user.Passwd,
			"age":      user.Age,
		})
	})

	r.POST("/login1", func(c *gin.Context) {
		var user User
		err := c.Bind(&user)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"username": user.Username,
			"passwd":   user.Passwd,
			"age":      user.Age,
		})
	})

	//render xml
	r.GET("/render", func(c *gin.Context) {
		contentType := c.DefaultQuery("contentType", "json")

		if contentType == "json" {
			c.JSON(http.StatusOK, gin.H{
				"username": "linheng",
				"passwd":   "test",
			})
		} else if contentType == "xml" {
			c.XML(http.StatusOK, gin.H{
				"username": "linheng",
				"passwd":   "test",
			})
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/h", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World.")
	})

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	//Querystring
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Lin")
		lastname := c.Query("lastname")
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	//POST Form
	//curl -X POST http://127.0.0.1:8000/post_form -H "Content-Type:application/x-www-form-urlencoded" -d "message=welcome&nick=jetlin"
	r.POST("/post_form", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
			"message": message,
			"nick":    nick,
		})
	})

	//put
	r.PUT("/put", func(c *gin.Context) {
		id := c.Query("id")
		page := c.Query("page")
		message := c.PostForm("message")
		nick := c.PostForm("nick")
		fmt.Printf("%s %s %s %s", id, page, message, nick)
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusOK,
		})
	})

	//upload html
	r.LoadHTMLGlob("template/*")
	r.GET("upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})

	//upload a file
	//curl -X POST http://127.0.0.1:8000/upload -F "upload=@/Users/jetlin/go/src/markdown/main.go" -H "Content-Type: multipart/form-data"
	r.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		fmt.Println(name)
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad Request")
			return
		}
		filename := header.Filename

		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		b, err := io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusCreated, "%d upload successful.\n", b)
	})

	//upload multiple files
	//curl -X POST http://127.0.0.1:8000/multi/upload -F "upload=@/Users/jetlin/go/src/markdown/main.go" \
	// 		-F "upload=@/Users/jetlin/go/src/urlrouting/httpr.go" -H "Content-Type: multipart/form-data"
	r.POST("/multi/upload", func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(20000)
		if err != nil {
			log.Fatal(err)
		}

		formdata := c.Request.MultipartForm

		files := formdata.File["upload"]

		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}

			out, err := os.Create(files[i].Filename)
			defer out.Close()
			if err != nil {
				log.Fatal(err)
			}

			b, err := io.Copy(out, file)
			if err != nil {
				log.Fatal(err)
			}
			c.String(http.StatusCreated, "%d upload successful.", b)
		}
	})

	r.GET("/redirect", Redirect)

	r.Use(Middleware())
	{
		r.GET("/middle", func(c *gin.Context) {
			request := c.MustGet("request").(string)
			req, _ := c.Get("request")
			c.JSON(http.StatusOK, gin.H{
				"middle_request": request,
				"request":        req,
			})
		})
	}

	r.GET("/before", Middleware(), func(c *gin.Context) {
		request := c.MustGet("request")
		c.JSON(http.StatusOK, gin.H{
			"middle_request": request,
		})
	})

	r.GET("/auth/signin", func(c *gin.Context) {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "123",
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		c.String(http.StatusOK, "login successful...")
	})

	r.GET("/home", AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": "Home",
		})
	})

	r.GET("/sync", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		log.Println("Done! in path:", c.Request.URL.Path)
	})

	r.GET("/async", func(c *gin.Context) {
		cCp := c.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			log.Println("Done! in path:", cCp.Request.URL.Path)
		}()
	})

	r.Run(":8000")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("session_id")
		if err != nil {
			log.Fatal(err)
		}
		value := cookie.Value
		if value == "123" {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}

func Redirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
}

//middleware
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request", "client_request")
		c.Next()
	}
}
