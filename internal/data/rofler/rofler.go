package rofler

import (
	"time"
)

const (
	OfTheWeek  = "of_the_week"
	OfTheMonth = "of_the_month"
	TopRofler  = "top_rofler"
	TopTaste   = "top_taste"
)

type Rofler struct {
	UserName   string     `firestore:"user_name"`
	TelegramId int64      `firestore:"telegram_id"`
	Accolades  []Accolade `firestore:"accolades"`
	UserPic    []byte     `firestore:"user_pic"`
}

type Accolade struct {
	Type      string    `firestore:"type"`
	TimeFrame string    `firestore:"time_frame"`
	WonOn     time.Time `firestore:"won_on"`
}
