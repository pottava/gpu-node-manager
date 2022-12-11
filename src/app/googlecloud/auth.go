package googlecloud

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/revel/revel"
	"google.golang.org/api/option"
)

var (
	ProjectID string
	Location  = "asia-northeast1"
	Zone      = "asia-northeast1-c"
)

func init() {
	if candidate, found := os.LookupEnv("GOOGLE_CLOUD_PROJECT"); found {
		ProjectID = candidate
	}
	if candidate, found := os.LookupEnv("GOOGLE_CLOUD_LOCATION"); found {
		Location = candidate
	}
}

func clientOption() option.ClientOption {
	co := option.WithHTTPClient(http.DefaultClient)
	if _, err := os.Stat("key.json"); !errors.Is(err, os.ErrNotExist) {
		co = option.WithCredentialsFile("key.json")
	}
	return co
}

func Auth(r *revel.Request) (*auth.Token, error) {
	app, err := firebase.NewApp(r.Context(), nil, clientOption())
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(r.Context())
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(r.Header.Get("Authorization"), " ")
	if len(tokens) != 2 {
		return nil, errors.New("invalid token")
	}
	token, err := client.VerifyIDToken(r.Context(), strings.TrimSpace(tokens[1]))
	if err != nil {
		return nil, err
	}
	if token.Firebase.SignInProvider != "password" {
		return nil, errors.New("unknown provider")
	}
	return token, nil
}

func VerifiedEmail(token *auth.Token) string {
	emails := []string{}
	bytes, err := json.Marshal(token.Firebase.Identities["email"])
	if err != nil {
		return ""
	}
	if err := json.Unmarshal(bytes, &emails); err != nil {
		return ""
	}
	return emails[0]
}
