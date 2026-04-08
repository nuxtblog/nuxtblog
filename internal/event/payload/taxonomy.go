package payload

// TaxonomyCreated is delivered when a taxonomy record (category/tag assignment) is created.
type TaxonomyCreated struct {
	TaxID    int64
	TermID   int64
	TermName string
	TermSlug string
	// Taxonomy type, e.g. "category" or "tag"
	Taxonomy string
}

// TaxonomyDeleted is delivered when a taxonomy record is removed.
type TaxonomyDeleted struct {
	TaxID    int64
	TermName string
	TermSlug string
	Taxonomy string
}

// TermCreated is delivered when a new term (category/tag definition) is created.
type TermCreated struct {
	TermID int64
	Name   string
	Slug   string
}

// TermDeleted is delivered when a term is deleted.
type TermDeleted struct {
	TermID int64
	Name   string
	Slug   string
}
