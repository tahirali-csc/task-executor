package secret

type Manager interface {
	NewSecretsFactory() (SecretsFactory, error)
}

//type ManagerFactory interface {
//	NewInstance(interface{}) (Manager, error)
//}
