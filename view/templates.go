package view

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
)

const (
	templateExtension = ".tmpl.html"
)

var (
	_, binpath, _, _ = runtime.Caller(0)
	basepath         = filepath.Dir(binpath)
	funcs            = template.FuncMap{
		"context": contextSignature,
		"dict":    dict,
	}
	templates = template.Must(template.New("picoshop").
			Funcs(funcs).
			ParseGlob(basepath + "/*" + templateExtension))
)

type Page struct {
	Title string
	Data  interface{}
}

func contextSignature(_ interface{}) interface{} {
	return nil
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 { // check that parameters are pairs (key, value)
		return nil, errors.New("invalid dict call")
	}

	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}

	return dict, nil
}

func Render(ctx context.Context, writer http.ResponseWriter, tpl string, page Page) {
	funcs := template.FuncMap{
		"context": ctx.Value,
		"dict":    dict,
	}

	err := templates.Funcs(funcs).ExecuteTemplate(writer, tpl+".tmpl.html", page)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
