package report

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestWhere(t *testing.T) {
  tests := []struct{
    name string
    attrsMap map[string]interface{}
    options *ReportOptions
    expect *whereData
    expectError bool
  } {
    {
      name: "empty_options",
      attrsMap: map[string]interface{}{"where": []interface{}{"event", "person_id"}},
      options: &ReportOptions{WhereValues: map[string]WhereValue{
      }},
      expect: &whereData{},
      expectError: false,
    },
    {
      name: "one_option",
      attrsMap: map[string]interface{}{"where": []interface{}{"event", "person_id"}},
      options: &ReportOptions{WhereValues: map[string]WhereValue{
        "event_id": {Op: "eq", Value: "ID"},
      }},
      expect: &whereData{
        expr: "event.id eq 'ID'",
        whereclause: " where event.id eq 'ID'",
        andclause: " && event.id eq 'ID'",
      },
      expectError: false,
    },
    {
      name: "invalid_option",
      attrsMap: map[string]interface{}{"where": []interface{}{"event", "person_id"}},
      options: &ReportOptions{WhereValues: map[string]WhereValue{
        "no_such_field": {Op: "eq", Value: "anything"},
      }},
      expectError: true,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := where(tc.attrsMap, tc.options)
      if tc.expectError {
        if err == nil {
          t.Fatalf("where: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("where: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got, cmp.AllowUnexported(whereData{})); diff != "" {
          t.Errorf("where mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}

func TestAttrsToWhereList(t *testing.T) {
  tests := []struct{
    name string
    input map[string]interface{}
    expect []string
    expectError bool
  } {
    {
      name: "no_where",
      input: map[string]interface{}{},
      expect: []string{},
      expectError: false,
    },
    {
      name: "simple",
      input: map[string]interface{}{"where": []interface{}{"a", "b", "c"}},
      expect: []string{"a", "b", "c"},
      expectError: false,
    },
    {
      name: "expanded",
      input: map[string]interface{}{"where": []interface{}{"a", "event", "c"}},
      expect: []string{"a", "event_id", "event_name", "event_number", "c"},
      expectError: false,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := attrsToWhereList(tc.input)
      if tc.expectError {
        if err == nil {
          t.Fatalf("attrsToWhereList: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("attrsToWhereList: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got); diff != "" {
          t.Errorf("attrsToWhereList mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}

func TestExtractWhereList(t *testing.T) {
  tests := []struct{
    name string
    input map[string]interface{}
    expect []string
    expectError bool
  } {
    {
      name: "no_where",
      input: map[string]interface{}{},
      expect: nil,
      expectError: false,
    },
    {
      name: "simple",
      input: map[string]interface{}{"where": []interface{}{"a", "b", "c"}},
      expect: []string{"a", "b", "c"},
      expectError: false,
    },
    {
      name: "where_not_array",
      input: map[string]interface{}{"where": "x"},
      expectError: true,
    },
    {
      name: "where_item_not_string",
      input: map[string]interface{}{"where": []interface{}{"a", 2, "c"}},
      expectError: true,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T){
      got, err := extractWhereList(tc.input)
      if tc.expectError {
        if err == nil {
          t.Fatalf("extractWhereList: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("extractWhereList: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got); diff != "" {
          t.Errorf("extractWhereList mismatch (-want +got):\n%s", diff)
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
    whereValues map[string]WhereValue
    expect []string
  } {
    {
      name: "empty",
      whereList: []string{},
      whereValues: map[string]WhereValue{},
      expect: []string{},
    },
    {
      name: "simple",
      whereList: []string{"a", "b", "c"},
      whereValues: map[string]WhereValue{
        "a": WhereValue{},
        "c": WhereValue{},
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
    whereValues map[string]WhereValue
    expect *whereData
    expectError bool
  } {
    {
      name: "noValues",
      whereMap: map[string]whereDetails{
        "a": whereDetails{table:"A", column:"c"},
      },
      whereListInUse: []string{},
      whereValues: map[string]WhereValue{
      },
      expect: &whereData{},
      expectError: false,
    },
    {
      name: "oneValue",
      whereMap: map[string]whereDetails{
        "a": whereDetails{table:"A", column:"c"},
      },
      whereListInUse: []string{"a"},
      whereValues: map[string]WhereValue{
        "a": WhereValue{Op: "eq", Value: 123},
      },
      expect: &whereData{
        expr: "A.c eq 123",
        whereclause: " where A.c eq 123",
        andclause: " && A.c eq 123",
      },
      expectError: false,
    },
    {
      name: "twoValues",
      whereMap: map[string]whereDetails{
        "a": whereDetails{table:"A", column:"c"},
        "b": whereDetails{table:"B", column:"d"},
      },
      whereListInUse: []string{"a", "b"},
      whereValues: map[string]WhereValue{
        "a": WhereValue{Op: "eq", Value: 123},
        "b": WhereValue{Op: "ne", Value: "xyz"},
      },
      expect: &whereData{
        expr: "A.c eq 123 && B.d ne 'xyz'",
        whereclause: " where A.c eq 123 && B.d ne 'xyz'",
        andclause: " && A.c eq 123 && B.d ne 'xyz'",
      },
      expectError: false,
    },
    {
      name: "invalidValue",
      whereMap: map[string]whereDetails{
        "a": whereDetails{table:"A", column:"c"},
      },
      whereListInUse: []string{"a","b"},
      whereValues: map[string]WhereValue{
        "a": WhereValue{Op: "eq", Value: 123},
        "b": WhereValue{Op: "ne", Value: "xyz"},
      },
      expect: &whereData{},
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
        if diff := cmp.Diff(want, got, cmp.AllowUnexported(whereData{})); diff != "" {
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
