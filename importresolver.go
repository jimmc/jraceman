package main

import (
    "errors"
    "net/http"
    "regexp"
    "strings"

    "github.com/golang/glog"
)

var (
    importPattern = regexp.MustCompile(`import *(?:({[^}]*}|[?:(\* * as *)a-zA-Z0-9]+) *from *)?("[^\./][^"]*"|'[^\./][^']*') *;`)
    importFirstMatchIndex = 4   // Index of match value at start of what we want to change.
    importLastMatchIndex = 5    // Index of match value at end of what we want to change.
    exportPattern = regexp.MustCompile(`export *\* *from *("[^\./][^"]*"|'[^\./][^']*') *;`)

    defaultReplacements = map[string]string {
        "@lit/reactive-element": "@lit/reactive-element/reactive-element.js",
        "lit": "lit/index.js",
        "lit-element": "lit-element/index.js",
        "lit-element/development": "lit-element/development/index.js",
        "lit-html": "lit-html/lit-html.js",
    }

    errInvalidWrite = errors.New("invalid write result")
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
    glog.V(3).Infof("in importResolver with request=%v", req)
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
    glog.V(3).Infof("in resolveResponseWriter.Header with header=%v", hdr)
    return hdr
}

func (rrw resolverResponseWriter) WriteHeader(status int) {
    glog.V(2).Infof("in resolveResponseWriter.WriteHeader with status=%v", status)
    // We might change the size of the content we are returning, so don't send
    // that value, and instead set chunked transfer.
    rrw.Header().Del("Content-length")
    rrw.Header().Set("Transfer-Encoding", "chunked")
    rrw.rw.WriteHeader(status)
}

func (rrw resolverResponseWriter) Write(p []byte) (int, error) {
    nr := len(p)        // Number of bytes passed to us.
    glog.V(3).Infof("in resolveResponseWriter.Write with bytes=%v", p)
    glog.V(2).Infof("in resolveResponseWriter.Write for %q, content length before replacement is %d",
      rrw.req.URL, len(p))
    hdr := rrw.rw.Header()
    contentType := hdr.Get("content-type");
    glog.V(2).Infof("in resolveResponseWriter.Write content-type=%v", contentType)
    if strings.HasPrefix(contentType, "text/javascript") ||
      strings.HasPrefix(contentType, "application/javascript") {
        glog.V(2).Infof("Checking javascript from URL=%v", rrw.req.URL)
        p = []byte(rrw.fixImports(string(p)))
        p = []byte(rrw.fixExports(string(p)))
    }
    glog.V(2).Infof("in resolveResponseWriter.Write for %q, content length after replacement is %d",
      rrw.req.URL, len(p))
    nw, err := rrw.rw.Write(p)
    // We are being called from the io.Copy function (io.go:copyBuffer), which checks the
    // return value of bytes written to make sure it matches the value for bytes read from
    // the copy source. If we return a value other than the length of p, our caller treats it
    // as an error and aborts the write. So we need to return nr rather than nw.
    // This bug only shows up with files larger than 32KB, because copyBuffer uses a loop
    // with a buffer size defined as "size := 32 * 1024" bytes.
    if nw != len(p) {
      return nw, errInvalidWrite
    }
    return nr, err
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
    // Our pattern contains three submatches (parens), so each matched import has 6 indexes.
    matchCount := len(importLocations)
    for m:=0; m<matchCount; m++ {
        matchM := importLocations[m]
        // We have one subpattern we are looking for in order to replace that text.
        // Our subpattern includes the quotes, so we adjust the locations inwards to skip those.
        start := matchM[importFirstMatchIndex] + 1
        end := matchM[importLastMatchIndex] - 1
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
