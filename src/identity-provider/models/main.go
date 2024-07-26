//package models
//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"identity-provider/objects"
//	"identity-provider/utils"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"strconv"
//	"strings"
//)
//
//
//func (model *AuthM) Create(request *objects.UserCreateRequest) error {
//	request.GroupIds = []string{utils.Config.Okta.ClientGroup}
//	request.Profile.UserType = utils.User.String()
//
//	req_body, err := json.Marshal(request)
//	if err != nil {
//		return err
//	}
//	req, _ := http.NewRequest(
//		"POST",
//		fmt.Sprintf("%s/api/v1/users/?activate=true", utils.Config.Okta.Endpoint),
//		bytes.NewBuffer(req_body),
//	)
//	req.Header.Add("Content-Type", "application/json")
//	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", utils.Config.Okta.SSWSToken))
//
//	resp, err := model.client.Do(req)
//	if err != nil {
//		return err
//	} else if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("auth failed with code %d", resp.StatusCode)
//	} else {
//		return nil
//	}
//}
//
//func (model *AuthM) GetAuthCode() string {
//	authRequest := url.Values{}
//	authRequest.Set("scope", "openid")
//	authRequest.Set("response_type", "code")
//	authRequest.Set("response_mode", "query")
//	authRequest.Set("redirect_uri", "http://localhost:8040/api/v1/callback")
//	authRequest.Set("client_id", utils.Config.Okta.ClientId)
//	authRequest.Set("state", "ikjs0la6d8b")
//	authRequest.Set("nonce", "f72eosg89cs")
//	encodedData := authRequest.Encode()
//
//	req, _ := http.NewRequest(
//		"GET",
//		fmt.Sprintf("%s/oauth2/default/v1/authorize", utils.Config.Okta.Endpoint),
//		strings.NewReader(encodedData),
//	)
//
//	resp, err := model.client.Do(req)
//	if err != nil {
//		return ""
//	}
//	defer resp.Body.Close()
//
//	return ""
//}
//
//func (model *AuthM) Auth(login string, password string) (*objects.AuthResponse, error) {
//
//	authRequest := url.Values{}
//	authRequest.Set("scope", "openid")
//	authRequest.Set("grant_type", "password")
//	authRequest.Set("login", login)
//	authRequest.Set("password", password)
//	authRequest.Set("client_id", utils.Config.Okta.ClientId)
//	authRequest.Set("client_secret", utils.Config.Okta.ClientSecret)
//	encodedData := authRequest.Encode()
//
//	req, _ := http.NewRequest(
//		"POST",
//		fmt.Sprintf("%s/oauth2/default/v1/token", utils.Config.Okta.Endpoint),
//		strings.NewReader(encodedData),
//	)
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Content-Length", strconv.Itoa(len(authRequest.Encode())))
//
//	resp, err := model.client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	data := &objects.AuthResponse{}
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("auth failed, code: %d", resp.StatusCode)
//	}
//
//	err = json.Unmarshal(body, data)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}

package models

import (
	"identity-provider/objects"
	"identity-provider/repositories"

	"github.com/jinzhu/gorm"
)

type Models struct {
	User *UserModel
}

func InitModels(db *gorm.DB) *Models {
	models := new(Models)
	models.User = NewUserModel(repositories.NewPostgresUserRepository(db))

	return models
}

type UserModel struct {
	repository repositories.UserRepository
}

func NewUserModel(repository repositories.UserRepository) *UserModel {
	return &UserModel{
		repository: repository,
	}
}

func (model *UserModel) RegisterUser(firstName, lastName, login, password, email string) (*objects.User, error) {
	user := &objects.User{
		FirstName: firstName,
		LastName:  lastName,
		Login:     login,
		Password:  password,
		Email:     email,
	}

	createdUser, err := model.repository.Create(user)

	return createdUser, err
}

func (model *UserModel) GetUser(login string) (*objects.User, error) {
	user, err := model.repository.Find(login)
	if err != nil {
		return nil, err
	}

	return user, nil
}
