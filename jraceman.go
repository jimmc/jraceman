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

func main() {
  config := &config{}

  flag.IntVar(&config.port, "port", 8080, "port on which to listen for connections")
  flag.StringVar(&config.uiRoot, "uiroot", "", "location of ui root (build/default)")

  flag.Parse()

  area := domain.Area{}         // Just to use the domain package
  ph := app.Placeholder{}       // Just to use the app package
  log.Printf("area is %v, ph is %v", area, ph)

  mux := http.NewServeMux()

  uiFileHandler := http.FileServer(http.Dir(config.uiRoot))
  apiHandler := api.NewHandler(&api.Config{
    Prefix: "/api/",
  })
  mux.Handle("/ui/", http.StripPrefix("/ui/", uiFileHandler))
  mux.Handle("/api/", apiHandler)
  // mux.Handle("/api/", authHandler.RequireAuth(apiHandler))
  // mux.Handle("/auth/", authHandler.ApiHandler)
  mux.HandleFunc("/", redirectToUi)

  fmt.Printf("jraceman serving on port %v\n", config.port)
  log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.port), mux))
}

func redirectToUi(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
}
