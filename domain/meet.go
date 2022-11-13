package domain

// MeetRepo describes how Meet records are loaded and saved.
type MeetRepo interface {
  FindByID(ID string) (*Meet, error)
  List(offset, limit int) ([]*Meet, error)
  Save(*Meet) (string, error)
  UpdateByID(ID string, oldMeet, newMeet *Meet, diffs Diffs) error
  DeleteByID(ID string) error
}

// Meet describes a sporting meet or regatta.
type Meet struct {
  ID string
  Name string
  ShortName string
  SiteID string
  StartDate *string
  EndDate *string
  AgeDate *string
  WebReportsDirectory *string
  TransferDirectory *string
  LabelImageLeft *string
  LabelImageRight *string
  ScoringSystemID *string
}

// MeetMeta provides funcions related to the Meet struct.
type MeetMeta struct {}

func (m *MeetMeta) EntityTypeName() string {
  return "meet"
}

func (m *MeetMeta) EntityGroupName() string {
  return "regatta"
}

func (m *MeetMeta) NewEntity() interface{} {
  return &Meet{}
}
