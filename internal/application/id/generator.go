package id

type Generator interface {
	Generate(prefix EntityPrefix) string
}
