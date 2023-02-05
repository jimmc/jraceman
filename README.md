# jraceman

JRaceman is an open-source canoe and kayak regatta management system
(Jimmc's RACE MANager).

The original version of JRaceman was written in Java to run as a standalone
program or as a server with a thick client and first released in 2002.
This is version 2, rewritten using Go and Lit as a web server that
requires no local installation for client users of the system.

## Quick start

### Install Go compiler

See if the go compiler is already installed on your system:

    go version

If not, [install Go](https://go.dev/doc/install).

### Get the jraceman sources

Once the Go compiler is installed, it can do this for you automatically:

    go install github.com/jimmc/jraceman@latest

### Build

#### Compile the Go code

Change your working directory to the jraceman directory:

    cd ~/go/src/github.com/jimmc/jraceman

Compile jraceman:

    go build

This creates the executable `jraceman` in the jraceman directory.

#### Test the Go code

Run the unit tests:

    go test ./...

If you want to check the unit test coverage:

    go test ./... -coverprofile=cover.out
    go tool cover -html=cover.out -o cover.html
    # Open cover.html in your browser

### Build the UI pages

Follow the instructions in the [\_ui](./_ui) directory.

### Make your database

JRaceman uses the [glog](https://github.com/golang/glog)
logging package, which by default sends output
to files in /tmp. During setup, it is typically simpler to direct this
output to stderr so that it comes directly to the terminal. To do this,
add the `-logtostderr` command line option to all of the `jraceman`
commands in this section.

#### Choose a location for your database

Select a location for your database, such a `$HOME/jrdb`, and pass that
value to  the `-db` option to `jraceman` when you run it. The
commands in this section assume you have set the `JRDB` environment
variable to point to the location of your database. If your database
is located at `$HOME/jrdb`, you can use the following line to `sh`:

  export JRDB $HOME/jrdb

#### Create a new empty database

Select a location for your database, for example $HOME/jrdb, then run the jraceman binary
specifying that database:

    ./jraceman -db sqlite3:$JRDB -create

You can import a jraceman data file. For example, if you have the jraceman v1
source files in $JRACEMAN, you can load the USACK sports definition:

    ./jraceman -db sqlite3:$JRDB -import $JRACEMAN/data/usack-sports.txt

#### Upgrade the database

If you did not use the latest JRaceman v2 data file to create your database,
upgrade it so that it includes all the tables it needs. You can start with
a dry run:

    ./jraceman -db sqlite3:$JRDB -checkupgrade

Then do the upgrade:

    ./jraceman -db sqlite3:$JRDB -upgrade

#### Add a user

You need at least one user in order to log in. Add one:

    ./jraceman -db sqlite3:$JRDB -updatepassword user1

This will prompt you for a password and ask again for confirmation, then
create or update the named user with the new password. Once you have one
user, you can then log in and use the Auth Setup tabs to add more users.
You can also use the command line option `-password` to give the password
on the command line rather than typing it in twice.

#### Add standard permissions

API calls require permission, so you will need to add permissions to your
database if they are not included in your imported file. To add the standard
permissions and a "guru" role that has all permissions, import the standard
permissions file:

    ./jraceman -db sqlite3:$JRDB -logtostderr -import ./_data/stdperms.txt

#### Grant the "guru" role to your user

Assuming you loaded the standard permissions file, the role ID for the guru
role is R1. Assign that role to your user:

    ./jraceman -db sqlite3:$JRDB -logtostderr -sql \
        "INSERT INTO userrole(id,userid,roleid) VALUES('UR1','U1','R1')"

You can check the IDs of the user and role before running this command by
doing some SQL queries, for example:

    ./jraceman -db sqlite3:$JRDB -logtostderr -sql "SELECT * FROM user"

### Run the server

    ./jraceman -db sqlite3:$JRDB -logtostderr

## Documentation

To view the go documentation in your web browser:

    godoc -http=":6060"

Then open [localhost:6060](http://localhost:6060/) in your browser.

## Development

For information about developing additional code for JRaceman,
see [README-dev.md](./README-dev.md).
