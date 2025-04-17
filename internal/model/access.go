package model

type Permission struct {
	Id          int64  `json:"id" db:"permission_id"`
	Endpoint    string `json:"endpoint" db:"resource_name"`
	MinPriority int64  `json:"minPriority" db:"min_role_priority"`
}

type Role struct {
	Id       int64  `db:"role_id"`
	Name     string `db:"role_name"`
	Priority int64  `db:"priority"`
}
