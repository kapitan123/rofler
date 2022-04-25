package model

import (
	"time"
)

type Rofler struct {
	UserName   string     `firestore:"user_name"`
	TelegramId int64      `firestore:"telegram_id"`
	Accolades  []Accolade `firestore:"accolades"`
	UserPic    []byte     `firestore:"user_pic"`
}

type Accolade struct {
	Type      string    `firestore:"type"`       // rofler, worker, taste
	TimeFrame string    `firestore:"time_frame"` // week, month, year
	WonOn     time.Time `firestore:"won_on"`
}
