package dbrepo

import (
  "errors"
  "fmt"
  "strconv"
  "strings"
  "unicode"
)

const (
  TokenErr = iota
  TokenComma
  TokenNull
  TokenInt
  TokenBool
  TokenString
)

type QuotedScanner struct {
  runes []rune
  pos int
  nextToken *QuotedToken
  err *QuotedToken
}

type QuotedToken struct {
  Type int              // One of the Token* types
  Pos int               // The number of runes preceeding this token in the string
  Source string         // The original string for this field
  Value interface{}     // The converted value of this field
  Err error             // nil if no parsing error
}

func NewQuotedScanner(line string) *QuotedScanner {
  q := &QuotedScanner{}
  q.runes = []rune(strings.TrimSpace(line))
  q.pos = 0
  return q
}

// CommaSeparatedValues scans the entire line and returns the values
// between the commas. If there are two commas with no value between,
// a null value is assumed for that position. If there are two values
// without a comma between, an error is returned. When an error is
// returned, the values that have been scanned so far are also returned.
func (q *QuotedScanner) CommaSeparatedTokens() ([]*QuotedToken, error) {
  values := make([]*QuotedToken, 0)
  commaNext := false    // True when we expect to see a comma token
  for q.Next() {
    t := q.Token()
    if t.Err != nil {
      return values, t.Err
    }
    if t.Type == TokenComma {
      if !commaNext {
        // If we weren't expecting a comma, then we have a blank field
        values = append(values, &QuotedToken{
          Type: TokenNull,
          Pos: t.Pos,
        })
      }
      commaNext = false;
    } else {
      if commaNext {
        // We need a comma between every value
        values = append(values, t)      // Return that token so caller can see the problem.
        return values, fmt.Errorf("expected comma at pos=%d before token %v", t.Pos, t.Source)
      }
      values = append(values, t)
      commaNext = true
    }
  }
  return values, nil
}

// TokensToValues extracts the value field from each token and returns a slice
// of interface{} the same length as the given slice of tokens.
func (q *QuotedScanner) TokensToValues(tokens []*QuotedToken) []interface{} {
  values := make([]interface{}, len(tokens))
  for n, t := range tokens {
    values[n] = t.Value
  }
  return values
}

func (q *QuotedScanner) Next() bool {
  if q.nextToken != nil {
    // Caller must call Token() to read the previous token before
    // we scan for the next token.
    return true
  }
  if (q.err != nil) {
    return false        // We don't try to scan after getting an error.
  }

  if q.pos >= len(q.runes) {
    return false
  }
  start := q.pos
  r := q.runes[start]
  switch {
  case r == ',':
    q.nextToken = &QuotedToken{
      Type: TokenComma,
      Pos: start,
      Source: string(q.runes[start]),
    }
    q.pos++
  case unicode.IsDigit(r):
    end := start + 1
    for end < len(q.runes) && unicode.IsDigit(q.runes[end]) {
      end++
    }
    // TODO - look for a . for floating point; for now assume it's an int.
    source := string(q.runes[start:end])
    n, err := strconv.Atoi(source)
    if err != nil {
      q.nextToken = &QuotedToken{
        Type: TokenErr,
        Pos: start,
        Err: err,
      }
    } else {
      q.nextToken = &QuotedToken{
        Type: TokenInt,
        Pos: start,
        Source: source,
        Value: n,
      }
    }
    q.pos = end
  case r == '"':
    // Quoted string
    end := start + 1
    unescapedRunes := []rune{}
    for end < len(q.runes) && q.runes[end] != '"' && q.nextToken == nil {
      if q.runes[end] == '\\' {
        end++
        if end >= len(q.runes) {
          q.nextToken = &QuotedToken{
            Type: TokenErr,
            Pos: q.pos,
            Source: string(q.runes[start:]),
            Err: errors.New("backslash at end of quoted string"),
          }
          break
        }
        switch q.runes[end] {
        case '"', '\\':
          unescapedRunes = append(unescapedRunes, q.runes[end])
        case 'n':
          unescapedRunes = append(unescapedRunes, '\n')
        case 't':
          unescapedRunes = append(unescapedRunes, '\t')
        default:
          q.nextToken = &QuotedToken{
            Type: TokenErr,
            Pos: q.pos,
            Source: string(q.runes[start:end]),
            Err: fmt.Errorf("invalid quoted character '%v'", q.runes[end]),
          }
        }
      } else {
        unescapedRunes = append(unescapedRunes, q.runes[end])
      }
      end++
    }
    if q.nextToken == nil && end < len(q.runes) && q.runes[end] == '"' {
      end++     // Include the closing quote in the source string
    }
    q.pos = end
    if q.nextToken == nil {
      q.nextToken = &QuotedToken{
        Type: TokenString,
        Pos: start,
        Source: string(q.runes[start:end]),
        Value: string(unescapedRunes),
      }
    }
  default:
    // The only other thing we support is boolean values or null
    end := start + 1
    for end < len(q.runes) && q.runes[end] != ',' {
      end++
    }
    source := string(q.runes[start:end])
    if source == "null" {
      q.nextToken = &QuotedToken{
        Type: TokenNull,
        Pos: start,
        Source: source,
      }
    } else if source == "true" {
      q.nextToken = &QuotedToken{
        Type: TokenBool,
        Pos: start,
        Source: source,
        Value: true,
      }
    } else if source == "false" {
      q.nextToken = &QuotedToken{
        Type: TokenBool,
        Pos: start,
        Source: source,
        Value: false,
      }
    } else {
      q.nextToken = &QuotedToken{
        Type: TokenErr,
        Pos: start,
        Source: source,
        Err: fmt.Errorf("unrecognized unquoted token: %v", source),
      }
    }
    q.pos = end
  }
  if q.nextToken.Type == TokenErr {
    q.err = q.nextToken
  }
  return true
}

func (q *QuotedScanner) Token() *QuotedToken {
  if q.nextToken == nil {
    return &QuotedToken{
      Err: errors.New("no token available"),
    }
  }
  t := q.nextToken
  q.nextToken = nil
  return t
}
