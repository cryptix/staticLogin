package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/cryptix/staticLogin/app/models"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

// private helpers
func (c App) loggedIn() *models.User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c App) getUser(username string) *models.User {
	// revel.INFO.Printf("Looking for user: %s\n", username)
	for _, user := range Users {
		// revel.INFO.Printf("checking against user: %s\n", user)
		if user.Username == username {
			revel.INFO.Printf("found user %s\n", user.Username)
			return user
		}
	}
	return nil
}

// public
//

// landing page
func (c App) Index() revel.Result {
	return c.Render()
}

// serve register form
func (c App) Register() revel.Result {
	return c.Render()
}

// restricted area
func (c App) RestrictedIndex() revel.Result {
	// check if user is logged in
	if user := c.loggedIn(); user == nil {
		// send him away if not
		c.Flash.Error("Sorry - You don't have access!")
		return c.Redirect(App.Index)
	} else {
		return c.Render(user)
	}
}

// check if login is correct
func (c App) ProcessLogin(username, password string) revel.Result {
	user := c.getUser(username)

	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			revel.INFO.Printf("Password matched for user %s\n", username)
			c.Session["user"] = username
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(App.RestrictedIndex)
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(App.Index)
}

// create new user
func (c App) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Register)
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	Users = append(Users, &user)

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(App.RestrictedIndex)
}

// destroy session
func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(App.Index)
}
