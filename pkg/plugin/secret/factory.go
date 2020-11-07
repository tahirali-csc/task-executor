package secret

type Secret struct {
	Metadata map[string]interface{}
	Data     interface{}
}

type SecretsFactory interface {
	NewSecrets() Secrets
}

type Secrets interface {
	Get(string) (*Secret, error)
}
