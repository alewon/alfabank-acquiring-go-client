package alfabank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type ErrorCode string

type FlexibleString string

func (s *FlexibleString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*s = ""
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = FlexibleString(str)
		return nil
	}

	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*s = FlexibleString(n.String())
		return nil
	}

	return fmt.Errorf("unsupported string payload: %s", string(data))
}

type FlexibleInt64 int64

func (n *FlexibleInt64) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*n = 0
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*n = FlexibleInt64(i)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		parsed, err := strconv.ParseInt(num.String(), 10, 64)
		if err != nil {
			return err
		}
		*n = FlexibleInt64(parsed)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" {
			*n = 0
			return nil
		}
		parsed, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*n = FlexibleInt64(parsed)
		return nil
	}

	return fmt.Errorf("unsupported int64 payload: %s", string(data))
}

func (c *ErrorCode) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*c = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*c = ErrorCode(s)
		return nil
	}

	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*c = ErrorCode(n.String())
		return nil
	}

	return fmt.Errorf("unsupported errorCode payload: %s", string(data))
}

type OrderStatus int

const (
	OrderStatusRegistered        OrderStatus = 0
	OrderStatusAuthorized        OrderStatus = 1
	OrderStatusDeposited         OrderStatus = 2
	OrderStatusReversed          OrderStatus = 3
	OrderStatusRefunded          OrderStatus = 4
	OrderStatusACSAuthentication OrderStatus = 5
	OrderStatusDeclined          OrderStatus = 6
	OrderStatusAwaitingPayment   OrderStatus = 7
	OrderStatusPartialDeposit    OrderStatus = 8
)

type PaymentState string

const (
	PaymentStateCreated   PaymentState = "CREATED"
	PaymentStateApproved  PaymentState = "APPROVED"
	PaymentStateDeposited PaymentState = "DEPOSITED"
	PaymentStateDeclined  PaymentState = "DECLINED"
	PaymentStateReversed  PaymentState = "REVERSED"
	PaymentStateRefunded  PaymentState = "REFUNDED"
)

type PaymentWay string

const (
	PaymentWayCard                 PaymentWay = "CARD"
	PaymentWayCardBinding          PaymentWay = "CARD_BINDING"
	PaymentWayAlfaClick            PaymentWay = "ALFA_ALFACLICK"
	PaymentWayAlfaPay              PaymentWay = "ALFAPAY"
	PaymentWayAlfaPayBinding       PaymentWay = "ALFAPAY_BINDING"
	PaymentWayCardMOTO             PaymentWay = "CARD_MOTO"
	PaymentWayP2P                  PaymentWay = "P2P"
	PaymentWayApplePay             PaymentWay = "APPLE_PAY"
	PaymentWayApplePayBinding      PaymentWay = "APPLE_PAY_BINDING"
	PaymentWayGooglePayCard        PaymentWay = "GOOGLE_PAY_CARD"
	PaymentWayGooglePayCardBinding PaymentWay = "GOOGLE_PAY_CARD_BINDING"
	PaymentWayGooglePayTokenized   PaymentWay = "GOOGLE_PAY_TOKENIZED_BINDING"
	PaymentWaySBPC2B               PaymentWay = "SBP_C2B"
	PaymentWaySBPC2BBinding        PaymentWay = "SBP_C2B_BINDING"
	PaymentWayYandexPayCard        PaymentWay = "YANDEX_PAY_CARD"
	PaymentWayYandexPayBinding     PaymentWay = "YANDEX_PAY_CARD_BINDING"
)

type OrderFeature string

const (
	OrderFeatureAutoPayment     OrderFeature = "AUTO_PAYMENT"
	OrderFeatureVerify          OrderFeature = "VERIFY"
	OrderFeatureSBPBinding      OrderFeature = "SBP_BINDING"
	OrderFeatureForceTDS        OrderFeature = "FORCE_TDS"
	OrderFeatureForceSSL        OrderFeature = "FORCE_SSL"
	OrderFeatureForceFullTDS    OrderFeature = "FORCE_FULL_TDS"
	OrderFeatureForceCreateBind OrderFeature = "FORCE_CREATE_BINDING"
)

type TaxSystem int

