package actions

func GetAll() (map[string]map[string]string, error) {

	return store.storage.GetAll(), nil
}
