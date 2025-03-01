package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates
var emailTemplatesFS embed.FS

func (app *application) SendMail(from, to, subject, tmpl string, data interface{}) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl")

	t, err := template.New("email-html").ParseFS(emailTemplatesFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	formatedMessage := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplatesFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	plainMessage := tpl.String()

	return nil
}
