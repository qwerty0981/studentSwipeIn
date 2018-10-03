package delivery

type DriverInitializer interface {
	Title() string
	Configure() (Driver, error)
}

type Driver interface {
	Input(map[string]string) error
}
