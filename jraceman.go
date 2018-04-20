package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "strconv"

  "github.com/jimmc/jracemango/domain"
  "github.com/jimmc/jracemango/api"
  "github.com/jimmc/jracemango/app"
)

type config struct {
  port int
  uiRoot string
}

type SiteRepoTest struct {}
func (r *SiteRepoTest) FindById(ID string) (domain.Site, error) {
  return domain.Site{
    ID: ID,
    Name: "Site" + ID,
  }, nil
}
func (r *SiteRepoTest) Save(site domain.Site) error {
  return nil
}

func main() {
  config := &config{}

  flag.IntVar(&config.port, "port", 8080, "port on which to listen for connections")
  flag.StringVar(&config.uiRoot, "uiroot", "", "location of ui root (build/default)")

  flag.Parse()

  domainRepos := &domain.Repos{
    Site: &SiteRepoTest{},
  }
  ph := app.Placeholder{}       // Just to use the app package
  log.Printf("ph is %v", ph)

  mux := http.NewServeMux()

  uiFileHandler := http.FileServer(http.Dir(config.uiRoot))
  apiPrefix := "/api/"
  apiHandler := api.NewHandler(&api.Config{
    Prefix: apiPrefix,
    DomainRepos: domainRepos,
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
