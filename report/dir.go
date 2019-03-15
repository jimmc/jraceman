package report

import (
  "fmt"

  "github.com/jimmc/gtrepgen/gen"
)

/* ReadTemplateAttrs loads the attributes from all the template files in
 * the given directory.
 */
func ReadTemplateAttrs(templateDir string) ([]map[string]interface{}, error) {
  attrs, err := gen.ReadDirFilesAttributes(templateDir)
  if attrs == nil {
    return nil, err
  }
  attrMaps := []map[string]interface{}{}
  for _, fattrs := range attrs {
    if fattrs.Err != nil {
      return nil, fmt.Errorf("for template %q received error %v", fattrs.Name, fattrs.Err)
    }
    fmap, ok := fattrs.Attributes.(map[string]interface{})
    if !ok {
      return nil, fmt.Errorf("invalid data type for template %q", fattrs.Name)
    }
    attrMaps = append(attrMaps, fmap)
  }
  return attrMaps, err
}
