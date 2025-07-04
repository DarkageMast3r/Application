package api

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

// Mock voor de CatalogusRepository interface
type MockCatalogusRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCatalogusRepositoryMockRecorder
}

type MockCatalogusRepositoryMockRecorder struct {
	mock *MockCatalogusRepository
}

func NewMockCatalogusRepository(ctrl *gomock.Controller) *MockCatalogusRepository {
	mock := &MockCatalogusRepository{ctrl: ctrl}
	mock.recorder = &MockCatalogusRepositoryMockRecorder{mock}
	return mock
}

func (m *MockCatalogusRepository) EXPECT() *MockCatalogusRepositoryMockRecorder {
	return m.recorder
}

func (m *MockCatalogusRepository) Healthcheck(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Healthcheck", c)
}

func (m *MockCatalogusRepository) MaakZorgTechProduct(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MaakZorgTechProduct", c)
}

func (m *MockCatalogusRepository) WijzigZorgTechProduct(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WijzigZorgTechProduct", c)
}

func (m *MockCatalogusRepository) VerwijderZorgTechProduct(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VerwijderZorgTechProduct", c)
}

func (m *MockCatalogusRepository) VoegTechnischDetailToe(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VoegTechnischDetailToe", c)
}

func (m *MockCatalogusRepository) VerwijderTechnischDetail(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VerwijderTechnischDetail", c)
}

func (m *MockCatalogusRepository) GetProductById(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetProductById", c)
}

func (m *MockCatalogusRepository) FindByCategorie(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FindByCategorie", c)
}

func (m *MockCatalogusRepository) ListAlleProducten(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListAlleProducten", c)
}

func (m *MockCatalogusRepository) ZoekOpNaam(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ZoekOpNaam", c)
}

func (mr *MockCatalogusRepositoryMockRecorder) Healthcheck(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Healthcheck", reflect.TypeOf((*catalogusRepository)(nil).Healthcheck), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) MaakZorgTechProduct(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaakZorgTechProduct", reflect.TypeOf((*catalogusRepository)(nil).MaakZorgTechProduct), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) WijzigZorgTechProduct(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WijzigZorgTechProduct", reflect.TypeOf((*catalogusRepository)(nil).WijzigZorgTechProduct), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) VerwijderZorgTechProduct(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerwijderZorgTechProduct", reflect.TypeOf((*catalogusRepository)(nil).VerwijderZorgTechProduct), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) VoegTechnischDetailToe(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VoegTechnischDetailToe", reflect.TypeOf((*catalogusRepository)(nil).VoegTechnischDetailToe), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) VerwijderTechnischDetail(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerwijderTechnischDetail", reflect.TypeOf((*catalogusRepository)(nil).VerwijderTechnischDetail), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) GetProductById(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductById", reflect.TypeOf((*catalogusRepository)(nil).GetProductById), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) FindByCategorie(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCategorie", reflect.TypeOf((*catalogusRepository)(nil).FindByCategorie), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) ListAlleProducten(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAlleProducten", reflect.TypeOf((*catalogusRepository)(nil).ListAlleProducten), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) ZoekOpNaam(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ZoekOpNaam", reflect.TypeOf((*catalogusRepository)(nil).ZoekOpNaam), c)
}
