{{/* This fragment calculates and return the orderBy values.
     On entry, dot must be set to the default value for the
     orderBy key in the attributes. */ -}}
{{ $orderByKey := options "orderby" -}}
{{ if not $orderByKey }}{{ $orderByKey = . }}{{end -}}
{{ $orderBySql := attrs "orderby" $orderByKey "sql" -}}
{{ return (mkmap "key" $orderByKey "sql" $orderBySql) -}}
