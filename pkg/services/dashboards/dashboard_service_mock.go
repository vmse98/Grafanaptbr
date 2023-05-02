// Code generated by mockery v2.23.2. DO NOT EDIT.

package dashboards

import (
	context "context"

	folder "github.com/grafana/grafana/pkg/services/folder"
	mock "github.com/stretchr/testify/mock"

	model "github.com/grafana/grafana/pkg/services/search/model"

	user "github.com/grafana/grafana/pkg/services/user"
)

// FakeDashboardService is an autogenerated mock type for the DashboardService type
type FakeDashboardService struct {
	mock.Mock
}

type FakeDashboardService_Expecter struct {
	mock *mock.Mock
}

func (_m *FakeDashboardService) EXPECT() *FakeDashboardService_Expecter {
	return &FakeDashboardService_Expecter{mock: &_m.Mock}
}

// BuildSaveDashboardCommand provides a mock function with given fields: ctx, dto, shouldValidateAlerts, validateProvisionedDashboard
func (_m *FakeDashboardService) BuildSaveDashboardCommand(ctx context.Context, dto *SaveDashboardDTO, shouldValidateAlerts bool, validateProvisionedDashboard bool) (*SaveDashboardCommand, error) {
	ret := _m.Called(ctx, dto, shouldValidateAlerts, validateProvisionedDashboard)

	var r0 *SaveDashboardCommand
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *SaveDashboardDTO, bool, bool) (*SaveDashboardCommand, error)); ok {
		return rf(ctx, dto, shouldValidateAlerts, validateProvisionedDashboard)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *SaveDashboardDTO, bool, bool) *SaveDashboardCommand); ok {
		r0 = rf(ctx, dto, shouldValidateAlerts, validateProvisionedDashboard)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SaveDashboardCommand)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *SaveDashboardDTO, bool, bool) error); ok {
		r1 = rf(ctx, dto, shouldValidateAlerts, validateProvisionedDashboard)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_BuildSaveDashboardCommand_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BuildSaveDashboardCommand'
type FakeDashboardService_BuildSaveDashboardCommand_Call struct {
	*mock.Call
}

// BuildSaveDashboardCommand is a helper method to define mock.On call
//   - ctx context.Context
//   - dto *SaveDashboardDTO
//   - shouldValidateAlerts bool
//   - validateProvisionedDashboard bool
func (_e *FakeDashboardService_Expecter) BuildSaveDashboardCommand(ctx interface{}, dto interface{}, shouldValidateAlerts interface{}, validateProvisionedDashboard interface{}) *FakeDashboardService_BuildSaveDashboardCommand_Call {
	return &FakeDashboardService_BuildSaveDashboardCommand_Call{Call: _e.mock.On("BuildSaveDashboardCommand", ctx, dto, shouldValidateAlerts, validateProvisionedDashboard)}
}

func (_c *FakeDashboardService_BuildSaveDashboardCommand_Call) Run(run func(ctx context.Context, dto *SaveDashboardDTO, shouldValidateAlerts bool, validateProvisionedDashboard bool)) *FakeDashboardService_BuildSaveDashboardCommand_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*SaveDashboardDTO), args[2].(bool), args[3].(bool))
	})
	return _c
}

func (_c *FakeDashboardService_BuildSaveDashboardCommand_Call) Return(_a0 *SaveDashboardCommand, _a1 error) *FakeDashboardService_BuildSaveDashboardCommand_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_BuildSaveDashboardCommand_Call) RunAndReturn(run func(context.Context, *SaveDashboardDTO, bool, bool) (*SaveDashboardCommand, error)) *FakeDashboardService_BuildSaveDashboardCommand_Call {
	_c.Call.Return(run)
	return _c
}

