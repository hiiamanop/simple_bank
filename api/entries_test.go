package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/hiiamanop/simple_bank/db/mock"
	db "github.com/hiiamanop/simple_bank/db/sqlc"
	"github.com/hiiamanop/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	entry := randomEntry()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"account_id": entry.AccountID,
				"amount":     entry.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateEntriesParams{
					AccountID: entry.AccountID,
					Amount:    entry.Amount,
				}

				store.EXPECT().
					CreateEntries(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(entry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchEntry(t, recorder.Body, entry)
			},
		},
		{
			name: "InvalidAccountID",
			body: gin.H{
				"account_id": 0,
				"amount":     entry.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateEntries(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"account_id": entry.AccountID,
				"amount":     entry.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateEntries(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/entries"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetEntry(t *testing.T) {
	entry := randomEntry()

	testCases := []struct {
		name          string
		entryID       int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetEntries(gomock.Any(), gomock.Eq(entry.ID)).
					Times(1).
					Return(entry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchEntry(t, recorder.Body, entry)
			},
		},
		{
			name:    "NotFound",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetEntries(gomock.Any(), gomock.Eq(entry.ID)).
					Times(1).
					Return(db.Entry{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetEntries(gomock.Any(), gomock.Eq(entry.ID)).
					Times(1).
					Return(db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/entries/%d", tc.entryID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListEntries(t *testing.T) {
	n := 5
	entries := make([]db.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = randomEntry()
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListEntriesParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListEntries(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchEntries(t, recorder.Body, entries)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListEntries(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/api/v1/entries"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateEntry(t *testing.T) {
	entry := randomEntry()

	testCases := []struct {
		name          string
		entryID       int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			entryID: entry.ID,
			body: gin.H{
				"amount": entry.Amount + 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateEntriesParams{
					ID:     entry.ID,
					Amount: entry.Amount + 100,
				}

				store.EXPECT().
					UpdateEntries(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(entry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			entryID: entry.ID,
			body: gin.H{
				"amount": entry.Amount + 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateEntries(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Entry{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/entries/%d", tc.entryID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	entry := randomEntry()

	testCases := []struct {
		name          string
		entryID       int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteEntries(gomock.Any(), gomock.Eq(entry.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteEntries(gomock.Any(), gomock.Eq(entry.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/entries/%d", tc.entryID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomEntry() db.Entry {
	return db.Entry{
		ID:        int64(util.RandomInt(1, 1000)),
		AccountID: int64(util.RandomInt(1, 1000)),
		Amount:    int64(util.RandomMoney()),
	}
}

func requireBodyMatchEntry(t *testing.T, body *bytes.Buffer, entry db.Entry) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotEntry db.Entry
	err = json.Unmarshal(data, &gotEntry)
	require.NoError(t, err)
	require.Equal(t, entry, gotEntry)
}

func requireBodyMatchEntries(t *testing.T, body *bytes.Buffer, entries []db.Entry) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotEntries []db.Entry
	err = json.Unmarshal(data, &gotEntries)
	require.NoError(t, err)
	require.Equal(t, entries, gotEntries)
}
