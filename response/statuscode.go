package response

type StatusCode uint32

const (
	StatusOK StatusCode = iota
	StatusBindError
	StatusSaveDiskError
	StatusAddError
	StatusGetError
	StatusDelError
	StatusModError

	StatusAlreadyExist
)
