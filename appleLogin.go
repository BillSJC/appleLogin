package appleLogin

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/url"
	"time"
)

// AppleConfig Main struct of the package
type AppleConfig struct {
	TeamID      string      //Your Apple Team ID
	ClientID    string      //Your Service which enable sign-in-with-apple service
	RedirectURI string      //Your RedirectURI config in apple website
	KeyID       string      //Your Secret Key ID
	AESCert     interface{} //Your Secret Key Created By X509 package
}

// AppleAuthToken main response of apple REST-API
type AppleAuthToken struct {
	AccessToken  string `json:"access_token"`  //AccessToken
	ExpiresIn    int64  `json:"expires_in"`    //Expires in
	IDToken      string `json:"id_token"`      //ID token
	RefreshToken string `json:"refresh_token"` //RF token
	TokenType    string `json:"token_type"`    //Token Type
}

const AppleAuthURL = "https://appleid.apple.com/auth/token" //the auth URL of apple
const AppleGrantType = "authorization_code"                 //the grant type of apple auth

//LoadP8CertByByte use x509.ParsePKCS8PrivateKey to Parse cert file
func (a *AppleConfig) LoadP8CertByByte(str []byte) (err error) {
	block, _ := pem.Decode([]byte(str))
	cert, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	a.AESCert = cert
	return nil
}

//LoadP8CertByFile load file and Parse it
func (a *AppleConfig) LoadP8CertByFile(path string) (err error) {
	b, err := ioutil.ReadFile("cert")
	if err != nil {
		return
	}
	return a.LoadP8CertByByte([]byte(b))
}

//InitAppleConfig init a new Client of this pkg
func InitAppleConfig(teamID string, clientID string, redirectURI string, keyID string) *AppleConfig {
	return &AppleConfig{
		teamID,
		clientID,
		redirectURI,
		keyID,
		nil,
	}
}

//CreateCallbackURL create a callback URL for frontend
func (a *AppleConfig) CreateCallbackURL(state string) string {
	u := url.Values{}
	u.Add("response_type", "code")
	u.Add("redirect_uri", a.RedirectURI)
	u.Add("client_id", a.ClientID)
	u.Add("state", state)
	u.Add("scope", "name email")
	return "https://appleid.apple.com/auth/authorize?" + u.Encode()
}

//input your code and expire-time to get AccessToken of apple
func (a *AppleConfig) GetAppleToken(code string, expireTime int64) (*AppleAuthToken, error) {
	//test cert
	if a.AESCert == nil {
		return nil, errors.New("missing cert")
	}
	//set jwt
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": a.TeamID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + expireTime,
		"aud": AppleAuthURL,
		"sub": a.ClientID,
	})
	//set JWT header
	token.Header = map[string]interface{}{
		"kid": a.KeyID,
		"alg": "ES256",
	}
	//make JWT sign
	tokenString, _ := token.SignedString(a.AESCert)
	v := url.Values{}
	v.Set("client_id", a.ClientID)
	v.Set("client_secret", tokenString)
	v.Set("code", code)
	v.Set("grant_type", "authorization_code")
	v.Set("redirect_uri", a.RedirectURI)
	fmt.Println(tokenString)
	vs := v.Encode()
	//send request
	resp, body, err2 := gorequest.New().Post("https://appleid.apple.com/auth/token").Type("urlencoded").Send(vs).End()
	if err2 != nil {
		return nil, fmt.Errorf(fmt.Sprint(err2))
	}
	//check response
	if resp.StatusCode != 200 {
		fmt.Println(body)
		panic(errors.New("post failed : resp code is not 200"))
	}
	t := new(AppleAuthToken)
	err := json.Unmarshal([]byte(body), t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
