package main

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Users -
type Users struct {
	gorm.Model
	Username string `gorm:"type:varchar(25);unique_index" json:"username"`
	Password string `gorm:"type:varchar(256)" json:"password"`
	Role     int    `gorm:"default:0" json:"role"`
}

// Receives -
type Receives struct {
	gorm.Model
	Name     string `json:"name"`
	Type     string
	Header   string
	Keyword  string
	Body     string
	Variable string
}

// Pushers -
type Pushers struct {
	gorm.Model
	URL      string // with key
	Name     string
	Vendor   string
	Template string
}

// Templates -
type Templates struct {
	gorm.Model
	URL    string
	Vendor string // Receive Vendor
	Name   string // Template Name
	Body   string // Template Body
}

// Relations -
type Relations struct {
	gorm.Model
	Status    bool `gorm:"default:true" json:"status"`
	UserID    uint `json:"user_id"`
	PusherID  uint `json:"pusher_id"`
	ReceiveID uint `json:"receive_id"`
}

// Logs -
type Logs struct {
	gorm.Model
	ReceivesID uint
	PushersID  uint
}

func debugDB() {
	user, _ := addUser(Users{Username: "virink", Password: MD5("123456"), Role: 9})
	logger.Debugln(user)
	r, _ := addReceive(Receives{
		Name:     "chamd5_ti_dingding",
		Type:     "dingding",
		Keyword:  "dingding",
		Variable: `markdown.text,markdown.title`,
		Body:     `{"msgtype":"markdown","markdown":{"title":"xxx","text":"xxx"}}`,
	})
	logger.Debugln(r)

	// t1, _ := addTemplate(Templates{
	// 	URL:    "https://oapi.dingtalk.com/robot/send?access_token=$key",
	// 	Vendor: "dingding",
	// 	Name:   "dingding_text",
	// 	Body:   `{"msgtype":"text","text":{"content":"key.subkey"}`,
	// })
	// t2, _ := addTemplate(Templates{
	// 	URL:    "https://open.feishu.cn/open-apis/bot/hook/$key",
	// 	Vendor: "feishu",
	// 	Name:   "feishu_text",
	// 	Body:   `{"msgtype":"text","text":{"content":"key.subkey"}`,
	// })
	// logger.Debugln(t1)
	// logger.Debugln(t2)
	p1, _ := addPusher(Pushers{
		URL:      "https://oapi.dingtalk.com/robot/send?access_token=e38ed299a78666082c96034ed47efae6aa8e812ed3aa9b4264d26a89ec534288",
		Name:     "chamd5_ti_dingding_text",
		Vendor:   "dingding",
		Template: `{"msgtype":"markdown","markdown":{"title":"${markdown.title}","text":"${markdown.text}"}}`,
	})
	p2, _ := addPusher(Pushers{
		URL:      "https://open.feishu.cn/open-apis/bot/hook/e3fe90555f7b47289cabb3ab2b5a239f",
		Name:     "chamd5_ti_feishu_text",
		Vendor:   "feishu",
		Template: `{"title":"${markdown.title}","text":"${markdown.text}"}`,
	})
	logger.Debugln(p1)
	logger.Debugln(p2)

	addRelation(Relations{
		Status:    true,
		UserID:    user.ID,
		PusherID:  p1.ID,
		ReceiveID: r.ID,
	})
	addRelation(Relations{
		Status:    true,
		UserID:    user.ID,
		PusherID:  p2.ID,
		ReceiveID: r.ID,
	})
}

func initDatabase() {
	db.DropTableIfExists(&Users{}, &Pushers{}, &Receives{}, &Relations{}, &Templates{})
	db.CreateTable(&Users{}, &Pushers{}, &Receives{}, &Relations{}, &Templates{})
	addUser(Users{Username: "virink", Password: MD5("123456"), Role: 9})
	// debugDB()
	// db.CreateTable(&Templates{})
	// db.Model(&Relations{}).AddForeignKey(field string, dest string, onDelete string, onUpdate string)
}

func addUser(user Users) (out Users, err error) {
	if db.First(&user, Users{Username: user.Username}).RecordNotFound() {
		if err = db.Create(&user).Error; err != nil {
			return
		}
		return user, nil
	}
	return out, errors.New("Username is exists")
}

func addRelation(relation Relations) (out Relations, err error) {
	logger.Debug(relation)
	if err = db.FirstOrCreate(&out, &relation).Error; err != nil {
		return out, err
	}
	return out, nil
}

func addPusher(pusher Pushers) (out Pushers, err error) {
	if err = db.FirstOrCreate(&out, pusher).Error; err != nil {
		return out, err
	}
	return out, nil
}

func addReceive(receive Receives) (out Receives, err error) {
	if err = db.FirstOrCreate(&out, receive).Error; err != nil {
		return out, err
	}
	return out, nil
}

func addTemplate(template Templates) (out Templates, err error) {
	logger.Debug(template)
	if err = db.FirstOrCreate(&out, &template).Error; err != nil {
		return out, err
	}
	return out, nil
}

func findPushers(username string) {
	user := Users{
		Username: username,
	}
	var puserhs Pushers
	db.Model(&user).Related(&puserhs)
}

func findRecevices() (receives []*Receives, err error) {
	if err = db.Find(&receives).Error; err != nil {
		return
	}
	return receives, nil
}

func findPusherByRecevice(rid uint) (pushers []*Pushers, err error) {
	if err = db.Joins("JOIN relations ON relations.pusher_id = pushers.id").
		Joins("JOIN users ON users.id = relations.user_id").
		Joins("JOIN receives ON receives.id = relations.receive_id").
		Where("receives.id = ?", rid).Find(&pushers).Error; err != nil {
		return
	}
	return pushers, nil
}

func findUsersByUsername(username string) (user Users, err error) {
	if err = db.First(&user, Users{Username: username}).Error; err != nil {
		return
	}
	return user, nil
}

func findUsers(username, password string) (user Users, err error) {
	if err = db.First(&user, Users{Username: username, Password: password}).Error; err != nil {
		return
	}
	return user, nil
}

func findUsersByID(id string) (users []*Users, err error) {
	logger.Debugln(id)
	stmp := db.New()
	if id != ":id" && id != "" {
		stmp = stmp.Where("id = ?", id)
	}
	if err = stmp.Find(&users).Error; err != nil {
		return
	}
	return users, nil
}

func findReceivesByID(id string) (receives []*Receives, err error) {
	stmp := db.New()
	if id != ":id" && id != "" {
		stmp = stmp.Where("id = ?", id)
	}
	if err = stmp.Find(&receives).Error; err != nil {
		return
	}
	return receives, nil
}

func findPushersByID(id string) (pushers []*Pushers, err error) {
	stmp := db.New()
	if id != ":id" && id != "" {
		stmp = stmp.Where("id = ?", id)
	}
	if err = stmp.Find(&pushers).Error; err != nil {
		return
	}
	return pushers, nil
}

func findRelationsByID(id string) (relations []*Relations, err error) {
	stmp := db.New()
	if id != ":id" && id != "" {
		stmp = stmp.Where("id = ?", id)
	}
	if err = stmp.Find(&relations).Error; err != nil {
		return
	}
	return relations, nil
}

func findTemplatesByID(id string) (templates []*Templates, err error) {
	stmp := db.New()
	if id != ":id" && id != "" {
		stmp = stmp.Where("id = ?", id)
	}
	if err = stmp.Find(&templates).Error; err != nil {
		return
	}
	return templates, nil
}
