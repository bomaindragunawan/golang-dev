package helper

import "go-login-api/models"

// ToUser mengonversi RegisterRequest menjadi User dengan password yang sudah di-hash
func ToUser(req models.RegisterRequest, hashedPassword string) models.User {
	return models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}
}

// ToUserResponse mengonversi User menjadi UserResponse
func ToUserResponse(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
