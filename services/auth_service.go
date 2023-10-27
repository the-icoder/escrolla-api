package services

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"escrolla-api/config"
	"escrolla-api/db"
	apiError "escrolla-api/errors"
	"escrolla-api/models"
	"escrolla-api/services/jwt"
	"fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

type AuthService interface {
	LoginUser(request *models.LoginRequest) (*models.LoginResponse, *apiError.Error)
	SignupUser(request *models.User) (*models.User, *apiError.Error)
	FacebookSignInUser(token string) (*string, *apiError.Error)
	VerifyEmail(token string) error
	SendEmailForPasswordReset(user *models.ForgotPassword) *apiError.Error
	ResetPassword(user *models.ResetPassword, token string) *apiError.Error
	DeleteUserByEmail(userEmail string) *apiError.Error
}

// authService struct
type authService struct {
	Config   *config.Config
	authRepo db.AuthRepository
	mail     Mailer
}

// NewAuthService instantiate an authService
func NewAuthService(authRepo db.AuthRepository, conf *config.Config, mailer Mailer) AuthService {
	return &authService{
		Config:   conf,
		authRepo: authRepo,
		mail:     mailer,
	}
}

func (a *authService) SignupUser(user *models.User) (*models.User, *apiError.Error) {
	err := a.authRepo.IsEmailExist(user.Email)
	if err != nil {
		return nil, apiError.New("email already exist", http.StatusBadRequest)
	}
	err = a.authRepo.IsPhoneExist(user.Phone)
	if err != nil {
		return nil, apiError.New("phone already exist", http.StatusBadRequest)
	}
	//if user.ReferralCode != "" {
	//	err = a.authRepo.IsReferralCodeExist(user.ReferralCode)
	//	if err != nil {
	//		return nil, apiError.New("invalid referral code", http.StatusBadRequest)
	//	}
	//}
	user.PasswordHash, err = GenerateHashPassword(user.Password)
	if err != nil {
		log.Printf("error generating password hash: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}

	token, err := jwt.GenerateToken(user.Email, a.Config.JWTSecret, jwt.ConfirmEmailValidity)
	if err != nil {
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	go func() {
		err = a.sendVerifyEmail(token, user.Email)
		if err != nil {
			log.Printf("unable to send email: %v", err)
			return
		}
	}()

	user.Password = ""
	user.EmailConfirmed = false
	user, err = a.authRepo.CreateUser(user)

	if err != nil {
		log.Printf("unable to create user: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}

	return user, nil
}

func (a *authService) sendVerifyEmail(token, email string) *apiError.Error {
	link := fmt.Sprintf("%s/verifyEmail/%s", a.Config.BaseUrl, token)
	value := map[string]interface{}{}
	value["link"] = link
	subject := "Verify your email"
	body := "Please Click the link below to verify your email"
	templateName := "emailverification"
	err := a.mail.SendMail(email, subject, body, templateName, value)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return apiError.New("Internal server error", http.StatusInternalServerError)
	}
	return nil
}

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (a *authService) LoginUser(loginRequest *models.LoginRequest) (*models.LoginResponse, *apiError.Error) {
	foundUser, err := a.authRepo.FindUserByEmail(loginRequest.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apiError.New("invalid email", http.StatusUnprocessableEntity)
		} else {
			log.Printf("error from database: %v", err)
			return nil, apiError.ErrInternalServerError
		}
	}

	if foundUser.EmailConfirmed == false {
		return nil, apiError.New("email not verified", http.StatusUnauthorized)
	}

	if err := foundUser.VerifyPassword(loginRequest.Password); err != nil {
		return nil, apiError.ErrInvalidPassword
	}

	accessToken, err := jwt.GenerateToken(foundUser.Email, a.Config.JWTSecret, jwt.AccessTokenValidity)
	if err != nil {
		log.Printf("error generating token %s", err)
		return nil, apiError.ErrInternalServerError
	}

	return foundUser.LoginUserToDto(accessToken), nil
}

func (a *authService) VerifyEmail(token string) error {
	claims, err := jwt.ValidateAndGetClaims(token, a.Config.JWTSecret)
	if err != nil {
		return apiError.New("invalid link", http.StatusUnauthorized)
	}
	email := claims["email"].(string)
	err = a.authRepo.VerifyEmail(email, token)
	return err
}

func (a *authService) GoogleSignInUser(token string) (*string, *apiError.Error) {

	googleUserDetails, googleUserDetailsError := GetUserInfoFromGoogle(token)

	if googleUserDetailsError != nil {
		return nil, apiError.New(fmt.Sprintf("unable to get user details from google: %v", googleUserDetailsError), http.StatusUnauthorized)
	}

	authToken, authTokenError := a.GetGoogleSignInToken(googleUserDetails)

	if authTokenError != nil {
		return nil, apiError.New(fmt.Sprintf("unable sign in user: %v", authTokenError), http.StatusUnauthorized)
	}
	return &authToken, nil
}

// GetUserInfoFromGoogle will return information of user which is fetched from Google
func GetUserInfoFromGoogle(token string) (*models.GoogleUser, error) {
	var googleUserDetails *models.GoogleUser

	url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token
	googleUserDetailsRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", err)
	}

	googleUserDetailsResponse, googleDetailsResponseError := http.DefaultClient.Do(googleUserDetailsRequest)
	if googleDetailsResponseError != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", googleDetailsResponseError)
	}

	body, err := ioutil.ReadAll(googleUserDetailsResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", err)
	}
	defer googleUserDetailsResponse.Body.Close()

	err = json.Unmarshal(body, &googleUserDetails)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", err)
	}

	return googleUserDetails, nil
}

