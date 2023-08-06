package runtime

type Configuration[R any] func(R) R

type Constructor[R any, I any] func(...Configuration[R]) I
