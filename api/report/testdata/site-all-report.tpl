{{/*GT: {
  "display": "Sites",
  "description": "List all sites",
  "permission": "view_regatta",
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
This is the test sites report, ordered by {{ (computed).OrderBy.Display }}.
{{range rows (printf "SELECT name, city, phone from site ORDER BY %s" (computed).OrderBy.Expr) -}}
name={{.name}}, city={{.city}}, phone={{.phone}}
{{- end}}
