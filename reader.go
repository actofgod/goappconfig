package goappconfig

import (
	"io"
)

type ByteArrayDecoderFunc func(bytes []byte, value any) error

type Decoder interface {
	Decode(v any) error
}

type DecoderConstructor func(reader io.Reader) Decoder

type bufferedDecoder struct {
	reader      io.Reader
	unmarshal   ByteArrayDecoderFunc
	maxFileSize int64
}

func NewBufferedDecoder(reader io.Reader, unmarshal ByteArrayDecoderFunc) Decoder {
	return &bufferedDecoder{
		reader:      reader,
		unmarshal:   unmarshal,
		maxFileSize: defaultMaxConfigFileSize,
	}
}

func (d *bufferedDecoder) Decode(v any) error {
	buf, err := io.ReadAll(d.reader)
	if err != nil {
		return err
	}
	return d.unmarshal(buf, v)
}
