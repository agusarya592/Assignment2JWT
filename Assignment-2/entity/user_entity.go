package entity

type User struct {
	UserID   uint64
	Email    string
	Username string
	Password string
}

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ur *UserRegistrationRequest) ToEntity() (u User) {
	u = User{
		Username: ur.Username,
		Email:    ur.Email,
		Password: ur.Password,
	}
	return
}

func (ul *UserLoginRequest) ToEntity() (us User) {
	us = User{
		Username: ul.Username,
		Password: ul.Password,
	}
	return
}

type UserLoginResponse struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

func NewUserLoginResponse(user *User, ac string) (res *UserLoginResponse, err error) {
	res = &UserLoginResponse{
		Username:    user.Username,
		AccessToken: ac,
	}
	return
}
