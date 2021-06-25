package contact

import (
	"encoding/json"
	"fmt"
	"github.com/timerzz/itchatgo/clients/base"
	"github.com/timerzz/itchatgo/enum"
	"github.com/timerzz/itchatgo/model"
	"github.com/timerzz/itchatgo/utils"
	"io/ioutil"
	"strconv"
	"time"
)

type Client struct {
	base.Client
}

func NewClient(base *base.Client) *Client {
	return &Client{
		*base,
	}
}

func (c *Client) GetAllContact() (contactMap map[string]*model.User, err error) {
	var getContact = func(seq int) (users []*model.User, reSeq int, err error) {
		urlMap := enum.InitParaEnum
		urlMap[enum.PassTicket] = c.LoginInfo.PassTicket
		urlMap[enum.R] = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
		urlMap["seq"] = strconv.Itoa(seq)
		urlMap[enum.SKey] = c.LoginInfo.BaseRequest.SKey

		url := fmt.Sprintf("%s/webwxgetcontact", c.LoginInfo.Url)
		resp, err := c.HttpClient.Get(url+utils.GetURLParams(urlMap), nil)
		if err != nil {
			return nil, 0, err
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}
		contactList := model.ContactList{}
		err = json.Unmarshal(bodyBytes, &contactList)
		if err != nil {
			return nil, 0, err
		}
		return contactList.MemberList, seq, err
	}
	contactMap = make(map[string]*model.User)
	var seq = 0
	for {
		var users []*model.User
		users, seq, err = getContact(seq)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			contactMap[u.UserName] = u
		}
		if seq == 0 {
			break
		}
	}
	return contactMap, nil
}