// CountInFolder provides a mock function with given fields: ctx, orgID, uid, _a3
func (_m *FakeDashboardService) CountInFolder(ctx context.Context, orgID int64, uid string, _a3 *user.SignedInUser) (int64, error) {
	ret := _m.Called(ctx, orgID, uid, _a3)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, *user.SignedInUser) (int64, error)); ok {
		return rf(ctx, orgID, uid, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, *user.SignedInUser) int64); ok {
		r0 = rf(ctx, orgID, uid, _a3)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string, *user.SignedInUser) error); ok {
		r1 = rf(ctx, orgID, uid, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_CountInFolder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountInFolder'
type FakeDashboardService_CountInFolder_Call struct {
	*mock.Call
}

// CountInFolder is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID int64
//   - uid string
//   - _a3 *user.SignedInUser
func (_e *FakeDashboardService_Expecter) CountInFolder(ctx interface{}, orgID interface{}, uid interface{}, _a3 interface{}) *FakeDashboardService_CountInFolder_Call {
	return &FakeDashboardService_CountInFolder_Call{Call: _e.mock.On("CountInFolder", ctx, orgID, uid, _a3)}
}

func (_c *FakeDashboardService_CountInFolder_Call) Run(run func(ctx context.Context, orgID int64, uid string, _a3 *user.SignedInUser)) *FakeDashboardService_CountInFolder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string), args[3].(*user.SignedInUser))
	})
	return _c
}

func (_c *FakeDashboardService_CountInFolder_Call) Return(_a0 int64, _a1 error) *FakeDashboardService_CountInFolder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_CountInFolder_Call) RunAndReturn(run func(context.Context, int64, string, *user.SignedInUser) (int64, error)) *FakeDashboardService_CountInFolder_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteACLByUser provides a mock function with given fields: ctx, userID
func (_m *FakeDashboardService) DeleteACLByUser(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FakeDashboardService_DeleteACLByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteACLByUser'
type FakeDashboardService_DeleteACLByUser_Call struct {
	*mock.Call
}

// DeleteACLByUser is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int64
func (_e *FakeDashboardService_Expecter) DeleteACLByUser(ctx interface{}, userID interface{}) *FakeDashboardService_DeleteACLByUser_Call {
	return &FakeDashboardService_DeleteACLByUser_Call{Call: _e.mock.On("DeleteACLByUser", ctx, userID)}
}

func (_c *FakeDashboardService_DeleteACLByUser_Call) Run(run func(ctx context.Context, userID int64)) *FakeDashboardService_DeleteACLByUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *FakeDashboardService_DeleteACLByUser_Call) Return(_a0 error) *FakeDashboardService_DeleteACLByUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FakeDashboardService_DeleteACLByUser_Call) RunAndReturn(run func(context.Context, int64) error) *FakeDashboardService_DeleteACLByUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteDashboard provides a mock function with given fields: ctx, dashboardId, orgId
func (_m *FakeDashboardService) DeleteDashboard(ctx context.Context, dashboardId int64, orgId int64) error {
	ret := _m.Called(ctx, dashboardId, orgId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, dashboardId, orgId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FakeDashboardService_DeleteDashboard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteDashboard'
type FakeDashboardService_DeleteDashboard_Call struct {
	*mock.Call
}

// DeleteDashboard is a helper method to define mock.On call
//   - ctx context.Context
//   - dashboardId int64
//   - orgId int64
func (_e *FakeDashboardService_Expecter) DeleteDashboard(ctx interface{}, dashboardId interface{}, orgId interface{}) *FakeDashboardService_DeleteDashboard_Call {
	return &FakeDashboardService_DeleteDashboard_Call{Call: _e.mock.On("DeleteDashboard", ctx, dashboardId, orgId)}
}

func (_c *FakeDashboardService_DeleteDashboard_Call) Run(run func(ctx context.Context, dashboardId int64, orgId int64)) *FakeDashboardService_DeleteDashboard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(int64))
	})
	return _c
}

func (_c *FakeDashboardService_DeleteDashboard_Call) Return(_a0 error) *FakeDashboardService_DeleteDashboard_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FakeDashboardService_DeleteDashboard_Call) RunAndReturn(run func(context.Context, int64, int64) error) *FakeDashboardService_DeleteDashboard_Call {
	_c.Call.Return(run)
	return _c
}

// FindDashboards provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) FindDashboards(ctx context.Context, query *FindPersistedDashboardsQuery) ([]DashboardSearchProjection, error) {
	ret := _m.Called(ctx, query)

	var r0 []DashboardSearchProjection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *FindPersistedDashboardsQuery) ([]DashboardSearchProjection, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *FindPersistedDashboardsQuery) []DashboardSearchProjection); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]DashboardSearchProjection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *FindPersistedDashboardsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_FindDashboards_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindDashboards'
type FakeDashboardService_FindDashboards_Call struct {
	*mock.Call
}

