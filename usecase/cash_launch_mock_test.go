package usecase_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	mock_repository "github.com/CharlesSchiavinato/minsait-challenge-backend/mock/repository"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	repository_in_memory "github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository/in_memory"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/stretchr/testify/assert"
)

func TestMockCashLaunchInsert(t *testing.T) {
	type test struct {
		name            string
		inputCashLaunch *model.CashLaunch
		wantCashLaunch  *model.CashLaunch
		wantError       error
		mockOn          func(mockRepositoryCashLaunch *mock_repository.MockCashLaunch, resultCashLaunch *model.CashLaunch, err error)
		assert          func(t *testing.T, tt *test, resultCashLaunch *model.CashLaunch, err error)
	}

	tests := []test{
		{
			name: "ReferenceDateEmptyError",
			inputCashLaunch: &model.CashLaunch{
				Type:        modelCashLaunchDefault.Type,
				Description: modelCashLaunchDefault.Description,
				Value:       modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageReferenceDateEmptyError},
		},
		{
			name: "ReferenceDateBeforeError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, -1),
				Type:          modelCashLaunchDefault.Type,
				Description:   modelCashLaunchDefault.Description,
				Value:         modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageReferenceDateBetweenError},
		},
		{
			name: "ReferenceDateAfterError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: usecase.CashLaunchReferenceDateMax.AddDate(0, 0, 1),
				Type:          modelCashLaunchDefault.Type,
				Description:   modelCashLaunchDefault.Description,
				Value:         modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageReferenceDateBetweenError},
		},
		{
			name: "TypeEmptyError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: modelCashLaunchDefault.ReferenceDate,
				Type:          "",
				Description:   modelCashLaunchDefault.Description,
				Value:         modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageTypeEmptyError},
		},
		{
			name: "TypeInvalidError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: modelCashLaunchDefault.ReferenceDate,
				Type:          "CC",
				Description:   modelCashLaunchDefault.Description,
				Value:         modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageTypeInvalidError},
		},
		{
			name: "DescriptionEmptyError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: modelCashLaunchDefault.ReferenceDate,
				Type:          modelCashLaunchDefault.Type,
				Value:         modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageDescriptionEmptyError},
		},
		{
			name: "DescriptionSizeLessError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: modelCashLaunchDefault.ReferenceDate,
				Type:          modelCashLaunchDefault.Type,
				Description:   fmt.Sprintf("%v ", modelCashLaunchDefault.Description[0:usecase.CashLaunchDescriptionMinLen-1]),
				Value:         modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageDescriptionSizeError},
		},
		{
			name: "DescriptionSizeGreaterError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: modelCashLaunchDefault.ReferenceDate,
				Type:          modelCashLaunchDefault.Type,
				Description:   fmt.Sprintf("%v ", strings.Repeat(modelCashLaunchDefault.Description, 10)[0:usecase.CashLaunchDescriptionMaxLen+1]),
				Value:         modelCashLaunchDefault.Value,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageDescriptionSizeError},
		},
		{
			name: "ValueZeroError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: modelCashLaunchDefault.ReferenceDate,
				Type:          modelCashLaunchDefault.Type,
				Description:   modelCashLaunchDefault.Description,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageValueError},
		},
		{
			name: "ValueLessZeroError",
			inputCashLaunch: &model.CashLaunch{
				ReferenceDate: modelCashLaunchDefault.ReferenceDate,
				Type:          modelCashLaunchDefault.Type,
				Description:   modelCashLaunchDefault.Description,
				Value:         -0.01,
			},
			wantError: usecase.ErrModelValidate{Message: usecase.CashLaunchMessageValueError},
		},
		{
			name:            "Success",
			inputCashLaunch: modelCashLaunchDefault,
			mockOn: func(mockRepositoryCashLaunch *mock_repository.MockCashLaunch, resultCashLaunch *model.CashLaunch, err error) {
				mockRepositoryCashLaunch.On("Insert").Return(resultCashLaunch, err)
			},
			assert: func(t *testing.T, tt *test, resultCashLaunch *model.CashLaunch, err error) {
				if !reflect.DeepEqual(err, tt.wantError) {
					t.Errorf("Insert() got error = %v, want %v", err, tt.wantError)
				}

				assert.NotNil(t, resultCashLaunch)
				assert.NotEqual(t, tt.inputCashLaunch.ID, resultCashLaunch.ID)
				assert.Equal(t, tt.inputCashLaunch.ReferenceDate, resultCashLaunch.ReferenceDate)
				assert.Equal(t, tt.inputCashLaunch.Type, resultCashLaunch.Type)
				assert.Equal(t, tt.inputCashLaunch.Description, resultCashLaunch.Description)
				assert.Equal(t, tt.inputCashLaunch.Value, resultCashLaunch.Value)
				assert.NotEqual(t, tt.inputCashLaunch.UpdatedAt, resultCashLaunch.UpdatedAt)
				assert.NotEqual(t, tt.inputCashLaunch.CreatedAt, resultCashLaunch.CreatedAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepositoryCashLaunch := new(mock_repository.MockCashLaunch)

			if tt.mockOn != nil {
				tt.mockOn(mockRepositoryCashLaunch, tt.inputCashLaunch, tt.wantError)
			}

			usecaseCashLaunch := usecase.NewCashLaunch(mockRepositoryCashLaunch)

			modelCashLaunch := *tt.inputCashLaunch

			resultCashLaunch, err := usecaseCashLaunch.Insert(&modelCashLaunch)

			if tt.assert != nil {
				tt.assert(t, &tt, resultCashLaunch, err)
			} else {
				if !reflect.DeepEqual(err, tt.wantError) {
					t.Errorf("Insert() got error = %v, want = %v.", err, tt.wantError)
				}

				if !reflect.DeepEqual(resultCashLaunch, tt.wantCashLaunch) {
					t.Errorf("Insert() got result = %v, want = %v.", resultCashLaunch, tt.wantCashLaunch)
				}
			}
		})
	}
}

