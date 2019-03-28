{{/*GT: {
  "where": [ "site_id" ]
} */ -}}
This is the test sites report.
{{ $where := where -}}
{{with row (printf "SELECT name, city, phone from site%s" $where.WhereClause) -}}
name={{.name}}, city={{.city}}, phone={{.phone}}
{{- end}}
