/*
Package domain defines the domain entities used in jraceman.
These are the Go structs that are used for all of the app-specific
algorithms, and define the canonical forms of the data.

Field names in these structs follow the Go naming conventions.

Fields which are defined as pointer fields are optional and may
contain nil pointers.

Note that these structs intentionally do not include tags for such
things as json field names or database field names. In the architecture
used here, that information should be included in the packages that
deal with that aspect of the design.

The Repo interfaces define the functions necessary to get and save
data from and to the repository (database).
*/
package domain
