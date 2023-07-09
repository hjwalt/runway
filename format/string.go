package format

type StringFormat struct {
}

func (helper StringFormat) Default() string {
	return ""
}

func (helper StringFormat) Marshal(value string) ([]byte, error) {
	return []byte(value), nil
}

func (helper StringFormat) Unmarshal(value []byte) (string, error) {
	return string(value), nil
}

func String() Format[string] {
	return StringFormat{}
}
