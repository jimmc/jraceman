package dbrepo

import (
  "reflect"
  "testing"
)

func TestQuotedScannerEmpty(t *testing.T) {
  q := NewQuotedScanner("")
  if q.Next() {
    t.Errorf("Empty string should have no tokens")
  }
}

func TestQuotedScannerBasic(t *testing.T) {
  q := NewQuotedScanner("123,true,,\"abc\",456")

  expectedTokens := []*QuotedToken{
    &QuotedToken{
      Type: TokenInt,
      Pos: 0,
      Source: "123",
      Value: 123,
    },
    &QuotedToken{
      Type: TokenComma,
      Pos: 3,
      Source: ",",
    },
    &QuotedToken{
      Type: TokenBool,
      Pos: 4,
      Source: "true",
      Value: true,
    },
    &QuotedToken{
      Type: TokenComma,
      Pos: 8,
      Source: ",",
    },
    &QuotedToken{
      Type: TokenComma,
      Pos: 9,
      Source: ",",
    },
    &QuotedToken{
      Type: TokenString,
      Pos: 10,
      Source: "\"abc\"",
      Value: "abc",
    },
    &QuotedToken{
      Type: TokenComma,
      Pos: 15,
      Source: ",",
    },
    &QuotedToken{
      Type: TokenInt,
      Pos: 16,
      Source: "456",
      Value: 456,
    },
  }

  for n, xt := range expectedTokens {
    if !q.Next() {
      t.Fatalf("Should have token #%d", n)
    }
    if got, want := q.Token(), xt; !reflect.DeepEqual(got, want) {
      t.Fatalf("Token #%d: got %v, want %v", n, got, want)
    }
  }

  if q.Next() {
    xt := q.Token()
    t.Fatalf("Extra token: %v", xt)
  }
}

func TestQuotedScannerQuotes(t *testing.T) {
  q := NewQuotedScanner(`"a\"b","c\\d","e\n\t"`)

  expectedTokens := []*QuotedToken{
    &QuotedToken{
      Type: TokenString,
      Pos: 0,
      Source: `"a\"b"`,
      Value: `a"b`,
    },
    &QuotedToken{
      Type: TokenComma,
      Pos: 6,
      Source: ",",
    },
    &QuotedToken{
      Type: TokenString,
      Pos: 7,
      Source: `"c\\d"`,
      Value: `c\d`,
    },
    &QuotedToken{
      Type: TokenComma,
      Pos: 13,
      Source: ",",
    },
    &QuotedToken{
      Type: TokenString,
      Pos: 14,
      Source: `"e\n\t"`,
      Value: "e\n\t",
    },
  }

  for n, xt := range expectedTokens {
    if !q.Next() {
      t.Fatalf("Should have token #%d", n)
    }
    if got, want := q.Token(), xt; !reflect.DeepEqual(got, want) {
      t.Fatalf("Token #%d: got %v, want %v", n, got, want)
    }
  }

  if q.Next() {
    xt := q.Token()
    t.Fatalf("Extra token: %v", xt)
  }
}

// TODO - add test for bad strings, test for bad calling sequences