// FindDashboards is a helper method to define mock.On call
//   - ctx context.Context
//   - query *FindPersistedDashboardsQuery
func (_e *FakeDashboardService_Expecter) FindDashboards(ctx interface{}, query interface{}) *FakeDashboardService_FindDashboards_Call {
	return &FakeDashboardService_FindDashboards_Call{Call: _e.mock.On("FindDashboards", ctx, query)}
}

func (_c *FakeDashboardService_FindDashboards_Call) Run(run func(ctx context.Context, query *FindPersistedDashboardsQuery)) *FakeDashboardService_FindDashboards_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*FindPersistedDashboardsQuery))
	})
	return _c
}

func (_c *FakeDashboardService_FindDashboards_Call) Return(_a0 []DashboardSearchProjection, _a1 error) *FakeDashboardService_FindDashboards_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_FindDashboards_Call) RunAndReturn(run func(context.Context, *FindPersistedDashboardsQuery) ([]DashboardSearchProjection, error)) *FakeDashboardService_FindDashboards_Call {
	_c.Call.Return(run)
	return _c
}

// GetDashboard provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) GetDashboard(ctx context.Context, query *GetDashboardQuery) (*Dashboard, error) {
	ret := _m.Called(ctx, query)

	var r0 *Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardQuery) (*Dashboard, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardQuery) *Dashboard); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_GetDashboard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDashboard'
type FakeDashboardService_GetDashboard_Call struct {
	*mock.Call
}

// GetDashboard is a helper method to define mock.On call
//   - ctx context.Context
//   - query *GetDashboardQuery
func (_e *FakeDashboardService_Expecter) GetDashboard(ctx interface{}, query interface{}) *FakeDashboardService_GetDashboard_Call {
	return &FakeDashboardService_GetDashboard_Call{Call: _e.mock.On("GetDashboard", ctx, query)}
}

func (_c *FakeDashboardService_GetDashboard_Call) Run(run func(ctx context.Context, query *GetDashboardQuery)) *FakeDashboardService_GetDashboard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*GetDashboardQuery))
	})
	return _c
}

func (_c *FakeDashboardService_GetDashboard_Call) Return(_a0 *Dashboard, _a1 error) *FakeDashboardService_GetDashboard_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_GetDashboard_Call) RunAndReturn(run func(context.Context, *GetDashboardQuery) (*Dashboard, error)) *FakeDashboardService_GetDashboard_Call {
	_c.Call.Return(run)
	return _c
}

// GetDashboardACLInfoList provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) GetDashboardACLInfoList(ctx context.Context, query *GetDashboardACLInfoListQuery) ([]*DashboardACLInfoDTO, error) {
	ret := _m.Called(ctx, query)

	var r0 []*DashboardACLInfoDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardACLInfoListQuery) ([]*DashboardACLInfoDTO, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardACLInfoListQuery) []*DashboardACLInfoDTO); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*DashboardACLInfoDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardACLInfoListQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_GetDashboardACLInfoList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDashboardACLInfoList'
type FakeDashboardService_GetDashboardACLInfoList_Call struct {
	*mock.Call
}

// GetDashboardACLInfoList is a helper method to define mock.On call
//   - ctx context.Context
//   - query *GetDashboardACLInfoListQuery
func (_e *FakeDashboardService_Expecter) GetDashboardACLInfoList(ctx interface{}, query interface{}) *FakeDashboardService_GetDashboardACLInfoList_Call {
	return &FakeDashboardService_GetDashboardACLInfoList_Call{Call: _e.mock.On("GetDashboardACLInfoList", ctx, query)}
}

func (_c *FakeDashboardService_GetDashboardACLInfoList_Call) Run(run func(ctx context.Context, query *GetDashboardACLInfoListQuery)) *FakeDashboardService_GetDashboardACLInfoList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*GetDashboardACLInfoListQuery))
	})
	return _c
}

func (_c *FakeDashboardService_GetDashboardACLInfoList_Call) Return(_a0 []*DashboardACLInfoDTO, _a1 error) *FakeDashboardService_GetDashboardACLInfoList_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_GetDashboardACLInfoList_Call) RunAndReturn(run func(context.Context, *GetDashboardACLInfoListQuery) ([]*DashboardACLInfoDTO, error)) *FakeDashboardService_GetDashboardACLInfoList_Call {
	_c.Call.Return(run)
	return _c
}

