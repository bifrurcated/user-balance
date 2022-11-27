package balance

import (
	"bytes"
	"context"
	"encoding/json"
	balance2 "github.com/bifrurcated/user-balance/internal/balance"
	"github.com/bifrurcated/user-balance/test/testdata"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	server := testdata.GetTestServer()
	code := m.Run()
	err := testdata.ExecuteSQLScript(context.TODO(), server.Store, "../../test/testdata/drop.sql")
	if err != nil {
		panic(err)
	}
	defer server.Test.Close()
	os.Exit(code)
}

func TestAddMoney(t *testing.T) {
	server := testdata.GetTestServer()
	convey.Convey("Test API POST add money to user", t, func() {
		userBalanceDTO := balance2.CreateUserBalanceDTO{
			UserID: 1,
			Amount: 1000,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(userBalanceDTO)
		if err != nil {
			t.Fatal(err)
		}
		r, err := http.NewRequest(http.MethodPost, server.Test.URL+"/api/v1/add-money", &buf)
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		convey.So(response.StatusCode, convey.ShouldEqual, 204)
		convey.Reset(func() {
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM balance WHERE user_id=$1`, userBalanceDTO.UserID)
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}

func TestGetBalance(t *testing.T) {
	server := testdata.GetTestServer()
	convey.Convey("Test API Get user balance", t, func() {
		_, err := server.Store.Exec(context.TODO(), `INSERT INTO balance (user_id, amount) VALUES (1,100)`)
		if err != nil {
			t.Fatal(err)
		}
		userBalance := balance2.UserBalance{UserID: 1}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(userBalance)
		if err != nil {
			log.Fatal(err)
		}
		r, err := http.NewRequest(http.MethodGet, server.Test.URL+"/api/v1/balance", &buf)
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		var amount float32
		err = json.NewDecoder(response.Body).Decode(&amount)
		if err != nil {
			t.Fatal(err)
		}
		convey.So(amount, convey.ShouldEqual, 100)
		convey.Reset(func() {
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM balance WHERE user_id=1`)
			if err != nil {
				t.Fatal(err)
			}
			response.Body.Close()
		})
	})
}
