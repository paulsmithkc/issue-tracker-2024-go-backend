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

func ResSuccess(w http.ResponseWriter, data any) (int, error) {
  var status int = http.StatusOK
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(status)
  json.NewEncoder(w).Encode(ResData { true, data, nil })
  log.Printf("%d\n", 200);
  return status, nil
}

func ResError(w http.ResponseWriter, status int, err error, data any) (int, error) {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(status)
  json.NewEncoder(w).Encode(ResData { false, data, fmt.Sprintf("%v", err) })
  log.Printf("%d %v\n", status, err);
  return status, nil
}

func routeIssueApi(w http.ResponseWriter, r *http.Request) {
  var err error
  var status int = http.StatusInternalServerError

  method := r.Method
  slug := r.PathValue("slug")
  log.Printf("%s /issue/%s\n", method, slug);

  switch slug {
  case "list":
    switch method {
    case "GET": status, err = listIssues(w, r)
    default: status, err = 405, fmt.Errorf("unsupported method: %s /issue/%s", method, slug)
    }
  case "new":
    switch method {
      case "POST": status, err = insertIssue(w, r, slug)
      default: status, err = 405, fmt.Errorf("unsupported method: %s /issue/%s", method, slug)
    }
  default:
    switch method {
    case "GET": status, err = getIssue(w, r, slug)
    case "POST": status, err = updateIssue(w, r, slug)
    case "DELETE": status, err = deleteIssue(w, r, slug)
    default: status, err = 405, fmt.Errorf("unsupported method: %s /issue/%s", method, slug)
    }
  }

  if (err != nil) {
    ResError(w, status, err, nil)
  }
}

func listIssues(w http.ResponseWriter, r *http.Request) (int, error) {
  db, err := OpenDB();
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Open: %w", err), nil)
  }

  query := "SELECT id, title, description, created_at FROM issues ORDER BY created_at desc LIMIT 100"
  rows, err := db.Query(query)
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Query: %w", err), nil)
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
      return ResError(w, 500, fmt.Errorf("Database.Scan: %w", err), nil)
    }
    items = append(items, Issue { id, title, description, created_at })
  }

  return ResSuccess(w, items)
}

func getIssue(w http.ResponseWriter, r *http.Request, slug string) (int, error) {
  db, err := OpenDB();
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Open: %w", err), nil)
  }

  var id string
  var title string
  var description string;
  var created_at string;

  query := "SELECT id, title, description, created_at FROM issues WHERE id = $1 LIMIT 1"
  err = db.QueryRow(query, slug).Scan(&id, &title, &description, &created_at)
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Query: %w", err), nil)
  }

  item := Issue { id, title, description, created_at }
  return ResSuccess(w, item)
}

func insertIssue(w http.ResponseWriter, r *http.Request, slug string) (int, error) {
  db, err := OpenDB();
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Open: %w", err), nil)
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
    return ResError(w, 500, fmt.Errorf("Database.Exec: %w", err), nil)
  }

  item := Issue { d.Id, d.Title, d.Description, "" }
  return ResSuccess(w, item)
}

func updateIssue(w http.ResponseWriter, r *http.Request, slug string) (int, error) {
  db, err := OpenDB();
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Open: %w", err), nil)
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
    return ResError(w, 500, fmt.Errorf("Database.Exec: %w", err), nil)
  }

  item := Issue { slug, d.Title, d.Description, "" }
  return ResSuccess(w, item)
}

func deleteIssue(w http.ResponseWriter, r *http.Request, slug string) (int, error) {
  db, err := OpenDB();
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Open: %w", err), nil)
  }

  query := "DELETE FROM issues WHERE id = $1"
  _, err = db.Exec(query, slug)
  if err != nil {
    return ResError(w, 500, fmt.Errorf("Database.Exec: %w", err), nil)
  }

  item := Issue { slug, "", "", "" }
  return ResSuccess(w, item)
}
