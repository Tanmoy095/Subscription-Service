package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	stringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]interface{}
	Flash         string
	Warning       string
	Error         string
	Authenticated int
	Now           time.Time
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {

	//Define the path to template files
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}
	// Initialize a slice to hold all template files
	var templateSlice []string
	// Add the main template file to the slice
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, t))

	//Add all the partial template file to the slice
	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}
	// Check if TemplateData is nil; if so, create a new empty TemplateData
	if td == nil {
		td = &TemplateData{}

	}
	// Parse all the template files into a single *template.Template object
	temp, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // Send an internal server error response
		return
	}

	// Execute the template with the provided TemplateData
	if err := temp.Execute(w, app.AppDefaultData(td, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // Send an internal server error response
		return
	}

}
func (app *Config) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}

func (app *Config) AppDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")     // Get and remove "flash" message from session
	td.Warning = app.Session.PopString(r.Context(), "warning") // Get and remove "warning" message from session
	td.Error = app.Session.PopString(r.Context(), "error")     // Get and remove "error" message from session

	if app.IsAuthenticated(r) {
		td.Authenticated = 1 // Set to 1 if user is authenticated
	}

	return td // Return the updated TemplateData
}
