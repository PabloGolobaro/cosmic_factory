package record

import "time"

type UserRecord struct {
	UUID         string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
