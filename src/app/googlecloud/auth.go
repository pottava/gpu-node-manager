package googlecloud

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/revel/revel"
	"google.golang.org/api/option"
)

func clientOption() option.ClientOption {
	co := option.WithGRPCConnectionPool(10)
	if _, err := os.Stat("key.json"); !errors.Is(err, os.ErrNotExist) {
		co = option.WithCredentialsFile("key.json")
	}
	return co
}

func Auth(r *revel.Request) (*auth.Token, error) {
	tokens := strings.Split(r.Header.Get("Authorization"), " ")
	if len(tokens) != 2 {
		return nil, errors.New("invalid token")
	}
	client, err := authClient(r)
	if err != nil {
		return nil, err
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

func authClient(r *revel.Request) (*auth.Client, error) {
	app, err := firebaseApp(r)
	if err != nil {
		return nil, err
	}
	return app.Auth(r.Context())
}

func firebaseApp(r *revel.Request) (*firebase.App, error) {
	if _, err := os.Stat("key.json"); !errors.Is(err, os.ErrNotExist) {
		return firebase.NewApp(r.Context(), nil, option.WithCredentialsFile("key.json"))
	}
	return firebase.NewApp(r.Context(), nil)
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

func PasswordResetLink(r *revel.Request, email string) (string, error) {
	client, err := authClient(r)
	if err != nil {
		return "", err
	}
	return client.PasswordResetLink(r.Context(), email)
}
