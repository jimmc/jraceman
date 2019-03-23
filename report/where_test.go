package report

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestExpandStringList(t *testing.T) {
  emptyExpansions := map[string][]string{}
  tests := []struct{
    name string
    in []string
    expansions map[string][]string
    expect []string
    expectError bool
  } {
    {
      name: "empty",
      in: []string{},
      expansions: emptyExpansions,
      expect: []string{},
      expectError: false,
    },
    {
      name: "none",
      in: []string{"a", "b", "c"},
      expansions: emptyExpansions,
      expect: []string{"a", "b", "c"},
      expectError: false,
    },
    {
      name: "simple",
      in: []string{"a", "b", "c"},
      expansions: map[string][]string{"b": {"x", "y"}},
      expect: []string{"a", "x", "y", "c"},
      expectError: false,
    },
    {
      name: "nested",
      in: []string{"a", "b", "c"},
      expansions: map[string][]string{"a": {"x", "y"}, "y": {"j", "k"}, "j": {"p", "q"}},
      expect: []string{"x", "p", "q", "k", "b", "c"},
      expectError: false,
    },
    {
      name: "duplicates",
      in: []string{"a", "b", "b", "c"},
      expansions: map[string][]string{"b": {"c", "d"}},
      expect: []string{"a", "c", "d", "c", "d", "c"},
      expectError: false,
    },
    {
      name: "cycle",
      in: []string{"a", "b", "c"},
      expansions: map[string][]string{"a": {"x", "y"}, "y": {"z"}, "z": {"b", "a"}},
      expect: []string{},
      expectError: true,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := expandStringList([]string{}, []string{}, tc.in, tc.expansions)
      if err != nil {
        if !tc.expectError {
          t.Fatalf("expandStringList: unexpected error: %v", err)
        }
      } else {
       want := tc.expect
       if diff := cmp.Diff(want, got); diff != "" {
         t.Errorf("expandStringList mismatch (-want +got):\n%s", diff)
       }
      }
    })
  }
}

func TestWhereString(t *testing.T) {
  tests := []struct{
    name string
    field whereDetails
    value WhereValue
    expect string
  } {
    {
      name: "int",
      field: whereDetails{table: "tbl", column: "col"},
      value: WhereValue{Op: "=", Value: 123},
      expect: "tbl.col = 123",
    },
    {
      name: "string",
      field: whereDetails{table: "tbl", column: "col"},
      value: WhereValue{Op: "=", Value: "abc"},
      expect: "tbl.col = 'abc'",
    },
    {
      name: "field override",
      field: whereDetails{table: "tbl", column: "col", field: "fff"},
      value: WhereValue{Op: "=", Value: 123},
      expect: "fff = 123",
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      expr, err := whereString(tc.field, tc.value)
      if err != nil {
        t.Fatalf("whereString: unexpected error: %v", err)
      }
      if got, want := expr, tc.expect; got != want {
        t.Errorf("whereString: got %q, want %q", got, want)
      }
    })
  }
}

func TestSqlQuotedString(t *testing.T) {
  tests := []struct{
    name string
    input string
    expected string
  } {
    { "empty", "", "''" },
    { "simple", "abc", "'abc'" },
    { "double quotes", `a"b"c`, `'a"b"c'` },
    { "single quotes", `a'b'c`, `'a''b''c'` },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      if got, want := sqlQuotedString(tc.input), tc.expected; got != want {
        t.Errorf("sqlQuotedString(%q) = %q, wanted %q", tc.input, got, want)
      }
    })
  }
}