const (
	TaxSystemCommon              TaxSystem = 0
	TaxSystemSimplifiedIncome    TaxSystem = 1
	TaxSystemSimplifiedOutcome   TaxSystem = 2
	TaxSystemUnifiedImputed      TaxSystem = 3
	TaxSystemUnifiedAgricultural TaxSystem = 4
	TaxSystemPatent              TaxSystem = 5
)

type TaxType int

const (
	TaxTypeNoVAT    TaxType = 0
	TaxTypeVAT0     TaxType = 1
	TaxTypeVAT10    TaxType = 2
	TaxTypeVAT10110 TaxType = 4
	TaxTypeVAT20    TaxType = 6
	TaxTypeVAT20120 TaxType = 7
	TaxTypeVAT5     TaxType = 10
	TaxTypeVAT5105  TaxType = 11
	TaxTypeVAT7     TaxType = 12
	TaxTypeVAT7107  TaxType = 13
	TaxTypeVAT22    TaxType = 14
	TaxTypeVAT22122 TaxType = 15
)

type PaymentMethod int

const (
	PaymentMethodFullPrepayment         PaymentMethod = 1
	PaymentMethodPartialPrepayment      PaymentMethod = 2
	PaymentMethodAdvance                PaymentMethod = 3
	PaymentMethodFullPayment            PaymentMethod = 4
	PaymentMethodPartialCredit          PaymentMethod = 5
	PaymentMethodTransferWithoutPayment PaymentMethod = 6
	PaymentMethodCreditPayment          PaymentMethod = 7
)

type PaymentObject int

const (
	PaymentObjectCommodity          PaymentObject = 1
	PaymentObjectExcise             PaymentObject = 2
	PaymentObjectWork               PaymentObject = 3
	PaymentObjectService            PaymentObject = 4
	PaymentObjectBet                PaymentObject = 5
	PaymentObjectGamblingWin        PaymentObject = 6
	PaymentObjectLotteryTicket      PaymentObject = 7
	PaymentObjectLotteryWin         PaymentObject = 8
	PaymentObjectRID                PaymentObject = 9
	PaymentObjectPayment            PaymentObject = 10
	PaymentObjectAgentFee           PaymentObject = 11
	PaymentObjectComposite          PaymentObject = 12
	PaymentObjectOther              PaymentObject = 13
	PaymentObjectPropertyRight      PaymentObject = 14
	PaymentObjectNonOperatingIncome PaymentObject = 15
	PaymentObjectInsuranceContrib   PaymentObject = 16
	PaymentObjectTradeFee           PaymentObject = 17
	PaymentObjectResortFee          PaymentObject = 18
)

type AgentType int

const (
	AgentTypeBankPaymentAgent    AgentType = 1
	AgentTypeBankPaymentSubagent AgentType = 2
	AgentTypePaymentAgent        AgentType = 3
	AgentTypePaymentSubagent     AgentType = 4
	AgentTypeAttorney            AgentType = 5
	AgentTypeCommissionAgent     AgentType = 6
	AgentTypeOther               AgentType = 7
)

type APISuccess struct {
	Success *bool `json:"success,omitempty"`
}

