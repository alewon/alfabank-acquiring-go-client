package alfabank

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newTestClient(t *testing.T, fn roundTripFunc) *Client {
	t.Helper()
	client := NewClient("api-user", "api-pass", "api-token", &http.Client{Transport: fn})
	client.BaseURL = "https://test.local"
	return client
}

func TestRegisterUsesFormEncodingAndAppliesClientAuth(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
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
			t.Fatalf("parse query: %v", err)
		}

		if got := values.Get("userName"); got != "api-user" {
			t.Fatalf("unexpected userName: %s", got)
		}
		if got := values.Get("password"); got != "api-pass" {
			t.Fatalf("unexpected password: %s", got)
		}
		if got := values.Get("token"); got != "api-token" {
			t.Fatalf("unexpected token: %s", got)
		}
		if got := values.Get("orderNumber"); got != "order-1" {
			t.Fatalf("unexpected orderNumber: %s", got)
		}
		if got := values.Get("extra"); got != "value" {
			t.Fatalf("unexpected extra value: %s", got)
		}

		var bundle OrderBundle
		if err := json.Unmarshal([]byte(values.Get("orderBundle")), &bundle); err != nil {
			t.Fatalf("unmarshal orderBundle: %v", err)
		}
		if len(bundle.Loyalties) != 1 || bundle.Loyalties[0].LoyaltyProgramName != "bonus" {
			t.Fatalf("unexpected loyalties payload: %+v", bundle.Loyalties)
		}
		if len(bundle.CartItems.Items) != 1 {
			t.Fatalf("unexpected cart items count: %d", len(bundle.CartItems.Items))
		}
		item := bundle.CartItems.Items[0]
		if item.ItemDetails.ItemDetailsParams.Name != "fes_truCode" {
			t.Fatalf("unexpected itemDetailsParams: %+v", item.ItemDetails.ItemDetailsParams)
		}
		if len(item.ItemAttributes.Attributes) != 2 {
			t.Fatalf("unexpected item attributes: %+v", item.ItemAttributes.Attributes)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"errorCode":"0","formUrl":"https://pay.example/form","orderId":"order-id-1"}`)),
		}, nil
	})

	req := &RegisterRequest{
		OrderNumber: "order-1",
		Amount:      12345,
		ReturnURL:   "https://merchant.example/success",
		OrderBundle: OrderBundle{
			Loyalties: []Loyalty{
				{PositionID: 1, LoyaltyProgramName: "bonus", BonusAmountRefunded: "0"},
			},
			CartItems: CartItems{
				Items: []CartItem{
					{
						PositionID: 1,
						Name:       "item",
						Quantity:   Quantity{Value: "1", Measure: "0"},
						ItemCode:   "sku-1",
						ItemDetails: ItemDetails{
							ItemDetailsParams: ItemDetailsParams{Name: "fes_truCode", Value: "329921120.06002020100000000643"},
						},
						ItemAttributes: ItemAttributes{
							Attributes: []NameValuePair{
								{Name: "paymentMethod", Value: "1"},
								{Name: "paymentObject", Value: "1"},
							},
						},
					},
				},
			},
		},
		Values: url.Values{"extra": {"value"}},
	}

	result, err := client.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if result.Response.OrderID != "order-id-1" {
		t.Fatalf("unexpected order id: %s", result.Response.OrderID)
	}
}

func TestTemplatesCreateTemplateUsesJSONContentType(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content-type: %s", got)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		var req TemplatesCreateTemplateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal request: %v", err)
		}
		if req.UserName != "api-user" || req.Password != "api-pass" {
			t.Fatalf("auth was not applied: %+v", req)
		}
		if req.Type != "SBP_QR" {
			t.Fatalf("unexpected type: %s", req.Type)
		}
		if req.QRTemplate.QRHeight != "150" || req.QRTemplate.QRWidth != "150" {
			t.Fatalf("unexpected qrTemplate: %+v", req.QRTemplate)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"templateId":"tpl-1","name":"Template","type":"SBP_QR","status":"ACTIVE","qrTemplate":{"payload":"https://qr.example/tpl-1"}}`)),
		}, nil
	})

	result, err := client.TemplatesCreateTemplate(context.Background(), &TemplatesCreateTemplateRequest{
		Name:       "Template",
		Type:       "SBP_QR",
		QRTemplate: QRTemplate{QRHeight: "150", QRWidth: "150"},
	})
	if err != nil {
		t.Fatalf("create template: %v", err)
	}
	if result.Response.TemplateID != "tpl-1" {
		t.Fatalf("unexpected template id: %s", result.Response.TemplateID)
	}
}

