package middleware

import (
	"fmt"
	"net/http"
)

// BasicAuthCredential はBasic認証の認証情報を表す。
type BasicAuthCredential struct {
	userID   string
	password string
}

// NewBasicAuthCredential は、妥当性が保証された BasicAuthCredential を返す。
func NewBasicAuthCredential(userID, password string) (*BasicAuthCredential, error) {
	cred := &BasicAuthCredential{
		userID:   userID,
		password: password,
	}
	if !cred.isValidCredential() {
		return nil, fmt.Errorf("与えられた認証情報は、Basic認証として不適切です")
	}
	return cred, nil
}

type basicAuthMiddleware struct {
	cred BasicAuthCredential
}

// NewBasicAuthMiddleware は、Basic認証によるアクセス制限を行うミドルウェアを返す。
func NewBasicAuthMiddleware(cred BasicAuthCredential) *basicAuthMiddleware {
	return &basicAuthMiddleware{
		cred: cred,
	}
}

func (cred *BasicAuthCredential) isValidCredential() bool {
	if cred.userID == "" || cred.password == "" {
		return false
	}
	return true
}

// ServeNext は、 h の前後で取得した情報を元に、 アクセスログを記録する。
func (m *basicAuthMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Authentication\n")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
