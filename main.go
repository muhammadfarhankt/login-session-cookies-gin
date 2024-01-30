package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge:   600,
		HttpOnly: true,
	})
	router.Use(sessions.Sessions("mysession", store))
	//fmt.Println("store : ",store)

	router.LoadHTMLGlob("templates/*.html")

	//login get
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		if username != nil {
			c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
			c.Redirect(http.StatusSeeOther, "/home")
			return
		}
		c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
		c.HTML(http.StatusOK, "login.html", nil)
		// c.HTML(http.StatusOK, "login.html", gin.H{
		// "message" : "Hello World from Go GIN",
		// })
		//c.String(200, "Hello WOrld")
		// })
	})

	//login post
	router.POST("/", func(c *gin.Context) {
		// c.String(200, "Login POST")

		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "admin" && password == "pass" {
			session := sessions.Default(c)
			session.Set("username", username)
			// session.Set("")
			session.Set("lastActivity", time.Now().Unix())
			session.Save()
			c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
			c.Redirect(http.StatusSeeOther, "/home")
		} else {
			c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "Invalid login credentials"})
		}
	})

	//home page  loading after login
	router.GET("/home", func(c *gin.Context) {
		// c.String(200, "Home load")

		session := sessions.Default(c)
		username := session.Get("username")
		lastActivity := session.Get("lastActivity")
		//lastActivity := time.Unix(lastActivity.(int64), 0)
		fmt.Println("last acitivity : ", lastActivity)
		fmt.Println("username : ", username)
		if lastActivity == nil || username == nil {
			fmt.Println("inside if lastActivity == nil || username == nil")
			//c.Redirect(http.StatusSeeOther, "/")
			c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
			//c.Redirect(http.StatusSeeOther, "/?message=Must%20Login")
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "Must login before accessing User home page"})
			return
		}
		// if username == nil {
		// 	//c.Redirect(http.StatusSeeOther, "/")
		// 	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
		// 	//c.Redirect(http.StatusSeeOther, "/?message=Must%20Login")
		// 	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "Must login before accessing User home page"})
		// 	return
		// }
		// timeoutDuration := 1 * time.Minute // Timeout duration
		// fmt.Println("Time out duration : ",timeoutDuration)
		// if time.Now().Unix()-lastActivity.(int64) > int64(timeoutDuration.Seconds()) {
		// 	fmt.Println("if case - session time out")
		// 	// Session has exceeded the timeout, clear the session
		// 	session.Clear()
		// 	session.Save()
		// 	// c.Redirect(http.StatusSeeOther, "/")
		// 	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
		// 	//c.Redirect(http.StatusSeeOther, "/login?message=Session%20Time%20Out")
		// 	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "Session Time Out"})
		// 	return
		// }
		lastActivityTime := time.Unix(lastActivity.(int64), 0)
		fmt.Println("lastActivityTime : ", lastActivityTime)
		// Check if the session has expired
		if time.Since(lastActivityTime) > time.Minute*1 {
			// Session has expired, redirect to login
			fmt.Println("inside if condition")
			session.Clear()
			session.Save()
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "Session Time Out"})
			return
		}

		// Update last activity timestamp
		session.Set("lastActivity", time.Now())
		session.Save()
		fmt.Println("session time updated")

		c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
		c.HTML(http.StatusOK, "home.html", gin.H{"username": username})
	})

	//logout
	router.GET("/logout", func(c *gin.Context) {
		// c.String(200, "Logout GET")
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8080")
}
