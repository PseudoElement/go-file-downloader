package services_models

type RandomValueCreatorParams struct {
	/* from sql_constants BOOL/STRING/NUMBER/AUTO_INCREMENT */
	ValueType   string
	Min         int
	Max         int
	IncrementFn func() int
}
