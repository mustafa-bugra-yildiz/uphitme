package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/mustafa-bugra-yildiz/uphitme/env"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const cookieName = "session"

type Auth struct{ userRepo user.Repo }
type Claims struct{ jwt.RegisteredClaims }

func New(userRepo user.Repo) *Auth {
	return &Auth{userRepo: userRepo}
}

func (a *Auth) Verify(r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		cookie.Value,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return env.JWTsecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claims, nil
}

func (a *Auth) Login(w http.ResponseWriter, u user.User) error {
	now := jwt.NumericDate{Time: time.Now()}
	expiresAt := jwt.NumericDate{Time: time.Now().Add(time.Hour)}
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "uphit.me",
			Subject:   u.ID.String(),
			Audience:  nil,
			ExpiresAt: &expiresAt,
			NotBefore: &now,
			IssuedAt:  &now,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(env.JWTsecret)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   env.Env == env.Prod,
		Expires:  expiresAt.Time,
	})

	return nil
}

func (a *Auth) Logout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   env.Env == env.Prod,
		Expires:  time.Unix(0, 0),
	})
}
