# Jraceman Lit UI

This directory contains the UI for the web-based version of jraceman.
The UI is implemented in Lit + Typescript.

## Building the UI

### Install tools

If you don't already have npm and tsc, you will have to install them.
This may differ depending on what version of Linux you are running,
but may require commands such as the following (for Ubuntu or Debian):

    curl -fsSL https://deb.nodesource.com/setup_19.x | sudo -E bash -
    sudo apt-get install -y nodejs
    curl -L https://npmjs.org/install.sh | sudo sh
    sudo npm install -g typescript

### Install Lit

The first time only, install a local copy of the Lit package:

    npm install lit

### Build

Compile the typescript code into the build directory:

    tsc

### Run

Run the jracemango program with appropriate command line arguments,
such as this command, which specifies a database, this dir as the UI root dir,
and debugging for importresolver.go:

    ./jracemango --db sqlite3:_private/jrdb-test --uiroot _ui2 -logtostderr -vmodule=importresolver=2

Hard-reload your browser page to load the newly compiled code.
