{{ $sectionSuffix := `
    (CASE WHEN SectionsPerRound.sectionCount>1
        THEN (CASE WHEN stage.isFinal
             THEN substring('ABCDEFGHIJK',race.section,1)
             ELSE race.section END)
        ELSE '' END) as sectionSuffix` -}}
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
	    ((CASE WHEN event.scratched THEN '<strike>' ELSE '' END) ||
	    COALESCE(('Event #' || event.number || ': '),'') ||
	    COALESCE(event.name,(
		competition.name || ' ' || level.name || ' ' || gender.name)) ||
	    (CASE WHEN event.scratched THEN '</strike>' ELSE ''END))
	as eventInfo` -}}
{{ $raceNumberInfo := `(CASE WHEN COALESCE(race.number,0)=0 THEN '0' ELSE
	(
	    (CASE WHEN race.scratched THEN '<strike>' ELSE '' END) ||
	    'Race #' || (CASE WHEN race.number=ROUND(race.number) THEN
		ROUND(race.number) ELSE race.number END) || ', ' ||
	    (CASE WHEN race.scratched THEN '</strike>' ELSE '' END)
	)
    END) as raceNumberInfo` -}}
{{ $raceColumns := printf `%s,
    race.id as raceid, race.number, stage.name, race.round, race.section, race.eventid,
    SectionsPerRound.sectionCount,%s,%s`
    $eventInfo $sectionSuffix $raceNumberInfo -}}
{{ $eventTitleInfo := `(
	    COALESCE(('Event #' || event.number || ': '),'') ||
	    COALESCE(event.name,(
		competition.name || ' ' || level.name || ' ' || gender.name)) ||
	    (CASE WHEN event.scratched THEN ' (SCR)' ELSE '' END)
	)` -}}
{{ $raceNumberTitleInfo := `
    (CASE WHEN COALESCE(race.number,0)=0 THEN '0' ELSE
	(
	    'Race #' ||
	    (CASE WHEN race.number=ROUND(race.number) THEN
		ROUND(race.number) ELSE race.number END) ||
	    (CASE WHEN race.scratched THEN ' (SCR)' ELSE '' END) ||
	    ', '
	)
    END)` -}}
{{ $raceEventCommentInfo := `(
        COALESCE(race.raceComment,'') ||
        (CASE WHEN race.raceComment is not null AND event.eventComment is not null THEN
            '<br>' ELSE '' END) ||
        COALESCE(event.eventComment,'')
    )` -}}
{{ $raceTitleTemplate := `
    <h3>{{.eventInfo}}<br/>{{.raceNumberInfo}}</h3>` -}}
{{ return (mkmap
    "stageSubCountTable" $stageSubCountTable
    "raceTables" $raceTables
    "raceColumns" $raceColumns
    "eventInfo" $eventInfo
    "eventTitleInfo" $eventTitleInfo
    "raceNumberInfo" $raceNumberInfo
    "raceNumberTitleInfo" $raceNumberTitleInfo
    "raceEventCommentInfo" $raceEventCommentInfo
    "raceTitleTemplate" $raceTitleTemplate
    ) -}}
