package actions

import "context"

func (r *Repo) GetAll(ctx context.Context) (map[string]map[string]string, error) {
	return r.storage.GetAll(ctx)
}
