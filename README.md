# jracemango

Jracemango is a rewrite of JRaceman as a web app with a golang server.

## Quick start

### Install Go compiler

See if the go compiler is already installed on your system:

    go version

If not, [install Go](https://go.dev/doc/install).

### Get the jracemango sources

Once the Go compiler is installed, it can do this for you automatically:

    go get github.com/jimmc/jracemango

### Build

Change your working directory to the jracemango directory:

    cd ~/go/src/github.com/jimmc/jracemango

Compile jraceman:

    go build

This creates the executable `jracemango` in the jracemango directory.

### Test

Run the unit tests:

    go test ./...

If you want to check the unit test coverage:

    go test ./... -coverprofile=cover.out
    go tool cover -html=cover.out -o cover.html
    # Open cover.html in your browser

### Build the UI pages

Follow the instructions in the [\_ui](./_ui) directory.

### Create a new database

Select a location for your database, for example $HOME/jrdb, then run the jraceman binary
specifying that database:

    mkdir _private      # Put files here you don't want to checkin to git
    ./jracemango -db sqlite3:_private/jrdb-test -create

You can import a jraceman data file. For example, if you have the jraceman v1
source files in $JRACEMAN, you can load the USACK sports definition:

    ./jracemango -db sqlite3:_private/jrdb-test -import $JRACEMAN/data/usack-sports.txt

### Run the server

    ./jracemango -db sqlite3:_private/jrdb-test

## Documentation

To view the go documentation in your web browser:

    godoc -http=":6060"

Then open [localhost:6060](http://localhost:6060/) in your browser.

## Development

### Logging

jracemango uses [glog](https://github.com/golang/glog) for logging.
For details, see the [User Guide](https://github.com/google/glog#user-guide)
or the [glog source](https://github.com/golang/glog/blob/master/glog.go).

* Log messages are written to files in `/tmp`, with filenames starting with `jracemango`,
  divided into separate files per level, date-time, and pid
* For convenience, the symlinks `/tmp/jracemango.{INFO,WARNING,ERROR,FATAL}` point to the latest
  log files for each level
* Log messages at severity levels ERROR and FATAL are also sent to stderr

You can change the behavior of logging by specifying the appropriate command line option
when starting jracemango:

* To send all messages to stderr instead of the log files, use `--logtostderr`
* To send all messages to stderr in addition to the log files, use `--alsologtostderr`
* To enable verbose debugging, use `--v=N`, where `N` is a verbosity level such as 1 or 2
* To enable verbose debugging for some packages, use `--vmodule=pattern1=N1,pattern2=N2,pattern3=N3`,
  where `pattern1` and the others are source file names (without the directory or `.go` extension)
  or prefixes with an asterisk,
  and the `N` numbers are the verbosity levels for the matching files

### Sources

The sources for jracemango are available on github in multiple repositories
under [jimmc](http://github.com/jimmc):

* [jracemango](http://github.com/jimmc/jracemango) (this repo)
* [golden](http://github.com/jimmc/golden) - Support for unit tests using golden reference files
* [gtrepgen](http://github.com/jimmc/gtrepgen) - Go-Template REport GENerator

If you want to make changes to golden or gtrepgen while working on jracemango,
add these lines to the go.mod file here in jraceman.go:

```
replace github.com/jimmc/golden => ../golden
replace github.com/jimmc/gtrepgen => ../gtrepgen
```
