{{/*GT: {
  "display": "Medal Count",
  "description": "The number of medals received by each team",
  "permission": "view_regatta",
  "where": [ "meet", "event", "team", "person" ],
  "orderby": [
    {
      "name": "team",
      "display": "Team",
      "sql": "team.shortname"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $teamRows := rows (printf `
  SELECT 
    Team.name as TeamName,
    Team.shortName as TeamAbbrev,
     count(CASE WHEN scorePlace=1 THEN 1 ELSE null END) as GoldMedals,
     count(CASE WHEN scorePlace=2 THEN 1 ELSE null END) as SilverMedals,
     count(CASE WHEN scorePlace=3 THEN 1 ELSE null END) as BronzeMedals,
     count(CASE WHEN scorePlace IN (1,2,3) THEN 1 ELSE null END) as TotalNatMedals,
     count(CASE WHEN (scorePlace is null AND place=1) THEN 1 ELSE null END) as IntlGoldMedals,
     count(CASE WHEN (scorePlace is null AND place=2) THEN 1 ELSE null END) as IntlSilverMedals,
     count(CASE WHEN (scorePlace is null AND place=3) THEN 1 ELSE null END) as IntlBronzeMedals,
     count(CASE WHEN scorePlace is null AND place IN (1,2,3) THEN 1 ELSE null END) as TotalIntlMedals,
     count(CASE WHEN scorePlace IN (1,2,3) or place in (1,2,3) THEN 1 ELSE null END) as TotalMedals
    from Lane
    LEFT JOIN Race on Lane.raceId=Race.id
    LEFT JOIN Event on Race.eventId=Event.id
    LEFT JOIN Entry on Lane.entryId=Entry.id
    LEFT JOIN Person on Entry.personId=Person.id
    LEFT JOIN Team on Person.teamId=Team.id
    LEFT JOIN Stage on Race.stageId=Stage.id
    LEFT JOIN Meet on Event.meetId=Meet.id
    %s
    group by Team.shortName
    %s`
    $comp.Where.WhereClause $comp.OrderBy.Clause) . -}}
{{ $totals := row (printf `
  SELECT count(distinct Team.id) as TeamCount,
     count(CASE WHEN scorePlace=1 THEN 1 ELSE null END) as GoldMedals,
     count(CASE WHEN scorePlace=2 THEN 1 ELSE null END) as SilverMedals,
     count(CASE WHEN scorePlace=3 THEN 1 ELSE null END) as BronzeMedals,
     count(CASE WHEN scorePlace IN (1,2,3) THEN 1 ELSE null END) as TotalNatMedals,
     count(CASE WHEN scorePlace is null AND place=1 THEN 1 ELSE null END) as IntlGoldMedals,
     count(CASE WHEN scorePlace is null AND place=2 THEN 1 ELSE null END) as IntlSilverMedals,
     count(CASE WHEN scorePlace is null AND place=3 THEN 1 ELSE null END) as IntlBronzeMedals,
     count(CASE WHEN scorePlace is null AND place IN (1,2,3) THEN 1 ELSE null END) as TotalIntlMedals,
     count(CASE WHEN scorePlace IN (1,2,3) or place in (1,2,3) THEN 1 ELSE null END) as TotalMedals
    from Lane
    LEFT JOIN Race on Lane.raceId=Race.id
    LEFT JOIN Event on Race.eventId=Event.id
    LEFT JOIN Entry on Lane.entryId=Entry.id
    LEFT JOIN Person on Entry.personId=Person.id
    LEFT JOIN Team on Person.teamId=Team.id
    LEFT JOIN Stage on Race.stageId=Stage.id
    LEFT JOIN Meet on Event.meetId=Meet.id
    %s
` $comp.Where.WhereClause) -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Medal Count Per Team ordered by {{ $comp.OrderBy.Display }}</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Team</th><th>Gold Medals</th><th>Silver Medals</th>
      <th>Bronze Medals</th><th>Total G/S/B</th>
      <th>Intl Gold Medals</th><th>Intl Silver Medals</th>
      <th>Intl Bronze Medals</th><th>Total Intl G/S/B</th>
      <th>Total Medals</th>
    </tr>
{{ range $teamRows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.TeamName}} ({{.TeamAbbrev}})</td>
      <td>{{.GoldMedals}}</td><td>{{.SilverMedals}}</td>
      <td>{{.BronzeMedals}}</td><td>{{.TotalNatMedals}}</td>
      <td>{{.IntlGoldMedals}}</td><td>{{.IntlSilverMedals}}</td>
      <td>{{.IntlBronzeMedals}}</td><td>{{.TotalIntlMedals}}</td>
      <td>{{.TotalMedals}}</td>
    </tr>
{{ end }}
{{ with $totals }}
    <tr class="rowTotal">
      <td><b>Totals</b> #Teams: {{.TeamCount}}</td>
      <td>{{.GoldMedals}}</td><td>{{.SilverMedals}}</td>
      <td>{{.BronzeMedals}}</td><td>{{.TotalNatMedals}}</td>
      <td>{{.IntlGoldMedals}}</td><td>{{.IntlSilverMedals}}</td>
      <td>{{.IntlBronzeMedals}}</td><td>{{.TotalIntlMedals}}</td>
      <td>{{.TotalMedals}}</td>
    </tr>
{{ end }}
  </table>
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
