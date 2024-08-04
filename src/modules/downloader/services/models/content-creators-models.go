package services_models

type RandomValueCreatorParams struct {
	/* from sql_constants BOOL/STRING/NUMBER/AUTO_INCREMENT/FIRST_NAME/LAST_NAME/COUNTRY/CAR */
	ValueType   string
	Min         int
	Max         int
	IncrementFn func() int
}