// GetDashboardTags provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) GetDashboardTags(ctx context.Context, query *GetDashboardTagsQuery) ([]*DashboardTagCloudItem, error) {
	ret := _m.Called(ctx, query)

	var r0 []*DashboardTagCloudItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardTagsQuery) ([]*DashboardTagCloudItem, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardTagsQuery) []*DashboardTagCloudItem); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*DashboardTagCloudItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardTagsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_GetDashboardTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDashboardTags'
type FakeDashboardService_GetDashboardTags_Call struct {
	*mock.Call
}

// GetDashboardTags is a helper method to define mock.On call
//   - ctx context.Context
//   - query *GetDashboardTagsQuery
func (_e *FakeDashboardService_Expecter) GetDashboardTags(ctx interface{}, query interface{}) *FakeDashboardService_GetDashboardTags_Call {
	return &FakeDashboardService_GetDashboardTags_Call{Call: _e.mock.On("GetDashboardTags", ctx, query)}
}

func (_c *FakeDashboardService_GetDashboardTags_Call) Run(run func(ctx context.Context, query *GetDashboardTagsQuery)) *FakeDashboardService_GetDashboardTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*GetDashboardTagsQuery))
	})
	return _c
}

func (_c *FakeDashboardService_GetDashboardTags_Call) Return(_a0 []*DashboardTagCloudItem, _a1 error) *FakeDashboardService_GetDashboardTags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_GetDashboardTags_Call) RunAndReturn(run func(context.Context, *GetDashboardTagsQuery) ([]*DashboardTagCloudItem, error)) *FakeDashboardService_GetDashboardTags_Call {
	_c.Call.Return(run)
	return _c
}

// GetDashboardUIDByID provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) GetDashboardUIDByID(ctx context.Context, query *GetDashboardRefByIDQuery) (*DashboardRef, error) {
	ret := _m.Called(ctx, query)

	var r0 *DashboardRef
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardRefByIDQuery) (*DashboardRef, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardRefByIDQuery) *DashboardRef); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DashboardRef)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardRefByIDQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_GetDashboardUIDByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDashboardUIDByID'
type FakeDashboardService_GetDashboardUIDByID_Call struct {
	*mock.Call
}

// GetDashboardUIDByID is a helper method to define mock.On call
//   - ctx context.Context
//   - query *GetDashboardRefByIDQuery
func (_e *FakeDashboardService_Expecter) GetDashboardUIDByID(ctx interface{}, query interface{}) *FakeDashboardService_GetDashboardUIDByID_Call {
	return &FakeDashboardService_GetDashboardUIDByID_Call{Call: _e.mock.On("GetDashboardUIDByID", ctx, query)}
}

func (_c *FakeDashboardService_GetDashboardUIDByID_Call) Run(run func(ctx context.Context, query *GetDashboardRefByIDQuery)) *FakeDashboardService_GetDashboardUIDByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*GetDashboardRefByIDQuery))
	})
	return _c
}

func (_c *FakeDashboardService_GetDashboardUIDByID_Call) Return(_a0 *DashboardRef, _a1 error) *FakeDashboardService_GetDashboardUIDByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_GetDashboardUIDByID_Call) RunAndReturn(run func(context.Context, *GetDashboardRefByIDQuery) (*DashboardRef, error)) *FakeDashboardService_GetDashboardUIDByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetDashboards provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) GetDashboards(ctx context.Context, query *GetDashboardsQuery) ([]*Dashboard, error) {
	ret := _m.Called(ctx, query)

	var r0 []*Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardsQuery) ([]*Dashboard, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardsQuery) []*Dashboard); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_GetDashboards_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDashboards'
type FakeDashboardService_GetDashboards_Call struct {
	*mock.Call
}

// GetDashboards is a helper method to define mock.On call
//   - ctx context.Context
//   - query *GetDashboardsQuery
func (_e *FakeDashboardService_Expecter) GetDashboards(ctx interface{}, query interface{}) *FakeDashboardService_GetDashboards_Call {
	return &FakeDashboardService_GetDashboards_Call{Call: _e.mock.On("GetDashboards", ctx, query)}
}

func (_c *FakeDashboardService_GetDashboards_Call) Run(run func(ctx context.Context, query *GetDashboardsQuery)) *FakeDashboardService_GetDashboards_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*GetDashboardsQuery))
	})
	return _c
}

