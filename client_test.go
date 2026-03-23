package alfabank

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClientFormRequests(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		call         func(context.Context, *Client) error
		assertValues func(*testing.T, url.Values)
	}{
		{
			name: "register",
			path: "/register.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Register(ctx, RegisterRequest{
					OrderNumber: "ord-1",
					Amount:      1200,
					ReturnURL:   "https://example.com/ok",
					Features:    []OrderFeature{OrderFeatureVerify, OrderFeatureForceTDS},
					OrderBundle: &OrderBundle{CustomerDetails: &CustomerDetails{Email: "user@example.com"}},
				})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if values.Get("orderNumber") != "ord-1" {
					t.Fatalf("unexpected orderNumber: %q", values.Get("orderNumber"))
				}
				if values.Get("amount") != "1200" {
					t.Fatalf("unexpected amount: %q", values.Get("amount"))
				}
				if got := values["features"]; len(got) != 2 || got[0] != "VERIFY" || got[1] != "FORCE_TDS" {
					t.Fatalf("unexpected features: %#v", got)
				}
				if !strings.Contains(values.Get("orderBundle"), `"email":"user@example.com"`) {
					t.Fatalf("orderBundle was not json encoded: %q", values.Get("orderBundle"))
				}
			},
		},
		{
			name: "register preauth",
			path: "/registerPreAuth.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.RegisterPreAuth(ctx, RegisterPreAuthRequest{
					OrderNumber: "ord-2",
					Amount:      2200,
					ReturnURL:   "https://example.com/ok",
				})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if values.Get("orderNumber") != "ord-2" {
					t.Fatalf("unexpected orderNumber: %q", values.Get("orderNumber"))
				}
			},
		},
		{
			name: "status",
			path: "/getOrderStatusExtended.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.GetOrderStatusExtended(ctx, GetOrderStatusExtendedRequest{OrderID: "md-1", Language: "ru"})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if values.Get("orderId") != "md-1" {
					t.Fatalf("unexpected orderId: %q", values.Get("orderId"))
				}
			},
		},
		{
			name: "deposit",
			path: "/deposit.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Deposit(ctx, DepositRequest{
					OrderID: "md-2",
					Amount:  "500",
					DepositItems: &ItemsContainer{
						Items: []Item{{PositionID: 1, Name: "item"}},
					},
				})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if values.Get("amount") != "500" {
					t.Fatalf("unexpected amount: %q", values.Get("amount"))
				}
				if !strings.Contains(values.Get("depositItems"), `"positionId":1`) {
					t.Fatalf("depositItems was not json encoded: %q", values.Get("depositItems"))
				}
			},
		},
		{
			name: "reverse",
			path: "/reverse.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Reverse(ctx, ReverseRequest{OrderID: "md-3", Amount: "300"})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if values.Get("amount") != "300" {
					t.Fatalf("unexpected amount: %q", values.Get("amount"))
				}
			},
		},
		{
			name: "refund",
			path: "/refund.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Refund(ctx, RefundRequest{
					OrderID:    "md-4",
					Amount:     "700",
					JSONParams: map[string]string{"penalty": "1300"},
				})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if !strings.Contains(values.Get("jsonParams"), `"penalty":"1300"`) {
					t.Fatalf("jsonParams was not json encoded: %q", values.Get("jsonParams"))
				}
			},
		},
		{
			name: "decline",
			path: "/decline.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Decline(ctx, DeclineRequest{OrderID: "md-7", OrderNumber: "merchant-ord-1", Language: "ru"})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if values.Get("orderNumber") != "merchant-ord-1" {
					t.Fatalf("unexpected orderNumber: %q", values.Get("orderNumber"))
				}
			},
		},
		{
			name: "raw sum refund",
			path: "/processRawSumRefund.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ProcessRawSumRefund(ctx, ProcessRawSumRefundRequest{
					OrderID:  "md-5",
					Language: "ru",
					Amount:   "900",
					OFDParams: &OFDParams{Name: "ofd"},
				})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if !strings.Contains(values.Get("ofdParams"), `"name":"ofd"`) {
					t.Fatalf("ofdParams was not json encoded: %q", values.Get("ofdParams"))
				}
			},
		},
		{
			name: "raw position refund",
			path: "/processRawPositionRefund.do",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ProcessRawPositionRefund(ctx, ProcessRawPositionRefundRequest{
					OrderID:    "md-6",
					Language:   "ru",
					Amount:     "1100",
					PositionID: 3,
				})
				return err
			},
			assertValues: func(t *testing.T, values url.Values) {
				t.Helper()
				if values.Get("positionId") != "3" {
					t.Fatalf("unexpected positionId: %q", values.Get("positionId"))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				if r.Method != http.MethodPost {
					t.Fatalf("unexpected method: %s", r.Method)
				}
				if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
					t.Fatalf("unexpected content-type: %s", got)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}

				values, err := url.ParseQuery(string(body))
				if err != nil {
					t.Fatalf("parse body: %v", err)
				}

				if values.Get("userName") != "demo-user" || values.Get("password") != "demo-pass" {
					t.Fatalf("auth was not injected: %#v", values)
				}

				tc.assertValues(t, values)

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"errorCode":"0","errorMessage":"Success","orderId":"ok","formUrl":"https://example.com/form","orderStatus":2}`))
			}))
			defer srv.Close()

			client := NewClient(Config{BaseURL: srv.URL, UserName: "demo-user", Password: "demo-pass"})
			if err := tc.call(context.Background(), client); err != nil {
				t.Fatalf("call failed: %v", err)
			}
		})
	}
}

func TestGetOrderStatusExtendedResponseDecoding(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"errorCode": 0,
			"orderNumber": "merchant-order-1",
			"orderStatus": 2,
			"actionCode": -100,
			"actionCodeDescription": "Requested action has been successfully completed",
			"amount": 2000,
			"currency": 810,
			"date": 1740392720718,
			"depositedDate": 1740392720718,
			"orderDescription": "test order",
			"ip": "127.0.0.1",
			"paymentWay": "CARD",
			"chargeback": false,
			"authDateTime": 1740392720718,
			"terminalId": "123456",
			"attributes": [{"name":"mdOrder","value":"md-order"}],
			"paymentAmountInfo": {
				"paymentState": "DEPOSITED",
				"approvedAmount": 2000,
				"depositedAmount": 2000,
				"refundedAmount": 0
			},
			"refunds": {
				"date": "2025-02-24T10:25:20",
				"externalRefundId": "rf-1",
				"approvalCode": "ABC123",
				"actionCode": 0,
				"referenceNumber": "123456789012",
				"amount": 300
			},
			"cardAuthInfo": {
				"maskedPan": "411111**1111",
				"expiration": "203412",
				"cardholderName": "TEST CARDHOLDER",
				"approvalCode": "123456",
				"paymentSystem": "VISA",
				"product": "",
				"productCategory": "DEBIT",
				"corporateCard": "false",
				"secureAuthInfo": {
					"eci": "05",
					"authTypeIndicator": "1",
					"threeDSProtocolVersion": "2.2.0",
					"paResStatus": "Y"
				}
			},
			"bindingInfo": {
				"clientId": "client-1",
				"bindingId": "binding-1",
				"externalCreated": true
			},
			"pluginInfo": {
				"name": "plugin-name",
				"params": {"p1":"v1"}
			},
			"orderBundle": {
				"customerDetails": {"email": "buyer@example.com"},
				"cartItems": {
					"items": [{"positionId": "1", "name": "Item", "itemAmount": 2000, "itemPrice": 2000, "depositedItemAmount": 2000}]
				}
			}
		}`))
	}))
	defer srv.Close()

	client := NewClient(Config{BaseURL: srv.URL, UserName: "demo-user", Password: "demo-pass"})
	resp, err := client.GetOrderStatusExtended(context.Background(), GetOrderStatusExtendedRequest{OrderID: "md-order"})
	if err != nil {
		t.Fatalf("status call failed: %v", err)
	}
	if resp.ErrorCode != "0" {
		t.Fatalf("unexpected errorCode: %q", resp.ErrorCode)
	}
	if resp.ActionCode != "-100" {
		t.Fatalf("unexpected actionCode: %q", resp.ActionCode)
	}
	if resp.Currency != "810" {
		t.Fatalf("unexpected currency: %q", resp.Currency)
	}
	if resp.OrderStatus != OrderStatusDeposited {
		t.Fatalf("unexpected orderStatus: %d", resp.OrderStatus)
	}
	if resp.PaymentAmountInfo == nil || resp.PaymentAmountInfo.PaymentState != PaymentStateDeposited {
		t.Fatalf("paymentAmountInfo was not decoded: %#v", resp.PaymentAmountInfo)
	}
	if len(resp.Refunds) != 1 || resp.Refunds[0].ExternalRefundID != "rf-1" || resp.Refunds[0].ActionCode != "0" {
		t.Fatalf("refunds were not decoded: %#v", resp.Refunds)
	}
	if resp.CardAuthInfo == nil || resp.CardAuthInfo.SecureAuthInfo == nil || resp.CardAuthInfo.SecureAuthInfo.PaResStatus != "Y" {
		t.Fatalf("cardAuthInfo was not decoded: %#v", resp.CardAuthInfo)
	}
	if len(resp.Attributes) != 1 || resp.Attributes[0].Value != "md-order" {
		t.Fatalf("attributes were not decoded: %#v", resp.Attributes)
	}
	if resp.OrderBundle == nil || resp.OrderBundle.CartItems == nil || len(resp.OrderBundle.CartItems.Items) != 1 {
		t.Fatalf("orderBundle was not decoded: %#v", resp.OrderBundle)
	}
	if resp.OrderBundle.CartItems.Items[0].PositionID != 1 {
		t.Fatalf("positionId was not decoded: %#v", resp.OrderBundle.CartItems.Items[0])
	}
	if resp.OrderBundle.CartItems.Items[0].DepositedItemAmount != "2000" {
		t.Fatalf("depositedItemAmount was not decoded: %#v", resp.OrderBundle.CartItems.Items[0])
	}
}

func TestClientHTTPAndDecodeErrors(t *testing.T) {
	t.Run("http error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadGateway) }))
		defer srv.Close()

		client := NewClient(Config{BaseURL: srv.URL, UserName: "u", Password: "p"})
		if _, err := client.Register(context.Background(), RegisterRequest{OrderNumber: "1", Amount: 1, ReturnURL: "https://example.com"}); err == nil {
			t.Fatal("expected HTTP error")
		}
	})

	t.Run("decode error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{`))
		}))
		defer srv.Close()

		client := NewClient(Config{BaseURL: srv.URL, UserName: "u", Password: "p"})
		if _, err := client.Register(context.Background(), RegisterRequest{OrderNumber: "1", Amount: 1, ReturnURL: "https://example.com"}); err == nil {
			t.Fatal("expected decode error")
		}
	})
}