package server

type Config struct {
	Ports Ports
}

type Ports struct {
	Admin   uint16
	Meta    uint16
	Network uint16
}
