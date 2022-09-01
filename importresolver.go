package main

import (
    "net/http"
    "regexp"
    "strings"

    "github.com/golang/glog"
)

var (
    importPattern = regexp.MustCompile(`import *(?:{[^}]*} *from *)?("[^\./][^"]*"|'[^\./][^']*') *;`)
    exportPattern = regexp.MustCompile(`export *\* *from *("[^\./][^"]*"|'[^\./][^']*') *;`)

    defaultReplacements = map[string]string {
        "@lit/reactive-element": "@lit/reactive-element/reactive-element.js",
        "lit": "lit/index.js",
        "lit-element": "lit-element/index.js",
        "lit-element/development": "lit-element/development/index.js",
        "lit-html": "lit-html/lit-html.js",
    }
)

/** importResolver looks for .js files in an http request, then intercepts
 * the content, looks for bare module import statements, and resolves them
 * by modifying the import path to include a prefix and .js extension
 * as needed.
 */
type importResolver struct {
    fileHandler http.Handler
}

func newImportResolver(fileHandler http.Handler) http.Handler {
    return importResolver{fileHandler}
}

func (ir importResolver) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    responseWriterWrapper := resolverResponseWriter{
        req: req,
        rw: rw,
        customReplacements: defaultReplacements,
    }
    // glog.V(1).Infof("in importResolver with request=%v", req)
    ir.fileHandler.ServeHTTP(responseWriterWrapper, req)
}

/** resolveResponseWriter intercepts the data stream of the response so that
 * we can look for import statements that are bare module imports and rewrite
 * them to be valid ES module imports.
 */
type resolverResponseWriter struct {
    req *http.Request
    rw http.ResponseWriter
    customReplacements map[string]string
}

func (rrw resolverResponseWriter) Header() http.Header {
    hdr := rrw.rw.Header()
    // glog.V(2).Infof("in resolveResponseWriter.Header with header=%v", hdr)
    return hdr
}

func (rrw resolverResponseWriter) WriteHeader(status int) {
    // glog.V(2).Infof("in resolveResponseWriter.WriteHeader with status=%v", status)
    // We might change the size of the content we are returning, so don't send
    // that value, and instead set chunked transfer.
    rrw.Header().Del("Content-length")
    rrw.Header().Set("Transfer-Encoding", "chunked")
    rrw.rw.WriteHeader(status)
}

func (rrw resolverResponseWriter) Write(p []byte) (int, error) {
    // glog.V(2).Infof("in resolveResponseWriter.Write with bytes=%v", p)
    hdr := rrw.rw.Header()
    contentType := hdr.Get("content-type");
    // glog.V(2).Infof("in resolveResponseWriter.Write content-type=%v", contentType)
    if strings.HasPrefix(contentType, "text/javascript") {
        glog.V(2).Infof("Checking javascript from URL=%v", rrw.req.URL)
        p = []byte(rrw.fixImports(string(p)))
        p = []byte(rrw.fixExports(string(p)))
    }
    n, err := rrw.rw.Write(p)
    return n, err
}

/** fixImports looks for base module imports and rewrites them.
 */
func (rrw resolverResponseWriter) fixImports(p string) string {
    // Use regexp to look for import patterns
    // import "filename"; 
    //   or
    // import { tags } from "filename";
    //   with single or double-quotes around filename
    //   and where filename does not start with . or /
    importLocations := importPattern.FindAllStringSubmatchIndex(p, -1)
    if importLocations == nil {
        return p        // No matching import statements found
    }
    carat := 0
    var sb strings.Builder
    // Our pattern contains one submatch, so each matched import has 4 indexes.
    matchCount := len(importLocations)
    for m:=0; m<matchCount; m++ {
        matchM := importLocations[m]
        // We have one subpattern, so we use indexes 2 and 3.
        // Our subpattern includes the quotes, so we adjust the locations inwards to skip those.
        start := matchM[2] + 1
        end := matchM[3] - 1
        filename := p[start:end]
        resolvedFilename := rrw.resolveFilename(filename)
        glog.V(2).Infof("Replaced %q by %q", filename, resolvedFilename)
        sb.WriteString(p[carat:start])
        sb.WriteString(resolvedFilename)
        carat = end
    }
    sb.WriteString(p[carat:])
    // Look for
    //   export * from "filename";
    //   where filename does not start with . or /
    // Do the same filename processing as for import
    return sb.String()
}

/** fixImports looks for base module exports and rewrites them.
 */
func (rrw resolverResponseWriter) fixExports(p string) string {
    // Use regexp to look for export patterns
    // Look for
    //   export * from "filename";
    //   with single or double-quotes around filename
    //   and where filename does not start with . or /
    // Do the same filename processing as for import
    exportLocations := exportPattern.FindAllStringSubmatchIndex(p, -1)
    if exportLocations == nil {
        return p        // No matching export statements found
    }
    carat := 0
    var sb strings.Builder
    // Our pattern contains one submatch, so each matched export has 4 indexes.
    matchCount := len(exportLocations)
    for m:=0; m<matchCount; m++ {
        matchM := exportLocations[m]
        // We have one subpattern, so we use indexes 2 and 3.
        // Our subpattern includes the quotes, so we adjust the locations inwards to skip those.
        start := matchM[2] + 1
        end := matchM[3] - 1
        filename := p[start:end]
        resolvedFilename := rrw.resolveFilename(filename)
        glog.V(2).Infof("Replaced %q by %q", filename, resolvedFilename)
        sb.WriteString(p[carat:start])
        sb.WriteString(resolvedFilename)
        carat = end
    }
    sb.WriteString(p[carat:])
    return sb.String()
}

func (rrw resolverResponseWriter) resolveFilename(filename string) string {
    if customFix, ok := rrw.customReplacements[filename]; ok {
        filename = customFix
    }
    // Prepend node_modules/ to filename.
    return "/ui/node_modules/" + filename
}
