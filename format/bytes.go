package format

type BytesFormat struct {
}

func (helper BytesFormat) Default() []byte {
	return []byte{}
}

func (helper BytesFormat) Marshal(value []byte) ([]byte, error) {
	return value, nil
}

func (helper BytesFormat) Unmarshal(value []byte) ([]byte, error) {
	return value, nil
}

func (helper BytesFormat) Mask(value []byte) ([]byte, error) {
	return value, nil
}

func (helper BytesFormat) Unmask(value []byte) ([]byte, error) {
	return value, nil
}

func Bytes() Format[[]byte] {
	return BytesFormat{}
}

func Plain() Mask {
	return BytesFormat{}
}
