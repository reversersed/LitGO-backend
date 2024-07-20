package app

type app struct {
}

var server *app

func Run() error {
	server = &app{}

	return nil
}
func Close() error {
	return nil
}
