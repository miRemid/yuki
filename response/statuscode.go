package response

type StatusCode uint32

const (
	StatusOK StatusCode = iota
	StatusBindError
	StatusSaveDiskError
)
