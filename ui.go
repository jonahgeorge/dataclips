package sqlexplorer

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed ui.html
var uiHTML string

var t = template.Must(template.New("test").Parse(uiHTML))

func (s *Server) uiHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sql := r.Form.Get("sql")
	vars := map[string]interface{}{"sql": sql}
	vars["out"] = func(sql string) map[string]interface{} {
		if sql == "" {
			return nil
		}

		out, err := s.queryRPC(r.Context(), &queryRequest{SQL: sql})
		if err != nil {
			return nil
		}

		rows := out.Rows
		if len(rows) == 0 {
			// TODO
		}

		return map[string]interface{}{
			"headers": []string{},
			"rows":    [][]string{},
		}
	}(sql)

	if sql != "" {
	}

	err := t.Execute(w, vars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
