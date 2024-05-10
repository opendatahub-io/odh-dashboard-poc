package data

// Models struct is a single convenient container to hold and represent all our models.
type Models struct {
	ModelRegistry ModelRegistryModel
}

func NewModels() Models {
	return Models{
		ModelRegistry: ModelRegistryModel{},
	}
}
