# Field Mapper
## Description
A utility API written in Go to match custom defined fields to a set of provided titles by comparing them to a field name and using some fuzzy matching. The same field is allowed to be present multiple times ONLY IF the refinement are different in the tags of a given field. If there is an refinement for a given field, then it is used to compare to the title. Otherwise the label is used. Each field/refinement combination should only map a single title and each single title should be mapped to a single field.

## Libraries Used
- [Gin-Gonic](https://github.com/gin-gonic/gin)
- [Excelize](https://github.com/qax-os/excelize)
- [fuzzysearch](https://github.com/lithammer/fuzzysearch)