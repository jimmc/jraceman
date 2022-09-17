package main

import (
  "flag"
  "fmt"
  "net/http"
  "os"
  "strconv"

  "github.com/jimmc/jracemango/api"
  "github.com/jimmc/jracemango/app"
  "github.com/jimmc/jracemango/dbrepo"

  "github.com/golang/glog"

  // _ "github.com/go-sql-driver/mysql"
  _ "github.com/mattn/go-sqlite3"       // driver name: sqlite3
)

type config struct {
  // configuration
  port int
  reportRoot string
  uiRoot string
  db string

  // actions
  create bool
  exportFile string
  importFile string
}

func main() {
  os.Exit(doMain())
}

// doMain returns 0 if no errors.
func doMain() int {
  config := &config{}

  // Configuration flags
  flag.IntVar(&config.port, "port", 8080, "port on which to listen for connections")
  flag.StringVar(&config.reportRoot, "reportroot", "report/template", "location of report templates")
  flag.StringVar(&config.uiRoot, "uiroot", "_ui/", "location of ui root")
  flag.StringVar(&config.db, "db", "", "location of database, in the form driver:name")

  // Action flags
  flag.BoolVar(&config.create, "create", false, "true to create the database specified by -db")
  flag.StringVar(&config.exportFile, "export", "", "export the database to a text file")
  flag.StringVar(&config.importFile, "import", "", "import a text file to the database")

  flag.Parse()

  if config.db == "" {
    glog.Error("-db is required")
    return 1
  }

  dbRepos, err := dbrepo.Open(config.db)
  if err != nil {
    glog.Errorf("Failed to open repository: %v", err)
    return 1
  }
  defer dbRepos.Close()

  actionTaken := false

  if config.create {
    glog.Info("Creating database tables")
    err = dbRepos.CreateTables()
    if err != nil {
      glog.Errorf("Failed to create repository tables: %v", err)
      return 1
    }
    actionTaken = true
  }

  if config.exportFile != "" {
    if err := exportFile(config, dbRepos); err != nil {
      glog.Error(err.Error())
      return 1
    }
    actionTaken = true;
  }

  if config.importFile != "" {
    if err := importFile(config, dbRepos); err != nil {
      glog.Error(err.Error())
      return 1
    }
    actionTaken = true;
  }

  _ = app.Placeholder{}       // Just to use the app package

  if actionTaken {
    return 0
  }

  runHttpServer(config, dbRepos)
  return 1      // runHttpServer shouldn't return.
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
  glog.Infof("Exporting to %s\n", config.exportFile)
  if err := dbRepos.Export(outFile); err != nil {
    return fmt.Errorf("error exporting to %s: %v", config.exportFile, err)
  }
  return nil
}

func importFile(config *config, dbRepos *dbrepo.Repos) error {
  glog.Infof("Importing from %s\n", config.importFile)
  counts, err := dbRepos.ImportFile(config.importFile)
  if err != nil {
    return fmt.Errorf("error importing from %s: %v", config.importFile, err)
  }
  glog.Infof("Import done: inserted %d, updated %d, unchanged %d records\n",
      counts.Inserted(), counts.Updated(), counts.Unchanged())
  return nil
}

func runHttpServer(config *config, dbRepos *dbrepo.Repos) {
  mux := http.NewServeMux()

  uiFileHandler := newImportResolver(http.FileServer(http.Dir(config.uiRoot)))
  apiPrefix := "/api/"
  apiHandler := api.NewHandler(&api.Config{
    Prefix: apiPrefix,
    DomainRepos: dbRepos,
    ReportRoots: []string{config.reportRoot},
  })
  mux.Handle("/ui/", http.StripPrefix("/ui/", uiFileHandler))
  mux.Handle(apiPrefix, apiHandler)
  // mux.Handle(apiPrefix, authHandler.RequireAuth(apiHandler))
  // mux.Handle("/auth/", authHandler.ApiHandler)
  mux.HandleFunc("/", redirectToUi)

  fmt.Printf("jraceman serving on port %v\n", config.port)
  glog.Error(http.ListenAndServe(":"+strconv.Itoa(config.port), mux))
  // If we return, that's an error.
}

func redirectToUi(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
}
