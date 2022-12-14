{{/*GT: {
  "display": "Entries",
  "description": "Entries ordered as selected.",
  "permission": "view_regatta",
  "where": [ "event", "team", "person" ],
  "orderby": [
    {
      "name": "team",
      "display": "Team, Person, Event",
      "sql": "team.shortname, person.lastname, person.firstname, event.number"
    },
    {
      "name": "person",
      "display": "Person, Event",
      "sql": "person.lastname, person.firstname, event.number, team.shortname"
    },
    {
      "name": "eventTeam",
      "display": "Event, Team, Person",
      "sql": "event.number, team.shortname, person.lastname, person.firstname"
    },
    {
      "name": "eventPerson",
      "display": "Event, Person",
      "sql": "event.number, person.lastname, person.firstname, team.shortname"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $rows := rows (printf `
    SELECT
      team.shortname as Team,
      person.lastName || ', ' || person.firstName as Person,
      '#' || event.number || ' ' ||
      COALESCE(event.name,
        competition.name || ' ' || level.name || ' ' || gender.name)
        as Event,
      entry.groupname as EGroup,
      CASE WHEN COALESCE(entry.alternate,false) THEN 'Yes' ELSE '' END as Alternate
    FROM entry
      LEFT JOIN person on entry.personid=person.id
      LEFT JOIN team on person.teamid=team.id
      LEFT JOIN event on entry.eventid=event.id
      LEFT JOIN competition on event.competitionid=competition.id
      LEFT JOIN level on event.levelid=level.id
      LEFT JOIN gender on event.genderid=gender.id
      LEFT JOIN meet on event.meetid=meet.id
      LEFT JOIN area on event.areaid=area.id
    WHERE (NOT COALESCE(event.scratched,false) AND NOT COALESCE(entry.scratched,false))%s
    %s
` $comp.Where.AndClause $comp.OrderBy.Clause) . -}}
{{ $totals := row (printf `
  SELECT
    count(distinct team.id) as TeamCount,
    count(distinct person.id) as PeopleCount,
    count(distinct entry.id) as EntriesCount,
    count(distinct event.id) as EventsCount,
    count(distinct entry.groupname ||
      CASE WHEN entry.groupname is null THEN null ELSE entry.eventid END) as GroupsCount,
    count(distinct CASE WHEN COALESCE(entry.alternate,0)=0 THEN null ELSE entry.id END) as AlternateCount
  FROM entry
    LEFT JOIN person on entry.personid=person.id
    LEFT JOIN team on person.teamid=team.id
    LEFT JOIN event on entry.eventid=event.id
    LEFT JOIN meet on event.meetid=meet.id
  WHERE (NOT COALESCE(event.scratched,false) AND NOT COALESCE(entry.scratched,false))%s
` $comp.Where.AndClause) -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Entries By {{ $comp.OrderBy.Display }}</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Team</th>
      <th>Person</th>
      <th>Event</th>
      <th>Group</th>
      <th>Alternate?</th>
    </tr>
{{ range $rows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.Team}}</td>
      <td>{{.Person}}</td>
      <td>{{.Event}}</td>
      <td>{{.EGroup}}</td>
      <td>{{.Alternate}}</td>
    </tr>
{{ end }}
{{ with $totals }}
    <tr class="rowTotal">
      <td><b>Totals</b></td>
      <td><b>#Teams: {{.TeamCount}}, &nbsp; #People: {{.PeopleCount}}</b></td>
      <td><b>#Entries: {{.EntriesCount}}, &nbsp; #Events: {{.EventsCount}}</b></td>
      <td><b>#Groups: {{.GroupsCount}}</b></td>
      <td><b>#Alt: {{.AlternateCount}}</b></td>
    </tr>
{{ end }}
  </table>
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
