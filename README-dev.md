# JRaceman Development

This file discusses development tips for JRaceman. The target audience
is software developers who intend to make code changes or additions
to JRaceman.

If you are only interested in getting JRaceman running on your system,
please see the instructions in [README.md](./README.md).

## Sources

The sources for JRaceman are available on github in multiple repositories
under [jimmc](http://github.com/jimmc):

* [jraceman](http://github.com/jimmc/jraceman) (this repo)
* [auth](http://github.com/jimmc/auth) - Support for authentication (login) and authorization (permissions)
* [golden](http://github.com/jimmc/golden) - Support for unit tests using golden reference files
* [gtrepgen](http://github.com/jimmc/gtrepgen) - Go-Template REport GENerator

If you want to make changes to auth, golden, or gtrepgen while working on JRaceman,
clone that repository from github.com and
add one or more of these lines to the go.mod file here:

```
replace github.com/jimmc/auth => ../auth
replace github.com/jimmc/golden => ../golden
replace github.com/jimmc/gtrepgen => ../gtrepgen
```

## Logging

JRaceman uses [glog](https://github.com/golang/glog) for logging.
For details, see the [glog User Guide](https://github.com/google/glog#user-guide)
or the [glog source](https://github.com/golang/glog/blob/master/glog.go).

* Log messages are written to files in `/tmp`, with filenames starting with `jraceman`,
  divided into separate files per level, date-time, and pid
* For convenience, the symlinks `/tmp/jraceman.{INFO,WARNING,ERROR,FATAL}` point to the latest
  log files for each level
* Log messages at severity levels ERROR and FATAL are also sent to stderr

You can change the behavior of logging by specifying the appropriate command line option
when starting JRaceman:

* To send all messages to stderr instead of the log files, use `--logtostderr`
* To send all messages to stderr in addition to the log files, use `--alsologtostderr`
* To enable verbose debugging, use `--v=N`, where `N` is a verbosity level such as 1 or 2
* To enable verbose debugging for some packages, use `--vmodule=pattern1=N1,pattern2=N2,pattern3=N3`,
  where `pattern1` and the others are source file names (without the directory or `.go` extension)
  or prefixes with an asterisk,
  and the `N` numbers are the verbosity levels for the matching files.

## Adding a database table

JRaceman provides extensive automatic support for database tables. With a minumum
set of changes, JRaceman will automatically support initial creation, import, export, 
marshaling and unmarshaling to the client, and client-side query and edit forms,
including foreign keys.

### Design your table

When selecting table and column names, follow these conventions:

* Use the singular name of the data items being stored (example: person, not people)
* The primary key column must be named "ID" and be a string
* Foreign key columns must have a name that is the foreign key table name followed by id
  (example: Area.SiteID is a foreign key from the Area table to the Site table; see domain/area.go)
  * If this is not the case, see the note below about using the `UpdateColumnInfos` method

### Add server code

When adding a new database table, select an existing table to use as a template, such
as the area table (if your table name is one word) or the laneorder table (if your
table name is two or more words, to get capitalization right). In each of the following
directories, copy that file to the corresponding file with your new table name, and edit
the new file appropriately.

* `api/crud`
* `api/query`
  * Look at the `SummaryQuery` method and create an appropriate summary string
    for your new type
* `dbrepo`
  * When editing this file, look at the `Save` method and select an appropriate
    ID prefix in the call to `structsql.UniqueID`. Look at some of the other
    tables for examples.
  * If you have any special case columns, such as a foreign key column that does
    not match the standard pattern, define the `UpdateColumnInfos` method in
    your Repo class, and call the `WithUpdate` versions of `CreateTable` and
    `UpgradeTable` from your functions of those names.
    See `dbrepo/simplanrule.go` for an example.
* `domain`
  * This is where you define the columns of your table
  * If you want a column value to be optional, make it a pointer in this struct

In addition, edit the following files, look for the table name you are using
as your source template, then copy that line for your new table:

* `api/crud/handler.go`
* `api/query/handler.go`
* `dbrepo/repos.go` (four places)
* `domain/repos.do`

### Add UI code

Many tables are appropriate to display directly in the UI and allow an admin
user to edit. To do this, look in the `_ui/src` directory. Decide which Setup
tab you want to add your new table to, and edit that file. For example, if you
want your new table to go in the same Setup tab as the area table, you would edit
`venue-setup.ts`. Each table has two lines, one with the tab label and one with
the tablename for the `table-queryedit` component. Copy a pair of those lines to
the position with the Setup tab where you want your new tab, and edit the tab
label ane tableName parameter to refer to your new table.

### Build and test

That's it! Recompile the go code and the typescript code, restart your jraceman
server, and reload your web page, and you should see the query and edit forms for
your new tab.
