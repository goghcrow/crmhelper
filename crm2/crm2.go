package crm2

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	URL_SSO_A = "http://private"
	URL_SSO_B = "http://private"
	URL_SSO_C = "http://private"
)

var (
	errPwd    = errors.New("pwd needed")
	errCookie = errors.New("cookie is empty")
	errLogin  = errors.New("Login failed")
)

type Crm2 struct {
	bufferSize   int
	n            int
	chBuf        chan *bytes.Buffer
	User         *url.Userinfo
	Cookie_SSO_A []*http.Cookie
	Cookie_SSO_B []*http.Cookie
	Cookie_SSO_C []*http.Cookie
	HttpClient   *http.Client
}

func New(erp, pwd string, n, bufferSize int) *Crm2 {
	bufpool := make(chan *bytes.Buffer, n<<1)
	for i := 0; i < n; i++ {
		bufpool <- bytes.NewBuffer(make([]byte, 0, bufferSize))
	}

	cookieJar, _ := cookiejar.New(nil)
	return &Crm2{
		n:          n,
		bufferSize: bufferSize,
		chBuf:      bufpool,
		User:       url.UserPassword(erp, pwd),
		HttpClient: &http.Client{Jar: cookieJar},
	}
}

// 设置代理
func (self *Crm2) SetProxy(proxyUrl *url.URL) {
	self.HttpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}
