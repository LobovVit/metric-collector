package actions

func GetSingleVal(tp string, name string) (string, error) {

	return store.storage.GetSingle(tp, name)
}
