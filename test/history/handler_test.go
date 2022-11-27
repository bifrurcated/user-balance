package history

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bifrurcated/user-balance/internal/history"
	"github.com/bifrurcated/user-balance/test/testdata"
	"github.com/smartystreets/goconvey/convey"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	server := testdata.GetTestServer(testdata.History)
	code := m.Run()
	err := testdata.ExecuteSQLScript(context.TODO(), server.Store, "../../test/testdata/drop.sql")
	if err != nil {
		panic(err)
	}
	defer server.Test.Close()
	os.Exit(code)
}

func TestGetUserTransactions(t *testing.T) {
	server := testdata.GetTestServer(testdata.History)
	convey.Convey("Test API Get user transaction", t, func() {
		q := `INSERT INTO history_operations (sender_user_id, user_id, amount, type) VALUES  (1,2,200,'перевод')`
		_, err := server.Store.Exec(context.TODO(), q)
		if err != nil {
			t.Fatal(err)
		}
		historyDTO := history.UserHistoryDTO{UserID: 1}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(historyDTO)
		if err != nil {
			t.Fatal(err)
		}
		r, err := http.NewRequest(http.MethodGet, server.Test.URL+"/api/v1/history", &buf)
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		var histories history.UserHistoriesDTO
		err = json.NewDecoder(response.Body).Decode(&histories)
		if err != nil {
			t.Fatal(err)
		}
		convey.So(histories.NextPage, convey.ShouldBeNil)
		convey.So(len(histories.Histories), convey.ShouldEqual, 1)
		convey.So(*histories.Histories[0].SenderUserID, convey.ShouldEqual, 1)
		convey.So(histories.Histories[0].UserID, convey.ShouldEqual, 2)
		convey.So(histories.Histories[0].ServiceID, convey.ShouldBeNil)
		convey.So(histories.Histories[0].Amount, convey.ShouldEqual, 200)
		convey.So(histories.Histories[0].Type, convey.ShouldEqual, "перевод")

		convey.Reset(func() {
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM history_operations WHERE id=$1`, histories.Histories[0].ID)
			if err != nil {
				t.Fatal(err)
			}
			response.Body.Close()
		})
	})
}
