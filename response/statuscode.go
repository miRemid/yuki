package response

type StatusCode uint32

const (
	StatusOK StatusCode = iota
	StatusBindError
	StatusAddDiskError
	StatusDelDiskError
	StatusGetDiskError
	StatusModDiskError
	StatusAddError
	StatusGetError
	StatusDelError
	StatusModError
	StatusAlreadyExist
	StatusNotExist
	StatusRegCompileError
)
