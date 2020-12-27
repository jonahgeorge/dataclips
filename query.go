package sqlexplorer

import (
	"encoding/json"
	"net/http"
)

// accepts SQL, returns JSON?
type queryRequest struct {
	SQL string
}

type queryResponse struct {
	Rows []map[string]interface{}
}

func (s *Server) queryHandler(w http.ResponseWriter, r *http.Request) {
	var in queryRequest

	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := s.db.QueryxContext(r.Context(), in.SQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}

	buf, err := json.Marshal(queryResponse{Rows: results})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(buf)
}
