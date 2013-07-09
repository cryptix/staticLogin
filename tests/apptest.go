package tests

import (
	"github.com/robfig/revel"
	"net/url"
)

type AppTest struct {
	revel.TestSuite
}

func (t AppTest) Before() {
	println("Set up")
}

func (t AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestThatRestrictedPageDoesntWork() {
	t.Get("/App/RestrictedIndex")
	t.AssertStatus(401)
}

func (t AppTest) TestCorrectCredentials() {
	v := url.Values{}
	v.Set("username", "cryptix")
	v.Set("password", "12345")
	t.PostForm("/login", v)
	t.AssertOk()
	t.AssertContains("Welcome")
}

func (t AppTest) TestIncorrectCredentials() {
	v := url.Values{}
	v.Set("username", "cryptix")
	v.Set("password", "4444")
	t.PostForm("/login", v)
	t.AssertOk()
	t.AssertContains("Login failed")
}

func (t AppTest) TestThatRegisterPageWorks() {
	t.Get("/register")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) After() {
	println("Tear down")
}
