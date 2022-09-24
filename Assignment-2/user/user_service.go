package user

import (
	"assignment2/config"
	"assignment2/entity"
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	NewUser(ctx context.Context, data *entity.UserRegistrationRequest) (userID int64, err error)
	LoginUser(ctx context.Context, data *entity.UserLoginRequest) (res *entity.UserLoginResponse, err error)
}

type authServiceImpl struct {
	repo AuthRepository
}

func ProvideAuthService(repo AuthRepository) *authServiceImpl {
	return &authServiceImpl{repo: repo}
}

func (a authServiceImpl) NewUser(ctx context.Context, data *entity.UserRegistrationRequest) (userID int64, err error) {
	userInfo := data.ToEntity()

	hashedPassword, err := encryptPassword(userInfo.Password)
	if err != nil {
		log.Printf("failed to encrypt the password")
		return 0, err
	}

	userInfo.Password = string(hashedPassword)
	userID, err = a.repo.StoreUser(ctx, userInfo)
	if err != nil {
		return
	}

	return
}

func (a authServiceImpl) LoginUser(ctx context.Context, data *entity.UserLoginRequest) (res *entity.UserLoginResponse, err error) {
	user := data.ToEntity()
	userInfo, err := a.repo.GetUserByEmail(ctx, user.Username)
	if err != nil {
		log.Printf("[LoginUser] failed to get the data by email, err => %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(user.Password))
	if err != nil {
		log.Printf("[LoginUser] failed to validate the password, err => %v", err)
		return nil, err
	}

	token, err := a.newAccessToken(userInfo)
	if err != nil {
		log.Printf("[LoginUser] failed to generate access token, err => %v", err)
		return nil, err
	}
	return entity.NewUserLoginResponse(userInfo, token)
}

func (a authServiceImpl) newAccessToken(user *entity.User) (string, error) {
	cfg := config.GetConfig()
	claims := a.newUserClaim(user.Email, user.Username, cfg.JWT_EXP)
	accessToken := jwt.NewWithClaims(cfg.JWT_SIGNING_METHOD, claims)

	signed, err := accessToken.SignedString([]byte(cfg.JWT_SECRET_KEY))
	if err != nil {
		log.Printf("[newAccessToken] failed to signed the token, err => %v", err)
		return "", err
	}

	return signed, err
}

func (a authServiceImpl) newUserClaim(email string, username string, exp time.Duration) *jwt.MapClaims {
	return &jwt.MapClaims{
		"exp": time.Now().Add(exp).Unix(),
		"data": map[string]interface{}{
			"email":    email,
			"username": username,
		},
	}
}

func encryptPassword(password string) (string, error) {
	bytepass := []byte(password)
	hashed, err := bcrypt.GenerateFromPassword(bytepass, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to generate hashed password, err => %v", err)
		return "", err
	}
	return string(hashed), nil
}
