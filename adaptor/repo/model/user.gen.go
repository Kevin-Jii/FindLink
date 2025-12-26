package model

import (
	"time"
)

const TableNameUser = "user"

// User C端用户表
type User struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Nickname  string     `gorm:"column:nickname;type:varchar(64);not null;default:'';comment:昵称" json:"nickname"`
	Avatar    string     `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像" json:"avatar"`
	Mobile    string     `gorm:"column:mobile;type:varchar(20);not null;uniqueIndex;comment:手机号" json:"mobile"`
	Password  string     `gorm:"column:password;type:varchar(128);not null;default:'';comment:密码" json:"password"`
	OpenID    string     `gorm:"column:open_id;type:varchar(64);not null;default:'';index;comment:微信OpenID" json:"open_id"`
	UnionID   string     `gorm:"column:union_id;type:varchar(64);not null;default:'';index;comment:微信UnionID" json:"union_id"`
	Gender    int32      `gorm:"column:gender;not null;default:0;comment:性别 0未知 1男 2女" json:"gender"`
	Status    int32      `gorm:"column:status;not null;default:1;comment:状态 1正常 -1禁用" json:"status"`
	CreateAt  time.Time  `gorm:"column:create_at;autoCreateTime" json:"create_at"`
	UpdateAt  time.Time  `gorm:"column:update_at;autoUpdateTime" json:"update_at"`
	LastLogin *time.Time `gorm:"column:last_login" json:"last_login"`
}

func (*User) TableName() string {
	return TableNameUser
}
