package tmpl

// TemplateBundleList is a list of template bundles
type TemplateBundleList []*TemplateBundle

// Generate generates the template bundle list
func (t *TemplateBundleList) Generate(outputDir string, data any) error {
	for _, tb := range *t {
		if err := tb.Generate(outputDir, data); err != nil {
			return err
		}
	}
	return nil
}

// Add adds or replaces template bundle by output file
func (t *TemplateBundleList) _add(tb *TemplateBundle) *TemplateBundleList {
	for i, v := range *t {
		if v.OutputFile == tb.OutputFile {
			(*t)[i] = tb
			return t
		}
	}
	*t = append(*t, tb)
	return t
}

// Add adds or replaces template bundle by output file
func (t *TemplateBundleList) Add(tb ...*TemplateBundle) *TemplateBundleList {
	for _, v := range tb {
		t._add(v)
	}
	return t
}

// Remove removes template bundle by output file
func (t *TemplateBundleList) Remove(outputFile string) *TemplateBundleList {
	for i, v := range *t {
		if v.OutputFile == outputFile {
			*t = append((*t)[:i], (*t)[i+1:]...)
			return t
		}
	}
	return t
}
