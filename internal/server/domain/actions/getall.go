package actions

func (r *Repo) GetAll() (map[string]map[string]string, error) {
	return r.storage.GetAll(), nil
}
