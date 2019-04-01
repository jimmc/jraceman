package ixport

import (
  "fmt"

  "github.com/golang/glog"
)

type tableColumnsInfo struct {
  v1TableName string
  v1ColumnNames []string
  v2TableName string
  v2ColumnNames []string
}

/* tableMapValue is what we put in a table map that we build from the tableColumnsInfo. */
type tableMapValue struct {
  tableName string
  columnMap map[string]string
}

var tableMaps []tableColumnsInfo = []tableColumnsInfo{
  tableColumnsInfo{"Areas", []string{"id", "name", "siteId", "lanes", "extraLanes"},
      "area", []string{"id", "name", "siteid", "lanes", "extralanes"}},
  tableColumnsInfo{"Challenges", []string{"id","name"},
      "challenge", []string{"id","name"}},
  tableColumnsInfo{"Competitions", []string{"id","name","groupSize","maxAlternates","scheduledDuration"},
      "competition", []string{"id","name","groupsize","maxalternates","scheduledduration"}},
  tableColumnsInfo{"Complans",[]string{"id","system","plan","minEntries","maxEntries","planOrder"},
      "complan", []string{"id","system","plan","minentries","maxentries","planorder"}},
  tableColumnsInfo{"ComplanRules", []string{"id","complanId","fromRound","fromSection","fromPlace","toRound","toSection","toLane"},
      "complanrule", []string{"id","complanid","fromround","fromsection","fromplace","toround","tosection","tolane"}},
  tableColumnsInfo{"ComplanStages", []string{"id","complanId","stageId","round","sectionCount","fillOrder"},
      "complanstage", []string{"id","complanid","stageid","round","sectioncount","fillorder"}},
  tableColumnsInfo{"ContextOptions", []string{"id","name","value","host","webContext","meetId"},
      "contextoption", []string{"id","name","value","host","webcontext","meetId"}},
  tableColumnsInfo{"Entries", []string{"id","personId","group","alternate","scratched","eventId"},
      "entry", []string{"id","personid","groupname","alternate","scratched","eventid"}},
  tableColumnsInfo{"Events", []string{"id","meetId","number","name","competitionId","levelId","genderId","areaId","seedingPlanId","progressionId","progressionState","scoringSystemId","scratched","eventComment"},
      "event", []string{"id","meetid","number","name","competitionid","levelid","genderid","areaid","seedingplanid","progressionid","progressionstate","scoringsystemid","scratched","eventcomment"}},
  tableColumnsInfo{"Exceptions", []string{"id","name","shortName","resultAllowedRequired"},
      "exception", []string{"id","name","shortname","resultallowedrequired"}},
  tableColumnsInfo{"Genders", []string{"id","name"},
      "gender", []string{"id","name"}},
  tableColumnsInfo{"LaneOrder", []string{"id","areaId","lane","order"},
      "laneorder", []string{"id","areaId","lane","ordering"}},
  tableColumnsInfo{"Lanes", []string{"id","entryId","raceId","lane","place","result","exceptionId","scorePlace","score"},
      "lane", []string{"id","entryid","raceid","lane","place","result","exceptionid","scoreplace","score"}},
  tableColumnsInfo{"Levels", []string{"id","name","minEntryAge","minAge","maxAge","maxEntryAge","useGroupAverage"},
      "level", []string{"id","name","minentryage","minage","maxage","maxentryage","usegroupaverage"}},
  tableColumnsInfo{"Meets", []string{"id","name","shortName","siteId","startDate","endDate","ageDate","webReportsDirectory","transferDirectory","labelImageLeft","labelImageRight","scoringSystemId"},
      "meet", []string{"id","name","shortname","siteid","startdate","enddate","agedate","webreportsdirectory","transferdirectory","labelimageleft","labelimageright","scoringsystemid"}},
  tableColumnsInfo{"Options", []string{"name", "value"},
      "option", []string{"name", "value"}},
  tableColumnsInfo{"People", []string{"id","lastName","firstName","title","teamId","birthday","genderId","membership","membershipExpiration","teamEdit","street","street2","city","state","country","zip","phone","fax","email","webPassword"},
      "person", []string{"id","lastname","firstname","title","teamid","birthday","genderid","membership","membershipexpiration","","street","street2","city","state","country","zip","phone","","email",""}},
      // Note we drop some columns: teamEdit, webPassword.
  tableColumnsInfo{"Progressions", []string{"id","name","class","parameters"},
      "progression", []string{"id","name","class","parameters"}},
  tableColumnsInfo{"Races", []string{"id","eventId","number","stageId","round","section","areaId","scheduledStart","scheduledDuration","actualStart","scratched","raceComment"},
      "race", []string{"id","eventid","number","stageid","round","section","areaid","scheduledstart","scheduledduration","actualstart","scratched","racecomment"}},
  tableColumnsInfo{"Registrations", []string{"id","meetId","personId","amountCharged","surcharge","discount","amountPaid","waiverSigned","paymentNotes"},
      "registration", []string{"id","meetid","personid","amountcharged","surcharge","discount","amountpaid","waiversigned","paymentnotes"}},
  tableColumnsInfo{"RegistrationFees", []string{"id","meetId","eventCount","amountCharged"},
      "registrationfee", []string{"id","meetid","eventcount","amountcharged"}},
  tableColumnsInfo{"ScoringRules", []string{"id","scoringSystemId","rule","value","points"},
      "scoringrule", []string{"id","scoringSystemId","rule","value","points"}},
  tableColumnsInfo{"ScoringSystems", []string{"id","name","class","parameters"},
      "scoringsystem", []string{"id","name","class","parameters"}},
  tableColumnsInfo{"SeedingLists", []string{"id","seedingPlanId","rank","personId"},
      "seedinglist", []string{"id","seedingplanid","rank","personid"}},
  tableColumnsInfo{"SeedingPlans", []string{"id","name","seedingOrder"},
      "seedingplan", []string{"id","name","seedingorder"}},
  tableColumnsInfo{"Simplans", []string{"id","system","plan","minEntries","maxEntries"},
      "simplan", []string{"id","system","plan","minentries","maxentries"}},
  tableColumnsInfo{"SimplanRules", []string{"id","simplanId","fromStageId","toStageId","thruPlace","nextBestTimes"},
      "simplanrule", []string{"id","simplanid","fromstageid","tostageid","thruplace","nextbesttimes"}},
  tableColumnsInfo{"SimplanStages", []string{"id","simplanId","stageId","sectionCount","fillOrder"},
      "simplanstage", []string{"id","simplanid","stageid","sectioncount","fillorder"}},
  tableColumnsInfo{"Sites",[]string{"id","name","street","street2","city","state","country","zip","phone","fax"},
      "site", []string{"id","name","street","street2","city","state","country","zip","phone","fax"}},
  tableColumnsInfo{"Stages", []string{"id","name","number","isFinal"},
      "stage", []string{"id","name","number","isfinal"}},
  tableColumnsInfo{"Teams", []string{"id","shortName","name","number","challengeId","nonScoring","street","street2","city","state","country","zip","phone","fax"},
      "team", []string{"id","shortname","name","number","challengeid","nonscoring","street","street2","city","state","country","zip","phone","fax"}},
}