type BaseResponse struct {
	APISuccess
	ErrorCode    ErrorCode `json:"errorCode,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

type RegisterResponse struct {
	BaseResponse
	FormURL string `json:"formUrl,omitempty"`
	OrderID string `json:"orderId,omitempty"`
}

type RegisterPreAuthResponse = RegisterResponse

type OperationResponse struct {
	BaseResponse
}

type RegisterRequest struct {
	OrderNumber                             string              `form:"orderNumber"`
	Amount                                  int64               `form:"amount"`
	Currency                                string              `form:"currency,omitempty"`
	ReturnURL                               string              `form:"returnUrl"`
	FailURL                                 string              `form:"failUrl,omitempty"`
	DynamicCallbackURL                      string              `form:"dynamicCallbackUrl,omitempty"`
	Description                             string              `form:"description,omitempty"`
	Language                                string              `form:"language,omitempty"`
	IP                                      string              `form:"ip,omitempty"`
	ClientID                                string              `form:"clientId,omitempty"`
	MerchantLogin                           string              `form:"merchantLogin,omitempty"`
	CardholderName                          string              `form:"cardholderName,omitempty"`
	JSONParams                              map[string]string   `form:"jsonParams,omitempty,json"`
	SessionTimeoutSecs                      int64               `form:"sessionTimeoutSecs,omitempty"`
	ExpirationDate                          string              `form:"expirationDate,omitempty"`
	BindingID                               string              `form:"bindingId,omitempty"`
	Features                                []OrderFeature      `form:"features,omitempty"`
	PostAddress                             string              `form:"postAddress,omitempty"`
	OrderBundle                             *OrderBundle        `form:"orderBundle,omitempty,json"`
	AdditionalOfdParams                     *AdditionalOFDParams `form:"additionalOfdParams,omitempty,json"`
	TaxSystem                               TaxSystem           `form:"taxSystem,omitempty"`
	FeeInput                                int64               `form:"feeInput,omitempty"`
	Email                                   string              `form:"email,omitempty"`
	MerchantInn                             string              `form:"merchantInn,omitempty"`
	BillingPayerData                        *BillingPayerData   `form:"billingPayerData,omitempty,json"`
	ShippingPayerData                       *ShippingPayerData  `form:"shippingPayerData,omitempty,json"`
	PreOrderPayerData                       *PreOrderPayerData  `form:"preOrderPayerData,omitempty,json"`
	OrderPayerData                          *OrderPayerData     `form:"orderPayerData,omitempty,json"`
	BillingAndShippingAddressMatchIndicator string              `form:"billingAndShippingAddressMatchIndicator,omitempty"`
	MobileNumber                            string              `form:"mobile_number,omitempty"`
	WalletNumber                            string              `form:"wallet_number,omitempty"`
	WalletOperatorINN                       string              `form:"wallet_operator_inn,omitempty"`
}

type RegisterPreAuthRequest = RegisterRequest

type GetOrderStatusExtendedRequest struct {
	OrderID       string `form:"orderId,omitempty"`
	OrderNumber   string `form:"orderNumber,omitempty"`
	Language      string `form:"language,omitempty"`
	MerchantLogin string `form:"merchantLogin,omitempty"`
}

type DepositRequest struct {
	OrderID      string            `form:"orderId"`
	Amount       string            `form:"amount"`
	DepositItems *ItemsContainer   `form:"depositItems,omitempty,json"`
	Language     string            `form:"language,omitempty"`
	Currency     string            `form:"currency,omitempty"`
	JSONParams   map[string]string `form:"jsonParams,omitempty,json"`
}

type ReverseRequest struct {
	OrderID     string `form:"orderId"`
	Language    string `form:"language,omitempty"`
	OrderNumber string `form:"orderNumber,omitempty"`
	Amount      string `form:"amount,omitempty"`
	Currency    string `form:"currency,omitempty"`
}

type RefundRequest struct {
	OrderID             string               `form:"orderId"`
	Amount              string               `form:"amount"`
	Language            string               `form:"language,omitempty"`
	JSONParams          map[string]string    `form:"jsonParams,omitempty,json"`
	RefundItems         *ItemsContainer      `form:"refundItems,omitempty,json"`
	AdditionalOfdParams *AdditionalOFDParams `form:"additionalOfdParams,omitempty,json"`
}

type DeclineRequest struct {
	MerchantLogin string `form:"merchantLogin,omitempty"`
	Language      string `form:"language,omitempty"`
	OrderID       string `form:"orderId"`
	OrderNumber   string `form:"orderNumber"`
}

type ProcessRawSumRefundRequest struct {
	OrderID             string               `form:"orderId"`
	Language            string               `form:"language"`
	Amount              string               `form:"amount"`
	AdditionalOfdParams *AdditionalOFDParams `form:"additionalOfdParams,omitempty,json"`
	OFDParams           *OFDParams           `form:"ofdParams,omitempty,json"`
}

type ProcessRawPositionRefundRequest struct {
	OrderID             string               `form:"orderId"`
	Language            string               `form:"language"`
	Amount              string               `form:"amount"`
	PositionID          int64                `form:"positionId,omitempty"`
	AdditionalOfdParams *AdditionalOFDParams `form:"additionalOfdParams,omitempty,json"`
}

type GetOrderStatusExtendedResponse struct {
	BaseResponse
	OrderNumber           string                 `json:"orderNumber,omitempty"`
	OrderStatus           OrderStatus            `json:"orderStatus,omitempty"`
	ActionCode            FlexibleString         `json:"actionCode,omitempty"`
	ActionCodeDescription string                 `json:"actionCodeDescription,omitempty"`
	Amount                int64                  `json:"amount,omitempty"`
	Currency              FlexibleString         `json:"currency,omitempty"`
	Date                  int64                  `json:"date,omitempty"`
	DepositedDate         int64                  `json:"depositedDate,omitempty"`
	OrderDescription      string                 `json:"orderDescription,omitempty"`
	IP                    string                 `json:"ip,omitempty"`
	AuthRefNum            string                 `json:"authRefNum,omitempty"`
	RefundedDate          int64                  `json:"refundedDate,omitempty"`
	ReversedDate          int64                  `json:"reversedDate,omitempty"`
	PaymentWay            PaymentWay             `json:"paymentWay,omitempty"`
	AVSCode               string                 `json:"avsCode,omitempty"`
	Chargeback            *bool                  `json:"chargeback,omitempty"`
	AuthDateTime          int64                  `json:"authDateTime,omitempty"`
	TerminalID            string                 `json:"terminalId,omitempty"`
	OrderBundle           *OrderBundle           `json:"orderBundle,omitempty"`
	PaymentAmountInfo     *PaymentAmountInfo     `json:"paymentAmountInfo,omitempty"`
	Refunds               RefundCollection       `json:"refunds,omitempty"`
	CardAuthInfo          *CardAuthInfo          `json:"cardAuthInfo,omitempty"`
	TransactionAttributes []TransactionAttribute `json:"transactionAttributes,omitempty"`
	PrepaymentMdOrder     string                 `json:"prepaymentMdOrder,omitempty"`
	PartpaymentMdOrders   []string               `json:"partpaymentMdOrders,omitempty"`
	FeUtrnno              int64                  `json:"feUtrnno,omitempty"`
	BindingInfo           *BindingInfo           `json:"bindingInfo,omitempty"`
	EfectyOrderInfo       *EfectyOrderInfo       `json:"efectyOrderInfo,omitempty"`
	PluginInfo            *PluginInfo            `json:"pluginInfo,omitempty"`
	DisplayErrorMessage   string                 `json:"displayErrorMessage,omitempty"`
	TII                   string                 `json:"tii,omitempty"`
	UsedPsdIndicatorValue string                 `json:"usedPsdIndicatorValue,omitempty"`
	PayerData             *PayerData             `json:"payerData,omitempty"`
	LoyaltyInfo           *LoyaltyInfo           `json:"loyaltyInfo,omitempty"`
	MerchantOrderParams   []KeyValue             `json:"merchantOrderParams,omitempty"`
	Attributes            []KeyValue             `json:"attributes,omitempty"`
	BankInfo              *BankInfo              `json:"bankInfo,omitempty"`
}

type RefundCollection []Refund

func (c *RefundCollection) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*c = nil
		return nil
	}
	if data[0] == '[' {
		var items []Refund
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		*c = items
		return nil
	}
	var item Refund
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	*c = []Refund{item}
	return nil
}

type Refund struct {
	Date             string         `json:"date,omitempty"`
	ExternalRefundID string         `json:"externalRefundId,omitempty"`
	ApprovalCode     string         `json:"approvalCode,omitempty"`
	ActionCode       FlexibleString `json:"actionCode,omitempty"`
	ReferenceNumber  string         `json:"referenceNumber,omitempty"`
	Amount           int64          `json:"amount,omitempty"`
	Attributes       []KeyValue     `json:"attributes,omitempty"`
}

type PaymentAmountInfo struct {
	ApprovedAmount  int64        `json:"approvedAmount,omitempty"`
	DepositedAmount int64        `json:"depositedAmount,omitempty"`
	RefundedAmount  int64        `json:"refundedAmount,omitempty"`
	PaymentState    PaymentState `json:"paymentState,omitempty"`
	TotalAmount     int64        `json:"totalAmount,omitempty"`
}

type BankInfo struct {
	BankName        string `json:"bankName,omitempty"`
	BankCountryCode string `json:"bankCountryCode,omitempty"`
	BankCountryName string `json:"bankCountryName,omitempty"`
}

type PayerData struct {
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	PostAddress string `json:"postAddress,omitempty"`
}

type LoyaltyInfo struct {
	AwardBonus   *AwardBonus   `json:"awardBonus,omitempty"`
	PaymentBonus *PaymentBonus `json:"paymentBonus,omitempty"`
	LoyaltyName  string        `json:"loyaltyName,omitempty"`
}

type AwardBonus struct {
	ApprovedAmountAward  int64  `json:"approvedAmountAward,omitempty"`
	DepositedAmountAward int64  `json:"depositedAmountAward,omitempty"`
	RefundedAmountAward  int64  `json:"refundedAmountAward,omitempty"`
	PCID                 string `json:"pcId,omitempty"`
	Successful           string `json:"successful,omitempty"`
	PaymentOperation     string `json:"paymentOperation,omitempty"`
}

type PaymentBonus struct {
	ApprovedAmountBonus  int64  `json:"approvedAmountBonus,omitempty"`
	DepositedAmountBonus int64  `json:"depositedAmountBonus,omitempty"`
	RefundedAmountBonus  int64  `json:"refundedAmountBonus,omitempty"`
	PCID                 string `json:"pcId,omitempty"`
	Successful           string `json:"successful,omitempty"`
	PaymentOperation     string `json:"paymentOperation,omitempty"`
}

type EfectyOrderInfo struct {
	ReferenceNumber int64  `json:"referenceNumber,omitempty"`
	ReferenceDate   int64  `json:"referenceDate,omitempty"`
	ReferenceStatus string `json:"referenceStatus,omitempty"`
	ReferenceTerm   int64  `json:"referenceTerm,omitempty"`
	NetworkID       int64  `json:"networkID,omitempty"`
	NetworkName     string `json:"networkName,omitempty"`
}

type PluginInfo struct {
	Name   string            `json:"name,omitempty"`
	Params map[string]string `json:"params,omitempty"`
}

type KeyValue struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type TransactionAttribute struct {
	Name                    string `json:"name,omitempty"`
	Value                   string `json:"value,omitempty"`
	SubscriptionServiceName string `json:"subscriptionServiceName,omitempty"`
	SubscriptionServiceID   string `json:"subscriptionServiceId,omitempty"`
}

type CardAuthInfo struct {
	MaskedPan       string          `json:"maskedPan,omitempty"`
	Expiration      string          `json:"expiration,omitempty"`
	CardholderName  string          `json:"cardholderName,omitempty"`
	ApprovalCode    string          `json:"approvalCode,omitempty"`
	PaymentSystem   string          `json:"paymentSystem,omitempty"`
	Product         string          `json:"product,omitempty"`
	ProductCategory string          `json:"productCategory,omitempty"`
	CorporateCard   string          `json:"corporateCard,omitempty"`
	SecureAuthInfo  *SecureAuthInfo `json:"secureAuthInfo,omitempty"`
}

type SecureAuthInfo struct {
	ECI                    string `json:"eci,omitempty"`
	AuthTypeIndicator      string `json:"authTypeIndicator,omitempty"`
	CAVV                   string `json:"cavv,omitempty"`
	XID                    string `json:"xid,omitempty"`
	ThreeDSProtocolVersion string `json:"threeDSProtocolVersion,omitempty"`
	RReqTransStatus        string `json:"rreqTransStatus,omitempty"`
	AResTransStatus        string `json:"aresTransStatus,omitempty"`
	PaResStatus            string `json:"paResStatus,omitempty"`
	VeResStatus            string `json:"veResStatus,omitempty"`
	PaResCheckStatus       string `json:"paResCheckStatus,omitempty"`
}

type BindingInfo struct {
	ClientID        string `json:"clientId,omitempty"`
	BindingID       string `json:"bindingId,omitempty"`
	AuthDateTime    int64  `json:"authDateTime,omitempty"`
	AuthRefNum      string `json:"authRefNum,omitempty"`
	TerminalID      string `json:"terminalId,omitempty"`
	ExternalCreated *bool  `json:"externalCreated,omitempty"`
}

type BillingPayerData struct {
	BillingCity         string `json:"billingCity,omitempty"`
	BillingCountry      string `json:"billingCountry,omitempty"`
	BillingAddressLine1 string `json:"billingAddressLine1,omitempty"`
	BillingAddressLine2 string `json:"billingAddressLine2,omitempty"`
	BillingAddressLine3 string `json:"billingAddressLine3,omitempty"`
	BillingPostalCode   string `json:"billingPostalCode,omitempty"`
	BillingState        string `json:"billingState,omitempty"`
}

type ShippingPayerData struct {
	ShippingCity            string `json:"shippingCity,omitempty"`
	ShippingCountry         string `json:"shippingCountry,omitempty"`
	ShippingAddressLine1    string `json:"shippingAddressLine1,omitempty"`
	ShippingAddressLine2    string `json:"shippingAddressLine2,omitempty"`
	ShippingAddressLine3    string `json:"shippingAddressLine3,omitempty"`
	ShippingPostalCode      string `json:"shippingPostalCode,omitempty"`
	ShippingState           string `json:"shippingState,omitempty"`
	ShippingMethodIndicator string `json:"shippingMethodIndicator,omitempty"`
	DeliveryTimeframe       string `json:"deliveryTimeframe,omitempty"`
	DeliveryEmail           string `json:"deliveryEmail,omitempty"`
}

type PreOrderPayerData struct {
	PreOrderDate        string `json:"preOrderDate,omitempty"`
	PreOrderPurchaseInd string `json:"preOrderPurchaseInd,omitempty"`
	ReorderItemsInd     string `json:"reorderItemsInd,omitempty"`
}

type OrderPayerData struct {
	HomePhone   string `json:"homePhone,omitempty"`
	WorkPhone   string `json:"workPhone,omitempty"`
	MobilePhone string `json:"mobilePhone,omitempty"`
}

type OrderBundle struct {
	OrderCreationDate string           `json:"orderCreationDate,omitempty"`
	CustomerDetails   *CustomerDetails `json:"customerDetails,omitempty"`
	CartItems         *ItemsContainer  `json:"cartItems,omitempty"`
	Agent             *Agent           `json:"agent,omitempty"`
	SupplierPhones    []string         `json:"supplierPhones,omitempty"`
	Loyalties         *Loyalties       `json:"loyalties,omitempty"`
}

type Agent struct {
	AgentType              AgentType `json:"agentType,omitempty"`
	PayingOperation        string    `json:"payingOperation,omitempty"`
	PayingPhones           []string  `json:"payingPhones,omitempty"`
	PaymentsOperatorPhones []string  `json:"paymentsOperatorPhones,omitempty"`
	MTOperatorPhones       []string  `json:"MTOperatorPhones,omitempty"`
	MTOperatorName         string    `json:"MTOperatorName,omitempty"`
	MTOperatorAddress      string    `json:"MTOperatorAddress,omitempty"`
	MTOperatorInn          string    `json:"MTOperatorInn,omitempty"`
}

type Loyalties struct {
	BonusAmountForCredit string `json:"bonusAmountForCredit,omitempty"`
	BonusAmountForDebit  string `json:"bonusAmountForDebit,omitempty"`
	BonusAmountRefunded  string `json:"bonusAmountRefunded,omitempty"`
}

type CustomerDetails struct {
	Email        string        `json:"email,omitempty"`
	Phone        string        `json:"phone,omitempty"`
	Contact      string        `json:"contact,omitempty"`
	FullName     string        `json:"fullName,omitempty"`
	Passport     string        `json:"passport,omitempty"`
	DeliveryInfo *DeliveryInfo `json:"deliveryInfo,omitempty"`
	INN          string        `json:"inn,omitempty"`
}

type DeliveryInfo struct {
	DeliveryType string `json:"deliveryType,omitempty"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	PostAddress  string `json:"postAddress,omitempty"`
}

