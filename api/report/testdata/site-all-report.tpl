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
{{ $orderBy := include "org.jimmc.jraceman.orderBy" "name" -}}
This is the test sites report, ordered by {{ attrs "orderby" $orderBy.key "display" }}.
{{range rows (printf "SELECT name, city, phone from site ORDER BY %s" $orderBy.sql) -}}
name={{.name}}, city={{.city}}, phone={{.phone}}
{{- end}}
