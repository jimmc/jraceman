This is the test sites report.
{{with row "SELECT name, city, phone from site where id = ?" . -}}
name={{.name}}, city={{.city}}, phone={{.phone}}
{{- end}}
