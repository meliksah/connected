package model

type DataType string

const (
	DataTypeOCR   DataType = "ocr"
	DataTypeError DataType = "error"
)

type Data struct {
	Type DataType `json:"type"`
	Data string   `json:"data"`
}
