package ixport_test

import (
  "reflect"
  "testing"

  "github.com/jimmc/jracemango/dbrepo/ixport"
)

func TestQuotedScannerEmpty(t *testing.T) {
  q := ixport.NewQuotedScanner("")
  if q.Next() {
    t.Errorf("Empty string should have no tokens")
  }
}

func TestQuotedScannerBasic(t *testing.T) {
  q := ixport.NewQuotedScanner("123,true,,\"abc\",null,456")

  expectedTokens := []*ixport.QuotedToken{
    &ixport.QuotedToken{
      Type: ixport.TokenInt,
      Pos: 0,
      Source: "123",
      Value: 123,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenComma,
      Pos: 3,
      Source: ",",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenBool,
      Pos: 4,
      Source: "true",
      Value: true,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenComma,
      Pos: 8,
      Source: ",",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenComma,
      Pos: 9,
      Source: ",",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenString,
      Pos: 10,
      Source: "\"abc\"",
      Value: "abc",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenComma,
      Pos: 15,
      Source: ",",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenNull,
      Pos: 16,
      Source: "null",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenComma,
      Pos: 20,
      Source: ",",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenInt,
      Pos: 21,
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
  q := ixport.NewQuotedScanner(`"a\"b","c\\d","e\n\t"`)

  expectedTokens := []*ixport.QuotedToken{
    &ixport.QuotedToken{
      Type: ixport.TokenString,
      Pos: 0,
      Source: `"a\"b"`,
      Value: `a"b`,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenComma,
      Pos: 6,
      Source: ",",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenString,
      Pos: 7,
      Source: `"c\\d"`,
      Value: `c\d`,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenComma,
      Pos: 13,
      Source: ",",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenString,
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

func TestQuotedScannerCommaSeparatedTokens(t *testing.T) {
  q := ixport.NewQuotedScanner("123,true,,\"abc\",456")

  expectedTokens := []*ixport.QuotedToken{
    &ixport.QuotedToken{
      Type: ixport.TokenInt,
      Pos: 0,
      Source: "123",
      Value: 123,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenBool,
      Pos: 4,
      Source: "true",
      Value: true,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenNull,
      Pos: 9,
      Source: "",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenString,
      Pos: 10,
      Source: "\"abc\"",
      Value: "abc",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenInt,
      Pos: 16,
      Source: "456",
      Value: 456,
    },
  }

  tokens, err := q.CommaSeparatedTokens()
  if err != nil {
    t.Fatalf("from CommaSeparatedTokens: %v", err)
  }
  if got, want := tokens, expectedTokens; !reflect.DeepEqual(got, want) {
    t.Fatalf("Tokens: got %v, want %v", got, want)
  }

  if q.Next() {
    xt := q.Token()
    t.Fatalf("Extra token: %v", xt)
  }
}

func TestQuotedScannerTokensToValues(t *testing.T) {
  tokens := []*ixport.QuotedToken{
    &ixport.QuotedToken{
      Type: ixport.TokenInt,
      Pos: 0,
      Source: "123",
      Value: 123,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenBool,
      Pos: 4,
      Source: "true",
      Value: true,
    },
    &ixport.QuotedToken{
      Type: ixport.TokenNull,
      Pos: 9,
      Source: "",
    },
    &ixport.QuotedToken{
      Type: ixport.TokenString,
      Pos: 10,
      Source: "\"abc\"",
      Value: "abc",
    },
  }
  expectedValues := []interface{}{
    123,
    true,
    nil,
    "abc",
  }

  q := ixport.NewQuotedScanner("")
  if got, want := q.TokensToValues(tokens), expectedValues; !reflect.DeepEqual(got, want) {
    t.Fatalf("Values: got %v, want %v", got, want)
  }
}
