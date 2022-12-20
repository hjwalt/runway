package inverse

type QualifierContext struct {
	Name string
}

func Qualifier(name string) QualifierContext {
	return QualifierContext{Name: name}
}
