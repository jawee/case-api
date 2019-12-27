package controllers

import (
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"

  _ "github.com/jinzhu/gorm/dialects/mysql"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  _ "github.com/jinzhu/gorm/dialects/sqlite"

  "github.com/jawee/case-api/api/models"
)

type Server struct {
  DB *gorm.DB
  Router *mux.Router
}

func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
  var err error

  if DbDriver == "mysql" {
    DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
    server.DB, err = gorm.Open(DbDriver, DBURL)
    if err != nil {
      fmt.Printf("Cannot connect to %s database", DbDriver)
      log.Fatal("This is the error:", err)
    } else {
      fmt.Printf("We are connected to the %s database", DbDriver)
    }
  }

  if DbDriver == "Postgres" {
    DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
    server.DB, err = gorm.Open(DbDriver, DBURL)
    if err != nil {
      fmt.Printf("Cannot connect to %s database", DbDriver)
      log.Fatal("This is the error:", err)
    } else {
      fmt.Printf("We are connected to the %s database", DbDriver)
    }
  }

  if DbDriver == "sqlite3" {
    server.DB, err = gorm.Open(DbDriver, DbName)
    if err != nil {
      fmt.Printf("Cannot connect to %s database", DbDriver)
      log.Fatal("This is the error:", err)
    } else {
      fmt.Printf("We are connected to the %s database", DbDriver)
    }
  }

  server.DB.Debug().AutoMigrate(&models.User{}, &models.CasePriority{}, &models.CaseStatus{}, &models.Case{})

  server.Router = mux.NewRouter()

  server.initializeRoutes()

}

func (server *Server) Run(addr string) {
  fmt.Println("Listening to port 8080")
  log.Fatal(http.ListenAndServe(addr, server.Router))
}