type ItemsContainer struct {
	Items []Item `json:"items,omitempty"`
}

type Item struct {
	PositionID          FlexibleInt64   `json:"positionId,omitempty"`
	Name                string          `json:"name,omitempty"`
	ItemDetails         *ItemDetails    `json:"itemDetails,omitempty"`
	Quantity            *Quantity       `json:"quantity,omitempty"`
	ItemAmount          int64           `json:"itemAmount,omitempty"`
	ItemPrice           int64           `json:"itemPrice,omitempty"`
	DepositedItemAmount FlexibleString  `json:"depositedItemAmount,omitempty"`
	ItemCurrency        int64           `json:"itemCurrency,omitempty"`
	ItemCode            string          `json:"itemCode,omitempty"`
	Tax                 *Tax            `json:"tax,omitempty"`
	ItemAttributes      *ItemAttributes `json:"itemAttributes,omitempty"`
	Nomenclature        string          `json:"nomenclature,omitempty"`
	MarkQuantity        *MarkQuantity   `json:"markQuantity,omitempty"`
	UserData            string          `json:"userData,omitempty"`
	AgentInfo           *AgentInfo      `json:"agent_info,omitempty"`
	SupplierInfo        *SupplierInfo   `json:"supplier_info,omitempty"`
}

type ItemDetails struct {
	ItemDetailsParams []KeyValue `json:"itemDetailsParams,omitempty"`
}

