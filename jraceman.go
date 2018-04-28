package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "strconv"

  "github.com/jimmc/jracemango/api"
  "github.com/jimmc/jracemango/app"
  "github.com/jimmc/jracemango/dbrepo"

  // _ "github.com/go-sql-driver/mysql"
  _ "github.com/proullon/ramsql/driver" // driver name: ramsql
  _ "github.com/mattn/go-sqlite3"       // driver name: sqlite3
)

type config struct {
  port int
  uiRoot string
  db string
  create bool
}

func main() {
  config := &config{}

  flag.IntVar(&config.port, "port", 8080, "port on which to listen for connections")
  flag.StringVar(&config.uiRoot, "uiroot", "", "location of ui root (build/default)")
  flag.StringVar(&config.db, "db", "", "location of database, in the form driver:name")
  flag.BoolVar(&config.create, "create", false, "true to create the database specified by -db")

  flag.Parse()

  if config.db == "" {
    log.Fatal("-db is required")
  }

  dbRepos, err := dbrepo.Open(config.db)
  if err != nil {
    log.Fatalf("Failed to open repository: %v", err)
  }
  defer dbRepos.Close()

  if config.create {
    log.Printf("Creating database tables")
    err = dbRepos.CreateTables()
    if err != nil {
      log.Fatalf("Failed to create repository tables: %v", err)
    }
  }

  ph := app.Placeholder{}       // Just to use the app package
  log.Printf("ph is %v", ph)

  mux := http.NewServeMux()

  uiFileHandler := http.FileServer(http.Dir(config.uiRoot))
  apiPrefix := "/api/"
  apiHandler := api.NewHandler(&api.Config{
    Prefix: apiPrefix,
    DomainRepos: dbRepos,
  })
  mux.Handle("/ui/", http.StripPrefix("/ui/", uiFileHandler))
  mux.Handle(apiPrefix, apiHandler)
  // mux.Handle(apiPrefix, authHandler.RequireAuth(apiHandler))
  // mux.Handle("/auth/", authHandler.ApiHandler)
  mux.HandleFunc("/", redirectToUi)

  fmt.Printf("jraceman serving on port %v\n", config.port)
  log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.port), mux))
}

func redirectToUi(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
}
