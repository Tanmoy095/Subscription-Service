package main

import "net/http"

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)

}
func (app *Config) LogInPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)

}
func (app *Config) postLoginPage(w http.ResponseWriter, r *http.Request) {
	// Validate and process form data
}
func (app *Config) LogOutPage(w http.ResponseWriter, r *http.Request) {
}
func (app *Config) registerPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}
func (app *Config) postRegisterPage(w http.ResponseWriter, r *http.Request) {

}
func (app *Config) activateAccount(w http.ResponseWriter, r *http.Request) {

}
