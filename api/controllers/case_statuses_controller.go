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

func (server *Server) CreateCaseStatus(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  c := models.CaseStatus{}
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
  caseStatusCreated, err := c.SaveCaseStatus(server.DB)
  if err != nil {
    formattedError := formaterror.FormatError(err.Error())
    responses.ERROR(w, http.StatusInternalServerError, formattedError)
    return
  }
  w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, caseStatusCreated.ID))
  responses.JSON(w, http.StatusCreated, caseStatusCreated)
}

func (server *Server) GetCaseStatuses(w http.ResponseWriter, r *http.Request) {

  c := models.CaseStatus{}
  caseStatuses, err := c.FindAllCaseStatuses(server.DB)
  if err != nil {
    responses.ERROR(w, http.StatusInternalServerError, err)
    return
  }
  responses.JSON(w, http.StatusOK, caseStatuses)
}

func (server *Server) GetCaseStatus(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  cid, err := strconv.ParseUint(vars["id"], 10, 64)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }
  c := models.CaseStatus{}

  caseStatusReceived, err := c.FindCaseStatusByID(server.DB, cid)
  if err != nil {
    responses.ERROR(w, http.StatusInternalServerError, err)
    return
  }
  responses.JSON(w, http.StatusOK, caseStatusReceived)
}

func (server *Server) UpdateCaseStatus(w http.ResponseWriter, r *http.Request) {
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

  c := models.CaseStatus{}
  err = server.DB.Debug().Model(models.CaseStatus{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    responses.ERROR(w, http.StatusNotFound, errors.New("Case not found"))
    return
  }

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  caseStatusUpdate := models.CaseStatus{}
  err = json.Unmarshal(body, &caseStatusUpdate)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  caseStatusUpdate.Prepare()
  err = caseStatusUpdate.Validate()
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  caseStatusUpdate.ID = c.ID

  caseStatusUpdated, err := caseStatusUpdate.UpdateACaseStatus(server.DB)

  if err != nil {
    formattedError := formaterror.FormatError(err.Error())
    responses.ERROR(w, http.StatusInternalServerError, formattedError)
    return
  }
  responses.JSON(w, http.StatusOK, caseStatusUpdated)
}

func (server *Server) DeleteCaseStatus(w http.ResponseWriter, r *http.Request) {
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

  c := models.CaseStatus{}
  err = server.DB.Debug().Model(models.CaseStatus{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
    return
  }

  _, err = c.DeleteACaseStatus(server.DB, cid)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }
  w.Header().Set("Entity", fmt.Sprintf("%d", cid))
  responses.JSON(w, http.StatusNoContent, "")
}

