package converter

import (
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/pkg/user_v1"
)

func ToUserFromCreateRequest(user *user_v1.CreateRequest) *model.UserInfo {

	var role int

	switch user.Role.Number() {
	case 1:
		role = 2
		break
	default:
		role = 1
	}

	return &model.UserInfo{
		Name:     user.GetName(),
		Email:    user.GetEmail(),
		Role:     role,
		Password: user.GetPassword(),
	}
}

func ToUserFromUpdateRequest(user *user_v1.UpdateRequest) *model.User {
	return &model.User{
		Id:    user.GetId().GetValue(),
		Name:  user.GetName().String(),
		Email: user.GetEmail().String(),
	}
}
