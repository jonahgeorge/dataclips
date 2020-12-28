package dataclips

import (
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/Masterminds/sprig"
)

//go:embed ui.html
var uiHTML string

var t = template.Must(template.New("").
	Funcs(sprig.FuncMap()).
	Parse(uiHTML))

type results struct {
	Warning string
	Headers []string
	Rows    [][]interface{}
}

type templateData struct {
	SQL     string
	Results results
}

func (s *Server) uiHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	vars := func(sql string) templateData {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			sql = s.PlaceholderQuery
		}

		vars := templateData{SQL: sql}

		out, err := s.queryRPC(r.Context(), &queryRequest{SQL: sql})
		if err != nil {
			vars.Results.Warning = err.Error()
			return vars
		}

		if len(out.Rows) == 0 {
			vars.Results.Warning = "result set empty"
			return vars
		}

		var headers []string
		for k := range out.Rows[0] {
			headers = append(headers, k)
		}

		var rows [][]interface{}
		for _, row := range out.Rows {
			var cols []interface{}
			for _, h := range headers {
				cols = append(cols, row[h])
			}
			rows = append(rows, cols)
		}

		vars.Results.Headers = headers
		vars.Results.Rows = rows
		return vars
	}(r.Form.Get("sql"))

	err := t.Execute(w, vars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
