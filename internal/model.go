package internal

var _ Converter = (*Convert)(nil)

type Converter interface {
	Process()
}
