/*
The report package implements report generation for jraceman.

Families of structs represent different phases of processing of reports.
Report templates can include attributes at the start, which get loaded
into struct ReportAttributes.

Some fields are extracted from ReportAttributes and are passed to the
user in struct ReportFields.

When the user makes a request to generate a report, that request can include
ReportOptions as specified by the user.

Before invoking the templating engine, the code here precomputes some
variables and passed them into the template in struct ReportComputed.

Each of these structs can have contained structs with portions of the data.
These structs drop the "Report" part of the containing struct and add
another word for the contained data. Examples are AttributesWhere and
FieldsOrderBy.
*/
package report
