package main

import (
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Resp -
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func loginHandler(c *gin.Context) {
	var (
		user, json Users
		err        error
	)
	if err = c.BindJSON(&json); err != nil {
		c.JSON(400, Resp{Code: 0, Msg: err.Error()})
		return
	}
	if user, err = findUsersByUsername(json.Username); err != nil {
		c.JSON(401, Resp{Code: 0, Msg: err.Error()})
		return
	}
	if user.Password != MD5(json.Password) {
		c.JSON(401, Resp{Code: 0, Msg: "password is error"})
		return
	}
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("role", user.Role)
	session.Save()
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: user})
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.JSON(200, Resp{Code: 0, Msg: "ok"})
}

func webhookHandler(c *gin.Context) {
	var (
		err      error
		receives []*Receives
		data     []byte
	)
	// TODO: add receive log
	data, _ = ioutil.ReadAll(c.Request.Body)
	if receives, err = findRecevices(); err != nil {
		logger.Errorln(err.Error())
		// TODO: add err logs
	}
	headers := url.Values(c.Request.Header).Encode()
	for _, receive := range receives {
		if strings.Contains(headers, receive.Header) ||
			strings.Contains(string(data), receive.Keyword) {
			go parseDataAndPush(data, receive)
		}
	}
	c.String(200, "ok")
}

func addUserHandler(c *gin.Context) {
	var (
		user Users
		err  error
	)
	if err = c.BindJSON(&user); err != nil {
		c.JSON(400, Resp{Code: 0, Msg: err.Error()})
		return
	}
	user.Password = MD5(user.Password)
	if user, err = addUser(user); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: user})
}

func addReceiveHandler(c *gin.Context) {
	var (
		receive Receives
		err     error
	)
	if err = c.BindJSON(&receive); err != nil {
		c.JSON(400, Resp{Code: 0, Msg: err.Error()})
		return
	}
	if receive, err = addReceive(receive); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: receive})
}

func addPusherHandler(c *gin.Context) {
	var (
		pusher Pushers
		err    error
	)
	if err = c.BindJSON(&pusher); err != nil {
		c.JSON(400, Resp{Code: 0, Msg: err.Error()})
		return
	}
	if pusher, err = addPusher(pusher); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: pusher})
}

func addRelationHandler(c *gin.Context) {
	var (
		relation Relations
		err      error
	)
	if err = c.BindJSON(&relation); err != nil {
		c.JSON(400, Resp{Code: 0, Msg: err.Error()})
		return
	}
	if relation, err = addRelation(relation); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: relation})
}

func addTemplateHandler(c *gin.Context) {
	var (
		template Templates
		err      error
	)
	if err = c.BindJSON(&template); err != nil {
		c.JSON(400, Resp{Code: 0, Msg: err.Error()})
		return
	}
	if template, err = addTemplate(template); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: template})
}
func getUserHandler(c *gin.Context) {
	var (
		users []*Users
		err   error
	)
	if users, err = findUsersByID(c.Param("id")); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: users})
}

func getReceiveHandler(c *gin.Context) {
	var (
		receives []*Receives
		err      error
	)
	if receives, err = findReceivesByID(c.Param("id")); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: receives})
}

func getPusherHandler(c *gin.Context) {
	var (
		pushers []*Pushers
		err     error
	)
	if pushers, err = findPushersByID(c.Param("id")); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: pushers})
}

func getRelationHandler(c *gin.Context) {
	var (
		relations []*Relations
		err       error
	)
	if relations, err = findRelationsByID(c.Param("id")); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: relations})
}

func getTemplateHandler(c *gin.Context) {
	var (
		templates []*Templates
		err       error
	)
	if templates, err = findTemplatesByID(c.Param("id")); err != nil {
		c.JSON(503, Resp{Code: 0, Msg: err.Error()})
		return
	}
	c.JSON(200, Resp{Code: 0, Msg: "ok", Data: templates})
}

func sessionAuth(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("username") == nil || session.Get("username").(string) == "" {
		c.AbortWithStatusJSON(401, Resp{Code: -1, Msg: "Please login!"})
		return
	}
	if strings.HasPrefix(c.Request.RequestURI, "/user") {
		if session.Get("role") == nil || session.Get("username").(int) <= 0 {
			c.AbortWithStatusJSON(401, Resp{Code: -1, Msg: "Permission denied"})
			return
		}
	}
}

func newRouter() *gin.Engine {
	r := gin.Default()

	store := sessions.NewCookieStore([]byte(conf.Server.Secret))
	store.Options(sessions.Options{MaxAge: int(30 * time.Minute), Path: "/"})
	r.Use(sessions.Sessions("session", store))

	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// 使用 Logger 中间件
	r.Use(gin.Logger())
	// 使用 Recovery 中间件
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) { c.JSON(200, Resp{Code: 0, Msg: "pong"}) })

	r.POST("/webhook", webhookHandler)

	r.GET("/", func(c *gin.Context) { c.HTML(200, "index.html", nil) })

	// r.GET("/login", func(c *gin.Context) { c.HTML(200, "login.html", nil) })
	r.POST("/login", loginHandler)
	r.GET("/logout", logoutHandler)

	g := r.Group("/api", sessionAuth)
	{
		g.POST("/receive", addReceiveHandler)
		g.POST("/pusher", addPusherHandler)
		g.POST("/relation", addRelationHandler)
		g.POST("/template", addTemplateHandler)

		g.GET("/receive", getReceiveHandler)
		g.GET("/receive/:id", getReceiveHandler)
		g.GET("/pusher", getPusherHandler)
		g.GET("/pusher/:id", getPusherHandler)
		g.GET("/relation", getRelationHandler)
		g.GET("/relation/:id", getRelationHandler)
		g.GET("/template", getTemplateHandler)
		g.GET("/template/:id", getTemplateHandler)

		g.POST("/user", addUserHandler)
		g.GET("/user", getUserHandler)
		g.GET("/user/:id", getUserHandler)

	}

	r.POST("/debug/user", addUserHandler)

	return r
}
