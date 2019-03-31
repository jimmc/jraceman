{{/*GT: {
  "display": "Sites",
  "description": "List all sites",
  "orderby": [
    {
      "name": "name",
      "display": "Site name",
      "sql": "name"
    },
    {
      "name": "city",
      "display": "City",
      "sql": "city"
    }
  ]
} */ -}}
This is the test sites report, ordered by {{ (computed).OrderByDisplay }}.
{{range rows (printf "SELECT name, city, phone from site ORDER BY %s" (computed).OrderByExpr) -}}
name={{.name}}, city={{.city}}, phone={{.phone}}
{{- end}}
