package main

import (
	"JWT-AUTH-GIN/controllers"
	"JWT-AUTH-GIN/initializers"
	"bytes"
	"github.com/freshman-tech/news-demo/news"
	"github.com/gin-gonic/gin"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

var tpl = template.Must(template.ParseFiles("index.html"))

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

func indexHandler(c *gin.Context) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.Write(buf.Bytes())
}

func searchHandler(newsapi *news.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := url.Parse(c.Request.URL.String())
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults / newsapi.PageSize))),
			Results:    results,
		}

		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}

		// Render the search results using the template
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Write the rendered template to the response
		c.Writer.WriteString(buf.String())
	}
}

func isLoggedIn(r *http.Request) bool {
	// Check if user is logged in (implement your own logic here)
	// For example, you can check if a session cookie exists
	_, err := r.Cookie("Authorization")
	return err == nil
}

func authMiddleware(c *gin.Context) {
	if isLoggedIn(c.Request) {
		// User is logged in, proceed to the next handler
		c.Next()
	} else {
		// User is not logged in, render a message and a link to the login page
		c.AbortWithStatus(http.StatusUnauthorized)
		c.Writer.WriteString("You are not logged in. Please <a href=\"/\">login</a>.")
	}
}

func main() {
	r := gin.Default()
	r.Static("/css", "./templates/css")

	r.LoadHTMLGlob("templates/*.html")
	r.POST("/signup", controllers.Signup())
	r.POST("/login", controllers.Login())
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index1.html", map[string]string{"title": "Register"})
	})

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		println("Env: apiKey must be set")
		return
	}
	r.Use()
	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, 20)
	r.Static("/assets", "./assets")
	r.GET("/search", searchHandler(newsapi))
	r.GET("/main", authMiddleware, indexHandler)
	//r.GET("/validate", middleware.RequireAuth(), controllers.Validate)

	r.Run()
}
