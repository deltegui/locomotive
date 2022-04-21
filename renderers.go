package phoenix

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func parseTemplate(layout, view string) *template.Template {
	templateRoot := "web/templates"
	path := fmt.Sprintf("%s%s", templateRoot, view)
	layoutPath := fmt.Sprintf("%s%s", templateRoot, layout)
	if layout == "" {
		return template.Must(template.ParseFiles(path))
	} else {
		return template.Must(template.ParseFiles(layoutPath, path))
	}
}

type RequestMapper func(req *http.Request) interface{}

type ViewConfig struct {
	Layout, View, Name string
	MapRequest         RequestMapper
}

func RenderView(conf ViewConfig) http.HandlerFunc {
	tmpl := parseTemplate(conf.Layout, conf.View)
	return func(w http.ResponseWriter, req *http.Request) {
		tmpl.ExecuteTemplate(w, conf.Name, conf.MapRequest(req))
	}
}

type HTMLRenderer struct {
	view     string
	template *template.Template
	w        http.ResponseWriter
}

type HTMLConfig struct {
	Layout, View, Name string
}

func NewHTMLRenderer(conf HTMLConfig) HTMLRenderer {
	return HTMLRenderer{
		view:     conf.Name,
		template: parseTemplate(conf.Layout, conf.View),
	}
}

func (renderer *HTMLRenderer) Use(w http.ResponseWriter) {
	renderer.w = w
}

func (renderer HTMLRenderer) execute(viewmodel interface{}) {
	renderer.template.ExecuteTemplate(renderer.w, renderer.view, viewmodel)
}

func (renderer HTMLRenderer) Render(data interface{}) {
	renderer.execute(data)
}

// JSONPresenter is a presenter that renders your data in JSON format.
type JSONPresenter struct {
	Writer http.ResponseWriter
}

// NewJSONPresenter creates a presenter that renders your data in JSON format.
func NewJSONPresenter(writer http.ResponseWriter) JSONPresenter {
	return JSONPresenter{writer}
}

// Present data in JSON format
func (renderer JSONPresenter) Present(data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		return
	}
	renderer.Writer.Header().Set("Content-Type", "application/json")
	renderer.Writer.Write(response)
}

// PresentError renders a JSON with your error
func (renderer JSONPresenter) PresentError(caseError error) {
	renderer.Writer.WriteHeader(http.StatusBadRequest)
	renderer.Present(caseError)
}
