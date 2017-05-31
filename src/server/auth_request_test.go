package main

import (
	"testing"
)

func TestAuthRequest(t *testing.T)  {
	authInfos := []struct {
		email, password string
	}{
		{"xindongzhe@yahoo.com", "bonjour"},
		{"demo@yahoo.com", "123"},
	}

	auth_url = "http://php7.i.hrbbwx.com/api"

	for _, infoItem := range authInfos {
		authResult := authRequest(infoItem.email, infoItem.password)

		if !authResult {
			t.Errorf("error, %s", authResult);
		}
	}
}

func TestQueryUserInfoRequest(t *testing.T)  {
	emails := []struct {
		email string
	}{
		{"xindongzhe@yahoo.com"},
		{"demo@yahoo.com"},
	}

	auth_url = "http://php7.i.hrbbwx.com/api"
	auth_token = "token"

	for _, item := range emails {
		userInfo := queryUserInfoRequest(item.email)

		if nil == userInfo {
			t.Errorf("error, %s", userInfo);
		}
	}
}