func (a *authService) FacebookSignInUser(token string) (*string, *apiError.Error) {
	// rename function
	fbUserDetails, fbUserDetailsError := GetUserInfoFromFacebook(token)

	if fbUserDetailsError != nil {
		return nil, apiError.New(fmt.Sprintf("unable to get user details from facebook: %v", fbUserDetailsError), http.StatusUnauthorized)
	}

	authToken, authTokenError := a.GetFacebookSignInToken(fbUserDetails)
	if authTokenError != nil {
		return nil, apiError.New(fmt.Sprintf("unable sign in user: %v", authTokenError), http.StatusUnauthorized)
	}
	return &authToken, nil
}

// GetUserInfoFromFacebook will return information of user which is fetched from facebook
func GetUserInfoFromFacebook(token string) (*models.FacebookUser, error) {
	var fbUserDetails *models.FacebookUser
	facebookUserDetailsRequest, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=name,email&access_token="+token, nil)
	facebookUserDetailsResponse, facebookUserDetailsResponseError := http.DefaultClient.Do(facebookUserDetailsRequest)

	if facebookUserDetailsResponseError != nil {
		return nil, fmt.Errorf("error occurred while getting information from Facebook: %+v", facebookUserDetailsResponseError)
	}
	body, err := ioutil.ReadAll(facebookUserDetailsResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Facebook: %+v", err)
	}
	defer facebookUserDetailsResponse.Body.Close()
	err = json.Unmarshal(body, &fbUserDetails)

	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Facebook: %+v", err)
	}

	return fbUserDetails, nil
}

// GetGoogleSignInToken Used for Signing In the Users
func (a *authService) GetGoogleSignInToken(googleUserDetails *models.GoogleUser) (string, error) {
	var result *models.User

	if googleUserDetails == nil {
		return "", fmt.Errorf("error: google user details can't be empty")
	}

	if googleUserDetails.Email == "" {
		return "", fmt.Errorf("error: email can't be empty")
	}

	if googleUserDetails.Name == "" {
		return "", fmt.Errorf("error: name can't be empty")
	}

	result, err := a.authRepo.FindUserByEmail(googleUserDetails.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("error finding user: %+v", err)
	}

	if result == nil {
		result = &models.User{}
		result.Email = googleUserDetails.Email
		result.FirstName = googleUserDetails.Name
		result.EmailConfirmed = true
		_, err = a.authRepo.CreateUser(result)
		if err != nil {
			return "", fmt.Errorf("error occurred creating user: %+v", err)
		}
	}

	tokenString, err := jwt.GenerateToken(googleUserDetails.Email, a.Config.JWTSecret, jwt.AccessTokenValidity)

	if tokenString == "" {
		return "", fmt.Errorf("unable to generate Auth token: %+v", err)
	}

	return tokenString, nil
}

// GetFacebookSignInToken Used for Signing In the Users
func (a *authService) GetFacebookSignInToken(facebookUserDetails *models.FacebookUser) (string, error) {
	var result *models.User

	if facebookUserDetails == nil {
		return "", fmt.Errorf("error: facebook user details can't be empty")
	}

	if facebookUserDetails.Email == "" {
		return "", fmt.Errorf("error: email can't be empty")
	}

	if facebookUserDetails.Name == "" {
		return "", fmt.Errorf("error: name can't be empty")
	}

	result, err := a.authRepo.FindUserByEmail(facebookUserDetails.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("error finding user: %+v", err)
	}

	if result == nil {
		result = &models.User{}
		result.Email = facebookUserDetails.Email
		result.FirstName = facebookUserDetails.Name
		result.EmailConfirmed = true
		_, err = a.authRepo.CreateUser(result)
		if err != nil {
			return "", fmt.Errorf("error occurred creating user: %+v", err)
		}
	}

	tokenString, err := jwt.GenerateToken(result.Email, a.Config.JWTSecret, jwt.AccessTokenValidity)

	if tokenString == "" {
		return "", fmt.Errorf("unable to generate Auth token: %+v", err)
	}

	return tokenString, nil
}

func (a *authService) DeleteUserByEmail(userEmail string) *apiError.Error {
	err := a.authRepo.DeleteUserByEmail(userEmail)
	if err != nil {
		return apiError.ErrInternalServerError
	}
	return nil
}

func GenerateRandomString() (string, error) {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%X", b)
	return s, nil
}
