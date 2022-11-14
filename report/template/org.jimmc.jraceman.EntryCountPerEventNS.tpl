{{/*GT: {
  "display": "Entry Count Per Event non-scoring",
  "description": "The number of Entries and Teams entered, including non-scoring counts, for each Event.",
  "permission": "view_regatta",
  "where": [ "event", "meet", "person", "team" ],
  "orderby": [
    {
      "name": "event",
      "display": "Event",
      "sql": "event.number, event.name"
    },
    {
      "name": "teamcount",
      "display": "#Teams",
      "sql": "teamcount desc, event.number, event.name"
    },
    {
      "name": "entrycount",
      "display": "#Entries",
      "sql": "entrycount desc, event.number, event.name"
    },
    {
      "name": "nsentries",
      "display": "#NS-Entries",
      "sql": "nonscoringentriescount desc, event.number, event.name"
    },
    {
      "name": "groupcount",
      "display": "#Groups",
      "sql": "groupcount desc, event.number, event.name"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $rows := rows (printf `
    SELECT
      ( CASE WHEN Event.scratched THEN '<strike>' ELSE '' END ||
        Event.meetId || ' ' ||
        COALESCE('#' || Event.number || ': ', '') ||
        COALESCE(Event.name,
          (Competition.name || ' ' || Level.name || ' ' || Gender.name)) ||
        ' [' || Event.id || ']' ||
        CASE WHEN Event.scratched THEN '</strike>' ELSE '' END
      ) as Event,
      count(distinct Team.id) as TeamCount,
      count(distinct Entry.id) as EntryCount,
      count(distinct EntryNS.id) as NonScoringEntryCount,
      count(distinct Entry.groupname) as GroupCount
    FROM Entry
      LEFT JOIN Person on Entry.personId=Person.id
      LEFT JOIN Team on Person.teamId=Team.id
      LEFT JOIN Event on Entry.eventId=Event.id
      LEFT JOIN Competition on Event.competitionId=Competition.id
      LEFT JOIN Level on Event.levelId=Level.id
      LEFT JOIN Gender on Event.genderId=Gender.id
      LEFT JOIN Meet on Event.meetId=Meet.id
      LEFT JOIN Entry as EntryNS on Entry.id=EntryNS.id and Team.nonScoring
    WHERE (NOT COALESCE(event.scratched,false) AND NOT COALESCE(entry.scratched,false)) %s
    GROUP BY Event.id
    %s
` $comp.Where.AndClause $comp.OrderBy.Clause) . -}}
{{ $totals := row (printf `
    SELECT
      count(distinct Event.id) as EventCount,
      count(distinct Person.id) as PersonCount,
      count(distinct Team.id) as TeamCount,
      count(distinct Entry.id) as EntryCount,
      count(distinct EntryNS.id) as NonScoringEntryCount,
      count(distinct (Event.id || '-' || Entry.groupname)) as GroupCount
    FROM Entry
      LEFT JOIN Person on Entry.personId=Person.id
      LEFT JOIN Team on Person.teamId=Team.id
      LEFT JOIN Event on Entry.eventId=Event.id
      LEFT JOIN Meet on Event.meetId=Meet.id
      LEFT JOIN Entry as EntryNS on Entry.id=EntryNS.id and Team.nonScoring
    WHERE (NOT COALESCE(event.scratched,false) AND NOT COALESCE(entry.scratched,false)) %s
` $comp.Where.AndClause) -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Entries Per Event by {{ $comp.OrderBy.Display }} with Non-Scoring Counts</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Event</th>
      <th>#Teams</th>
      <th>#Entries</th>
      <th>#NS-Entries</th>
      <th>#Groups</th>
    </tr>
{{ range $rows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.Event}}</td>
      <td align=right>{{.TeamCount}}</td>
      <td align=right>{{.EntryCount}}</td>
      <td align=right>{{.NonScoringEntryCount}}</td>
      <td align=right>{{.GroupCount}}</td>
    </tr>
{{ end }}
{{ with $totals }}
    <tr class="rowTotal">
      <td><b>Totals -- #Events: {{.EventCount}} -- #People: {{.PersonCount}}</b></td>
      <td><b>#Teams: {{.TeamCount}}</b></td>
      <td><b>#Entries: {{.EntryCount}}</b></td>
      <td><b>#NS-Entries: {{.NonScoringEntryCount}}</b></td>
      <td><b>#Groups: {{.GroupCount}}</b></td>
    </tr>
{{ end }}
  </table>
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
