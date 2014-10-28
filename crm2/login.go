package crm2

import (
	"io/ioutil"
	"net/url"
	"strings"
)

// 同一个 HttpClient 不需要手动设置Cookie
func (self *Crm2) Login() error {
	err := self.SSO()
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Add("userName", self.User.Username())
	pwd, ok := self.User.Password()
	if !ok {
		return errPwd
	}
	v.Add("userPass", pwd)
	v.Add("telSystem", "00") //no soft phone
	v.Add("userTel", "")
	v.Add("empId", "")

	resp, err := self.HttpClient.PostForm(URL_SSO_C, v)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	cookie := resp.Cookies()
	if !strings.Contains(string(body), "【HOME】") {
		//fmt.Println(string(body))
		return errLogin
	}
	if len(cookie) == 0 {
		return errCookie
	}
	self.Cookie_SSO_C = cookie
	return nil
}

func (self *Crm2) SSO() error {
	if self.Cookie_SSO_A == nil {
		if err := self.SSO_A(); err != nil {
			return err
		}
	}
	if self.Cookie_SSO_B == nil {
		if err := self.SSO_B(); err != nil {
			return err
		}
	}
	return nil
}

// 获取SSO_A cookie
func (self *Crm2) SSO_A() error {
	v := url.Values{}
	v.Set("optiontype", "login")
	v.Set("userName_erp", strings.ToLower(self.User.Username()))
	pwd, ok := self.User.Password()
	if !ok {
		return errPwd
	}
	v.Set("userPassword_erp", pwd)
	buf := self.GetBuf()
	defer self.PutBuf(buf)
	buf.WriteString(URL_SSO_A)
	buf.WriteString("/newhrm/ssoResponse.aspx?")
	buf.WriteString(v.Encode())
	resp, err := self.HttpClient.Get(buf.String())
	if err != nil {
		return err
	}
	buf.Reset()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	cookie := resp.Cookies()

	if buf.String() != "('true')" {
		return errLogin
	}
	if len(cookie) == 0 {
		return errCookie
	}
	self.Cookie_SSO_A = cookie
	return nil
}

// 获取SSO_B cookie
func (self *Crm2) SSO_B() error {
	v := url.Values{}
	v.Set("optiontype", "login")
	v.Set("userName_erp", strings.ToLower(self.User.Username()))
	pwd, ok := self.User.Password()
	if !ok {
		return errPwd
	}
	v.Set("userPassword_erp", pwd)
	buf := self.GetBuf()
	defer self.PutBuf(buf)
	buf.WriteString(URL_SSO_B)
	buf.WriteString("/newhrm/ssoResponse.aspx?")
	buf.WriteString(v.Encode())
	resp, err := self.HttpClient.Get(buf.String())
	if err != nil {
		return err
	}
	buf.Reset()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	cookie := resp.Cookies()
	if buf.String() != "('true')" {
		return errLogin
	}
	if len(cookie) == 0 {
		return errCookie
	}
	self.Cookie_SSO_B = cookie
	return nil
}
