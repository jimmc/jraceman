package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "strconv"

  "github.com/jimmc/jracemango/api"
  "github.com/jimmc/jracemango/app"
  "github.com/jimmc/jracemango/dbrepo"

  // _ "github.com/go-sql-driver/mysql"
  _ "github.com/mattn/go-sqlite3"       // driver name: sqlite3
)

type config struct {
  // configuration
  port int
  uiRoot string
  db string

  // actions
  create bool
  exportFile string
  importFile string
}

func main() {
  config := &config{}

  // Configuration flags
  flag.IntVar(&config.port, "port", 8080, "port on which to listen for connections")
  flag.StringVar(&config.uiRoot, "uiroot", "_ui/build/default", "location of ui root (build/default)")
  flag.StringVar(&config.db, "db", "", "location of database, in the form driver:name")

  // Action flags
  flag.BoolVar(&config.create, "create", false, "true to create the database specified by -db")
  flag.StringVar(&config.exportFile, "export", "", "export the database to a text file")
  flag.StringVar(&config.importFile, "import", "", "import a text file to the database")

  flag.Parse()

  if config.db == "" {
    log.Fatal("-db is required")
  }

  dbRepos, err := dbrepo.Open(config.db)
  if err != nil {
    log.Fatalf("Failed to open repository: %v", err)
  }
  defer dbRepos.Close()

  actionTaken := false

  if config.create {
    log.Printf("Creating database tables")
    err = dbRepos.CreateTables()
    if err != nil {
      log.Fatalf("Failed to create repository tables: %v", err)
    }
    actionTaken = true
  }

  if config.exportFile != "" {
    if err := exportFile(config, dbRepos); err != nil {
      log.Fatalf(err.Error())
    }
    actionTaken = true;
  }

  if config.importFile != "" {
    if err := importFile(config, dbRepos); err != nil {
      log.Fatalf(err.Error())
    }
    actionTaken = true;
  }

  _ = app.Placeholder{}       // Just to use the app package

  if !actionTaken {
    runHttpServer(config, dbRepos)
    // Doesn't return.
  }
}

func exportFile(config *config, dbRepos *dbrepo.Repos) error {
  if _, err := os.Stat(config.exportFile); !os.IsNotExist(err) {
    return fmt.Errorf("output file %s exists, will not overwrite", config.exportFile)
  }
  outFile, err := os.Create(config.exportFile)
  if err != nil {
    return fmt.Errorf("error opening export output file %s: %v", config.exportFile, err)
  }
  defer outFile.Close()
  log.Printf("Exporting to %s\n", config.exportFile)
  if err := dbRepos.Export(outFile); err != nil {
    return fmt.Errorf("error exporting to %s: %v", config.exportFile, err)
  }
  return nil
}

func importFile(config *config, dbRepos *dbrepo.Repos) error {
  inFile, err := os.Open(config.importFile)
  if err != nil {
    return fmt.Errorf("error opening import imput file %s: %v", config.importFile, err)
  }
  defer inFile.Close()

  log.Printf("Importing from %s\n", config.importFile)
  insertCount, updateCount, unchangedCount, err := dbRepos.Import(inFile)
  if err != nil {
    return fmt.Errorf("error importing from %s: %v", config.importFile, err)
  }
  log.Printf("Import done: inserted %d, updated %d, unchanged %d records\n",
      insertCount, updateCount, unchangedCount)
  return nil
}

func runHttpServer(config *config, dbRepos *dbrepo.Repos) {
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
