package models

import (
  "errors"
  "html"
  "strings"
  "time"

  "github.com/jinzhu/gorm"
)

type CasePriority struct {
  ID uint32 `gorm:"primary_key;auto_increment" json:"id"`
  Name string `gorm:"size:255;not null; unique" json:"name"`
  CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
  UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *CasePriority) Prepare() {
  c.ID = 0
  c.Name = html.EscapeString(strings.TrimSpace(c.name)
  c.CreatedAt = time.Now()
  c.UpdatedAt = time.Now()
}

func (c *CasePriority) Validate() error {
  if c.Name == "" {
    return errors.New("Required Name")
  }
  return nil
}

func (c *CasePriority) SaveCasePriority(db *gorm.DB) (*CasePriority, error) {
  var err error
  err = db.Debug().Model(&CasePriority{}).Create(&c).Error
  if err != nil {
    return &CasePriority{}, err
  }
  return c, nil
}

func (c *CasePriority) FindAllCasePriorities(db *gorm.DB) (*[]CasePriority, error) {
  var err error
  casePriorities := []CasePriority{}
  err = db.Debug().Model(&CasePriority{}).Limit(100).Find(&casePriorities).Error
  if err != nil {
    return &[]CasePriority{}, err
  }
  return &casePriorities, nil
}

func (c *CasePriority) FindCasePriorityByID(db *gorm.DB, cid uint32) (*CasePriority, error) {
  var err error
  err = db.Debug().Model(&CasePriority{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    return &CasePriority{}, err
  }
  return c, nil
}

func (c *CasePriority) UpdateACasePriority(db *gorm.DB) (*CasePriority, error) {
  var err error

  err = db.Debug().Model(&CasePriority{}).Where("id = ?", c.ID)Updates(CasePriority{Name: c.Name, UpdatedAt: time.Now()}).Error
  if err != nil {
    return &CasePriority{}, err
  }
  return c, nil
}

func (c *CasePriority) DeleteACasePriority(db *gorm.DB, cid uint32) (int64, error) {

  db = db.Debug().Model(&CasePriority{}).Where("id = ?", cid).Take(&CasePriority{}).Delete(&CasePriority{})

  if db.Error != nil {
    if gorm.IsRecordNotFoundError(db.Error) {
      return 0, errors.New("CasePriority not found")
    }
    return 0, db.Error
  }

  return db.RowsAffected, nil
}

