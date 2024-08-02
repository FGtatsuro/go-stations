package model

type ErrNotFound struct{}

func (*ErrNotFound) Error() string {
	return "Not Found"
}
