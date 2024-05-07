package models

import (
	"time"

	"github.com/lib/pq"
)

type OrderStatus string

const (
	Normal      OrderStatus = "normal"      // 報名參加
	Attending   OrderStatus = "attending"   // 已報名
	Attended    OrderStatus = "attended"    // 之前的活動
	Recommended OrderStatus = "recommended" // 推薦參加
)

type Order struct {
	ID                  string         `json:"id" db:"id" example:"uuid"`
	CreatorID           string         `json:"creator_id" db:"creator_id" example:"uuid"`
	IsActive            bool           `json:"is_active" db:"is_active"`
	IsDeleted           bool           `json:"is_deleted" db:"is_deleted"`
	Status              OrderStatus    `json:"status" db:"-"`
	Title               string         `json:"title" db:"title"`
	Content             string         `json:"content" db:"content"`
	ParticipantCount    int            `json:"participant_count" db:"participant_count"`
	MaxParticipantCount *int           `json:"max_participant_count,omitempty" db:"max_participant_count"`
	Tags                pq.StringArray `json:"tags" db:"tags"`
	PictureURL          *string        `json:"picture_url,omitempty" db:"picture_url"`
	CreatedAt           time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at" db:"updated_at"`
	HostingAt           time.Time      `json:"hosting_at" db:"hosting_at"`
	EndingAt            *time.Time     `json:"ending_at" db:"ending_at"`
	Speaker             *string        `json:"speaker" db:"speaker"`
	Point               *string        `json:"point" db:"point"`
	Contact             *string        `json:"contact" db:"contact"`
	ApplyInfo           *string        `json:"apply_info" db:"apply_info"`
	ChargingFee         *string        `json:"charging_fee" db:"charging_fee"`
	Location            *string        `json:"location" db:"location"`
	Hoster              *string        `json:"hoster" db:"hoster"`
	URL                 *string        `json:"url" db:"url"`
	Coins               int            `json:"coins" db:"coins"`
	Awards              int            `json:"awards" db:"awards"`
}
