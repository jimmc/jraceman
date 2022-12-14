package main

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestImportRegexp(t *testing.T) {
  tests := []struct {
    name string
    str string
    expectResult [][]int
  } {
    {
      name: "no-match",
      str: "random string",
      expectResult: nil,        // nil when no match
    },
    {
      name: "import without filename",
      str: "import comment",
      expectResult: nil,
    },
    {
      name: "import file single quote",
      str: "import 'foo';",
      expectResult: [][]int{{0,13,-1,-1,7,12}},
    },
    {
      name: "import file relative path",
      str: "import './foo';",
      expectResult: nil,
    },
    {
      name: "import file double quote",
      str: `import "foo";`,
      expectResult: [][]int{{0,13,-1,-1,7,12}},
    },
    {
      name: "import file absolute path",
      str: `import "/foo";`,
      expectResult: nil,
    },
    {
      name: "two import files",
      str: `import "foo"; import "bar" ;`,
      expectResult: [][]int{{0,13,-1,-1,7,12},{14,28,-1,-1,21,26}},
    },
    {
      name: "import names",
      str: `import {a,b,c} from "foo";`,
      expectResult: [][]int{{0,26,7,14,20,25}},
    },
    {
      name: "import names no spaces",
      str: `import{a,b,c}from"foo";`,
      expectResult: [][]int{{0,23,6,13,17,22}},
    },
    {
      name: "default import",
      str: `import dflt from"foo";`,
      expectResult: [][]int{{0, 22, 7, 12, 16, 21}},
    },
    {
      name: "wildcard import",
      str: `import * as dflt from "foo";`,
      expectResult: [][]int{{0, 28, 7, 17, 22, 27}},
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got := importPattern.FindAllStringSubmatchIndex(tc.str, -1)
      want := tc.expectResult
      if diff := cmp.Diff(want, got); diff != "" {
          t.Errorf("regexp match (-want +got):\n%s", diff)
      }
    })
  }
}
