package structsql

import (
  "sort"
  "strconv"
  "strings"
  "unicode"

  "github.com/jimmc/jraceman/dbrepo/compat"
  "github.com/jimmc/jraceman/dbrepo/strsql"
)

// UniqueID returns a unique id within the table, following the
// id pattern given but using a different trailing number if necessary.
func UniqueID(db compat.DBorTx, table string, id string) string {
  n, err := strsql.QueryInt(db, "select count(*) from " + table + " where id=?", id)
  if err != nil {
    return id   // We will probably get an error when we try to save this record.
  }
  if n == 0 {
    return id   // This ID is not yet taken, so we can use it.
  }

  // String off the trailing digits to get a prefix to which we can add other digits.
  prefix := strings.TrimRightFunc(id, unicode.IsDigit)
  pattern := prefix + "%"
  ss, err := strsql.QueryStrings(db, "select id from " + table + " where id LIKE ?", pattern)
  if err != nil {
    return id   // Punt
  }

  // Take the list of matching IDs, strip off the known prefix part, then convert
  // the suffix to a number. If the suffix doesn't parse as an int, then it won't
  // conflict with the ID we generate (with an int suffix), so we ignore it.
  nn := make([]int, 0)
  for _, s := range ss {
    nstr := strings.TrimPrefix(s, prefix)
    if strings.HasPrefix(nstr, "-") {
      continue          // Ignore negative numbers.
    }
    n, err := strconv.Atoi(nstr)
    if err == nil {
      nn = append(nn, n)
    }
  }
  sort.Ints(nn)

  n = NumberNotIn(nn)  // Get a suffix number that is not in the list.

  // There is a potential race condition if two clients are asking for a unique
  // ID in the same table at the same time. If that happens, a save of a new
  // record may fail. In practice this is rare, so we won't worry about it.
  return prefix + strconv.Itoa(n)
}

// NumberNotIn returns a number which is not in the given array.
// The array must be sorted.
func NumberNotIn(nn []int) int {
  cn := len(nn)

  // If there are no numbers in the array, start at 1.
  if cn == 0 {
    return 1
  }

  // If there are available small numbers, pick the one just before
  // the lowest that is currently in use.
  if nn[0] > 1 {
    return nn[0] - 1
  }

  // If we have a compact list of numbers (which we think is a common case),
  // pick the next number after the largest number in use.
  if nn[cn - 1] - nn[0] == cn - 1 {
    return nn[cn - 1] + 1
  }

  // We have an array with holes in it, scan for the first hole and use that value.
  last := nn[0] - 1
  for i := 0; i < cn; i++ {
    if nn[i] > 1 && nn[i] > last + 1 {
      return nn[i] - 1
    }
    last = nn[i]
  }
  // We shouldn't get here, since we previously checked for a compact array.
  return nn[cn - 1] + 1
}
