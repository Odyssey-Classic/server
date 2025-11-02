package server

type Config struct {
	Ports   Ports
	DataDir string
}

type Ports struct {
	Admin   uint16
	Meta    uint16
	Network uint16
}
