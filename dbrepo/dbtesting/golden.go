package dbtesting

import (
  "bytes"
  "database/sql"
  "fmt"
  "io"
  "io/ioutil"
  "os"
)

func CompareOutToGolden(outfilename, goldenfilename string) error {
  outcontent, err := ioutil.ReadFile(outfilename)
  if err != nil {
    return fmt.Errorf("error reading back output file %s: %v", outfilename, err)
  }
  goldencontent, err := ioutil.ReadFile(goldenfilename)
  if err != nil {
    return fmt.Errorf("error reading golden file %s: %v", goldenfilename, err)
  }
  if !bytes.Equal(outcontent, goldencontent) {
    return fmt.Errorf("outfile %s does not match golden file %s", outfilename, goldenfilename)
  }
  return nil
}

// FromSetupToGolden loads a setup file into a fresh test database, runs the specified
// callback to produce a test output file, and compares it to the golden file.
func FromSetupToGolden(base string, callback func(*sql.DB, io.Writer) error) error {
  setupfilename := "testdata/" + base + ".setup"
  outfilename := "testdata/" + base + ".out"
  goldenfilename := "testdata/" + base + ".golden"

  db, err := DbWithSetupFile(string(setupfilename))
  if err != nil {
    return fmt.Errorf("error setting up for %s: %v", base, err)
  }
  defer db.Close()

  os.Remove(outfilename)
  outfile, err := os.Create(outfilename)
  if err != nil {
    return fmt.Errorf("error creating output file: %v", outfile)
  }

  // Run the specific test step.
  err = callback(db, outfile)

  if err != nil {
    return err
  }
  outfile.Close()
  return CompareOutToGolden(outfilename, goldenfilename)
}
