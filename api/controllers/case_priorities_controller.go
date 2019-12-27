package controllers

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "strconv"

  "github.com/gorilla/mux"
  "github.com/jawee/case-api/api/auth"
  "github.com/jawee/case-api/api/models"
  "github.com/jawee/case-api/api/responses"
  "github.com/jawee/case-api/api/utils/formaterror"
)

func (server *Server) CreateCasePriority(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  c := models.CasePriority{}
  err = json.Unmarshal(body, &c)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }
  c.Prepare()
  err = c.Validate()
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }
  uid, err := auth.ExtractTokenID(r)
  if err != nil {
    responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
    return
  }
  casePriorityCreated, err := c.SaveCasePriority(server.DB)
  if err != nil {
    formattedError := formaterror.FormatError(err.Error())
    responses.ERROR(w, http.StatusInternalServerError, formattedError)
    return
  }
  w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, casePriorityCreated.ID))
  responses.JSON(w, http.StatusCreated, casePriorityCreated)
}

func (server *Server) GetCasePriorities(w http.ResponseWriter, r *http.Request) {

  c := models.CasePriority{}
  casePriorities, err := c.FindAllCasePriorities(server.DB)
  if err != nil {
    responses.ERROR(w, http.StatusInternalServerError, err)
    return
  }
  responses.JSON(w, http.StatusOK, casePriorities)
}

func (server *Server) GetCasePriority(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  cid, err := strconv.ParseUint(vars["id"], 10, 64)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }
  c := models.CasePriority{}

  casePriorityReceived, err := c.FindCasePriorityByID(server.DB, cid)
  if err != nil {
    responses.ERROR(w, http.StatusInternalServerError, err)
    return
  }
  responses.JSON(w, http.StatusOK, casePriorityReceived)
}

func (server *Server) UpdateCasePriority(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  cid, err := strconv.ParseUint(vars["id"], 10, 64)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }

  uid, err := auth.ExtractTokenID(r)
  if err != nil {
    responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
    return
  }

  c := models.CasePriority{}
  err = server.DB.Debug().Model(models.CasePriority{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    responses.ERROR(w, http.StatusNotFound, errors.New("Case not found"))
    return
  }

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  casePriorityUpdate := models.CasePriority{}
  err = json.Unmarshal(body, &casePriorityUpdate)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  casePriorityUpdate.Prepare()
  err = casePriorityUpdate.Validate()
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  casePriorityUpdate.ID = c.ID

  casePriorityUpdated, err := casePriorityUpdate.UpdateACasePriority(server.DB)

  if err != nil {
    formattedError := formaterror.FormatError(err.Error())
    responses.ERROR(w, http.StatusInternalServerError, formattedError)
    return
  }
  responses.JSON(w, http.StatusOK, casePriorityUpdated)
}

func (server *Server) DeleteCasePriority(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  cid, err := strconv.ParseUint(vars["id"], 10, 64)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }

  uid, err := auth.ExtractTokenID(r)
  if err != nil {
    responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
    return
  }

  c := models.CasePriority{}
  err = server.DB.Debug().Model(models.CasePriority{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
    return
  }

  _, err = c.DeleteACasePriority(server.DB, cid)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }
  w.Header().Set("Entity", fmt.Sprintf("%d", cid))
  responses.JSON(w, http.StatusNoContent, "")
}

