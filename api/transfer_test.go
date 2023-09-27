package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mockdb "github.com/mateusribs/simple_bank/db/mock"
	db "github.com/mateusribs/simple_bank/db/sqlc"
	"github.com/mateusribs/simple_bank/token"
	"github.com/mateusribs/simple_bank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)


func TestCreateTransfer(t *testing.T){
	
	amount := int64(100)
	invalidID := int64(5430)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)
	account3 := randomAccount(user3.Username)

	account1.Currency = "BRL"
	account2.Currency = "BRL"
	account3.Currency = "EUR"

	testCases := []struct{
		name string
		body gin.H
		setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: amount,
				}

				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnauthorizedTransfer",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user2.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0	)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			body: gin.H{
				"from_account_id": 0,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": "MEX", 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": -amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFoundAccountID",
			body: gin.H{
				"from_account_id": invalidID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(invalidID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "MismatchCurrency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account3.ID,
				"amount": amount,
				"currency": util.BRL, 
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			//build stubs
			tc.buildStubs(store)

			//start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// get authorization
			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			//check responses
			tc.checkResponse(t, recorder)
		})
	}
}