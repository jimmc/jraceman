{{/*GT: {
  "display": "Genders",
  "description": "Gender mismatches between Person and Event",
  "where": [ "meet", "event", "team", "person" ],
  "orderby": [
    {
      "name": "team",
      "display": "Team",
      "sql": "team.shortname, person.lastname, person.firstname"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $fromTables := `
    from Entry
    LEFT JOIN Person on Entry.personId=Person.id
    LEFT JOIN Team on Person.teamId=Team.id
    LEFT JOIN Event on Entry.eventId=Event.id
    LEFT JOIN Competition on Event.competitionid=Competition.id
    LEFT JOIN Level on Event.levelid=Level.id
    LEFT JOIN Gender as EventGender on Event.genderid=EventGender.id
    LEFT JOIN Meet on Event.meetId=Meet.id
    LEFT JOIN Gender as PersonGender on Person.genderid=PersonGender.id
` -}}
{{ $where := `
    WHERE Event.genderid != 'X' AND Person.genderid != Event.genderid
       AND NOT COALESCE(Event.scratched,false)
       AND NOT COALESCE(Entry.scratched,false)
` -}}
{{ $teamRows := rows (printf `
  SELECT 
    Person.lastname || ', ' || Person.firstname as PersonName,
    Team.shortName as TeamAbbrev,
    Entry.id as EntryId,
    Event.name as EventSummary,
    Meet.shortname as MeetName,
    EventGender.name as EventGenderName,
    PersonGender.name as PersonGenderName
    %s %s %s %s`
    $fromTables $where $comp.Where.AndClause $comp.OrderBy.Clause) . -}}
{{ $countRow := row (printf `
  SELECT count(*) as Count
  %s %s %s`
  $fromTables $where $comp.Where.AndClause) . -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Gender mismatch errors ordered by {{ $comp.OrderBy.Display }}</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Person</th>
      <th>Team</th>
      <th>Meet</th>
      <th>Event</th>
      <th>Entry ID</th>
      <th>Event Gender</th>
      <th>Person Gender</th>
    </tr>
{{ range $teamRows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.TeamAbbrev}}</td>
      <td>{{.PersonName}}</td>
      <td>{{.MeetName}}</td>
      <td>{{.EventSummary}}</td>
      <td>{{.EntryId}}</td>
      <td>{{.EventGenderName}}</td>
      <td>{{.PersonGenderName}}</td>
    </tr>
{{ end }}
  </table>
  Count of gender errors: {{ $countRow.Count }}
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
