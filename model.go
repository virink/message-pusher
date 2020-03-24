package main

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Users -
type Users struct {
	gorm.Model
	Username string `gorm:"type:varchar(25);unique_index"`
	Password string `gorm:"type:varchar(256)" json:"password"`
}

// Receives -
type Receives struct {
	gorm.Model
	Name    string
	Type    string
	Header  string
	Keyword string
	Body    string
}

// Pushers -
type Pushers struct {
	gorm.Model
	URL        string
	Name       string
	Vendor     string
	TemplateID uint `json:"template_id"`
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
	// Users     Users     `gorm:"foreignkey:UserID"`
	// Pushers   []Pushers `gorm:"foreignkey:PusherID"`
	// Receives  Receives  `gorm:"foreignkey:ReceiveID"`
}

// `{"title": "Feishu","text": "Required"}`
// `{"name":"virink","sex":1,"text":"hello world","ignore":"keyword"}`

// Logs -
type Logs struct {
	gorm.Model
	ReceivesID uint
	PushersID  uint
}

func initDatabase() {
	db.DropTableIfExists(&Users{}, &Pushers{}, &Receives{}, &Relations{})
	db.CreateTable(&Users{}, &Pushers{}, &Receives{}, &Relations{})
	// db.Model(&Relations{}).AddForeignKey(field string, dest string, onDelete string, onUpdate string)
}

func addUser(username, password string) (user Users, err error) {
	if db.First(&user, Users{Username: username}).RecordNotFound() {
		user = Users{
			Username: username,
			Password: password,
		}
		if err = db.Create(&user).Error; err != nil {
			return
		}
		return user, nil
	}
	return user, errors.New("Username is exists")
	// if err = db.FirstOrCreate(&out, &user).Error; err != nil {
	// 	return out, err
	// }
	// return out, nil
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

func getPushers(username string) {
	user := Users{
		Username: username,
	}
	var puserhs Pushers
	db.Model(&user).Related(&puserhs)
}

func getRecevices() (receives []*Receives, err error) {
	if err = db.Find(&receives).Error; err != nil {
		return
	}
	return receives, nil
}

func getPusherByRecevice(rid uint) (pushers []*Pushers, err error) {
	if err = db.Joins("JOIN relations ON relations.pusher_id = pushers.id").
		Joins("JOIN users ON users.id = relations.user_id").
		Joins("JOIN receives ON receives.id = relations.receive_id").
		Where("relations.id = ?", rid).Find(&pushers).Error; err != nil {
		return
	}
	return pushers, nil
}
