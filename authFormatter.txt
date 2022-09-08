package formatter

import (
	"[project]/entity"
)

type UserFormatter struct {
	User entity.User `json:"profile"`
	Token   string         `json:"token"`
}

func FormatUser(user entity.User, token string) UserFormatter {
	formatter := UserFormatter{
		User: user,
		Token:   token,
	}
	return formatter
}
