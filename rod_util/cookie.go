package rod_util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type RodCookies []*proto.NetworkCookie

func (r RodCookies) ToHttpCookies() (cookies []*http.Cookie) {
	cookies = make([]*http.Cookie, 0, len(r))
	for _, obj := range r {
		cookies = append(cookies, &http.Cookie{
			Name:     obj.Name,
			Value:    obj.Value,
			Path:     obj.Path,
			Domain:   obj.Domain,
			Expires:  obj.Expires.Time(),
			Secure:   obj.Secure,
			HttpOnly: obj.HTTPOnly,
			SameSite: 0,
			Raw:      "",
			Unparsed: nil,
		})
	}
	return
}

func GetHttpCookies(page *rod.Page) (cookies []*http.Cookie, err error) {
	var rodCookies []*proto.NetworkCookie
	rodCookies, err = GetCookies(page)
	if err != nil {
		err = fmt.Errorf("GetCookies: %w", err)
		return
	}
	cookies = RodCookies(rodCookies).ToHttpCookies()

	return
}

func GetCookies(page *rod.Page) (cookies []*proto.NetworkCookie, err error) {
	cookies, err = page.Timeout(time.Second * 3).Cookies([]string{
		"https://e.qq.com",
		"https://sso.e.qq.com",
		"https://mp.weixin.qq.com",
		"https://ad.qq.com",
	})
	return
}

func GetCookieMap(page *rod.Page) (cookieMap map[string]*proto.NetworkCookie, err error) {
	cookies, err := GetCookies(page)
	if err != nil {
		err = fmt.Errorf("get cookies: %w", err)
		return
	}

	cookieMap = make(map[string]*proto.NetworkCookie)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie
	}
	return
}

func GetCookie(page *rod.Page, name string) (cookie *proto.NetworkCookie, err error) {
	var found bool
	var cookieMap map[string]*proto.NetworkCookie

	cookieMap, err = GetCookieMap(page)
	if err != nil {
		err = fmt.Errorf("48 r.GetCookieMap: %w", err)
		return
	}
	cookie, found = cookieMap[name]
	if !found {
		err = fmt.Errorf("not found cookie: %v", name)
		return
	}
	return
}

func CookieToParam(cookie *proto.NetworkCookie) (param *proto.NetworkCookieParam) {
	param = new(proto.NetworkCookieParam)
	param.Name = cookie.Name
	param.Value = cookie.Value
	param.Path = cookie.Path
	param.Domain = cookie.Domain
	param.SameSite = cookie.SameSite
	param.Secure = cookie.Secure
	param.Expires = cookie.Expires
	param.HTTPOnly = cookie.HTTPOnly
	return
}
