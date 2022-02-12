package db

import (
  "gorm.io/gorm"
  "time"
)

type Profile struct {
  gorm.Model
  UserType           string         `json:"usertype" gorm:"size:30"`
  Username           string         `json:"username" gorm:"size:20"`
  Password           string         `json:"password" gorm:"size:150"`
}

type Pupil struct {
  gorm.Model
  Profile            Profile        `json:"profile"`
  ProfileID          uint           `json:"profileid"`
  Firstname          string         `json:"firstname" gorm:"size:20"`
  Lastname           string         `json:"lastname" gorm:"size:20"`
  Phone              string         `json:"phone" gorm:"size:20"`
  UserCode           string         `json:"usercode" gorm:"size:20"`
  Deactivated        bool           `json:"deactivated" gorm:"default:false"`
}

type Teacher struct {
  gorm.Model
  Profile            Profile        `json:"profile"`
  ProfileID          uint           `json:"profileid"`
  Firstname          string         `json:"firstname" gorm:"size:20"`
  Lastname           string         `json:"lastname" gorm:"size:20"`
}

type Assignment struct {
    gorm.Model
    AssignmentListID  uint          `json:"assignmentlistid"`
    CharacterToDraw   string        `json:"charactertodraw" gorm:"size:2"`
}

type AssignmentList struct {
  gorm.Model
  Name               string         `json:"name" gorm:"size:20"`
  Expired            bool           `json:"expired" gorm:"default:false"`
  StartDate          time.Time      `json:"startdate"`
  EndDate            time.Time      `json:"enddate"`
  Assignments        []Assignment   `json:"assignments"`
}

type Solution struct {
  gorm.Model 
  PupilID                   uint    `json:"pupilid"`
  AssignmentID              uint    `json:"assignmentid"`
  AssignmentListID          uint    `json:"assignmentlistid"`
  AssignmentListName        string  `json:"assignmentlistname" gorm:"size:50"`
  AttemptTimeTakenInSeconds uint    `json:"attemptTimeTakenInSeconds"`
  CharacterToDraw           string  `json:"charactertodraw" gorm:"size:2"`
  AutoAttachedScore         float32 `json:"autoAttachedscore"`
  Comment                   string  `json:"comment"`
}

type ActivationRequest struct {
    gorm.Model
    PupilID           uint          `json:"pupilid" gorm:"size:5"`
    RequestMessage    string        `json:"requestmessage"`
}
  





