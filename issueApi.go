package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
  log.Println("Init issueApi")
  functions.HTTP("issue/{slug}", routeIssueApi)
}

type Issue struct {
  Id string `json:"id"`
  Title string `json:"title"`
  Description string `json:"description"`
  CreatedAt string `json:"created_at"`
}

type ResData struct {
  Success bool `json:"success"`
  Data any `json:"data"`
  Error any `json:"error"`
}

func ResSuccess(w http.ResponseWriter, data any) error {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(ResData { true, data, nil })
  log.Printf("%d\n", 200);
  return nil
}

func ResError(w http.ResponseWriter, status int, err error, data any) error {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(ResData { false, data, fmt.Sprintf("%v", err) })
  log.Printf("%d %v\n", status, err);
  return nil
}

func routeIssueApi(w http.ResponseWriter, r *http.Request) {
  var err error
  var status int

  method := r.Method
  slug := r.PathValue("slug")
  log.Printf("%s /issue/%s\n", method, slug);

  switch slug {
  case "list":
    switch method {
    case "GET": err, status = listIssues(w, r), 500
    default: err, status = fmt.Errorf("unsupported method: %s", method), 405
    }
  case "new":
    switch method {
      case "POST": err, status = insertIssue(w, r, slug), 500
      default: err, status = fmt.Errorf("unsupported method: %s", method), 405
    }
  default:
    switch method {
    case "GET": err, status = getIssue(w, r, slug), 500
    case "POST": err, status = updateIssue(w, r, slug), 500
    case "DELETE": err, status = deleteIssue(w, r, slug), 500
    default: err, status = fmt.Errorf("unsupported method: %s", method), 405
    }
  }

  if (err != nil) {
    ResError(w, status, err, nil)
  }
}

func listIssues(w http.ResponseWriter, r *http.Request) error {
  db, err := OpenDB();
  if err != nil {
    return err
  }

  query := "SELECT id, title, description, created_at FROM issues ORDER BY created_at desc LIMIT 100"
  rows, err := db.Query(query)
  if err != nil {
    return fmt.Errorf("Database.Query: %w", err)
  }

  defer rows.Close()

  items := []Issue{}
  var id string
  var title string
  var description string;
  var created_at string;

  for rows.Next() {
    err := rows.Scan(&id, &title, &description, &created_at);
    if err != nil {
      return fmt.Errorf("Database.Scan: %w", err)
    }
    items = append(items, Issue { id, title, description, created_at })
  }

  return ResSuccess(w, items)
}

func getIssue(w http.ResponseWriter, r *http.Request, slug string) error {
  db, err := OpenDB();
  if err != nil {
    return err
  }

  var id string
  var title string
  var description string;
  var created_at string;

  query := "SELECT id, title, description, created_at FROM issues WHERE id = $1 LIMIT 1"
  err = db.QueryRow(query, slug).Scan(&id, &title, &description, &created_at)
  if err != nil {
    return fmt.Errorf("Database.Query: %w", err)
  }

  item := Issue { id, title, description, created_at }
  return ResSuccess(w, item)
}

func insertIssue(w http.ResponseWriter, r *http.Request, slug string) error {
  db, err := OpenDB();
  if err != nil {
    return err
  }

  contentType := r.Header.Get("Content-Type")
  if (contentType != "application/json") {
    err = fmt.Errorf("unsupported content type: %s", contentType)
    return ResError(w, 400, err, nil)
  }

  var d struct {
    Id string `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
  }
  err = json.NewDecoder(r.Body).Decode(&d)
  if err != nil {
    err = fmt.Errorf("request body invalid: %w", err)
    return ResError(w, 400, err, nil)
  }
  if (d.Id == "" ||d.Title == "" || d.Description == "") {
    err = fmt.Errorf("missing required fields")
    return ResError(w, 400, err, d)
  }

  query := "INSERT INTO issues (id, title, description) VALUES ($1, $2, $3)"
  _, err = db.Exec(query, d.Id, d.Title, d.Description)
  if err != nil {
    return fmt.Errorf("Database.Exec: %w", err)
  }

  item := Issue { d.Id, d.Title, d.Description, "" }
  return ResSuccess(w, item)
}

func updateIssue(w http.ResponseWriter, r *http.Request, slug string) error {
  db, err := OpenDB();
  if err != nil {
    return err
  }

  contentType := r.Header.Get("Content-Type")
  if (contentType != "application/json") {
    err = fmt.Errorf("unsupported content type: %s", contentType)
    return ResError(w, 400, err, nil)
  }

  var d struct {
    Title string `json:"title"`
    Description string `json:"description"`
  }
  err = json.NewDecoder(r.Body).Decode(&d)
  if err != nil {
    err = fmt.Errorf("request body invalid: %w", err)
    return ResError(w, 400, err, nil)
  }
  if (d.Title == "" || d.Description == "") {
    err = fmt.Errorf("missing required fields")
    return ResError(w, 400, err, d)
  }

  query := "UPDATE issues SET title = $2, description = $3 WHERE id = $1"
  _, err = db.Exec(query, slug, d.Title, d.Description)
  if err != nil {
    return fmt.Errorf("Database.Exec: %w", err)
  }

  item := Issue { slug, d.Title, d.Description, "" }
  return ResSuccess(w, item)
}

func deleteIssue(w http.ResponseWriter, r *http.Request, slug string) error {
  db, err := OpenDB();
  if err != nil {
    return err
  }

  query := "DELETE FROM issues WHERE id = $1"
  _, err = db.Exec(query, slug)
  if err != nil {
    return fmt.Errorf("Database.Exec: %w", err)
  }

  item := Issue { slug, "", "", "" }
  return ResSuccess(w, item)
}
