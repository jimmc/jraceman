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

### Build the UI pages

Follow the instructions in the [\_ui](./_ui) directory.

### Create a new database

Select a location for your database, for example $HOME/jrdb, then run the jraceman binary
specifying that database:

    ./jracemango -db sqlite3:$HOME/jrdb -create

### Run the server

    ./jracemango -db sqlite3:$HOME/jrdb

## Documentation

To view the go documentation in your web browser:

    godoc -http=":6060"

Then open [localhost:6060](http://localhost:6060/) in your browser.

## Sources

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