func (_c *FakeDashboardService_GetDashboards_Call) Return(_a0 []*Dashboard, _a1 error) *FakeDashboardService_GetDashboards_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_GetDashboards_Call) RunAndReturn(run func(context.Context, *GetDashboardsQuery) ([]*Dashboard, error)) *FakeDashboardService_GetDashboards_Call {
	_c.Call.Return(run)
	return _c
}

// HasAdminPermissionInDashboardsOrFolders provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) HasAdminPermissionInDashboardsOrFolders(ctx context.Context, query *folder.HasAdminPermissionInDashboardsOrFoldersQuery) (bool, error) {
	ret := _m.Called(ctx, query)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *folder.HasAdminPermissionInDashboardsOrFoldersQuery) (bool, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *folder.HasAdminPermissionInDashboardsOrFoldersQuery) bool); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *folder.HasAdminPermissionInDashboardsOrFoldersQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasAdminPermissionInDashboardsOrFolders'
type FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call struct {
	*mock.Call
}

// HasAdminPermissionInDashboardsOrFolders is a helper method to define mock.On call
//   - ctx context.Context
//   - query *folder.HasAdminPermissionInDashboardsOrFoldersQuery
func (_e *FakeDashboardService_Expecter) HasAdminPermissionInDashboardsOrFolders(ctx interface{}, query interface{}) *FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call {
	return &FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call{Call: _e.mock.On("HasAdminPermissionInDashboardsOrFolders", ctx, query)}
}

func (_c *FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call) Run(run func(ctx context.Context, query *folder.HasAdminPermissionInDashboardsOrFoldersQuery)) *FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*folder.HasAdminPermissionInDashboardsOrFoldersQuery))
	})
	return _c
}

func (_c *FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call) Return(_a0 bool, _a1 error) *FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call) RunAndReturn(run func(context.Context, *folder.HasAdminPermissionInDashboardsOrFoldersQuery) (bool, error)) *FakeDashboardService_HasAdminPermissionInDashboardsOrFolders_Call {
	_c.Call.Return(run)
	return _c
}

// HasEditPermissionInFolders provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) HasEditPermissionInFolders(ctx context.Context, query *folder.HasEditPermissionInFoldersQuery) (bool, error) {
	ret := _m.Called(ctx, query)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *folder.HasEditPermissionInFoldersQuery) (bool, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *folder.HasEditPermissionInFoldersQuery) bool); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *folder.HasEditPermissionInFoldersQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_HasEditPermissionInFolders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasEditPermissionInFolders'
type FakeDashboardService_HasEditPermissionInFolders_Call struct {
	*mock.Call
}

// HasEditPermissionInFolders is a helper method to define mock.On call
//   - ctx context.Context
//   - query *folder.HasEditPermissionInFoldersQuery
func (_e *FakeDashboardService_Expecter) HasEditPermissionInFolders(ctx interface{}, query interface{}) *FakeDashboardService_HasEditPermissionInFolders_Call {
	return &FakeDashboardService_HasEditPermissionInFolders_Call{Call: _e.mock.On("HasEditPermissionInFolders", ctx, query)}
}

func (_c *FakeDashboardService_HasEditPermissionInFolders_Call) Run(run func(ctx context.Context, query *folder.HasEditPermissionInFoldersQuery)) *FakeDashboardService_HasEditPermissionInFolders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*folder.HasEditPermissionInFoldersQuery))
	})
	return _c
}

func (_c *FakeDashboardService_HasEditPermissionInFolders_Call) Return(_a0 bool, _a1 error) *FakeDashboardService_HasEditPermissionInFolders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_HasEditPermissionInFolders_Call) RunAndReturn(run func(context.Context, *folder.HasEditPermissionInFoldersQuery) (bool, error)) *FakeDashboardService_HasEditPermissionInFolders_Call {
	_c.Call.Return(run)
	return _c
}

