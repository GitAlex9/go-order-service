package id

type EntityPrefix string

const (
	ProductPrefix  EntityPrefix = "PROD"
	CustomerPrefix EntityPrefix = "CUS"
	OrderPrefix    EntityPrefix = "PED"
)

func (p EntityPrefix) String() string {
	return string(p)
}
