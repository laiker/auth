package model

import (
	"database/sql"
	"time"
)

// UserInfo
type UserInfo struct {
	Name     string
	Email    string
	Role     string
	Password string
}

type User struct {
	Id        int64
	Name      string
	Email     string
	Role      string
	UpdatedAt sql.NullTime
	CreatedAt time.Time
}
