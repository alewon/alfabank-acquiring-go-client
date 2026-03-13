package alfabank

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const DefaultBaseURL = "https://alfa.rbsuat.com"

type Client struct {
	BaseURL    string
	Username   string
	Password   string
	Token      string
	HTTPClient *http.Client
	Headers    http.Header
}

func NewClient(username, password, token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		BaseURL:    DefaultBaseURL,
		Username:   username,
		Password:   password,
		Token:      token,
		HTTPClient: httpClient,
		Headers:    make(http.Header),
	}
}

func (c *Client) baseURL() string {
	if c.BaseURL == "" {
		return DefaultBaseURL
	}
	return c.BaseURL
}

func (c *Client) newRequest(ctx context.Context, path string, request interface{}) (*http.Request, error) {
	values := encodeForm(request)
	endpoint := strings.TrimRight(c.baseURL(), "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	for key, items := range c.Headers {
		for _, item := range items {
			req.Header.Add(key, item)
		}
	}
	return req, nil
}

func (c *Client) newJSONRequest(ctx context.Context, path string, request interface{}) (*http.Request, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	endpoint := strings.TrimRight(c.baseURL(), "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for key, items := range c.Headers {
		for _, item := range items {
			req.Header.Add(key, item)
		}
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out interface{}) (int, http.Header, []byte, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, resp.Header.Clone(), nil, err
	}
	if len(body) > 0 && out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return resp.StatusCode, resp.Header.Clone(), body, err
		}
	}
	return resp.StatusCode, resp.Header.Clone(), body, nil
}

func encodeForm(input interface{}) url.Values {
	values := make(url.Values)
	if input == nil {
		return values
	}
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return values
		}
		v = v.Elem()
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return values
	}
	appendFormValues(values, v)
	return values
}

