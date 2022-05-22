package server

type Config struct {
	BindAddr       string
	MaxHeaderBytes int
	ReadTimeout    int
	WriteTimeout   int
}
