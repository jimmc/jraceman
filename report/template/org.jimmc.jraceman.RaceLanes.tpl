{{/*GT: {
  "display": "Race Lanes",
  "description": "Lane assignments for a Race",
  "where": [ "event", "race" ],
  "orderby": [
    {
      "name": "event",
      "display": "Event, Race",
      "sql": "event.number, race.number"
    },
    {
      "name": "race",
      "display": "Race",
      "sql": "race.number"
    }
  ]
} */ -}}
{{ $RaceInfo := include "org.jimmc.jraceman.RaceInfo" -}}
{{ $comp := computed -}}
{{ $raceRows := rows (printf `
    SELECT %s
    FROM %s
    WHERE event.number>0
    %s %s
` $RaceInfo.raceColumns $RaceInfo.raceTables $comp.Where.AndClause $comp.OrderBy.Clause) . -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
{{ range $raceRows -}}
  <div class="titleArea">
    {{ evalTemplate $RaceInfo.raceTitleTemplate . }}
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th width=60>Lane</th>
      <th width=300 alight=left>Person</th>
      <th width=100 alight=left>Team</th>
    </tr>
{{ $laneRows := rows (printf `
    SELECT lane.lane as Lane,
        GROUP_CONCAT((person.firstname || ' ' || person.lastname ||
            (CASE WHEN COALESCE(entry.alternate,false) THEN '(alt)' ELSE '' END)),',') as Person,
        team.shortname as Team
    FROM lane
    LEFT JOIN race on lane.raceid = race.id
    LEFT JOIN entry on lane.entryid = entry.id
    LEFT JOIN person on entry.personid = person.id
    LEFT JOIN team on person.teamid = team.id
    WHERE lane.raceid = '%s' AND lane.lane >= 0
    GROUP BY lane.lane
    ORDER BY lane.lane, Person` .raceid) -}}
{{ range $laneRows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td align=center>{{.Lane}}</td>
      <td>{{ range (split .Person ",")}}{{.}}<br/>{{end}}</td>
      <td>{{.Team}}</td>
    </tr>
{{ end }}
  </table>
{{ end }}
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
