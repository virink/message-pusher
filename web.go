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
	if c.DefaultPostForm("username", "") == "admin" && c.DefaultPostForm("password", "") == "123456" {
		session := sessions.Default(c)
		session.Set("username", "admin")
		session.Save()
		c.Redirect(302, "/manager")
		return
	}
	c.Redirect(302, "/login")
	return
}

func managerHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username != "admin" {
		c.Redirect(302, "/login")
		return
	}
	// TODO:
}

func webhookHandler(c *gin.Context) {
	var (
		err      error
		receives []*Receives
		data     []byte
	)
	// TODO: add receive log
	data, _ = ioutil.ReadAll(c.Request.Body)
	if receives, err = getRecevices(); err != nil {
		logger.Errorln(err.Error())
		// TODO: add err logs
	}

	headers := url.Values(c.Request.Header).Encode()
	for _, receive := range receives {
		// Header
		if strings.Contains(headers, receive.Header) ||
			strings.Contains(string(data), receive.Keyword) {
			go parseDataAbdPush(data, receive)
			continue
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
		c.JSON(400, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	if user, err = addUser(user.Username, user.Password); err != nil {
		c.JSON(503, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, Resp{
		Code: 0,
		Msg:  "ok",
		Data: user,
	})
}

func addReceiveHandler(c *gin.Context) {
	var (
		receive Receives
		err     error
	)
	if err = c.BindJSON(&receive); err != nil {
		c.JSON(400, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	if receive, err = addReceive(receive); err != nil {
		c.JSON(503, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, Resp{
		Code: 0,
		Msg:  "ok",
		Data: receive,
	})
}

func addPusherHandler(c *gin.Context) {
	var (
		pusher Pushers
		err    error
	)
	if err = c.BindJSON(&pusher); err != nil {
		c.JSON(400, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	if pusher, err = addPusher(pusher); err != nil {
		c.JSON(503, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, Resp{
		Code: 0,
		Msg:  "ok",
		Data: pusher,
	})
}

func addRelationHandler(c *gin.Context) {
	var (
		relation Relations
		err      error
	)
	if err = c.BindJSON(&relation); err != nil {
		c.JSON(400, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	if relation, err = addRelation(relation); err != nil {
		c.JSON(503, Resp{
			Code: 0,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, Resp{
		Code: 0,
		Msg:  "ok",
		Data: relation,
	})
}

func newRouter() *gin.Engine {
	r := gin.Default()

	store := sessions.NewCookieStore([]byte(conf.Server.Secret))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})

	r.Use(sessions.Sessions("session", store))

	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// 使用 Logger 中间件
	r.Use(gin.Logger())
	// 使用 Recovery 中间件
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/webhook", webhookHandler)

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.POST("/login", loginHandler)
	r.POST("/user/addreceive", addReceiveHandler)
	r.POST("/user/addpusher", addPusherHandler)
	r.POST("/user/addrelation", addRelationHandler)
	g := r.Group("/admin")
	{
		g.POST("/adduser", addUserHandler)
	}

	// g := r.Group("/api/v1")
	// {
	// }

	return r
}
