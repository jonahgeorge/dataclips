package sqlexplorer

import (
	"context"
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
	in := new(queryRequest)

	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := s.queryRPC(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(buf)
}

func (s *Server) queryRPC(ctx context.Context, in *queryRequest) (*queryResponse, error) {
	rows, err := s.db.QueryxContext(ctx, in.SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return &queryResponse{Rows: results}, nil
}
