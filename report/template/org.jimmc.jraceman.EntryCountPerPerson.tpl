{{/*GT: {
  "display": "Entry Count Per Person",
  "description": "The number of Entries entered for each Person.",
  "permission": "view_regatta",
  "where": [ "event", "meet", "person", "team" ],
  "orderby": [
    {
      "name": "person",
      "display": "Person",
      "sql": "Person.lastName, Person.firstName"
    },
    {
      "name": "entrycount",
      "display": "#Entries",
      "sql": "entrycount desc, event.number, event.name"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $rows := rows (printf `
    SELECT
      Team.shortName as Team,
      (Person.lastName || ', ' || Person.firstName) as Person,
      count(distinct Entry.id) as EntryCount,
      count(distinct Entry.groupname) as GroupCount
    FROM Entry
      LEFT JOIN Person on Entry.personId=Person.id
      LEFT JOIN Team on Person.teamId=Team.id
      LEFT JOIN Event on Entry.eventId=Event.id
      LEFT JOIN Competition on Event.competitionId=Competition.id
      LEFT JOIN Level on Event.levelId=Level.id
      LEFT JOIN Gender on Event.genderId=Gender.id
      LEFT JOIN Meet on Event.meetId=Meet.id
    WHERE (NOT COALESCE(event.scratched,false) AND NOT COALESCE(entry.scratched,false)) %s
    GROUP BY Person.id
    %s
` $comp.Where.AndClause $comp.OrderBy.Clause) . -}}
{{ $totals := row (printf `
    SELECT
      count(distinct Event.id) as EventCount,
      count(distinct Person.id) as PersonCount,
      count(distinct Team.id) as TeamCount,
      count(distinct Entry.id) as EntryCount,
      count(distinct (Event.id || '-' || Entry.groupname)) as GroupCount
    FROM Entry
      LEFT JOIN Person on Entry.personId=Person.id
      LEFT JOIN Team on Person.teamId=Team.id
      LEFT JOIN Event on Entry.eventId=Event.id
      LEFT JOIN Meet on Event.meetId=Meet.id
    WHERE (NOT COALESCE(event.scratched,false) AND NOT COALESCE(entry.scratched,false)) %s
` $comp.Where.AndClause) -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Entries Per Person by {{ $comp.OrderBy.Display }}</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Team</th>
      <th>Person</th>
      <th>#Entries</th>
      <th>#Groups</th>
    </tr>
{{ range $rows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.Team}}</td>
      <td>{{.Person}}</td>
      <td align=right>{{.EntryCount}}</td>
      <td align=right>{{.GroupCount}}</td>
    </tr>
{{ end }}
{{ with $totals }}
    <tr class="rowTotal">
      <td><b>#Teams: {{.TeamCount}}</b></td>
      <td><b>#People: {{.PersonCount}}</b></td>
      <td><b>#Entries: {{.EntryCount}}</b></td>
      <td><b>#Groups: {{.GroupCount}}</b></td>
    </tr>
{{ end }}
  </table>
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