/* InitTableMaps builds the maps from v1 table and column names to v2 table and column names.
 */
func (im *Importer) initTableMaps() {
  tableMap := make(map[string]*tableMapValue)
  for t, tableInfo := range tableMaps {
    columnMap := make(map[string]string)
    for c, columnName := range tableInfo.v1ColumnNames {
      columnMap[columnName] = tableInfo.v2ColumnNames[c]
    }
    tableMap[tableInfo.v1TableName] = &tableMapValue{
      tableName: tableMaps[t].v2TableName,
      columnMap: columnMap,
    }
  }
  im.v1v2map = tableMap
}

/* translateNames1To2 converts a table name and its column names from a v1 export file
 * to the corresponding table and columns used in v2.
 */
func (im *Importer) translateNames1To2(v1TableName string, v1ColumnNames []string) (string, []string, error) {
  tableMap := im.v1v2map[v1TableName]
  if tableMap == nil {
    return "", nil, fmt.Errorf("no translation to v2 from table name %s", v1TableName)
  }
  v2columnNames := make([]string, 0)
  for _, v1columnName := range v1ColumnNames {
    v2columnName := tableMap.columnMap[v1columnName]
    if v2columnName != "" {
      v2columnNames = append(v2columnNames, tableMap.columnMap[v1columnName])
    } else {
      glog.Errorf("No translation for v1 table %s column %s", v1TableName, v1columnName)
    }
  }
  return tableMap.tableName, v2columnNames, nil
}

/* translateValues1To2 takes a row of v1 data and fixes the values to be appropriate for v2.
 * In most cases this does nothing, but for selected tables it can make changes
 * such as converting a null value to a default value.
 * It may modify the values array passed in.
 */
func (im *Importer) translateValues1To2(values []interface{}) []interface{} {
  if im.tableName == "area" {
    return im.translateAreaValues1To2(values)
  }
  if im.tableName == "registration" {
    return im.translateRegistrationValues1To2(values)
  }
  if im.tableName == "stage" {
    return im.translateStageValues1To2(values)
  }
  if im.tableName == "team" {
    return im.translateTeamValues1To2(values)
  }
  // For all other tables, look for and remove values for columns we don't translate.
  return im.removeNontranslatedValues(values)
}

// In all of the following table data translation functions,
// we could optimize by setting up columns indexes in translateName1To2,
// but we'd have to save that state. We expect this import to be a
// relatively rare opertion, so we don't worry about that optimization.

func (im *Importer) translateAreaValues1To2(values []interface{}) []interface{} {
  for n, s := range im.columnNames {
    if s == "extralanes" {
      if values[n] == nil {
        values[n] = 0
      }
    }
  }
  return values
}

func (im *Importer) removeNontranslatedValues(values []interface{}) []interface{} {
  // We drop some column values.
  v2values := make([]interface{}, len(im.columnNames))
  c2 := 0
  for n, s := range im.v1ColumnNames {
    if im.v1v2map[im.v1TableName].columnMap[s] != "" {
      v2values[c2] = values[n]
      c2 += 1
    }
  }
  if c2 != len(im.columnNames) {
    glog.Errorf("Wrong number of values being returned from removeNontranslatedValues, got %d, want %d",
        c2, len(im.columnNames))
  }
  return v2values
}

func (im *Importer) translateRegistrationValues1To2(values []interface{}) []interface{} {
  for n, s := range im.columnNames {
    if s == "surcharge" || s == "discount" || s == "amountpaid" {
      if values[n] == nil {
        values[n] = 0
      }
    }
  }
  return values
}

func (im *Importer) translateStageValues1To2(values []interface{}) []interface{} {
  for n, s := range im.columnNames {
    if s == "isfinal" {
      if values[n] == nil {
        values[n] = false
      }
    }
    if s == "number" {
      if values[n] == nil {
        values[n] = 0
      }
    }
  }
  return values
}

func (im *Importer) translateTeamValues1To2(values []interface{}) []interface{} {
  for n, s := range im.columnNames {
    if s == "nonscoring" {
      if values[n] == nil {
        values[n] = false
      }
    }
  }
  return values
}
