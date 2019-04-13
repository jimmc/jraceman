{{/*GT: {
  "display": "Entry Count Per Level",
  "description": "The number of Entries and Teams entered for each Level.",
  "where": [ "event", "meet", "person", "team" ],
  "orderby": [
    {
      "name": "level",
      "display": "Level",
      "sql": "Level.minAge,Level.minEntryAge,Level.name"
    },
    {
      "name": "teamcount",
      "display": "#Teams",
      "sql": "teamcount desc, Level.minAge, Level.minEntryAge, Level.name"
    },
    {
      "name": "personcount",
      "display": "#People",
      "sql": "personcount desc, Level.minAge, Level.minEntryAge, Level.name"
    },
    {
      "name": "entrycount",
      "display": "#Entries",
      "sql": "entrycount desc, Level.minAge, Level.minEntryAge, Level.name"
    },
    {
      "name": "groupcount",
      "display": "#Groups",
      "sql": "groupcount desc, Level.minAge, Level.minEntryAge, Level.name"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $rows := rows (printf `
    SELECT
      Level.name as Level,
      count(distinct Person.id) as PersonCount,
      count(distinct Team.id) as TeamCount,
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
    GROUP BY Level.id
    %s
` $comp.Where.AndClause $comp.OrderBy.Clause) . -}}
{{ $totals := row (printf `
    SELECT
      count(distinct Level.id) as LevelCount,
      count(distinct Person.id) as PersonCount,
      count(distinct Team.id) as TeamCount,
      count(distinct Entry.id) as EntryCount,
      count(distinct (Event.id || '-' || Entry.groupname)) as GroupCount
    FROM Entry
      LEFT JOIN Person on Entry.personId=Person.id
      LEFT JOIN Team on Person.teamId=Team.id
      LEFT JOIN Event on Entry.eventId=Event.id
      LEFT JOIN Meet on Event.meetId=Meet.id
      LEFT JOIN Level on Event.levelId=Level.id
    WHERE NOT COALESCE(Entry.scratched, false) %s
` $comp.Where.AndClause) -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Entries Per Level by {{ $comp.OrderBy.Display }}</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Level</th>
      <th>#Teams</th>
      <th>#People</th>
      <th>#Entries</th>
      <th>#Groups</th>
    </tr>
{{ range $rows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.Level}}</td>
      <td align=right>{{.TeamCount}}</td>
      <td align=right>{{.PersonCount}}</td>
      <td align=right>{{.EntryCount}}</td>
      <td align=right>{{.GroupCount}}</td>
    </tr>
{{ end }}
{{ with $totals }}
    <tr class="rowTotal">
      <td><b>Totals -- #Levels: {{.LevelCount}}</b></td>
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
