{{/*GT: {
  "display": "Lane Assignments",
  "description": "Lane Assignments.",
  "permission": "view_regatta",
  "where": [ "event", "race", "team", "person", "event_race_entry_scratched" ],
  "orderby": [
    {
      "name": "event",
      "display": "Event, Race, Lane, Person",
      "sql": "event.number, race.number, lane.lane, person.lastname, person.firstname"
    },
    {
      "name": "race",
      "display": "Race, Lane, Person",
      "sql": "race.number, lane.lane, person.lastname, person.firstname"
    },
    {
      "name": "teamPerson",
      "display": "Team, Person",
      "sql": "team.shortname, person.lastname, person.firstname"
    },
    {
      "name": "teamRace",
      "display": "Team, Race, Person",
      "sql": "team.shortname, race.number, person.lastname, person.firstname"
    },
    {
      "name": "person",
      "display": "Person, Race",
      "sql": "person.lastname, person.firstname, race.number"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $RoundCountTable := `
 SELECT race.eventId AS eventId, race.round AS round, count(*) AS sectionCount
        FROM race JOIN stage on race.stageId=stage.id
        GROUP BY race.eventId,race.round
` -}}
{{ $rows := rows (printf `
  SELECT
      CASE WHEN event.scratched THEN '<strike>' ELSE '' END
        || event.name
        || CASE WHEN event.scratched THEN '</strike>' ELSE '' END
      as EventName,
      event.number as EventNumber,
      CASE WHEN race.scratched THEN '<strike>' ELSE '' END
        || stage.name,' '
        || CASE WHEN SectionsPerRound.sectionCount>1 THEN
            CASE WHEN stage.isFinal THEN substring('ABCDEFGHIJK',race.section,1) ELSE race.section END
        END
        || CASE WHEN race.scratched THEN '</strike>' ELSE '' END
      as RaceName,
      race.number as RaceNumber,
      lane.lane as LaneNumber,
      entry.groupname as GroupName,
      team.shortName as Team,
      CASE WHEN entry.scratched THEN '<strike>' ELSE '' END
        || person.lastname,', ',person.firstname
        || CASE WHEN entry.scratched THEN '</strike>' ELSE '' END
      as Person
  FROM lane 
      LEFT JOIN entry on lane.entryId=entry.id
      LEFT JOIN person on entry.personId=person.id
      LEFT JOIN team on person.teamId=team.id
      LEFT JOIN race on lane.raceId=race.id
      LEFT JOIN stage on race.stageId=stage.id
      LEFT JOIN event on entry.eventId=event.id
      LEFT JOIN meet on event.meetId=meet.id
      LEFT JOIN (%s) as SectionsPerRound
              on (race.round=SectionsPerRound.round and race.eventId=SectionsPerRound.eventId)
  WHERE (NOT COALESCE(event.scratched,false) AND
         NOT COALESCE(race.scratched,false) AND
         NOT COALESCE(entry.scratched,false))%s
  %s
` $RoundCountTable $comp.Where.AndClause $comp.OrderBy.Clause) . -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Lane Assignments By {{ $comp.OrderBy.Display }}</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Event Name</th>
      <th>Event Number</th>
      <th>Race Name</th>
      <th>Race Number</th>
      <th>Lane Number</th>
      <th>Group</th>
      <th>Team</th>
      <th>Person</th>
    </tr>
{{ range $rows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.EventName}}</td>
      <td>{{.EventNumber}}</td>
      <td>{{.RaceName}}</td>
      <td>{{.RaceNumber}}</td>
      <td>{{.LaneNumber}}</td>
      <td>{{.GroupName}}</td>
      <td>{{.Team}}</td>
      <td>{{.Person}}</td>
    </tr>
{{ end }}
  </table>
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
