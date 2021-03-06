package seed
import (
  "log"
  "github.com/jinzhu/gorm"
  "github.com/jawee/case-api/api/models"

)

var users = []models.User {
  models.User{
    Username: "Andreas Olsson",
    Email: "olsson.andreas@gmail.com",
    Password: "password",
  },
  models.User{
    Username: "Evelina Mattsson",
    Email: "some@example.com",
    Password: "password",
  },
}

var casePriorities = []models.CasePriority{
  models.CasePriority{
    Name: "High",
  },
  models.CasePriority{
    Name: "Medium",
  },
  models.CasePriority{
    Name: "Low",
  },
}

var caseStatuses = []models.CaseStatus{
  models.CaseStatus{
    Name: "Received",
  },
  models.CaseStatus{
    Name: "In Progress",
  },
  models.CaseStatus{
    Name: "Solved",
  },
  models.CaseStatus{
    Name: "Closed",
  },
}

var cases = []models.Case{
  models.Case{
    Title: "Case 1",
    Content: "Case 1 content",
  },
  models.Case{
    Title: "Case 2",
    Content: "Case 2 content",
  },
}

func Load(db *gorm.DB) {
  err := db.Debug().DropTableIfExists(&models.Case{}, &models.User{}, &models.CasePriority{}, &models.CaseStatus{}).Error
  if err != nil {
    log.Fatalf("Cannot drop table: %v", err)
  }
  err = db.Debug().AutoMigrate(&models.User{}, &models.CasePriority{}, &models.CaseStatus{}, &models.Case{}).Error
  if err != nil {
    log.Fatalf("cannot migrate table: %v", err)
  }

  // err = db.Debug().Model(&models.Case{}).AddForeignKey("created_by_id", "users(id)", "cascade", "cascade").Error
  // if err != nil {
  //   log.Fatalf("attaching foreign key error: %v", err)
  // }

  // err = db.Debug().Model(&models.Case{}).AddForeignKey("status_id", "case_statuses(id)", "cascade", "cascade").Error
  // if err != nil {
  //   log.Fatalf("attaching foreign key error: %v", err)
  // }

  // err = db.Debug().Model(&models.Case{}).AddForeignKey("priority_id", "case_priorities(id)", "cascade", "cascade").Error
  // if err != nil {
  //   log.Fatalf("attaching foreign key error: %v", err)
  // }

  for i, _ := range casePriorities {
    err = db.Debug().Model(&models.CasePriority{}).Create(&casePriorities[i]).Error
    if err != nil {
      log.Fatalf("cannot seed case priorities table: %v", err)
    }
  }

  for i, _ := range caseStatuses {
    err = db.Debug().Model(&models.CaseStatus{}).Create(&caseStatuses[i]).Error
    if err != nil {
      log.Fatalf("cannot seed case statuses table: %v", err)
    }
  }

  for i, _ := range users {
    err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
    if err != nil {
      log.Fatalf("cannot seed users table: %v", err)
    }
    cases[i].CreatedByID = users[i].ID
    // cases[i].ResponsibleID = users[i].ID
    cases[i].PriorityID = casePriorities[i].ID
    cases[i].StatusID = caseStatuses[i].ID
    err = db.Debug().Model(&models.Case{}).Create(&cases[i]).Error
    if err != nil {
      log.Fatalf("cannot seed cases table: %v", err)
    }
  }

  
}

