package adapter

import (
	"context"

	"vct-platform/backend/internal/domain/federation"
	"vct-platform/backend/internal/store"
)

// ── Province Repository (Postgres) ──

type pgProvinceRepo struct {
	*StoreAdapter[federation.Province]
}

func NewGenericProvinceRepo(ds store.DataStore) federation.ProvinceRepository {
	return &pgProvinceRepo{
		StoreAdapter: NewStoreAdapter[federation.Province](ds, "federation_provinces"),
	}
}

func (r *pgProvinceRepo) List(ctx context.Context) ([]federation.Province, error) {
	return r.StoreAdapter.List()
}

func (r *pgProvinceRepo) GetByID(ctx context.Context, id string) (*federation.Province, error) {
	return r.StoreAdapter.GetByID(id)
}

func (r *pgProvinceRepo) GetByCode(ctx context.Context, code string) (*federation.Province, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range items {
		if p.Code == code {
			return &p, nil
		}
	}
	return nil, federation.ErrNotFound
}

func (r *pgProvinceRepo) Create(ctx context.Context, p federation.Province) (*federation.Province, error) {
	return r.StoreAdapter.Create(p)
}

func (r *pgProvinceRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

func (r *pgProvinceRepo) ListByRegion(ctx context.Context, region federation.RegionCode) ([]federation.Province, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	var res []federation.Province
	for _, p := range items {
		if p.Region == region {
			res = append(res, p)
		}
	}
	return res, nil
}

// ── FederationUnit Repository (Postgres) ──

type pgUnitRepo struct {
	*StoreAdapter[federation.FederationUnit]
}

func NewGenericUnitRepo(ds store.DataStore) federation.FederationUnitRepository {
	return &pgUnitRepo{
		StoreAdapter: NewStoreAdapter[federation.FederationUnit](ds, "federation_units"),
	}
}

func (r *pgUnitRepo) List(ctx context.Context) ([]federation.FederationUnit, error) {
	return r.StoreAdapter.List()
}

func (r *pgUnitRepo) GetByID(ctx context.Context, id string) (*federation.FederationUnit, error) {
	return r.StoreAdapter.GetByID(id)
}

func (r *pgUnitRepo) Create(ctx context.Context, u federation.FederationUnit) (*federation.FederationUnit, error) {
	return r.StoreAdapter.Create(u)
}

func (r *pgUnitRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

func (r *pgUnitRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

func (r *pgUnitRepo) ListByType(ctx context.Context, uType federation.UnitType) ([]federation.FederationUnit, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	var res []federation.FederationUnit
	for _, u := range items {
		if u.Type == uType {
			res = append(res, u)
		}
	}
	return res, nil
}

func (r *pgUnitRepo) ListByParent(ctx context.Context, parentID string) ([]federation.FederationUnit, error) {
	items, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	var res []federation.FederationUnit
	for _, u := range items {
		if u.ParentID == parentID {
			res = append(res, u)
		}
	}
	return res, nil
}

// ── Personnel Repository (Postgres) ──

type pgPersonnelRepo struct {
	*StoreAdapter[federation.PersonnelAssignment]
}

func NewGenericPersonnelRepo(ds store.DataStore) federation.PersonnelRepository {
	return &pgPersonnelRepo{
		StoreAdapter: NewStoreAdapter[federation.PersonnelAssignment](ds, "federation_personnel"),
	}
}

func (r *pgPersonnelRepo) List(ctx context.Context, unitID string) ([]federation.PersonnelAssignment, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []federation.PersonnelAssignment
	for _, a := range items {
		if unitID == "" || a.UnitID == unitID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (r *pgPersonnelRepo) Create(ctx context.Context, a federation.PersonnelAssignment) error {
	_, err := r.StoreAdapter.Create(a)
	return err
}

func (r *pgPersonnelRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

func (r *pgPersonnelRepo) Deactivate(ctx context.Context, id string) error {
	_, err := r.StoreAdapter.Update(id, map[string]interface{}{"is_active": false})
	return err
}

func (r *pgPersonnelRepo) GetByUserAndUnit(ctx context.Context, userID, unitID string) (*federation.PersonnelAssignment, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	for _, a := range items {
		if a.UserID == userID && a.UnitID == unitID {
			return &a, nil
		}
	}
	return nil, federation.ErrNotFound
}
