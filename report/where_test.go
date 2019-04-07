package report

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestWhere(t *testing.T) {
  tests := []struct{
    name string
    attrs *ReportAttributes
    options *ReportOptions
    expect *ComputedWhere
    expectError bool
  } {
    {
      name: "empty_options",
      attrs: &ReportAttributes{
        Where: []string{"event", "person_id"},
      },
      options: &ReportOptions{Where: []OptionsWhereItem{}},
      expect: &ComputedWhere{},
      expectError: false,
    },
    {
      name: "one_option",
      attrs: &ReportAttributes{
        Where: []string{"event", "person_id"},
      },
      options: &ReportOptions{Where: []OptionsWhereItem{
        {Name: "event_id", Op: "eq", Value: "ID"},
      }},
      expect: &ComputedWhere{
        Expr: "event.id = 'ID'",
        WhereClause: " where event.id = 'ID'",
        AndClause: " AND event.id = 'ID'",
      },
      expectError: false,
    },
    {
      name: "invalid_option",
      attrs: &ReportAttributes{
        Where: []string{"event", "person_id"},
      },
      options: &ReportOptions{Where: []OptionsWhereItem{
        {Name: "no_such_field", Op: "eq", Value: "anything"},
      }},
      expectError: true,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := computeWhere(tc.attrs, tc.options)
      if tc.expectError {
        if err == nil {
          t.Fatalf("computeWhere: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("where: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got, cmp.AllowUnexported(ComputedWhere{})); diff != "" {
          t.Errorf("computeWhere mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}

func TestWhereListToMap(t *testing.T) {
  tests := []struct{
    name string
    input []string
    expect map[string]whereDetails
    expectError bool
  } {
    {
      name: "valid",
      input: []string{"event_id"},
      expect: map[string]whereDetails{
        "event_id": {display: "Event", table: "event", column: "id"},
      },
      expectError: false,
    },
    {
      name: "invalid_field",
      input: []string{"invalid"},
      expectError: true,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := whereListToMap(tc.input)
      if tc.expectError {
        if err == nil {
          t.Fatalf("whereListToMap: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("whereListToMap: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got, cmp.AllowUnexported(whereDetails{})); diff != "" {
          t.Errorf("whereListToMap mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}

func TestExtractWhereListInUse(t *testing.T) {
  tests := []struct{
    name string
    whereList []string
    whereValues []OptionsWhereItem
    expect []string
  } {
    {
      name: "empty",
      whereList: []string{},
      whereValues: []OptionsWhereItem{},
      expect: []string{},
    },
    {
      name: "simple",
      whereList: []string{"a", "b", "c"},
      whereValues: []OptionsWhereItem{
        {Name: "a"},
        {Name: "c"},
      },
      expect: []string{"a", "c"},
    },
  }
  for _, tc:= range tests {
    got, want := extractWhereListInUse(tc.whereList, tc.whereValues), tc.expect
    if diff := cmp.Diff(want, got); diff != "" {
      t.Errorf("extractWhereListInUse mismatch (-want +got):\n%s", diff)
    }
  }
}

func TestWhereMapToData(t *testing.T) {
  tests := []struct{
    name string
    whereMap map[string]whereDetails
    whereListInUse []string
    whereValues []OptionsWhereItem
    expect *ComputedWhere
    expectError bool
  } {
    {
      name: "noValues",
      whereMap: map[string]whereDetails{
        "a": {table:"A", column:"c"},
      },
      whereListInUse: []string{},
      whereValues: []OptionsWhereItem{},
      expect: &ComputedWhere{},
      expectError: false,
    },
    {
      name: "oneValue",
      whereMap: map[string]whereDetails{
        "a": {table:"A", column:"c"},
      },
      whereListInUse: []string{"a"},
      whereValues: []OptionsWhereItem{
        {Name: "a", Op: "eq", Value: 123},
      },
      expect: &ComputedWhere{
        Expr: "A.c = 123",
        WhereClause: " where A.c = 123",
        AndClause: " AND A.c = 123",
      },
      expectError: false,
    },
    {
      name: "twoValues",
      whereMap: map[string]whereDetails{
        "a": {table:"A", column:"c"},
        "b": {table:"B", column:"d"},
      },
      whereListInUse: []string{"a", "b"},
      whereValues: []OptionsWhereItem{
        {Name: "a", Op: "eq", Value: 123},
        {Name: "b", Op: "ne", Value: "xyz"},
      },
      expect: &ComputedWhere{
        Expr: "A.c = 123 AND B.d != 'xyz'",
        WhereClause: " where A.c = 123 AND B.d != 'xyz'",
        AndClause: " AND A.c = 123 AND B.d != 'xyz'",
      },
      expectError: false,
    },
    {
      name: "invalidValue",
      whereMap: map[string]whereDetails{
        "a": {table:"A", column:"c"},
      },
      whereListInUse: []string{"a","b"},
      whereValues: []OptionsWhereItem{
        {Name: "a", Op: "eq", Value: 123},
        {Name: "b", Op: "ne", Value: "xyz"},
      },
      expect: &ComputedWhere{},
      expectError: true,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := whereMapToData(tc.whereMap, tc.whereListInUse, tc.whereValues)
      if tc.expectError {
        if err == nil {
          t.Fatalf("whereMapToData: expected error but didn't get one")
        }
      } else if err != nil {
        t.Fatalf("whereMapToData: unexpected error: %v:", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got); diff != "" {
          t.Errorf("whereMapToData mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}

func TestExpandWhereList(t *testing.T) {
  input := []string{"event", "team_id"}
  want := []string{"event_id", "event_name", "event_number", "team_id"}
  got, err := expandWhereList(input)
  if err != nil {
    t.Fatalf("expandWhereList: unexpected error: %v", err)
  }
  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("expandWhereList mismatch (-want +got):\n%s", diff)
  }
}

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
      if tc.expectError {
        if err == nil {
          t.Fatalf("expandStringList: expected error but didn't get one")
        }
      } else if err != nil {
        t.Fatalf("expandStringList: unexpected error: %v", err)
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
    value OptionsWhereItem
    expect string
  } {
    {
      name: "int",
      field: whereDetails{table: "tbl", column: "col"},
      value: OptionsWhereItem{Op: "eq", Value: 123},
      expect: "tbl.col = 123",
    },
    {
      name: "string",
      field: whereDetails{table: "tbl", column: "col"},
      value: OptionsWhereItem{Op: "eq", Value: "abc"},
      expect: "tbl.col = 'abc'",
    },
    {
      name: "field override",
      field: whereDetails{table: "tbl", column: "col", field: "fff"},
      value: OptionsWhereItem{Op: "eq", Value: 123},
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
