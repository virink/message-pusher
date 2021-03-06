package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	// mysql driver

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	MySQL struct {
		Host    string `yaml:"host"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
		Name    string `yaml:"name"`
		Charset string `yaml:"charset"`
	} `yaml:"mysql"`
	Server struct {
		Debug  bool   `yaml:"debug"`
		Port   int    `yaml:"port"`
		Secret string `yaml:"secret"`
	} `yaml:"server"`
}

func initConnect() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.MySQL.User, conf.MySQL.Pass, conf.MySQL.Host, conf.MySQL.Name, conf.MySQL.Charset,
	)
	if db, err = gorm.Open("mysql", dsn); err != nil {
		logger.Errorln(err.Error())
		return nil, err
	}
	db.LogMode(conf.Server.Debug)
	// db.Debug()

	db.DB().SetConnMaxLifetime(100 * time.Second) // 最大连接周期，超过时间的连接就close
	db.DB().SetMaxOpenConns(100)                  // 设置最大连接数
	db.DB().SetMaxIdleConns(16)                   // 设置闲置连接数

	return
}

func loadConfig() (err error) {
	var (
		yamlFile []byte
	)
	if yamlFile, err = ioutil.ReadFile("config.yaml"); err != nil {
		logger.Errorln(err.Error())
		return err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		logger.Errorln(err.Error())
		return err
	}
	return nil
}

func initLogger(filename string, level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	if level == logrus.DebugLevel || level == logrus.InfoLevel {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:            true,
			DisableLevelTruncation: false,
			TimestampFormat:        "2006-01-02 15:04:05",
		})
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(logFile)
		}
	}
	return logger
}

func httpRequest(uri string, data string) (body []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	if req, err = http.NewRequest("POST", uri, strings.NewReader(data)); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return body, nil
}

func parseDataAndPush(data []byte, receive *Receives) {
	var (
		err    error
		varMap = make(map[string]string)
	)
	dataObj := gjson.ParseBytes(data)
	vars := strings.Split(receive.Variable, ",")
	for _, v := range vars {
		t := dataObj.Get(v)
		if t.Exists() {
			varMap[v] = t.String()
		}
	}

	var pushers []*Pushers
	if pushers, err = findPusherByRecevice(receive.ID); err != nil {
		logger.Errorln(err)
		return
	}
	for _, push := range pushers {
		go func(push *Pushers) {
			var (
				body   []byte
				result string
			)
			result = push.Template
			for k, v := range varMap {
				result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", k), v)
			}
			// Fix
			result = strings.ReplaceAll(result, "\n", "\\n")
			logger.Debugln(result)
			if body, err = httpRequest(push.URL, result); err != nil {
				logger.Errorln(err)
			}
			logger.Debugln(string(body))
		}(push)
	}
}

// MD5 -
func MD5(text string) string {
	ctx := md5.New()
	_, _ = ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
