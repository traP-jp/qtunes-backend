package router

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/gorilla/sessions"
	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
)

var baseURL, _ = url.Parse("https://q.trap.jp/api/v3")

// AuthResponse 認証の返答
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// PkceParams PKCE用のcode_challengeなど
type PkceParams struct {
	CodeChallenge       string `json:"codeChallenge,omitempty"`
	CodeChallengeMethod string `json:"codeChallengeMethod,omitempty"`
	ClientID            string `json:"clientID,omitempty"`
	ResponseType        string `json:"responseType,omitempty"`
}

// CallbackHandler GET /oauth/callbackのハンドラー
func CallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	if len(code) == 0 {
		return c.String(http.StatusBadRequest, "Code Is Null")
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err).Error())
	}

	codeVerifier := sess.Values["codeVerifier"].(string)
	res, err := getAccessToken(code, codeVerifier)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get access token: %w", err).Error())
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   res.ExpiresIn * 1000,
		HttpOnly: true,
	}
	sess.Values["accessToken"] = res.AccessToken
	sess.Values["refreshToken"] = res.RefreshToken
	user, err := getMe(res.AccessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Me: %w", err).Error())
	}

	err = model.CreateUser(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert user: %w", err))
	}

	sess.Values["id"] = user.ID
	sess.Values["name"] = user.Name

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 1000,
		HttpOnly: true,
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to save session: %w", err).Error())
	}

	return c.NoContent(http.StatusOK)
}

// PostLogoutHandler POST /oauth/logoutのハンドラー
func PostLogoutHandler(c echo.Context) error {
	accessToken := c.Get("accessToken").(string)

	path := *baseURL
	path.Path += "/oauth2/revoke"
	form := url.Values{}
	form.Set("token", accessToken)
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", path.String(), reqBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Making HTTP Request:%w", err).Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In HTTP Request:%w", err).Error())
	}
	if res.StatusCode != 200 {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Access Token:(Status:%d %s)", res.StatusCode, res.Status).Error())
	}

	err = s.RevokeSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to revoke session: %w", err))
	}

	return c.NoContent(http.StatusOK)
}

// PostGenerateCodeHandler POST /oauth/generate/codeのハンドラー
func PostGenerateCodeHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err).Error())
	}

	pkceParams := PkceParams{}

	pkceParams.ResponseType = "code"

	pkceParams.ClientID = clientID

	bytesCodeVerifier, err := randBytes(43)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to generate random bytes: %w", err).Error())
	}

	codeVerifier := string(bytesCodeVerifier)
	bytesCodeChallenge := sha256.Sum256(bytesCodeVerifier[:])
	codeChallenge := base64url.Encode(bytesCodeChallenge[:])
	pkceParams.CodeChallenge = codeChallenge
	sess.Values["codeVerifier"] = codeVerifier

	pkceParams.CodeChallengeMethod = "S256"

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 1000,
		HttpOnly: true,
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, pkceParams)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	max := new(big.Int)

	max.SetInt64(int64(len(letterBytes)))
	for i := range buf {
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random integer: %w", err)
		}

		buf[i] = letterBytes[r.Int64()]
	}

	return buf, nil
}

func getAccessToken(code string, codeVerifier string) (AuthResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", clientID)
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)
	reqBody := strings.NewReader(form.Encode())
	path := *baseURL
	path.Path += "/oauth2/token"
	req, err := http.NewRequest("POST", path.String(), reqBody)
	if err != nil {
		return AuthResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return AuthResponse{}, err
	}
	if res.StatusCode != 200 {
		return AuthResponse{}, fmt.Errorf("Failed In Getting Access Token:(Status:%d %s)", res.StatusCode, res.Status)
	}
	var authRes AuthResponse
	err = json.NewDecoder(res.Body).Decode(&authRes)
	if err != nil {
		return AuthResponse{}, err
	}
	return authRes, nil
}

func getMe(accessToken string) (*model.User, error) {
	path := *baseURL
	path.Path += "/users/me"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	var user model.User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
