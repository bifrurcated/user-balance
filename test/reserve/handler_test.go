package reserve

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bifrurcated/user-balance/internal/reserve"
	"github.com/bifrurcated/user-balance/test/testdata"
	"github.com/smartystreets/goconvey/convey"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	server := testdata.GetTestServer(testdata.Reserve)
	code := m.Run()
	err := testdata.ExecuteSQLScript(context.TODO(), server.Store, "../../test/testdata/drop.sql")
	if err != nil {
		panic(err)
	}
	defer server.Test.Close()
	os.Exit(code)
}

func TestReserveMoney(t *testing.T) {
	server := testdata.GetTestServer(testdata.Reserve)
	convey.Convey("Test API POST reserve money", t, func() {
		_, err := server.Store.Exec(context.TODO(), `INSERT INTO balance (user_id, amount) VALUES (1, 200)`)
		if err != nil {
			t.Fatal(err)
		}
		reserveDTO := reserve.CreateReserveDTO{
			UserID:    1,
			ServiceID: 1,
			OrderID:   1,
			Cost:      100,
			IsProfit:  false,
		}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(reserveDTO)
		if err != nil {
			t.Fatal(err)
		}
		r, err := http.NewRequest(http.MethodPost, server.Test.URL+"/api/v1/reserve", &buf)
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		var reserveMoney reserve.CreateReserveMoneyDTO
		err = json.NewDecoder(response.Body).Decode(&reserveMoney)
		if err != nil {
			t.Log(response.StatusCode)
			t.Fatal(err)
		}

		convey.So(reserveMoney, convey.ShouldResemble, reserve.CreateReserveMoneyDTO{
			UserID:    reserveDTO.UserID,
			ServiceID: reserveDTO.ServiceID,
			OrderID:   reserveDTO.OrderID,
			Cost:      reserveDTO.Cost,
		})
		convey.Reset(func() {
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM balance WHERE user_id=$1`, reserveDTO.UserID)
			if err != nil {
				t.Fatal(err)
			}
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM reserve WHERE user_id=$1 AND service_id=$2 AND order_id=$3`,
				reserveDTO.UserID, reserveDTO.ServiceID, reserveDTO.OrderID)
			if err != nil {
				t.Fatal(err)
			}
			response.Body.Close()
		})
	})
}

func TestReserveProfit(t *testing.T) {
	server := testdata.GetTestServer(testdata.Reserve)
	convey.Convey("Test API POST reserve profit", t, func() {
		_, err := server.Store.Exec(context.TODO(), `INSERT INTO reserve (user_id, service_id, order_id, cost) VALUES (1,1,1,100)`)
		if err != nil {
			t.Fatal(err)
		}
		profit := reserve.ProfitReserveDTO{
			UserID:    1,
			ServiceID: 1,
			OrderID:   1,
			Amount:    200,
		}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(profit)
		if err != nil {
			t.Fatal(err)
		}
		r, err := http.NewRequest(http.MethodPost, server.Test.URL+"/api/v1/reserve/profit", &buf)
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		var respProfit reserve.ProfitReserveDTO
		err = json.NewDecoder(response.Body).Decode(&respProfit)
		if err != nil {
			t.Log(response.StatusCode)
			t.Fatal(err)
		}

		convey.So(respProfit.UserID, convey.ShouldEqual, 1)
		convey.So(respProfit.ServiceID, convey.ShouldEqual, 1)
		convey.So(respProfit.OrderID, convey.ShouldEqual, 1)
		convey.So(respProfit.Cost, convey.ShouldEqual, 100)
		convey.So(respProfit.Amount, convey.ShouldEqual, 200)

		convey.Reset(func() {
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM reserve WHERE user_id=$1 AND service_id=$2 AND order_id=$3`,
				profit.UserID, profit.ServiceID, profit.OrderID)
			if err != nil {
				t.Fatal(err)
			}
			response.Body.Close()
		})
	})
}

func TestCancelReserve(t *testing.T) {
	server := testdata.GetTestServer(testdata.Reserve)
	convey.Convey("Test API POST cancel reserve", t, func() {
		_, err := server.Store.Exec(context.TODO(), `INSERT INTO balance (user_id) VALUES (1)`)
		if err != nil {
			t.Fatal(err)
		}
		_, err = server.Store.Exec(context.TODO(), `INSERT INTO reserve (user_id, service_id, order_id, cost) VALUES (1,1,1,100)`)
		if err != nil {
			t.Fatal(err)
		}
		reserveDTO := reserve.CancelReserveDTO{
			UserID:    1,
			ServiceID: 1,
			OrderID:   1,
		}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(reserveDTO)
		if err != nil {
			t.Fatal(err)
		}
		r, err := http.NewRequest(http.MethodPost, server.Test.URL+"/api/v1/reserve/cancel", &buf)
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		convey.So(response.StatusCode, convey.ShouldEqual, 204)
		convey.Reset(func() {
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM balance WHERE user_id=$1`, reserveDTO.UserID)
			if err != nil {
				t.Fatal(err)
			}
			_, err = server.Store.Exec(context.TODO(), `DELETE FROM reserve WHERE user_id=$1 AND service_id=$2 AND order_id=$3`,
				reserveDTO.UserID, reserveDTO.ServiceID, reserveDTO.OrderID)
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