// ImportDashboard provides a mock function with given fields: ctx, dto
func (_m *FakeDashboardService) ImportDashboard(ctx context.Context, dto *SaveDashboardDTO) (*Dashboard, error) {
	ret := _m.Called(ctx, dto)

	var r0 *Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *SaveDashboardDTO) (*Dashboard, error)); ok {
		return rf(ctx, dto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *SaveDashboardDTO) *Dashboard); ok {
		r0 = rf(ctx, dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *SaveDashboardDTO) error); ok {
		r1 = rf(ctx, dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_ImportDashboard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ImportDashboard'
type FakeDashboardService_ImportDashboard_Call struct {
	*mock.Call
}

// ImportDashboard is a helper method to define mock.On call
//   - ctx context.Context
//   - dto *SaveDashboardDTO
func (_e *FakeDashboardService_Expecter) ImportDashboard(ctx interface{}, dto interface{}) *FakeDashboardService_ImportDashboard_Call {
	return &FakeDashboardService_ImportDashboard_Call{Call: _e.mock.On("ImportDashboard", ctx, dto)}
}

func (_c *FakeDashboardService_ImportDashboard_Call) Run(run func(ctx context.Context, dto *SaveDashboardDTO)) *FakeDashboardService_ImportDashboard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*SaveDashboardDTO))
	})
	return _c
}

func (_c *FakeDashboardService_ImportDashboard_Call) Return(_a0 *Dashboard, _a1 error) *FakeDashboardService_ImportDashboard_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_ImportDashboard_Call) RunAndReturn(run func(context.Context, *SaveDashboardDTO) (*Dashboard, error)) *FakeDashboardService_ImportDashboard_Call {
	_c.Call.Return(run)
	return _c
}

// MakeUserAdmin provides a mock function with given fields: ctx, orgID, userID, dashboardID, setViewAndEditPermissions
func (_m *FakeDashboardService) MakeUserAdmin(ctx context.Context, orgID int64, userID int64, dashboardID int64, setViewAndEditPermissions bool) error {
	ret := _m.Called(ctx, orgID, userID, dashboardID, setViewAndEditPermissions)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int64, bool) error); ok {
		r0 = rf(ctx, orgID, userID, dashboardID, setViewAndEditPermissions)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FakeDashboardService_MakeUserAdmin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeUserAdmin'
type FakeDashboardService_MakeUserAdmin_Call struct {
	*mock.Call
}

// MakeUserAdmin is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID int64
//   - userID int64
//   - dashboardID int64
//   - setViewAndEditPermissions bool
func (_e *FakeDashboardService_Expecter) MakeUserAdmin(ctx interface{}, orgID interface{}, userID interface{}, dashboardID interface{}, setViewAndEditPermissions interface{}) *FakeDashboardService_MakeUserAdmin_Call {
	return &FakeDashboardService_MakeUserAdmin_Call{Call: _e.mock.On("MakeUserAdmin", ctx, orgID, userID, dashboardID, setViewAndEditPermissions)}
}

func (_c *FakeDashboardService_MakeUserAdmin_Call) Run(run func(ctx context.Context, orgID int64, userID int64, dashboardID int64, setViewAndEditPermissions bool)) *FakeDashboardService_MakeUserAdmin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(int64), args[3].(int64), args[4].(bool))
	})
	return _c
}

func (_c *FakeDashboardService_MakeUserAdmin_Call) Return(_a0 error) *FakeDashboardService_MakeUserAdmin_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FakeDashboardService_MakeUserAdmin_Call) RunAndReturn(run func(context.Context, int64, int64, int64, bool) error) *FakeDashboardService_MakeUserAdmin_Call {
	_c.Call.Return(run)
	return _c
}

// SaveDashboard provides a mock function with given fields: ctx, dto, allowUiUpdate
func (_m *FakeDashboardService) SaveDashboard(ctx context.Context, dto *SaveDashboardDTO, allowUiUpdate bool) (*Dashboard, error) {
	ret := _m.Called(ctx, dto, allowUiUpdate)

	var r0 *Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *SaveDashboardDTO, bool) (*Dashboard, error)); ok {
		return rf(ctx, dto, allowUiUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *SaveDashboardDTO, bool) *Dashboard); ok {
		r0 = rf(ctx, dto, allowUiUpdate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *SaveDashboardDTO, bool) error); ok {
		r1 = rf(ctx, dto, allowUiUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_SaveDashboard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveDashboard'
type FakeDashboardService_SaveDashboard_Call struct {
	*mock.Call
}

// SaveDashboard is a helper method to define mock.On call
//   - ctx context.Context
//   - dto *SaveDashboardDTO
//   - allowUiUpdate bool
func (_e *FakeDashboardService_Expecter) SaveDashboard(ctx interface{}, dto interface{}, allowUiUpdate interface{}) *FakeDashboardService_SaveDashboard_Call {
	return &FakeDashboardService_SaveDashboard_Call{Call: _e.mock.On("SaveDashboard", ctx, dto, allowUiUpdate)}
}

func (_c *FakeDashboardService_SaveDashboard_Call) Run(run func(ctx context.Context, dto *SaveDashboardDTO, allowUiUpdate bool)) *FakeDashboardService_SaveDashboard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*SaveDashboardDTO), args[2].(bool))
	})
	return _c
}

