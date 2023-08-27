package runtime

type Constructor[R any, I any] func(...Configuration[R]) I

func ConstructorFor[R any, I any](defaultValue func() R, casting func(R) I) Constructor[R, I] {
	return func(c ...Configuration[R]) I {
		r := defaultValue()
		for _, configuration := range c {
			r = configuration(r)
		}
		return casting(r)
	}
}
