package oauth

import "errors"

var (
	clientFactory    = make(map[string]ClientCreator)
	clientRepository = make(map[string]Client)
)

func Register(name string, creator ClientCreator) {
	clientFactory[name] = creator
}

func Resolve(driver string) (Client, error) {
	if client, exists := clientRepository[driver]; exists {
		return client, nil
	}

	if _, exists := clientFactory[driver]; !exists {
		return nil, errors.New("Driver " + driver + " is not supported")
	}

	clientRepository[driver] = clientFactory[driver]()

	return clientRepository[driver], nil
}
