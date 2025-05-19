package model

import (
	"database/sql"
	"time"
)

// UserInfo
type UserInfo struct {
	Name     string
	Email    string
	Role     int
	Password string
}

type User struct {
	Id        int64
	Name      string
	Email     string
	Role      string
	Password  string
	UpdatedAt sql.NullTime
	CreatedAt time.Time
}

type UserName struct {
	Id   int64
	Name string
}
