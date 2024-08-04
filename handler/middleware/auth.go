package middleware

import (
	"fmt"
	"net/http"
)

const (
	defaultRealm = "Authorization Required Area"
)

// BasicAuthCredential はBasic認証の認証情報を表す。
type BasicAuthCredential struct {
	userID   string
	password string
	realm    string
}

// NewBasicAuthCredential は、妥当性が保証された BasicAuthCredential を返す。
//
// レルムにはデフォルト値が指定される。
func NewBasicAuthCredential(userID, password string) (*BasicAuthCredential, error) {
	return NewBasicAuthCredentialWithRealm(userID, password, defaultRealm)
}

// NewBasicAuthCredentialWithRealm は、レルムを指定した BasicAuthCredential を返す。
func NewBasicAuthCredentialWithRealm(userID, password, realm string) (*BasicAuthCredential, error) {
	cred := &BasicAuthCredential{
		userID:   userID,
		password: password,
		realm:    realm,
	}
	if err := cred.validate(); err != nil {
		return nil, err
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

func (cred *BasicAuthCredential) validate() error {
	if cred.userID == "" || cred.password == "" {
		return fmt.Errorf("与えられた認証情報は、Basic認証として不適切です")
	}
	return nil
}

func (cred *BasicAuthCredential) authenticate(r *http.Request) error {
	uid, passwd, ok := r.BasicAuth()
	if !ok {
		return fmt.Errorf("ユーザからの認証情報が取得できません")
	}
	if cred.userID != uid || cred.password != passwd {
		return fmt.Errorf("認証に失敗しました")
	}
	return nil
}

// ServeNext は、 h の前後で取得した情報を元に、 アクセスログを記録する。
func (m *basicAuthMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := m.cred.authenticate(r); err != nil {
			w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, m.cred.realm))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
