package main

import (
  "flag"
  "fmt"
  "net/http"
  "os"
  "strconv"

  "github.com/jimmc/jracemango/api"
  "github.com/jimmc/jracemango/app"
  "github.com/jimmc/jracemango/auth"
  "github.com/jimmc/jracemango/dbrepo"
  "github.com/jimmc/jracemango/dbrepo/strsql"

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
  maxClockSkewSeconds int
  password string

  // actions
  checkUpgrade bool
  create bool
  exportFile string
  importFile string
  sql string
  updatePassword string
  upgrade bool
  version bool
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
  flag.IntVar(&config.maxClockSkewSeconds, "maxClockSkewSeconds", 5,
        "max seconds of clock skew allowed between client and server")
  flag.StringVar(&config.password, "password", "", "password for update (for testing)")

  // Action flags
  flag.BoolVar(&config.checkUpgrade, "checkUpgrade", false, "true to check for upgrade the database specified by -db")
  flag.BoolVar(&config.create, "create", false, "true to create the database specified by -db")
  flag.StringVar(&config.exportFile, "export", "", "export the database to a text file")
  flag.StringVar(&config.importFile, "import", "", "import a text file to the database")
  flag.StringVar(&config.sql, "sql", "", "execute sql statement")
  flag.StringVar(&config.updatePassword, "updatePassword", "", "update password for named user")
  flag.BoolVar(&config.upgrade, "upgrade", false, "true to upgrade the database specified by -db")
  flag.BoolVar(&config.version, "version", false, "print the version")

  flag.Parse()

  if config.version {
    fmt.Printf("jraceman version %s\n", Version)
    return 1
  }

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
    err = dbRepos.CreateTables()
    if err != nil {
      glog.Errorf("Failed to create repository tables: %v", err)
      return 1
    }
    actionTaken = true
  }

  if config.checkUpgrade {
    err = dbRepos.UpgradeTables(true/*dryrun*/)
    if err != nil {
      glog.Errorf("Error checking for upgrading database: %v", err)
      return 1
    }
    actionTaken = true
  }

  if config.upgrade {
    glog.Info("Upgrading database")
    err = dbRepos.UpgradeTables(false/*dryrun*/)
    if err != nil {
      glog.Errorf("Error upgrading database: %v", err)
      return 1
    }
    actionTaken = true
  }

  if config.exportFile != "" {
    if err := exportFile(config, dbRepos); err != nil {
      glog.Error(err.Error())
      return 1
    }
    actionTaken = true
  }

  if config.importFile != "" {
    if err := importFile(config, dbRepos); err != nil {
      glog.Error(err.Error())
      return 1
    }
    actionTaken = true
  }

  if config.updatePassword != "" {
    var err error
    authHandler := auth.NewHandler(dbRepos.DB(), config.maxClockSkewSeconds)
    if config.password == "" {
      err = authHandler.UpdateUserPassword(config.updatePassword)
    } else {
      err = authHandler.UpdatePassword(config.updatePassword, config.password)
    }
    if err != nil {
      glog.Errorf("Error updating password for %s: %v", config.updatePassword, err)
      return 1
    }
    fmt.Printf("Password updated for user %s\n", config.updatePassword)
    actionTaken = true
  }

  if config.sql != "" {
    db := dbRepos.DB()
    result, err := strsql.QueryStarAndCollect(db, config.sql)
    if err != nil {
      fmt.Printf("Error executing sql: %v\n", err)
    } else {
      fmt.Printf("Results:\n%v\n", result)
    }
    actionTaken = true
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
  authHandler := auth.NewHandler(dbRepos.DB(), config.maxClockSkewSeconds)

  uiFileHandler := newImportResolver(http.FileServer(http.Dir(config.uiRoot)))
  apiPrefix := "/api/"
  apiHandler := api.NewHandler(&api.Config{
    Prefix: apiPrefix,
    DomainRepos: dbRepos,
    ReportRoots: []string{config.reportRoot},
  })
  mux.Handle("/ui/", http.StripPrefix("/ui/", uiFileHandler))
  mux.Handle(apiPrefix, authHandler.RequireAuth(apiHandler))
  mux.Handle("/auth/", authHandler.ApiHandler)
  mux.HandleFunc("/", redirectToUi)
  mux.Handle("/api0/", api.NewHandler0("/api0/", Version))

  fmt.Printf("jraceman serving on port %v\n", config.port)
  glog.Error(http.ListenAndServe(":"+strconv.Itoa(config.port), mux))
  // If we return, that's an error.
}

func redirectToUi(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
}