func appendFormValues(values url.Values, v reflect.Value) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fv := v.Field(i)
		if field.Name == "Values" && field.Type == reflect.TypeOf(url.Values{}) {
			if fv.IsZero() {
				continue
			}
			for _, key := range fv.MapKeys() {
				items := fv.MapIndex(key).Interface().([]string)
				for _, item := range items {
					values.Add(key.String(), item)
				}
			}
			continue
		}
		name := field.Tag.Get("form")
		if name == "" || name == "-" {
			continue
		}
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}
		if !fv.IsValid() || fv.IsZero() {
			continue
		}
		switch fv.Kind() {
		case reflect.String:
			values.Add(name, fv.String())
		case reflect.Bool:
			values.Add(name, strconv.FormatBool(fv.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values.Add(name, strconv.FormatInt(fv.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			values.Add(name, strconv.FormatUint(fv.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			values.Add(name, strconv.FormatFloat(fv.Float(), 'f', -1, 64))
		default:
			payload, err := json.Marshal(fv.Interface())
			if err == nil {
				values.Add(name, string(payload))
			}
		}
	}
}

func applyClientAuth(c *Client, request interface{}) {
	if request == nil {
		return
	}
	v := reflect.ValueOf(request)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return
	}
	v = v.Elem()
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return
	}
	set := func(name, value string) {
		f := v.FieldByName(name)
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.String && f.String() == "" {
			f.SetString(value)
		}
	}
	set("UserName", c.Username)
	set("Password", c.Password)
	set("Token", c.Token)
}

func (c *Client) Register(ctx context.Context, request *RegisterRequest) (*RegisterResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/register.do", request)
	if err != nil {
		return nil, err
	}
	response := &RegisterResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &RegisterResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) RegisterPreAuth(ctx context.Context, request *RegisterPreAuthRequest) (*RegisterPreAuthResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/registerPreAuth.do", request)
	if err != nil {
		return nil, err
	}
	response := &RegisterPreAuthResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &RegisterPreAuthResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) PaymentOrder(ctx context.Context, request *PaymentOrderRequest) (*PaymentOrderResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/paymentorder.do", request)
	if err != nil {
		return nil, err
	}
	response := &PaymentOrderResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &PaymentOrderResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) InstantPayment(ctx context.Context, request *InstantPaymentRequest) (*InstantPaymentResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/instantPayment.do", request)
	if err != nil {
		return nil, err
	}
	response := &InstantPaymentResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &InstantPaymentResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) GetOrderStatusExtended(ctx context.Context, request *GetOrderStatusExtendedRequest) (*GetOrderStatusExtendedResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/getOrderStatusExtended.do", request)
	if err != nil {
		return nil, err
	}
	response := &GetOrderStatusExtendedResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &GetOrderStatusExtendedResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) Deposit(ctx context.Context, request *DepositRequest) (*DepositResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/deposit.do", request)
	if err != nil {
		return nil, err
	}
	response := &DepositResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &DepositResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) Reverse(ctx context.Context, request *ReverseRequest) (*ReverseResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/reverse.do", request)
	if err != nil {
		return nil, err
	}
	response := &ReverseResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &ReverseResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) Refund(ctx context.Context, request *RefundRequest) (*RefundResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/refund.do", request)
	if err != nil {
		return nil, err
	}
	response := &RefundResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &RefundResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) Decline(ctx context.Context, request *DeclineRequest) (*DeclineResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/decline.do", request)
	if err != nil {
		return nil, err
	}
	response := &DeclineResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &DeclineResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) ProcessRawSumRefund(ctx context.Context, request *ProcessRawSumRefundRequest) (*ProcessRawSumRefundResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/processRawSumRefund.do", request)
	if err != nil {
		return nil, err
	}
	response := &ProcessRawSumRefundResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &ProcessRawSumRefundResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) ProcessRawPositionRefund(ctx context.Context, request *ProcessRawPositionRefundRequest) (*ProcessRawPositionRefundResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/processRawPositionRefund.do", request)
	if err != nil {
		return nil, err
	}
	response := &ProcessRawPositionRefundResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &ProcessRawPositionRefundResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) ProcessRawPositionOrderRefund(ctx context.Context, request *ProcessRawPositionOrderRefundRequest) (*ProcessRawPositionOrderRefundResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/processRawPositionOrderRefund.do", request)
	if err != nil {
		return nil, err
	}
	response := &ProcessRawPositionOrderRefundResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &ProcessRawPositionOrderRefundResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) PaymentOrderBinding(ctx context.Context, request *PaymentOrderBindingRequest) (*PaymentOrderBindingResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/paymentOrderBinding.do", request)
	if err != nil {
		return nil, err
	}
	response := &PaymentOrderBindingResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &PaymentOrderBindingResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) GetBindings(ctx context.Context, request *GetBindingsRequest) (*GetBindingsResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/getBindings.do", request)
	if err != nil {
		return nil, err
	}
	response := &GetBindingsResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &GetBindingsResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) GetBindingsByCardOrID(ctx context.Context, request *GetBindingsByCardOrIDRequest) (*GetBindingsByCardOrIDResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/getBindingsByCardOrId.do", request)
	if err != nil {
		return nil, err
	}
	response := &GetBindingsByCardOrIDResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &GetBindingsByCardOrIDResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) UnBindCard(ctx context.Context, request *UnBindCardRequest) (*UnBindCardResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/unBindCard.do", request)
	if err != nil {
		return nil, err
	}
	response := &UnBindCardResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &UnBindCardResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) BindCard(ctx context.Context, request *BindCardRequest) (*BindCardResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/bindCard.do", request)
	if err != nil {
		return nil, err
	}
	response := &BindCardResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &BindCardResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) ExtendBinding(ctx context.Context, request *ExtendBindingRequest) (*ExtendBindingResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/extendBinding.do", request)
	if err != nil {
		return nil, err
	}
	response := &ExtendBindingResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &ExtendBindingResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) CreateBindingNoPayment(ctx context.Context, request *CreateBindingNoPaymentRequest) (*CreateBindingNoPaymentResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/createBindingNoPayment.do", request)
	if err != nil {
		return nil, err
	}
	response := &CreateBindingNoPaymentResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &CreateBindingNoPaymentResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) SBPC2BQRDynamicGet(ctx context.Context, request *SBPC2BQRDynamicGetRequest) (*SBPC2BQRDynamicGetResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/sbp/c2b/qr/dynamic/get.do", request)
	if err != nil {
		return nil, err
	}
	response := &SBPC2BQRDynamicGetResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &SBPC2BQRDynamicGetResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) SBPC2BQRStatus(ctx context.Context, request *SBPC2BQRStatusRequest) (*SBPC2BQRStatusResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/sbp/c2b/qr/status.do", request)
	if err != nil {
		return nil, err
	}
	response := &SBPC2BQRStatusResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &SBPC2BQRStatusResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) SBPC2BQRDynamicReject(ctx context.Context, request *SBPC2BQRDynamicRejectRequest) (*SBPC2BQRDynamicRejectResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/sbp/c2b/qr/dynamic/reject.do", request)
	if err != nil {
		return nil, err
	}
	response := &SBPC2BQRDynamicRejectResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &SBPC2BQRDynamicRejectResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) TemplatesCreateTemplate(ctx context.Context, request *TemplatesCreateTemplateRequest) (*TemplatesCreateTemplateResult, error) {
	applyClientAuth(c, request)
	req, err := c.newJSONRequest(ctx, "/payment/rest/templates/createTemplate.do", request)
	if err != nil {
		return nil, err
	}
	response := &TemplatesCreateTemplateResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &TemplatesCreateTemplateResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) TemplatesGetTemplateDetails(ctx context.Context, request *TemplatesGetTemplateDetailsRequest) (*TemplatesGetTemplateDetailsResult, error) {
	applyClientAuth(c, request)
	req, err := c.newJSONRequest(ctx, "/payment/rest/templates/getTemplateDetails.do", request)
	if err != nil {
		return nil, err
	}
	response := &TemplatesGetTemplateDetailsResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &TemplatesGetTemplateDetailsResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) TemplatesUpdateTemplate(ctx context.Context, request *TemplatesUpdateTemplateRequest) (*TemplatesUpdateTemplateResult, error) {
	applyClientAuth(c, request)
	req, err := c.newJSONRequest(ctx, "/payment/rest/templates/updateTemplate.do", request)
	if err != nil {
		return nil, err
	}
	response := &TemplatesUpdateTemplateResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &TemplatesUpdateTemplateResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) SBPC2BGetBindings(ctx context.Context, request *SBPC2BGetBindingsRequest) (*SBPC2BGetBindingsResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/sbp/c2b/getBindings.do", request)
	if err != nil {
		return nil, err
	}
	response := &SBPC2BGetBindingsResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &SBPC2BGetBindingsResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) SBPC2BUnBind(ctx context.Context, request *SBPC2BUnBindRequest) (*SBPC2BUnBindResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/sbp/c2b/unBind.do", request)
	if err != nil {
		return nil, err
	}
	response := &SBPC2BUnBindResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &SBPC2BUnBindResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) Finish3DSPayment(ctx context.Context, request *Finish3DSPaymentRequest) (*Finish3DSPaymentResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/finish3dsPayment.do", request)
	if err != nil {
		return nil, err
	}
	response := &Finish3DSPaymentResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &Finish3DSPaymentResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) Finish3DSVer2Payment(ctx context.Context, request *Finish3DSVer2PaymentRequest) (*Finish3DSVer2PaymentResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/finish3dsVer2Payment.do", request)
	if err != nil {
		return nil, err
	}
	response := &Finish3DSVer2PaymentResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &Finish3DSVer2PaymentResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) ThreeDSContinue(ctx context.Context, request *ThreeDSContinueRequest) (*ThreeDSContinueResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/3ds/continue.do", request)
	if err != nil {
		return nil, err
	}
	response := &ThreeDSContinueResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &ThreeDSContinueResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) VerifyCard(ctx context.Context, request *VerifyCardRequest) (*VerifyCardResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/verifyCard.do", request)
	if err != nil {
		return nil, err
	}
	response := &VerifyCardResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &VerifyCardResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) CloseOFDReceipt(ctx context.Context, request *CloseOFDReceiptRequest) (*CloseOFDReceiptResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/closeOfdReceipt.do", request)
	if err != nil {
		return nil, err
	}
	response := &CloseOFDReceiptResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &CloseOFDReceiptResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) GetReceiptStatus(ctx context.Context, request *GetReceiptStatusRequest) (*GetReceiptStatusResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/getReceiptStatus.do", request)
	if err != nil {
		return nil, err
	}
	response := &GetReceiptStatusResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &GetReceiptStatusResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) TemplatesCreate(ctx context.Context, request *TemplatesCreateRequest) (*TemplatesCreateResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/templates/create.do", request)
	if err != nil {
		return nil, err
	}
	response := &TemplatesCreateResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &TemplatesCreateResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) TemplatesGet(ctx context.Context, request *TemplatesGetRequest) (*TemplatesGetResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/templates/get.do", request)
	if err != nil {
		return nil, err
	}
	response := &TemplatesGetResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &TemplatesGetResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) TemplatesUpdate(ctx context.Context, request *TemplatesUpdateRequest) (*TemplatesUpdateResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/templates/update.do", request)
	if err != nil {
		return nil, err
	}
	response := &TemplatesUpdateResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &TemplatesUpdateResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) TemplatesGetList(ctx context.Context, request *TemplatesGetListRequest) (*TemplatesGetListResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/templates/getList.do", request)
	if err != nil {
		return nil, err
	}
	response := &TemplatesGetListResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &TemplatesGetListResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}

func (c *Client) Finish3DS(ctx context.Context, request *Finish3DSRequest) (*Finish3DSResult, error) {
	applyClientAuth(c, request)
	req, err := c.newRequest(ctx, "/payment/rest/finish3ds.do", request)
	if err != nil {
		return nil, err
	}
	response := &Finish3DSResponse{}
	statusCode, header, body, err := c.do(req, response)
	if err != nil {
		return nil, err
	}
	return &Finish3DSResult{StatusCode: statusCode, Header: header, Body: body, Response: response}, nil
}
