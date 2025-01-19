package converter

import (
	"github.com/laiker/auth/internal/model"
	"github.com/laiker/auth/pkg/user_v1"
)

func ToUserFromCreateRequest(user *user_v1.CreateRequest) *model.UserInfo {
	return &model.UserInfo{
		Name:     user.GetName(),
		Email:    user.GetEmail(),
		Role:     user.GetRole().String(),
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
