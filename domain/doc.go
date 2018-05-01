/*
Package domain defines the domain entities used in jraceman.
These are the Go structs that are used for all of the app-specific
algorithms, and define the canonical forms of the data.

Field names in these structs follow the Go naming conventions
for letter case.

Fields which are defined as pointer fields are optional and may
contain nil pointers.

A field named "ID" is assumed to be the key field for that type
of entity, and would, for example and by default, become the
primary key in a database table.

A field named FooID (i.e., ends with "ID") is assumed to refer to
the ID of an entity of type Foo, and would, for example, become
a foreign key referencing Foo(ID) in a database table.

Note that these structs intentionally do not include tags for such
things as json field names or database field names. In the architecture
used here, that information should be included in the packages that
deal with that aspect of the design. If, for uses such as json or
database, there are any fields that require translation of names or
other details beyond the assumptions noted in the preceeding paragraphs,
that information should be provided by a mechanism specific to the area
doing that transformation.

The Repo interfaces define the functions necessary to get and save
data from and to the repository (database).
*/
package domain
