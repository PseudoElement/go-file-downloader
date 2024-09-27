package services_models

type RandomValueCreatorParams struct {
	/* from sql_constants BOOL/STRING/NUMBER/AUTO_INCREMENT/FIRST_NAME/LAST_NAME/COUNTRY/CAR/WORK/DATE */
	ValueType   string
	Min         int64
	Max         int64
	IncrementFn func() int
}
