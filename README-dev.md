# JRaceman Development

This file discusses development tips for JRaceman. The target audience
is software developers who intend to make code changes or additions
to JRaceman.

If you are only interested in getting JRaceman running on your system,
please see the instructions in [README.md](./README.md).

## Sources

The sources for jraceman are available on github in multiple repositories
under [jimmc](http://github.com/jimmc):

* [jraceman](http://github.com/jimmc/jraceman) (this repo)
* [auth](http://github.com/jimmc/auth) - Support for authentication (login) and authorization (permissions)
* [golden](http://github.com/jimmc/golden) - Support for unit tests using golden reference files
* [gtrepgen](http://github.com/jimmc/gtrepgen) - Go-Template REport GENerator

If you want to make changes to auth, golden, or gtrepgen while working on jraceman,
clone that repository from github.com
add one or more of these lines to the go.mod file here:

```
replace github.com/jimmc/auth => ../auth
replace github.com/jimmc/golden => ../golden
replace github.com/jimmc/gtrepgen => ../gtrepgen
```

## Logging

jraceman uses [glog](https://github.com/golang/glog) for logging.
For details, see the [glog User Guide](https://github.com/google/glog#user-guide)
or the [glog source](https://github.com/golang/glog/blob/master/glog.go).

* Log messages are written to files in `/tmp`, with filenames starting with `jraceman`,
  divided into separate files per level, date-time, and pid
* For convenience, the symlinks `/tmp/jraceman.{INFO,WARNING,ERROR,FATAL}` point to the latest
  log files for each level
* Log messages at severity levels ERROR and FATAL are also sent to stderr

You can change the behavior of logging by specifying the appropriate command line option
when starting jraceman:

* To send all messages to stderr instead of the log files, use `--logtostderr`
* To send all messages to stderr in addition to the log files, use `--alsologtostderr`
* To enable verbose debugging, use `--v=N`, where `N` is a verbosity level such as 1 or 2
* To enable verbose debugging for some packages, use `--vmodule=pattern1=N1,pattern2=N2,pattern3=N3`,
  where `pattern1` and the others are source file names (without the directory or `.go` extension)
  or prefixes with an asterisk,
  and the `N` numbers are the verbosity levels for the matching files.
