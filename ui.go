package sqlexplorer

import (
	"html/template"
	"net/http"
)

func (s *Server) uiHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("test").Parse(uiHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var uiHTML = `<html>
	<head>
		<title>SQL Explorer</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet" 
			integrity="sha384-giJF6kkoqNQ00vy+HMDP7azOuL0xtbfIcaT9wjKHr8RbDVddVHyTfAAsrekwKmP1" crossorigin="anonymous">
	</head>
	<body>
		<div class="navbar navbar-dark bg-dark shadow-sm">
			<div class="container">
				<a href="#" class="navbar-brand d-flex align-items-center">
					<strong>SQL Explorer</strong>
				</a>
			</div>
		</div>

		<div class="flex-container">
			<form action="/query" method="POST">
				<textarea class="form-control" id="sql" rows="4"></textarea>
				<button value="submit">Query</button>
			</form>
		</div>

		<div class="flex-container">
		<table class="table">
			<thead>
				<tr>
				<th scope="col">#</th>
				<th scope="col">First</th>
				<th scope="col">Last</th>
				<th scope="col">Handle</th>
				</tr>
			</thead>
			<tbody>
				<tr>
				<th scope="row">1</th>
				<td>Mark</td>
				<td>Otto</td>
				<td>@mdo</td>
				</tr>
				<tr>
				<th scope="row">2</th>
				<td>Jacob</td>
				<td>Thornton</td>
				<td>@fat</td>
				</tr>
				<tr>
				<th scope="row">3</th>
				<td colspan="2">Larry the Bird</td>
				<td>@twitter</td>
				</tr>
			</tbody>
		</table>
		</div>
	</body>
</html>`
