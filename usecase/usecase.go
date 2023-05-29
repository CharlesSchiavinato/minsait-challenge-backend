package usecase

// ErrModelValidate denotes failing validate model.
type ErrModelValidate struct {
	Message string
}

// ErrModelValidate returns the model validation error.
func (emv ErrModelValidate) Error() string {
	return emv.Message
}

// ErrParamValidate denotes failing validate param.
type ErrParamValidate struct {
	Message string
}

// ErrParamValidate returns the param validation error.
func (epv ErrParamValidate) Error() string {
	return epv.Message
}