func TestGetOrderStatusExtendedUnmarshalsExpandedResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/payment/rest/getOrderStatusExtended.do") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{
			"errorCode":"0",
			"orderNumber":"7005",
			"orderStatus":2,
			"actionCode":"0",
			"actionCodeDescription":"success",
			"amount":2000,
			"currency":"810",
			"date":1617972915659,
			"paymentWay":"CARD_BINDING",
			"displayErrorMessage":"visible error",
			"merchantOrderParams":[{"name":"originalResponseCode","value":"00"}],
			"transactionAttributes":[{"name":"subscription","value":"pro","subscriptionServiceName":"svc","subscriptionServiceId":"svc-1"}],
			"refunds":[{"date":"2025-01-01","externalRefundId":"refund-1","approvalCode":"ABC123","amount":500}],
			"bindingInfo":{"clientId":"client-1","bindingId":"binding-1","externalCreated":true},
			"paymentAmountInfo":{"paymentState":"DEPOSITED","approvedAmount":2000,"depositedAmount":2000,"refundedAmount":500,"totalAmount":2100},
			"bankInfo":{"bankCountryCode":"RU","bankCountryName":"Russia","bankName":"Alfa"},
			"cardAuthInfo":{
				"maskedPan":"411111**1111",
				"expiration":"203412",
				"cardholderName":"TEST CARDHOLDER",
				"approvalCode":"123456",
				"paymentSystem":"VISA",
				"product":"Corporate",
				"productCategory":"DEBIT",
				"secureAuthInfo":{"eci":5,"threeDSProtocolVersion":"2.2.0","aresTransStatus":"Y","rreqTransStatus":"Y"}
			}
		}`)),
		}, nil
	})

	result, err := client.GetOrderStatusExtended(context.Background(), &GetOrderStatusExtendedRequest{OrderID: "order-id-1"})
	if err != nil {
		t.Fatalf("get order status extended: %v", err)
	}

	resp := result.Response
	if resp.OrderStatus != 2 || resp.PaymentWay != "CARD_BINDING" {
		t.Fatalf("unexpected status payload: %+v", resp)
	}
	if len(resp.Refunds) != 1 || resp.Refunds[0].ExternalRefundID != "refund-1" {
		t.Fatalf("unexpected refunds: %+v", resp.Refunds)
	}
	if resp.BindingInfo.BindingID != "binding-1" || !resp.BindingInfo.ExternalCreated {
		t.Fatalf("unexpected bindingInfo: %+v", resp.BindingInfo)
	}
	if resp.CardAuthInfo.SecureAuthInfo.ThreeDSProtocol != "2.2.0" {
		t.Fatalf("unexpected secureAuthInfo: %+v", resp.CardAuthInfo.SecureAuthInfo)
	}
	if len(resp.TransactionAttributes) != 1 || resp.TransactionAttributes[0].SubscriptionServiceID != "svc-1" {
		t.Fatalf("unexpected transactionAttributes: %+v", resp.TransactionAttributes)
	}
}

func TestGetReceiptStatusUnmarshalsReceiptDetails(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{
			"errorCode":"0",
			"orderNumber":"220170606034051002_177",
			"orderId":"abd60d0c-e096-42c3-8b17-6081c67db214",
			"receipt":[
				{
					"receiptStatus":1,
					"uuid":"790925e5-739c-430c-9e92-79d9f14481a4",
					"shift_number":"27",
					"fiscal_receipt_number":"21",
					"receipt_date_time":1499256900000,
					"fn_number":"9999078900006364",
					"ofd_receipt_url":"https://ofd.example/receipt",
					"OFD":{"inn":"1234567890","name":"OFD Name","website":"ofd.example"},
					"ofdReceiptParams":{"ofdReceiptStatus":{"queuedDoc":2,"ffdVersion":"1.2","internetSign":true},"vatAmount":100},
					"ofdOrderBundle":[
						{
							"name":"water",
							"itemAmount":111165,
							"itemPrice":7411,
							"taxType":1,
							"quantity":{"value":"15","measure":"0"},
							"itemAttributes":[{"paymentMethod":"1","paymentObject":"1"}]
						}
					]
				}
			]
		}`)),
		}, nil
	})

	result, err := client.GetReceiptStatus(context.Background(), &GetReceiptStatusRequest{OrderID: "order-id-1"})
	if err != nil {
		t.Fatalf("get receipt status: %v", err)
	}

	resp := result.Response
	if len(resp.Receipt) != 1 {
		t.Fatalf("unexpected receipt count: %d", len(resp.Receipt))
	}
	receipt := resp.Receipt[0]
	if receipt.UUID != "790925e5-739c-430c-9e92-79d9f14481a4" {
		t.Fatalf("unexpected receipt uuid: %s", receipt.UUID)
	}
	if int64(receipt.ReceiptDateTimeLegacy) != 1499256900000 {
		t.Fatalf("unexpected legacy receipt datetime: %d", receipt.ReceiptDateTimeLegacy)
	}
	if receipt.OFD.Name != "OFD Name" || receipt.OFDReceiptParams.OFDReceiptStatus.FFDVersion != "1.2" {
		t.Fatalf("unexpected ofd payload: %+v %+v", receipt.OFD, receipt.OFDReceiptParams)
	}
	if len(receipt.OFDOrderBundle) != 1 || len(receipt.OFDOrderBundle[0].ItemAttributes) != 1 {
		t.Fatalf("unexpected ofdOrderBundle: %+v", receipt.OFDOrderBundle)
	}
}
