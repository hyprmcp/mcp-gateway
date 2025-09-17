package dcr

type ClientIDSource interface {
	GetClientID() string
	GetClientSecret() string
}
