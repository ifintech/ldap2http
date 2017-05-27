package main

import (
	"testing"
)

func TestAuthRequest(t *testing.T)  {
	auth_url = "http://php7.i.hrbbwx.com/api"
	authResult := authRequest("xindongzhe@yahoo.com", "123456")

	if !authResult {
		t.Errorf("error, %s", authResult);
	}
}

func TestQueryUserInfoRequest(t *testing.T)  {
	auth_url = "http://php7.i.hrbbwx.com/api"
	authResult := queryUserInfoRequest("xindongzhe@yahoo.com")

	if nil == authResult {
		t.Errorf("error, %s", authResult);
	}
}