type Quantity struct {
	Value   float64 `json:"value,omitempty"`
	Measure string  `json:"measure,omitempty"`
}

type Tax struct {
	TaxType TaxType `json:"taxType,omitempty"`
	TaxSum  int64   `json:"taxSum,omitempty"`
}

type ItemAttributes struct {
	Attributes []KeyValue `json:"attributes,omitempty"`
}

type MarkQuantity struct {
	Numerator   int64 `json:"numerator,omitempty"`
	Denominator int64 `json:"denominator,omitempty"`
}

type AdditionalOFDParams struct {
	AgentInfo            *AgentInfo           `json:"agent_info,omitempty"`
	SupplierInfo         *SupplierInfo        `json:"supplier_info,omitempty"`
	Cashier              string               `json:"cashier,omitempty"`
	AdditionalCheckProps string               `json:"additional_check_props,omitempty"`
	AdditionalUserProps  *AdditionalUserProps `json:"additional_user_props,omitempty"`
	CashierINN           string               `json:"cashier_inn,omitempty"`
	Client               *OFDClient           `json:"client,omitempty"`
	OperatingCheckProps  *OperatingCheckProps `json:"operatingcheckprops,omitempty"`
	SectoralCheckProps   *SectoralCheckProps  `json:"sectoralcheckprops,omitempty"`
	Company              *Company             `json:"company,omitempty"`
	UseLegacyVAT         *bool                `json:"use_legacy_vat,omitempty"`
}

