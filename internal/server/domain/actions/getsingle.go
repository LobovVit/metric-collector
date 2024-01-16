package actions

func (r Repo) GetSingleVal(tp string, name string) (string, error) {
	return r.storage.GetSingle(tp, name)
}
