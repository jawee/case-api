package models

import (
  "errors"
  "html"
  "strings"
  "time"

  "github.com/jinzhu/gorm"
)

type CaseStatus struct {
  ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
  Name string `gorm:"size:255;not null; unique" json:"name"`
  CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
  UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *CaseStatus) Prepare() {
  c.ID = 0
  c.Name = html.EscapeString(strings.TrimSpace(c.Name))
  c.CreatedAt = time.Now()
  c.UpdatedAt = time.Now()
}

func (c *CaseStatus) Validate() error {
  if c.Name == "" {
    return errors.New("Required Name")
  }
  return nil
}

func (c *CaseStatus) SaveCaseStatus(db *gorm.DB) (*CaseStatus, error) {
  var err error
  err = db.Debug().Model(&CaseStatus{}).Create(&c).Error
  if err != nil {
    return &CaseStatus{}, err
  }
  return c, nil
}

func (c *CaseStatus) FindAllCaseStatuses(db *gorm.DB) (*[]CaseStatus, error) {
  var err error
  caseStatuses := []CaseStatus{}
  err = db.Debug().Model(&CaseStatus{}).Limit(100).Find(&caseStatuses).Error
  if err != nil {
    return &[]CaseStatus{}, err
  }
  return &caseStatuses, nil
}

func (c *CaseStatus) FindCaseStatusByID(db *gorm.DB, cid uint64) (*CaseStatus, error) {
  var err error
  err = db.Debug().Model(&CaseStatus{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    return &CaseStatus{}, err
  }
  return c, nil
}

func (c *CaseStatus) UpdateACaseStatus(db *gorm.DB) (*CaseStatus, error) {
  var err error

  err = db.Debug().Model(&CaseStatus{}).Where("id = ?", c.ID).Updates(CaseStatus{Name: c.Name, UpdatedAt: time.Now()}).Error
  if err != nil {
    return &CaseStatus{}, err
  }
  return c, nil
}

func (c *CaseStatus) DeleteACaseStatus(db *gorm.DB, cid uint64) (int64, error) {

  db = db.Debug().Model(&CaseStatus{}).Where("id = ?", cid).Take(&CaseStatus{}).Delete(&CaseStatus{})

  if db.Error != nil {
    if gorm.IsRecordNotFoundError(db.Error) {
      return 0, errors.New("CaseStatus not found")
    }
    return 0, db.Error
  }

  return db.RowsAffected, nil
}
