package view

import (
	"html/template"
	"io"
	"path/filepath"
	"time"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/util"
)

type View_ interface {
	Render(w io.Writer) error
}

type View interface {
	RenderTemplate(templateName string, out io.Writer, context map[string]interface{}) error
}

type view struct {
	urlBuilder   *util.UrlBuilder
	config       *config.Config
	templateName string
	context      map[string]interface{}
}

func New(urlBuilder *util.UrlBuilder, config_ *config.Config) View {
	return &view{urlBuilder, config_, "", nil}
}

func (v *view) RenderTemplate(templateName string, out io.Writer, context map[string]interface{}) error {
	ts := template.New("root").Funcs(template.FuncMap{
		"RenderMarkdown": util.ParseMarkdown,
		"ShowTime": func(t time.Time) string {
			return t.Format("Jan. 02, 2006, 3:04 PM")
		},
	})
	ts = template.Must(ts.ParseFiles(
		filepath.Join(v.config.TemplateDir, "layout.tmpl"),
		filepath.Join(v.config.TemplateDir, templateName),
	))
	context_ := map[string]interface{}{
		"SiteTitle":  v.config.SiteTitle,
		"SiteFooter": v.config.SiteFooter,
		"Urls":       v.urlBuilder,
	}
	for key, value := range context {
		context_[key] = value
	}
	if err := ts.ExecuteTemplate(out, "layout.tmpl", context_); err != nil {
		return err
	}
	return nil
}

func (v *view) Render(w io.Writer) error {
	return v.RenderTemplate(v.templateName, w, v.context)
}
