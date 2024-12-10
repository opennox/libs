package binenc

type Encoded interface {
	EncodeSize() int
	Encode(data []byte) (int, error)
	Decode(data []byte) (int, error)
}
