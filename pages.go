package main

import (
	"github.com/coderconvoy/htmq"
)

func IndexPage() *htmq.Tag {
	p, b := htmq.NewPage("Site Manager")

	//form
	f := htmq.QForm("login", []*htmq.Tag{
		htmq.QText("UserName"), htmq.QInput("text", "username"),
		htmq.QText("<br>Password"), htmq.QInput("password", "password"),
		htmq.QSubmit("Login"),
	})

	b.AddChildren(f)
	return p
}

func HomePage() *htmq.Tag {
	p, b := htmq.NewPage("Home")

}
