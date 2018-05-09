package structsql_test

import (
  "errors"
  "testing"

  "github.com/jimmc/jracemango/dbrepo/structsql"
)

type oneResultTester struct {
  num int64
  err error
}
func (o *oneResultTester) RowsAffected() (int64, error) {
  return o.num, o.err
}

func TestRequireOneResult(t *testing.T) {
  oZero := &oneResultTester{0, nil}
  oOne := &oneResultTester{1, nil}
  oTwo := &oneResultTester{2, nil}
  oErr := &oneResultTester{0, errors.New("Test error")}

  if got, want := structsql.RequireOneResult(oOne, nil, "Tested", "foo", "123"), error(nil); got != want {
    t.Errorf("Happy path: got %v, want %v", got, want)
  }
  if got, want := structsql.RequireOneResult(oOne, errors.New("Test error"), "Tested", "foo", "123"), errors.New("Test error"); got.Error() != want.Error() {
    t.Errorf("With passed-in error: got %v, want %v", got, want)
  }
  if got, want := structsql.RequireOneResult(oErr, nil, "Tested", "foo", "123"), errors.New("Test error"); got.Error() != want.Error() {
    t.Errorf("With sql error: got %v, want %v", got, want)
  }
  if got, want := structsql.RequireOneResult(oZero, nil, "Tested", "foo", "123"), errors.New("Wrong-count error"); got == nil {
    t.Errorf("With count==0: got %v, want %v", got, want)
  }
  if got, want := structsql.RequireOneResult(oTwo, nil, "Tested", "foo", "123"), errors.New("Wrong-count error"); got == nil {
    t.Errorf("With count==0: got %v, want %v", got, want)
  }
}
