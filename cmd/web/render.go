package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// TemplateData holds data that will be passed to templates
type templateData struct {
	StringMap            map[string]string
	IntMap               map[string]int
	FloatMap             map[string]float32
	Data                 map[string]interface{}
	CSRFToket            string
	Flash                string
	Warning              string
	Error                string
	IsAuthenticated      int
	API                  string
	CSSVersion           string
	StripeSecretKet      string
	StripePublishableKey string
}

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

func formatCurrency(n int) string {
	f := float32(n) / float32(100)
	return fmt.Sprintf("$%.2f", f)
}

// Embed the templates directory into the binary - THERE SHOULD BE NO SPACE BETWEEN "//" and "go:embed templates"
//
//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = app.config.api
	td.StripeSecretKet = app.config.stripe.secret
	td.StripePublishableKey = app.config.stripe.key

	if app.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = 1
	} else {
		td.IsAuthenticated = 0
	}

	return td
}

// renderTemplate renders a page
func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	// Check if template is in cache
	_, templateInMap := app.templateCahce[templateToRender]

	// Use cached template if in production
	if templateInMap {
		t = app.templateCahce[templateToRender]
	} else {
		// Otherwise, parse template
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	// Add default data
	if td == nil {
		td = &templateData{}
	}

	td = app.addDefaultData(td, r)

	// Execute the template
	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

// ParseTemplate loads templates and stores them in cache
func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	// build partials (if any)
	if len(partials) > 0 { // this if statement is unnecessary but used to make code more readable
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", x)
		}
	}

	// Load and parse the template
	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", strings.Join(partials, ","), templateToRender)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)
	}
	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	// Store template in cache
	app.templateCahce[templateToRender] = t
	return t, nil
}
