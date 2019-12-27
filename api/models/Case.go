package models

import (
  "errors"
  "html"
  "strings"
  "time"

  "github.com/jinzhu/gorm"
)

type Case struct {
  ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
  Title string `gorm:"size:255;not null;" json:"title"`
  Content string `gorm:"type:text;not null;" json:"content"`
  CreatedBy User `json:"created_by"`
  CreatedByID uint32 `gorm:"not null" json:"created_by_id"`
  Responsible User `json:"responsible"`
  ResponsibleID uint32 `gorm:"" json:"responsible_id"`
  CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
  UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
  Status CaseStatus `json:"status"`
  StatusID uint32 `gorm:"not null" json:"status_id"`
  Priority CasePriority `json:"priority"`
  PriorityID uint32 `gorm:"not null" json:"priority_id"`
}

func (c *Case) Prepare() {
  c.ID = 0
  c.Title = html.EscapeString(strings.TrimSpace(c.Title))
  c.Content = html.EscapeString(strings.TrimSpace(c.Content))
  c.CreatedBy = User{}
  c.Status = CaseStatus{}
  c.Priority = CasePriority{}
  c.CreatedAt = time.Now()
  c.UpdatedAt = time.Now()
  c.ResponsibleID = 0
  c.Responsible = User{}
}

func (c *Case) Validate() error {
  if c.Title == "" {
    return errors.New("Required Title")
  }
  if c.Content == "" {
    return errors.New("Required Content")
  }
  if c.CreatedByID < 1 {
    return errors.New("Required Author")
  }
  if c.PriorityID < 1 {
    return errors.New("Required Priority")
  }
  if c.StatusID < 1 {
    return errors.New("Required Status")
  }
  return nil
}

func (c *Case) SaveCase(db *gorm.DB) (*Case, error) {
  var err error
  err = db.Debug().Model(&Case{}).Create(&c).Error
  if err != nil {
    return &Case{}, err
  }
  if c.ID != 0 {
    err = db.Debug().Model(&User{}).Where("id = ?", c.CreatedByID).Take(&c.CreatedBy).Error
    if err != nil {
      return &Case{}, err
    }
    err = db.Debug().Model(&CaseStatus{}).Where("id = ?", c.StatusID).Take(&c.Status).Error
    if err != nil {
      return &Case{}, err
    }
    err = db.Debug().Model(&CasePriority{}).Where("id = ?", c.PriorityID).Take(&c.Priority).Error
    if err != nil {
      return &Case{}, err
    }
    if c.ResponsibleID != 0 {
      err = db.Debug().Model(&User{}).Where("id = ?", c.ResponsibleID).Take(&c.Responsible).Error
      if err != nil {
        return &Case{}, err
      }
    }
  }
  return c, nil
}

func (c *Case) FindAllCases(db *gorm.DB) (*[]Case, error) {
  var err error
  cases := []Case{}
  err = db.Debug().Model(&Case{}).Limit(100).Find(&cases).Error
  if err != nil {
    return &[]Case{}, err
  }

  if len(cases) > 0 {
    for i, _ := range cases {
      err := db.Debug().Model(&User{}).Where("id = ?", cases[i].CreatedByID).Take(&cases[i].CreatedBy).Error
      if err != nil {
        return &[]Case{}, err
      }
       err = db.Debug().Model(&CaseStatus{}).Where("id = ?", cases[i].StatusID).Take(&cases[i].Status).Error
      if err != nil {
        return &[]Case{}, err
      }
      err = db.Debug().Model(&CasePriority{}).Where("id = ?", cases[i].PriorityID).Take(&cases[i].Priority).Error
      if err != nil {
        return &[]Case{}, err
      }
      if cases[i].ResponsibleID != 0 {
        err = db.Debug().Model(&User{}).Where("id = ?", cases[i].ResponsibleID).Take(&cases[i].Responsible).Error
        if err != nil {
          return &[]Case{}, err
        }
      }
    }
  }
  return &cases, nil
}

func (c *Case) FindCaseByID(db *gorm.DB, cid uint64) (*Case, error) {
  var err error
  err = db.Debug().Model(&Case{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    return &Case{}, err
  }
  if c.ID != 0 {
    err = db.Debug().Model(&User{}).Where("id = ?", c.CreatedByID).Take(&c.CreatedBy).Error
    if err != nil {
      return &Case{}, err
    }
    err = db.Debug().Model(&CaseStatus{}).Where("id = ?", c.StatusID).Take(&c.Status).Error
    if err != nil {
      return &Case{}, err
    }
    err = db.Debug().Model(&CasePriority{}).Where("id = ?", c.PriorityID).Take(&c.Priority).Error
    if err != nil {
      return &Case{}, err
    }
    if c.ResponsibleID != 0 {
      err = db.Debug().Model(&User{}).Where("id = ?", c.ResponsibleID).Take(&c.Responsible).Error
      if err != nil {
        return &Case{}, err
      }
    }
  }
  return c, nil
}


func (c *Case) UpdateACase(db *gorm.DB) (*Case, error) {
  var err error

  err = db.Debug().Model(&Case{}).Where("id = ?", c.ID).Updates(Case{Title: c.Title, Content: c.Content, UpdatedAt: time.Now()}).Error
  if err != nil {
    return &Case{}, err
  }
  if c.ID != 0 {
    err = db.Debug().Model(&User{}).Where("id = ?", c.CreatedByID).Take(&c.CreatedBy).Error
    if err != nil {
      return &Case{}, err
    }
    err = db.Debug().Model(&CaseStatus{}).Where("id = ?", c.StatusID).Take(&c.Status).Error
    if err != nil {
      return &Case{}, err
    }
    err = db.Debug().Model(&CasePriority{}).Where("id = ?", c.PriorityID).Take(&c.Priority).Error
    if err != nil {
      return &Case{}, err
    }
    if c.ResponsibleID != 0 {
      err = db.Debug().Model(&User{}).Where("id = ?", c.ResponsibleID).Take(&c.Responsible).Error
      if err != nil {
        return &Case{}, err
      }
    }
  }
  return c, nil
}

func (c *Case) DeleteACase(db *gorm.DB, cid uint64) (int64, error) {

  db.Debug().Model(&Case{}).Where("id = ?", cid).Take(&Case{}).Delete(&Case{})

  if db.Error != nil {
    if gorm.IsRecordNotFoundError(db.Error) {
      return 0, errors.New("Case not found")
    }
    return 0, db.Error
  }
  return db.RowsAffected, nil
}
