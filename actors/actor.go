package actors

type Actor[T any] interface {
	Start() T
}
