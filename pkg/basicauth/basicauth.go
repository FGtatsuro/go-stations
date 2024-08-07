package basicauth

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

const defaultRealm = "Authorization Required Area"

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
	bai := &BasicAuthInfo{
		userID:   userID,
		password: password,
		realm:    realm,
	}
	if err := bai.validate(); err != nil {
		return nil, err
	}
	return bai, nil
}

// NOTE: サーバ起動失敗時に表示される情報であり、エラー詳細を含んでもユーザに見えないため問題ない。
func (bai *BasicAuthInfo) validate() error {
	if bai.userID == "" || bai.password == "" {
		return fmt.Errorf("Basic認証のユーザID/パスワードは、空文字以外を指定する必要があります")
	}
	if strings.Contains(bai.userID, ":") {
		return fmt.Errorf("Basic認証のユーザIDは、コロン(:)を含んではいけません")
	}
	if containsControl(bai.userID) || containsControl(bai.password) {
		return fmt.Errorf("Basic認証のユーザID/パスワードは、制御文字を含んではいけません")
	}
	return nil
}

// Authenticate は、ユーザから送られた情報を元に認証を実施する。
func (bai *BasicAuthInfo) Authenticate(r *http.Request) error {
	uid, passwd, ok := r.BasicAuth()
	if !ok {
		return fmt.Errorf("ユーザからの認証情報が取得できません")
	}
	if bai.userID != uid || bai.password != passwd {
		return fmt.Errorf("認証に失敗しました")
	}
	return nil
}

// Challenge は、Basic認証のチャレンジ文字列を返す。
func (bai *BasicAuthInfo) Challenge() string {
	return fmt.Sprintf(`Basic realm="%s"`, bai.realm)
}

func containsControl(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
