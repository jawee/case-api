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

func (server *Server) CreateCase(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  c := models.Case{}
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
  if uid != c.CreatedByID {
    responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
    return
  }
  caseCreated, err := c.SaveCase(server.DB)
  if err != nil {
    formattedError := formaterror.FormatError(err.Error())
    responses.ERROR(w, http.StatusInternalServerError, formattedError)
    return
  }
  w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, caseCreated.ID))
  responses.JSON(w, http.StatusCreated, caseCreated)
}

func (server *Server) GetCases(w http.ResponseWriter, r *http.Request) {

  c := models.Case{}
  cases, err := c.FindAllCases(server.DB)
  if err != nil {
    responses.ERROR(w, http.StatusInternalServerError, err)
    return
  }
  responses.JSON(w, http.StatusOK, cases)
}

func (server *Server) GetCase(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  cid, err := strconv.ParseUint(vars["id"], 10, 64)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }
  c := models.Case{}

  caseReceived, err := c.FindCaseByID(server.DB, cid)
  if err != nil {
    responses.ERROR(w, http.StatusInternalServerError, err)
    return
  }
  responses.JSON(w, http.StatusOK, caseReceived)
}

func (server *Server) UpdateCase(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  // Check if the case id is valid
  cid, err := strconv.ParseUint(vars["id"], 10, 64)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }

  //CHeck if the auth token is valid and get the user id from it
  uid, err := auth.ExtractTokenID(r)
  if err != nil {
    responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
    return
  }

  c := models.Case{}
  err = server.DB.Debug().Model(models.Case{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    responses.ERROR(w, http.StatusNotFound, errors.New("Case not found"))
    return
  }

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  caseUpdate := models.Case{}
  err = json.Unmarshal(body, &caseUpdate)
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  caseUpdate.Prepare()
  err = caseUpdate.Validate()
  if err != nil {
    responses.ERROR(w, http.StatusUnprocessableEntity, err)
    return
  }

  caseUpdate.ID = c.ID

  caseUpdated, err := caseUpdate.UpdateACase(server.DB)

  if err != nil {
    formattedError := formaterror.FormatError(err.Error())
    responses.ERROR(w, http.StatusInternalServerError, formattedError)
    return
  }
  responses.JSON(w, http.StatusOK, caseUpdated)
}

func (server *Server) DeleteCase(w http.ResponseWriter, r *http.Request) {
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

  c := models.Case{}
  err = server.DB.Debug().Model(models.Case{}).Where("id = ?", cid).Take(&c).Error
  if err != nil {
    responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
    return
  }

  _, err = c.DeleteACase(server.DB, cid)
  if err != nil {
    responses.ERROR(w, http.StatusBadRequest, err)
    return
  }
  w.Header().Set("Entity", fmt.Sprintf("%d", cid))
  responses.JSON(w, http.StatusNoContent, "")
}