func (_c *FakeDashboardService_SaveDashboard_Call) Return(_a0 *Dashboard, _a1 error) *FakeDashboardService_SaveDashboard_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_SaveDashboard_Call) RunAndReturn(run func(context.Context, *SaveDashboardDTO, bool) (*Dashboard, error)) *FakeDashboardService_SaveDashboard_Call {
	_c.Call.Return(run)
	return _c
}

// SearchDashboards provides a mock function with given fields: ctx, query
func (_m *FakeDashboardService) SearchDashboards(ctx context.Context, query *FindPersistedDashboardsQuery) (model.HitList, error) {
	ret := _m.Called(ctx, query)

	var r0 model.HitList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *FindPersistedDashboardsQuery) (model.HitList, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *FindPersistedDashboardsQuery) model.HitList); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.HitList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *FindPersistedDashboardsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FakeDashboardService_SearchDashboards_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchDashboards'
type FakeDashboardService_SearchDashboards_Call struct {
	*mock.Call
}

// SearchDashboards is a helper method to define mock.On call
//   - ctx context.Context
//   - query *FindPersistedDashboardsQuery
func (_e *FakeDashboardService_Expecter) SearchDashboards(ctx interface{}, query interface{}) *FakeDashboardService_SearchDashboards_Call {
	return &FakeDashboardService_SearchDashboards_Call{Call: _e.mock.On("SearchDashboards", ctx, query)}
}

func (_c *FakeDashboardService_SearchDashboards_Call) Run(run func(ctx context.Context, query *FindPersistedDashboardsQuery)) *FakeDashboardService_SearchDashboards_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*FindPersistedDashboardsQuery))
	})
	return _c
}

func (_c *FakeDashboardService_SearchDashboards_Call) Return(_a0 model.HitList, _a1 error) *FakeDashboardService_SearchDashboards_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FakeDashboardService_SearchDashboards_Call) RunAndReturn(run func(context.Context, *FindPersistedDashboardsQuery) (model.HitList, error)) *FakeDashboardService_SearchDashboards_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDashboardACL provides a mock function with given fields: ctx, uid, items
func (_m *FakeDashboardService) UpdateDashboardACL(ctx context.Context, uid int64, items []*DashboardACL) error {
	ret := _m.Called(ctx, uid, items)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, []*DashboardACL) error); ok {
		r0 = rf(ctx, uid, items)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FakeDashboardService_UpdateDashboardACL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDashboardACL'
type FakeDashboardService_UpdateDashboardACL_Call struct {
	*mock.Call
}

// UpdateDashboardACL is a helper method to define mock.On call
//   - ctx context.Context
//   - uid int64
//   - items []*DashboardACL
func (_e *FakeDashboardService_Expecter) UpdateDashboardACL(ctx interface{}, uid interface{}, items interface{}) *FakeDashboardService_UpdateDashboardACL_Call {
	return &FakeDashboardService_UpdateDashboardACL_Call{Call: _e.mock.On("UpdateDashboardACL", ctx, uid, items)}
}

func (_c *FakeDashboardService_UpdateDashboardACL_Call) Run(run func(ctx context.Context, uid int64, items []*DashboardACL)) *FakeDashboardService_UpdateDashboardACL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].([]*DashboardACL))
	})
	return _c
}

func (_c *FakeDashboardService_UpdateDashboardACL_Call) Return(_a0 error) *FakeDashboardService_UpdateDashboardACL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FakeDashboardService_UpdateDashboardACL_Call) RunAndReturn(run func(context.Context, int64, []*DashboardACL) error) *FakeDashboardService_UpdateDashboardACL_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewFakeDashboardService interface {
	mock.TestingT
	Cleanup(func())
}

// NewFakeDashboardService creates a new instance of FakeDashboardService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFakeDashboardService(t mockConstructorTestingTNewFakeDashboardService) *FakeDashboardService {
	mock := &FakeDashboardService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}