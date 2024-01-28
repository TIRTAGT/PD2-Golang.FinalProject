package controllers

import (
	"github.com/TIRTAGT/PD2-Golang.FinalProject/controllers/home"
	"github.com/TIRTAGT/PD2-Golang.FinalProject/controllers/submit_form"
	"github.com/TIRTAGT/PD2-Golang.FinalProject/server/controller/handlerstruct"
)

var ControllerMap = map[string]handlerstruct.ControllerStruct{
	"/": {
		GET: home.GET,
	},
	"/submit_form": {
		GET:  submit_form.GET,
		POST: submit_form.POST,
	},
}
