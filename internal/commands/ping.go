package commands

func Ping() (string, error) {
	resp := "Pong!"
	return resp, nil
}
