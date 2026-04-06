package adapter

import (
	"context"

	"vct-platform/backend/internal/domain/provincial"
	"vct-platform/backend/internal/store"
)

// ── Association ──────────────────────────────────────────────────────────────

type pgAssociationRepo struct {
	*StoreAdapter[provincial.Association]
}

func NewPgAssociationRepo(ds store.DataStore) provincial.AssociationRepository {
	return &pgAssociationRepo{
		StoreAdapter: NewStoreAdapter[provincial.Association](ds, "provincial_associations"),
	}
}
func (r *pgAssociationRepo) List(ctx context.Context, provinceID string) ([]provincial.Association, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.Association
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgAssociationRepo) GetByID(ctx context.Context, id string) (*provincial.Association, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgAssociationRepo) Create(ctx context.Context, item provincial.Association) (*provincial.Association, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgAssociationRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgAssociationRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Sub-Association ──────────────────────────────────────────────────────────

type pgSubAssociationRepo struct {
	*StoreAdapter[provincial.SubAssociation]
}

func NewPgSubAssociationRepo(ds store.DataStore) provincial.SubAssociationRepository {
	return &pgSubAssociationRepo{
		StoreAdapter: NewStoreAdapter[provincial.SubAssociation](ds, "provincial_sub_associations"),
	}
}
func (r *pgSubAssociationRepo) List(ctx context.Context, provinceID string) ([]provincial.SubAssociation, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.SubAssociation
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgSubAssociationRepo) ListByAssociation(ctx context.Context, assocID string) ([]provincial.SubAssociation, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.SubAssociation
	for _, it := range items {
		if it.AssociationID == assocID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgSubAssociationRepo) GetByID(ctx context.Context, id string) (*provincial.SubAssociation, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgSubAssociationRepo) Create(ctx context.Context, item provincial.SubAssociation) (*provincial.SubAssociation, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgSubAssociationRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgSubAssociationRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Club ─────────────────────────────────────────────────────────────────────

type pgClubRepo struct {
	*StoreAdapter[provincial.ProvincialClub]
}

func NewPgClubRepo(ds store.DataStore) provincial.ClubRepository {
	return &pgClubRepo{
		StoreAdapter: NewStoreAdapter[provincial.ProvincialClub](ds, "provincial_clubs"),
	}
}
func (r *pgClubRepo) List(ctx context.Context, provinceID string) ([]provincial.ProvincialClub, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialClub
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgClubRepo) GetByID(ctx context.Context, id string) (*provincial.ProvincialClub, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgClubRepo) Create(ctx context.Context, item provincial.ProvincialClub) (*provincial.ProvincialClub, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgClubRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgClubRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Athlete ──────────────────────────────────────────────────────────────────

type pgAthleteRepo struct {
	*StoreAdapter[provincial.ProvincialAthlete]
}

func NewPgAthleteRepo(ds store.DataStore) provincial.AthleteRepository {
	return &pgAthleteRepo{
		StoreAdapter: NewStoreAdapter[provincial.ProvincialAthlete](ds, "provincial_athletes"),
	}
}
func (r *pgAthleteRepo) List(ctx context.Context, provinceID string) ([]provincial.ProvincialAthlete, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialAthlete
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgAthleteRepo) ListByClub(ctx context.Context, clubID string) ([]provincial.ProvincialAthlete, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialAthlete
	for _, it := range items {
		if it.ClubID == clubID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgAthleteRepo) GetByID(ctx context.Context, id string) (*provincial.ProvincialAthlete, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgAthleteRepo) Create(ctx context.Context, item provincial.ProvincialAthlete) (*provincial.ProvincialAthlete, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgAthleteRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgAthleteRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Coach ────────────────────────────────────────────────────────────────────

type pgCoachRepo struct {
	*StoreAdapter[provincial.ProvincialCoach]
}

func NewPgCoachRepo(ds store.DataStore) provincial.CoachRepository {
	return &pgCoachRepo{
		StoreAdapter: NewStoreAdapter[provincial.ProvincialCoach](ds, "provincial_coaches"),
	}
}
func (r *pgCoachRepo) List(ctx context.Context, provinceID string) ([]provincial.ProvincialCoach, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialCoach
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgCoachRepo) GetByID(ctx context.Context, id string) (*provincial.ProvincialCoach, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgCoachRepo) Create(ctx context.Context, item provincial.ProvincialCoach) (*provincial.ProvincialCoach, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgCoachRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

// ── Referee ──────────────────────────────────────────────────────────────────

type pgRefereeRepo struct {
	*StoreAdapter[provincial.ProvincialReferee]
}

func NewPgRefereeRepo(ds store.DataStore) provincial.RefereeRepository {
	return &pgRefereeRepo{
		StoreAdapter: NewStoreAdapter[provincial.ProvincialReferee](ds, "provincial_referees"),
	}
}
func (r *pgRefereeRepo) List(ctx context.Context, provinceID string) ([]provincial.ProvincialReferee, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialReferee
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgRefereeRepo) GetByID(ctx context.Context, id string) (*provincial.ProvincialReferee, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgRefereeRepo) Create(ctx context.Context, item provincial.ProvincialReferee) (*provincial.ProvincialReferee, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgRefereeRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgRefereeRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Referee Cert ─────────────────────────────────────────────────────────────

type pgRefereeCertRepo struct {
	*StoreAdapter[provincial.RefereeCertificate]
}

func NewPgRefereeCertRepo(ds store.DataStore) provincial.RefereeCertificateRepository {
	return &pgRefereeCertRepo{
		StoreAdapter: NewStoreAdapter[provincial.RefereeCertificate](ds, "provincial_referee_certs"),
	}
}
func (r *pgRefereeCertRepo) ListByReferee(ctx context.Context, refereeID string) ([]provincial.RefereeCertificate, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.RefereeCertificate
	for _, it := range items {
		if it.RefereeID == refereeID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgRefereeCertRepo) Create(ctx context.Context, item provincial.RefereeCertificate) (*provincial.RefereeCertificate, error) {
	return r.StoreAdapter.Create(item)
}

// ── Committee ────────────────────────────────────────────────────────────────

type pgCommitteeRepo struct {
	*StoreAdapter[provincial.CommitteeMember]
}

func NewPgCommitteeRepo(ds store.DataStore) provincial.CommitteeRepository {
	return &pgCommitteeRepo{
		StoreAdapter: NewStoreAdapter[provincial.CommitteeMember](ds, "provincial_committee"),
	}
}
func (r *pgCommitteeRepo) List(ctx context.Context, provinceID string) ([]provincial.CommitteeMember, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.CommitteeMember
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgCommitteeRepo) GetByID(ctx context.Context, id string) (*provincial.CommitteeMember, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgCommitteeRepo) Create(ctx context.Context, item provincial.CommitteeMember) (*provincial.CommitteeMember, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgCommitteeRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

// ── Transfer ─────────────────────────────────────────────────────────────────

type pgTransferRepo struct {
	*StoreAdapter[provincial.ClubTransfer]
}

func NewPgTransferRepo(ds store.DataStore) provincial.TransferRepository {
	return &pgTransferRepo{
		StoreAdapter: NewStoreAdapter[provincial.ClubTransfer](ds, "provincial_transfers"),
	}
}
func (r *pgTransferRepo) List(ctx context.Context, provinceID string) ([]provincial.ClubTransfer, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ClubTransfer
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgTransferRepo) GetByID(ctx context.Context, id string) (*provincial.ClubTransfer, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgTransferRepo) Create(ctx context.Context, item provincial.ClubTransfer) (*provincial.ClubTransfer, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgTransferRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

// ── Club Class ───────────────────────────────────────────────────────────────

type pgClubClassRepo struct {
	*StoreAdapter[provincial.ClubClass]
}

func NewPgClubClassRepo(ds store.DataStore) provincial.ClubClassRepository {
	return &pgClubClassRepo{
		StoreAdapter: NewStoreAdapter[provincial.ClubClass](ds, "provincial_club_classes"),
	}
}
func (r *pgClubClassRepo) List(ctx context.Context, clubID string) ([]provincial.ClubClass, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ClubClass
	for _, it := range items {
		if it.ClubID == clubID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgClubClassRepo) GetByID(ctx context.Context, id string) (*provincial.ClubClass, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgClubClassRepo) Create(ctx context.Context, item provincial.ClubClass) (*provincial.ClubClass, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgClubClassRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgClubClassRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Club Member ──────────────────────────────────────────────────────────────

type pgClubMemberRepo struct {
	*StoreAdapter[provincial.ClubMember]
}

func NewPgClubMemberRepo(ds store.DataStore) provincial.ClubMemberRepository {
	return &pgClubMemberRepo{
		StoreAdapter: NewStoreAdapter[provincial.ClubMember](ds, "provincial_club_members"),
	}
}
func (r *pgClubMemberRepo) List(ctx context.Context, clubID string) ([]provincial.ClubMember, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ClubMember
	for _, it := range items {
		if it.ClubID == clubID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgClubMemberRepo) GetByID(ctx context.Context, id string) (*provincial.ClubMember, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgClubMemberRepo) Create(ctx context.Context, item provincial.ClubMember) (*provincial.ClubMember, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgClubMemberRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgClubMemberRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Club Finance ─────────────────────────────────────────────────────────────

type pgClubFinanceRepo struct {
	*StoreAdapter[provincial.ClubFinanceEntry]
}

func NewPgClubFinanceRepo(ds store.DataStore) provincial.ClubFinanceRepository {
	return &pgClubFinanceRepo{
		StoreAdapter: NewStoreAdapter[provincial.ClubFinanceEntry](ds, "provincial_club_finance"),
	}
}
func (r *pgClubFinanceRepo) List(ctx context.Context, clubID string) ([]provincial.ClubFinanceEntry, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ClubFinanceEntry
	for _, it := range items {
		if it.ClubID == clubID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgClubFinanceRepo) GetByID(ctx context.Context, id string) (*provincial.ClubFinanceEntry, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgClubFinanceRepo) Create(ctx context.Context, item provincial.ClubFinanceEntry) (*provincial.ClubFinanceEntry, error) {
	return r.StoreAdapter.Create(item)
}

// ── Vo Sinh ──────────────────────────────────────────────────────────────────

type pgVoSinhRepo struct {
	*StoreAdapter[provincial.VoSinh]
}

func NewPgVoSinhRepo(ds store.DataStore) provincial.VoSinhRepository {
	return &pgVoSinhRepo{
		StoreAdapter: NewStoreAdapter[provincial.VoSinh](ds, "provincial_vosinh"),
	}
}
func (r *pgVoSinhRepo) List(ctx context.Context, provinceID string) ([]provincial.VoSinh, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.VoSinh
	for _, it := range items {
		if provinceID == "" || it.ProvinceID == provinceID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgVoSinhRepo) ListByClub(ctx context.Context, clubID string) ([]provincial.VoSinh, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.VoSinh
	for _, it := range items {
		if it.ClubID == clubID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgVoSinhRepo) GetByID(ctx context.Context, id string) (*provincial.VoSinh, error) {
	return r.StoreAdapter.GetByID(id)
}
func (r *pgVoSinhRepo) Create(ctx context.Context, item provincial.VoSinh) (*provincial.VoSinh, error) {
	return r.StoreAdapter.Create(item)
}
func (r *pgVoSinhRepo) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
func (r *pgVoSinhRepo) Delete(ctx context.Context, id string) error {
	return r.StoreAdapter.Delete(id)
}

// ── Belt History ─────────────────────────────────────────────────────────────

type pgBeltHistoryRepo struct {
	*StoreAdapter[provincial.BeltHistory]
}

func NewPgBeltHistoryRepo(ds store.DataStore) provincial.BeltHistoryRepository {
	return &pgBeltHistoryRepo{
		StoreAdapter: NewStoreAdapter[provincial.BeltHistory](ds, "provincial_belt_history"),
	}
}
func (r *pgBeltHistoryRepo) ListByVoSinh(ctx context.Context, voSinhID string) ([]provincial.BeltHistory, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.BeltHistory
	for _, it := range items {
		if it.VoSinhID == voSinhID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (r *pgBeltHistoryRepo) Create(ctx context.Context, item provincial.BeltHistory) (*provincial.BeltHistory, error) {
	return r.StoreAdapter.Create(item)
}