type OFDParams struct {
	Name     string  `json:"name,omitempty"`
	ItemCode string  `json:"itemCode,omitempty"`
	TaxType  TaxType `json:"taxType,omitempty"`
}

type AgentInfo struct {
	Type             AgentType    `json:"type,omitempty"`
	Paying           *Paying      `json:"paying,omitempty"`
	PaymentsOperator *PhonesBlock `json:"paymentsOperator,omitempty"`
	MTOperator       *MTOperator  `json:"MTOperator,omitempty"`
}

type Paying struct {
	Operation string   `json:"operation,omitempty"`
	Phones    []string `json:"phones,omitempty"`
}

type PhonesBlock struct {
	Phones []string `json:"phones,omitempty"`
}

type MTOperator struct {
	Phones  []string `json:"phones,omitempty"`
	Name    string   `json:"name,omitempty"`
	Address string   `json:"address,omitempty"`
	INN     string   `json:"inn,omitempty"`
}

type SupplierInfo struct {
	Phones []string `json:"phones,omitempty"`
	Name   string   `json:"name,omitempty"`
	INN    string   `json:"inn,omitempty"`
}

type AdditionalUserProps struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type OFDClient struct {
	Address        string `json:"address,omitempty"`
	BirthDate      string `json:"birth_date,omitempty"`
	Citizenship    string `json:"citizenship,omitempty"`
	DocumentCode   string `json:"document_code,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	INN            string `json:"inn,omitempty"`
	Name           string `json:"name,omitempty"`
}

type OperatingCheckProps struct {
	Name      string `json:"name,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Value     string `json:"value,omitempty"`
}

type SectoralCheckProps struct {
	Date      string `json:"date,omitempty"`
	FederalID string `json:"federalid,omitempty"`
	Number    string `json:"number,omitempty"`
	Value     string `json:"value,omitempty"`
}

type Company struct {
	AutomatNumber  string `json:"automat_number,omitempty"`
	Location       string `json:"location,omitempty"`
	PaymentAddress string `json:"payment_address,omitempty"`
}