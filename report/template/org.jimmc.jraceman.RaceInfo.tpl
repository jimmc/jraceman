{{ $sectionSuffix := `
    (CASE WHEN SectionsPerRound.sectionCount>1
        THEN (CASE WHEN stage.isFinal
             THEN substring('ABCDEFGHIJK',race.section,1)
             ELSE race.section END)
        ELSE '' END)` -}}
{{ $stageSubCountTable := `
    SELECT race.eventid as eventid, race.round as round, count(*) as sectionCount
    FROM race JOIN stage on race.stageid=stage.id
    GROUP BY race.eventid, race.round` -}}
{{ $raceTables := printf `
    race
    LEFT JOIN event on race.eventid=event.id
    LEFT JOIN meet on event.meetid=meet.id
    LEFT JOIN competition on event.competitionid=competition.id
    LEFT JOIN level on event.levelid=level.id
    LEFT JOIN gender on event.genderid=gender.id
    LEFT JOIN stage on race.stageid=stage.id
    LEFT JOIN (%s) as SectionsPerRound
        on race.round=SectionsPerRound.round and race.eventid=SectionsPerRound.eventid` $stageSubCountTable -}}
{{ $eventInfo := `
    (COALESCE(('Event #' || event.number || ': '),'') ||
    COALESCE(event.name,(
        competition.name || ' ' || level.name || ' ' || gender.name)))
    as eventInfo` -}}
{{ $raceInfo := printf `
    (COALESCE(('Race #' || race.number || ': '),'') || stage.name || ' ' || %s)
    as raceInfo` $sectionSuffix -}}
{{ $raceColumns := printf `
    race.id as raceid, race.number, stage.name, race.round, race.section, race.eventid,
    race.scratched as raceScratched,
    event.scratched as eventScratched,
    SectionsPerRound.sectionCount,%s,%s`
    $eventInfo $raceInfo -}}
{{ $eventTitleInfo := `(
	    COALESCE(('Event #' || event.number || ': '),'') ||
	    COALESCE(event.name,(
		competition.name || ' ' || level.name || ' ' || gender.name)) ||
	    (CASE WHEN event.scratched THEN ' (SCR)' ELSE '' END)
	)` -}}
{{ $raceEventCommentInfo := `(
        COALESCE(race.raceComment,'') ||
        (CASE WHEN race.raceComment is not null AND event.eventComment is not null THEN
            '<br>' ELSE '' END) ||
        COALESCE(event.eventComment,'')
    )` -}}
{{ $raceTitleTemplate := `
    <h3>
        {{ if .eventScratched }}<strike>{{end}}{{.eventInfo}}{{if .eventScratched}}</strike>{{end}}<br/>
        {{ if .raceScratched }}<strike>{{end}}{{.raceInfo}}{{if .raceScratched}}</strike>{{end}}
    </h3>` -}}
{{ return (mkmap
    "stageSubCountTable" $stageSubCountTable
    "raceTables" $raceTables
    "raceColumns" $raceColumns
    "eventInfo" $eventInfo
    "eventTitleInfo" $eventTitleInfo
    "raceInfo" $raceInfo
    "raceEventCommentInfo" $raceEventCommentInfo
    "raceTitleTemplate" $raceTitleTemplate
    ) -}}
