{{/*GT: {
  "permission": "view_regatta",
  "where": [ "site_id" ]
} */ -}}
This is the test sites report.
{{ $comp := computed -}}
{{with row (printf "SELECT name, city, phone from site%s" $comp.Where.WhereClause) -}}
name={{.name}}, city={{.city}}, phone={{.phone}}
{{- end}}
