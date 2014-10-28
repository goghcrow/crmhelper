package crm2

import (
	"net/url"
	"strings"
)

func (self *Crm2) CaseInfo(id string) (*string, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("caseType", "")
	v.Add("callTel", "")
	v.Add("applyId", "")
	v.Add("source", "")
	v.Add("clueVal", "")
	v.Add("turnCallFlag", "V2.0")
	v.Add("V", "")

	buf := self.GetBuf()
	defer self.PutBuf(buf)

	buf.WriteString(URL_SSO_C)
	buf.WriteString("/caseInfo?")
	buf.WriteString(v.Encode())
	resp, err := self.HttpClient.Get(buf.String())
	if err != nil {
		return nil, err
	}
	buf.Reset()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := buf.String()
	return &ret, nil
}

func (self *Crm2) FilterTel(html string) (tel *string) {
	var prefix string
	prefix1 := `name="callTel" class="itxt1-12" value="`
	prefix2 := `name="callTel" class="itxt1-12  itxtnone" value=`
	postfix := `" type="text" maxlength="50" id="callTel"`
	preStart := strings.Index(html, prefix1)
	if preStart > 0 {
		prefix = prefix1
	} else {
		preStart = strings.Index(html, prefix2)
		if preStart > 0 {
			prefix = prefix2
		}
	}

	postStart := strings.Index(html, postfix)

	if preStart != -1 && postStart != -1 {
		testr := string([]byte(html)[(preStart + len(prefix)):postStart])
		tel = &testr
	}

	return
}
