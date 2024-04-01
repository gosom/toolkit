package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/gosom/toolkit/pkg/errorsext"
)

// Template stores the meta data for each template, and whether it uses a layout.
type Template struct {
	layout   string
	name     string
	template *template.Template
}

func (t *Template) Layout() string {
	return t.layout
}

func (t *Template) Name() string {
	return t.name
}

func (t *Template) Template() *template.Template {
	return t.template
}

// TemplateRenderer is a custom html/template renderer for Echo framework.
type TemplateRenderer struct {
	templates     map[string]*Template
	templateFuncs template.FuncMap
	localizers    map[string]*i18n.Localizer
}

// New setup a new template renderer.
func New() *TemplateRenderer {
	tr := TemplateRenderer{
		templates:     make(map[string]*Template),
		templateFuncs: make(template.FuncMap),
		localizers:    make(map[string]*i18n.Localizer),
	}

	return &tr
}

func (t *TemplateRenderer) String() string {
	sb := strings.Builder{}

	for k, v := range t.templates {
		sb.WriteString(fmt.Sprintf("key=%s layout=%s name=%s\n", k, v.Layout(), v.Name()))
	}

	return strings.TrimSpace(sb.String())
}

func (t *TemplateRenderer) AddLocalizer(lang string, localizer *i18n.Localizer) {
	t.localizers[lang] = localizer
}

func (t *TemplateRenderer) AddTemplateFunc(name string, fn any) {
	if _, ok := t.templateFuncs[name]; !ok {
		t.templateFuncs[name] = fn
	} else {
		panic(fmt.Sprintf("template function %s already exists", name))
	}
}

// AddWithLayout register one or more templates using the provided layout.
func (t *TemplateRenderer) AddWithLayout(fsys fs.FS, layout string, patterns ...string) error {
	filenames, err := readFileNames(fsys, patterns...)
	if err != nil {
		return errorsext.WithStack(fmt.Errorf("%w: failed to list using file pattern", err))
	}

	for _, f := range filenames {
		tname := path.Base(f)
		lname := path.Base(layout)

		tmp, err := template.New(tname).Funcs(t.templateFuncs).ParseFS(fsys, layout, f)
		if err != nil {
			return errorsext.WithStack(fmt.Errorf("%w: failed to parse template %s", err, f))
		}

		t.templates[tname] = &Template{
			layout:   lname,
			name:     tname,
			template: tmp,
		}
	}

	return nil
}

// AddWithLayoutAndIncludes register one or more templates using the provided layout and includes.
func (t *TemplateRenderer) AddWithLayoutAndIncludes(fsys fs.FS, layout, includes string, patterns ...string) error {
	filenames, err := readFileNames(fsys, patterns...)
	if err != nil {
		return errorsext.WithStack(fmt.Errorf("%w: failed to list using file pattern", err))
	}

	for _, f := range filenames {
		tname := path.Base(f)
		lname := path.Base(layout)

		tmp, err := template.New(tname).Funcs(t.templateFuncs).ParseFS(fsys, layout, includes, f)
		if err != nil {
			return errorsext.WithStack(fmt.Errorf("%w: failed to parse template %s", err, f))
		}

		t.templates[tname] = &Template{
			layout:   lname,
			name:     tname,
			template: tmp,
		}
	}

	return nil
}

// Add add a template to the registry.
func (t *TemplateRenderer) Add(fsys fs.FS, patterns ...string) error {
	filenames, err := readFileNames(fsys, patterns...)
	if err != nil {
		return errorsext.WithStack(fmt.Errorf("%w: failed to read file names using file pattern", err))
	}

	for _, f := range filenames {
		tname := path.Base(f)

		tmp, err := template.New(tname).Funcs(t.templateFuncs).ParseFS(fsys, f)
		if err != nil {
			return errorsext.WithStack(fmt.Errorf("%w: failed to parse template %s", err, f))
		}

		t.templates[tname] = &Template{
			name:     tname,
			template: tmp,
		}
	}

	return nil
}

// Render renders a template document.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	// use the name of the template, or layout if it exists
	execName := tmpl.name
	if tmpl.layout != "" {
		execName = tmpl.layout
	}

	err := tmpl.template.ExecuteTemplate(w, execName, data)
	if err != nil {
		return err
	}

	return nil
}

func (t *TemplateRenderer) GetTemplate(name string) (*Template, bool) {
	ans, ok := t.templates[name]

	return ans, ok
}

func readFileNames(fsys fs.FS, patterns ...string) ([]string, error) {
	var filenames []string

	for _, pattern := range patterns {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return nil, errorsext.WithStack(fmt.Errorf("%w: failed to list using file pattern", err))
		}

		if len(list) == 0 {
			return nil, errorsext.WithStack(fmt.Errorf("template: pattern matches no files: %#q", pattern))
		}

		filenames = append(filenames, list...)
	}

	return filenames, nil
}
