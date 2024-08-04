package middleware

import (
	"fmt"
	"net/http"
)

const (
	defaultRealm = "Authorization Required Area"
)

// BasicAuthInfo はBasic認証でサーバ側が保持する情報を表す。
type BasicAuthInfo struct {
	userID   string
	password string
	realm    string
}

// NewBasicAuthInfo は、妥当性が保証された BasicAuthInfo を返す。
//
// レルムにはデフォルト値が指定される。
func NewBasicAuthInfo(userID, password string) (*BasicAuthInfo, error) {
	return NewBasicAuthInfoWithRealm(userID, password, defaultRealm)
}

// NewBasicAuthInfoWithRealm は、レルムを指定した BasicAuthInfo を返す。
func NewBasicAuthInfoWithRealm(userID, password, realm string) (*BasicAuthInfo, error) {
	cred := &BasicAuthInfo{
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
	cred BasicAuthInfo
}

// NewBasicAuthMiddleware は、Basic認証によるアクセス制限を行うミドルウェアを返す。
func NewBasicAuthMiddleware(cred BasicAuthInfo) *basicAuthMiddleware {
	return &basicAuthMiddleware{
		cred: cred,
	}
}

// NOTE: サーバ起動失敗時に表示される情報であり、エラー詳細を含んでもユーザに見えないため問題ない。
func (cred *BasicAuthInfo) validate() error {
	if cred.userID == "" || cred.password == "" {
		return fmt.Errorf("Basic認証のユーザID/パスワードは、空文字以外を指定する必要があります")
	}
	return nil
}

func (cred *BasicAuthInfo) authenticate(r *http.Request) error {
	uid, passwd, ok := r.BasicAuth()
	if !ok {
		return fmt.Errorf("ユーザからの認証情報が取得できません")
	}
	if cred.userID != uid || cred.password != passwd {
		return fmt.Errorf("認証に失敗しました")
	}
	return nil
}

// ServeNext は、Basic認証によるアクセス制限を行う。
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
