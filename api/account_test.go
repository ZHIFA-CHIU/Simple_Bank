package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	mockdb "github.com/zhifaq/simple_bank/db/mock"
	db "github.com/zhifaq/simple_bank/sqlc"
	"github.com/zhifaq/simple_bank/utils"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountId     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			"OK",
			account.ID,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},

		{
			"NotFound",
			account.ID,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			"InternalError",
			account.ID,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			"InvalidID",
			0,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func TestCreateAccount(t *testing.T) {
	testCases := []struct {
		name          string
		body          io.Reader
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			"OK",
			bytes.NewBufferString(`{"owner":"test_user","currency":"USD"}`),
			func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    "test_user",
					Currency: "USD",
					Balance:  0,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Account{
					ID:       1,
					Owner:    arg.Owner,
					Currency: arg.Currency,
					Balance:  arg.Balance,
				}, nil)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var account db.Account
				err := json.Unmarshal(recorder.Body.Bytes(), &account)
				require.NoError(t, err)
				require.Equal(t, "test_user", account.Owner)
				require.Equal(t, "USD", account.Currency)
				require.Equal(t, int64(0), account.Balance)
			},
		},
		{
			"InternalError",
			bytes.NewBufferString(`{"owner":"test_user","currency":"USD"}`),
			func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    "test_user",
					Currency: "USD",
					Balance:  0,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			"InvalidJSON",
			bytes.NewBufferString(`{"owner":"test_user","currency":123}`),
			func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/accounts", tc.body)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