func TestMockCashLaunchList(t *testing.T) {
	type test struct {
		name             string
		wantCashLaunches model.CashLaunches
		wantError        error
		mockOn           func(mockRepositoryCashLaunch *mock_repository.MockCashLaunch, resultCashLaunches model.CashLaunches, err error)
		assert           func(t *testing.T, tt *test, resultCashLaunches model.CashLaunches, err error)
	}

	tests := []test{
		{
			name:             "Success",
			wantCashLaunches: repository_in_memory.InMemoryCashLaunches,
			mockOn: func(mockRepositoryCashLaunch *mock_repository.MockCashLaunch, resultCashLaunches model.CashLaunches, err error) {
				mockRepositoryCashLaunch.On("List").Return(resultCashLaunches, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepositoryCashLaunch := new(mock_repository.MockCashLaunch)

			if tt.mockOn != nil {
				tt.mockOn(mockRepositoryCashLaunch, tt.wantCashLaunches, tt.wantError)
			}

			usecaseCashLaunch := usecase.NewCashLaunch(mockRepositoryCashLaunch)

			resultCashLaunches, err := usecaseCashLaunch.List()

			if tt.assert != nil {
				tt.assert(t, &tt, resultCashLaunches, err)
			} else {
				if !reflect.DeepEqual(err, tt.wantError) {
					t.Errorf("List() got error = %v, want = %v.", err, tt.wantError)
				}

				if !reflect.DeepEqual(resultCashLaunches, tt.wantCashLaunches) {
					t.Errorf("List() got result = %v, want = %v.", resultCashLaunches, &tt.wantCashLaunches)
				}
			}
		})
	}
}

func TestMockCashLaunchGetByID(t *testing.T) {
	type test struct {
		name           string
		inputID        int64
		wantCashLaunch *model.CashLaunch
		wantError      error
		mockOn         func(mockRepositoryCashLaunch *mock_repository.MockCashLaunch, resultCashLaunch *model.CashLaunch, err error)
	}

	tests := []test{
		{
			name:      "NotFoundError",
			wantError: repository.ErrNotFound{Message: "not found"},
			mockOn: func(mockRepositoryCashLaunch *mock_repository.MockCashLaunch, resultCashLaunch *model.CashLaunch, err error) {
				mockRepositoryCashLaunch.On("GetByID").Return(resultCashLaunch, err)
			},
		},
		{
			name:           "Success",
			inputID:        1,
			wantCashLaunch: &repository_in_memory.InMemoryCashLaunches[0],
			mockOn: func(mockRepositoryCashLaunch *mock_repository.MockCashLaunch, resultCashLaunch *model.CashLaunch, err error) {
				mockRepositoryCashLaunch.On("GetByID").Return(resultCashLaunch, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepositoryCashLaunch := new(mock_repository.MockCashLaunch)

			if tt.mockOn != nil {
				tt.mockOn(mockRepositoryCashLaunch, tt.wantCashLaunch, tt.wantError)
			}

			usecaseCashLaunch := usecase.NewCashLaunch(mockRepositoryCashLaunch)

			resultCashLaunches, err := usecaseCashLaunch.GetByID(tt.inputID)

			if !reflect.DeepEqual(err, tt.wantError) {
				t.Errorf("GetByID() got error = %v, want = %v.", err, tt.wantError)
			}

			if !reflect.DeepEqual(resultCashLaunches, tt.wantCashLaunch) {
				t.Errorf("GetByID() got result = %v, want = %v.", resultCashLaunches, &tt.wantCashLaunch)
			}
		})
	}
}
