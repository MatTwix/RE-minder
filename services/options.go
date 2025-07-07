package services

type Operator string

const (
	Equal Operator = "="
)

type Options struct {
	Condition *Condition
}

type Condition struct {
	Field    string
	Operator Operator
	Value    any
}
