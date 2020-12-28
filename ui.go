package sqlexplorer

import (
	_ "embed"
	"html/template"
	"net/http"

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

	sql := r.Form.Get("sql")
	vars := templateData{SQL: sql}

	vars.Results = func(sql string) results {
		if sql == "" {
			return results{Warning: "empty sql"}
		}

		out, err := s.queryRPC(r.Context(), &queryRequest{SQL: sql})
		if err != nil {
			return results{Warning: err.Error()}
		}

		if len(out.Rows) == 0 {
			return results{Warning: "result set empty"}
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

		return results{
			Headers: headers,
			Rows:    rows,
		}
	}(sql)

	err := t.Execute(w, vars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
