package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/cryptix/staticLogin/app/models"
	"github.com/robfig/revel"
)

var (
	Users [](*models.User)
)

func Init() {
	revel.INFO.Print("Initing Users...")
	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("12345"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "Demo User", "cryptix", "clearpw-why?", bcryptPassword}
	Users = append(Users, demoUser)
}
