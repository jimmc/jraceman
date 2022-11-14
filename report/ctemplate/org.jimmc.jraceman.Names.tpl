{{/*GT: {
  "display": "Names",
  "description": "Problems with Person names",
  "permission": "view_roster",
  "where": [ "team", "person" ],
  "orderby": [
    {
      "name": "team",
      "display": "Team",
      "sql": "TeamAbbrev, LastName, FirstName"
    }
  ]
} */ -}}
{{ $comp := computed -}}
{{ $standardColumns := `
    Person.id as PersonId,
    Person.lastname as LastName,
    Person.firstname as FirstName,
    Team.shortName as TeamAbbrev
` -}}
{{ $fromTables := `
    Person
    LEFT JOIN Team on Person.teamId=Team.id
` -}}
{{ $selectSpaces := (printf `
  SELECT
    %s,
    'NameNotTrimmed' as Reason
    FROM %s
    WHERE Person.firstname != TRIM(Person.firstname)
       OR Person.lastname != TRIM(Person.lastname)
    %s
    `
    $standardColumns $fromTables $comp.Where.AndClause)
-}}
{{ $selectLower := (printf `
  SELECT
    %s,
    'NameLowerCase' as Reason
    FROM %s
    WHERE Person.firstname = LOWER(Person.firstname)
       OR Person.lastname = LOWER(Person.lastname)
    %s
    `
    $standardColumns $fromTables $comp.Where.AndClause)
-}}
{{ $selectUpper := (printf `
  SELECT
    %s,
    'NameUpperCase' as Reason
    FROM %s
    WHERE Person.firstname = UPPER(Person.firstname)
       OR Person.lastname = UPPER(Person.lastname)
    %s
    `
    $standardColumns $fromTables $comp.Where.AndClause)
-}}
{{ $selectBlank := (printf `
  SELECT
    %s,
    'NameBlank' as Reason
    FROM %s
    WHERE COALESCE(Person.firstname,'')=''
       OR COALESCE(Person.lastname,'')=''
    %s
    `
    $standardColumns $fromTables $comp.Where.AndClause)
-}}
{{ $selectDupSameTeam := (printf `
  SELECT
    %s,
    'DuplicateName' as Reason
    FROM %s
    JOIN Person as P2
      ON TRIM(LOWER(Person.lastName))=TRIM(LOWER(P2.lastName))
       AND TRIM(LOWER(Person.firstName))=TRIM(LOWER(P2.firstName))
       AND Person.teamId=P2.teamId
       AND Person.id!=P2.id 
    %s
    `
    $standardColumns $fromTables $comp.Where.WhereClause)
-}}
{{ $selectDupDifferentTeam := (printf `
  SELECT
    %s,
    'DuplicateNameDifferentTeam' as Reason
    FROM %s
    JOIN Person as P2
      ON TRIM(LOWER(Person.lastName))=TRIM(LOWER(P2.lastName))
       AND TRIM(LOWER(Person.firstName))=TRIM(LOWER(P2.firstName))
       AND Person.teamId!=P2.teamId
       AND Person.id!=P2.id 
    %s
    `
    $standardColumns $fromTables $comp.Where.WhereClause)
-}}
{{ $selectAll := (printf `
  SELECT * from (%s) UNION ALL
  SELECT * from (%s) UNION ALL
  SELECT * from (%s) UNION ALL
  SELECT * from (%s) UNION ALL
  SELECT * from (%s) UNION ALL
  SELECT * from (%s)
  %s`
  $selectSpaces $selectLower $selectUpper $selectBlank
  $selectDupSameTeam $selectDupDifferentTeam
  $comp.OrderBy.Clause)
-}}
{{ $Rows := rows $selectAll . -}}
{{ $countRow := row (printf `
  SELECT 
    (SELECT count(1) from (%s)) as SpacesCount,
    (SELECT count(1) from (%s)) as LowerCount,
    (SELECT count(1) from (%s)) as UpperCount,
    (SELECT count(1) from (%s)) as BlankCount,
    (SELECT count(1) from (%s)) as DupSameCount,
    (SELECT count(1) from (%s)) as DupDifferentCount,
    (SELECT count(1) from (%s)) as TotalCount
  `
  $selectSpaces $selectLower $selectUpper $selectBlank
  $selectDupSameTeam $selectDupDifferentTeam
  $selectAll) . -}}
<html>
<body>
  <div class="header">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
  <center>
  <div class="titleArea">
    <h3>Name errors ordered by {{ $comp.OrderBy.Display }}</h3>
  </div>
  <table class="main" border=1>
    <tr class="rowHeader">
      <th>Person ID</th>
      <th>Last Name</th>
      <th>First Name</th>
      <th>Team</th>
      <th>Reason</th>
    </tr>
{{ range $Rows }}
    <tr class="rowParity{{ evenodd .rowindex "0" "1" }}">
      <td>{{.PersonId}}</td>
      <td>{{.LastName}}</td>
      <td>{{.FirstName}}</td>
      <td>{{.TeamAbbrev}}</td>
      <td>{{.Reason}}</td>
    </tr>
{{ end }}
  </table>
  Error counts:
    NotTrimmed:{{ $countRow.SpacesCount }}
    LowerCase:{{ $countRow.LowerCount }}
    UpperCase:{{ $countRow.UpperCount }}
    Blank:{{ $countRow.BlankCount }}
    DupSame:{{ $countRow.DupSameCount }}
    DupDifferent:{{ $countRow.DupDifferentCount }}
    Total:{{ $countRow.TotalCount }}
  <div class="footer">
{{ include "org.jimmc.jraceman.datePrintedLine" }}
  </div>
</body>
</html>
