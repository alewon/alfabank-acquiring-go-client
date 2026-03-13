package alfabank

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type JSONParams map[string]string

type FlexibleInt64 int64

func (v *FlexibleInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*v = 0
		return nil
	}

	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*v = FlexibleInt64(num)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		*v = 0
		return nil
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*v = FlexibleInt64(num)
	return nil
}

type NameValuePair struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type TransactionAttribute struct {
	Name                    string `json:"name,omitempty"`
	Value                   string `json:"value,omitempty"`
	SubscriptionServiceName string `json:"subscriptionServiceName,omitempty"`
	SubscriptionServiceID   string `json:"subscriptionServiceId,omitempty"`
}

type Loyalty struct {
	BonusAmountForCredit string `json:"bonusAmountForCredit,omitempty"`
	BonusAmountForDebit  string `json:"bonusAmountForDebit,omitempty"`
	BonusAmountRefunded  string `json:"bonusAmountRefunded,omitempty"`
	LoyaltyProgramName   string `json:"loyaltyProgramName,omitempty"`
	PositionID           int64  `json:"positionId,omitempty"`
}

type OrderBundle struct {
	OrderCreationDate string          `json:"orderCreationDate,omitempty"`
	CustomerDetails   CustomerDetails `json:"customerDetails,omitempty"`
	CartItems         CartItems       `json:"cartItems,omitempty"`
	Agent             Agent           `json:"agent,omitempty"`
	SupplierPhones    []string        `json:"supplierPhones,omitempty"`
	Loyalties         []Loyalty       `json:"loyalties,omitempty"`
}

type CustomerDetails struct {
	Email        string       `json:"email,omitempty"`
	Phone        string       `json:"phone,omitempty"`
	Contact      string       `json:"contact,omitempty"`
	FullName     string       `json:"fullName,omitempty"`
	Passport     string       `json:"passport,omitempty"`
	DeliveryInfo DeliveryInfo `json:"deliveryInfo,omitempty"`
	INN          int64        `json:"inn,omitempty"`
}

type DeliveryInfo struct {
	DeliveryType string `json:"deliveryType,omitempty"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	PostAddress  string `json:"postAddress,omitempty"`
}

type CartItems struct {
	Items []CartItem `json:"items,omitempty"`
}

type CartItem struct {
	PositionID          int64          `json:"positionId,omitempty"`
	Name                string         `json:"name,omitempty"`
	ItemDetails         ItemDetails    `json:"itemDetails,omitempty"`
	Quantity            Quantity       `json:"quantity,omitempty"`
	ItemAmount          int64          `json:"itemAmount,omitempty"`
	ItemPrice           int64          `json:"itemPrice,omitempty"`
	DepositedItemAmount string         `json:"depositedItemAmount,omitempty"`
	ItemCurrency        int64          `json:"itemCurrency,omitempty"`
	ItemCode            string         `json:"itemCode,omitempty"`
	Tax                 Tax            `json:"tax,omitempty"`
	ItemAttributes      ItemAttributes `json:"itemAttributes,omitempty"`
	AgentInfo           AgentInfo      `json:"agent_info,omitempty"`
	SupplierInfo        SupplierInfo   `json:"supplier_info,omitempty"`
}

type Quantity struct {
	Value   string `json:"value,omitempty"`
	Measure string `json:"measure,omitempty"`
}

type ItemDetails struct {
	ItemDetailsParams ItemDetailsParams `json:"itemDetailsParams,omitempty"`
}

type ItemDetailsParams struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Tax struct {
	TaxType int64 `json:"taxType,omitempty"`
	TaxSum  int64 `json:"taxSum,omitempty"`
}

type ItemAttributes struct {
	Attributes []NameValuePair `json:"attributes,omitempty"`
}

type AgentInfo struct {
	Type             int64            `json:"type,omitempty"`
	Paying           Paying           `json:"paying,omitempty"`
	PaymentsOperator PaymentsOperator `json:"paymentsOperator,omitempty"`
	MTOperator       MTOperator       `json:"MTOperator,omitempty"`
}

type Paying struct {
	Operation string   `json:"operation,omitempty"`
	Phones    []string `json:"phones,omitempty"`
}

type PaymentsOperator struct {
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

type MarkQuantity struct {
	Numerator   int64 `json:"numerator,omitempty"`
	Denominator int64 `json:"denominator,omitempty"`
}

type Agent struct {
	AgentType              int64    `json:"agentType,omitempty"`
	PayingOperation        string   `json:"payingOperation,omitempty"`
	PayingPhones           []string `json:"payingPhones,omitempty"`
	PaymentsOperatorPhones []string `json:"paymentsOperatorPhones,omitempty"`
	MTOperatorPhones       []string `json:"MTOperatorPhones,omitempty"`
	MTOperatorName         string   `json:"MTOperatorName,omitempty"`
	MTOperatorAddress      string   `json:"MTOperatorAddress,omitempty"`
	MTOperatorINN          string   `json:"MTOperatorInn,omitempty"`
}

type AdditionalUserProps struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type AdditionalOFDParams struct {
	AgentInfoType                int64               `json:"agent_info.type,omitempty"`
	AgentInfoPayingOperation     string              `json:"agent_info.paying.operation,omitempty"`
	AgentInfoPayingPhones        []string            `json:"agent_info.paying.phones,omitempty"`
	AgentInfoPaymentsOperator    []string            `json:"agent_info.paymentsOperator.phones,omitempty"`
	AgentInfoMTOperatorAddress   string              `json:"agent_info.MTOperator.address,omitempty"`
	AgentInfoMTOperatorINN       string              `json:"agent_info.MTOperator.inn,omitempty"`
	AgentInfoMTOperatorName      string              `json:"agent_info.MTOperator.name,omitempty"`
	AgentInfoMTOperatorPhones    []string            `json:"agent_info.MTOperator.phones,omitempty"`
	SupplierInfoPhones           []string            `json:"supplier_info.phones,omitempty"`
	Cashier                      string              `json:"cashier,omitempty"`
	AdditionalCheckProps         string              `json:"additional_check_props,omitempty"`
	AdditionalUserProps          AdditionalUserProps `json:"additional_user_props,omitempty"`
	CashierINN                   string              `json:"cashier_inn,omitempty"`
	ClientAddress                string              `json:"client.address,omitempty"`
	ClientBirthDate              string              `json:"client.birth_date,omitempty"`
	ClientCitizenship            string              `json:"client.citizenship,omitempty"`
	ClientDocumentCode           string              `json:"client.document_code,omitempty"`
	ClientPassportNumber         string              `json:"client.passport_number,omitempty"`
	ClientEmail                  string              `json:"client.email,omitempty"`
	ClientPhone                  string              `json:"client.phone,omitempty"`
	ClientINN                    string              `json:"client.inn,omitempty"`
	ClientName                   string              `json:"client.name,omitempty"`
	OperatingCheckPropsName      string              `json:"operatingcheckprops.name,omitempty"`
	OperatingCheckPropsValue     string              `json:"operatingcheckprops.value,omitempty"`
	OperatingCheckPropsTimestamp string              `json:"operatingcheckprops.timestamp,omitempty"`
	SectoralCheckPropsDate       string              `json:"sectoralcheckprops.date,omitempty"`
	SectoralCheckPropsFederalID  string              `json:"sectoralcheckprops.federal_id,omitempty"`
	SectoralCheckPropsNumber     string              `json:"sectoralcheckprops.number,omitempty"`
	SectoralCheckPropsValue      string              `json:"sectoralcheckprops.value,omitempty"`
	CompanyAutomatNumber         string              `json:"company.automat_number,omitempty"`
	CompanyLocation              string              `json:"company.location,omitempty"`
	CompanyPaymentAddress        string              `json:"company.payment_address,omitempty"`
	UseLegacyVAT                 bool                `json:"use_legacy_vat,omitempty"`
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

type ClientBrowserInfo struct {
	UserAgent             string `json:"userAgent,omitempty"`
	OS                    string `json:"OS,omitempty"`
	OSVersion             string `json:"OSVersion,omitempty"`
	BrowserAcceptHeader   string `json:"browserAcceptHeader,omitempty"`
	BrowserIPAddress      string `json:"browserIpAddress,omitempty"`
	BrowserLanguage       string `json:"browserLanguage,omitempty"`
	BrowserTimeZone       string `json:"browserTimeZone,omitempty"`
	BrowserTimeZoneOffset string `json:"browserTimeZoneOffset,omitempty"`
	ColorDepth            string `json:"colorDepth,omitempty"`
	Fingerprint           string `json:"fingerprint,omitempty"`
	IsMobile              bool   `json:"isMobile,omitempty"`
	JavaEnabled           bool   `json:"javaEnabled,omitempty"`
	JavascriptEnabled     bool   `json:"javascriptEnabled,omitempty"`
	Plugins               string `json:"plugins,omitempty"`
	ScreenHeight          int64  `json:"screenHeight,omitempty"`
	ScreenWidth           int64  `json:"screenWidth,omitempty"`
	ScreenPrint           string `json:"screenPrint,omitempty"`
}

type DepositItems struct {
	Items []CartItem `json:"items,omitempty"`
}

type RefundItems struct {
	Items []CartItem `json:"items,omitempty"`
}

type OFDParams struct {
	Name     string `json:"name,omitempty"`
	ItemCode string `json:"itemCode,omitempty"`
	TaxType  int64  `json:"taxType,omitempty"`
}

type AdditionalParameters []TemplateAdditionalParameter

type TemplateAdditionalParameter struct {
	Mode        string `json:"mode,omitempty"`
	Label       string `json:"label,omitempty"`
	Name        string `json:"name,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	Regexp      string `json:"regexp,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Value       string `json:"value,omitempty"`
	Visible     bool   `json:"visible,omitempty"`
}

type QRTemplate struct {
	QRHeight       string `json:"qrHeight,omitempty"`
	QRWidth        string `json:"qrWidth,omitempty"`
	PaymentPurpose string `json:"paymentPurpose,omitempty"`
	QRCID          string `json:"qrcId,omitempty"`
	Payload        string `json:"payload,omitempty"`
	RenderedQR     string `json:"renderedQr,omitempty"`
}

type TemplateCommission struct {
	FeeMin        int64  `json:"feeMin,omitempty"`
	FeeMax        int64  `json:"feeMax,omitempty"`
	FixedAmount   string `json:"fixedAmount,omitempty"`
	FeePercentage string `json:"feePercentage,omitempty"`
}

type Template struct {
	TemplateID           string               `json:"templateId,omitempty"`
	Status               string               `json:"status,omitempty"`
	Type                 string               `json:"type,omitempty"`
	Name                 string               `json:"name,omitempty"`
	PreAuth              bool                 `json:"preAuth,omitempty"`
	Amount               int64                `json:"amount,omitempty"`
	IsFreeAmount         bool                 `json:"isFreeAmount,omitempty"`
	Currency             string               `json:"currency,omitempty"`
	Description          string               `json:"description,omitempty"`
	NameForClient        string               `json:"nameForClient,omitempty"`
	DescriptionForClient string               `json:"descriptionForClient,omitempty"`
	SinglePayment        bool                 `json:"singlePayment,omitempty"`
	QRTemplate           QRTemplate           `json:"qrTemplate,omitempty"`
	Commission           TemplateCommission   `json:"commission,omitempty"`
	AdditionalParams     AdditionalParameters `json:"additionalParams,omitempty"`
	DistributionChannel  string               `json:"distributionChannel,omitempty"`
	StartDate            string               `json:"startDate,omitempty"`
	EndDate              string               `json:"endDate,omitempty"`
	IsIndefinite         bool                 `json:"isIndefinite,omitempty"`
}

type TemplateError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Message     string `json:"message,omitempty"`
}

type Binding struct {
	BindingID   string `json:"bindingId,omitempty"`
	MaskedPan   string `json:"maskedPan,omitempty"`
	ExpiryDate  string `json:"expiryDate,omitempty"`
	ClientID    string `json:"clientId,omitempty"`
	BindingType string `json:"bindingType,omitempty"`
	PaymentWay  string `json:"paymentWay,omitempty"`
}

type SBPC2BInfo struct {
	BankName          string `json:"bankName,omitempty"`
	BIC               string `json:"bic,omitempty"`
	MaskedBankAccount string `json:"maskedBankAccount,omitempty"`
	BindingID         string `json:"bindingId,omitempty"`
}

type ReceiptOFD struct {
	INN     string `json:"inn,omitempty"`
	Name    string `json:"name,omitempty"`
	Website string `json:"website,omitempty"`
}

type ReceiptOFDStatus struct {
	QueuedDoc          int64  `json:"queuedDoc,omitempty"`
	QueuedFirstDocNum  int64  `json:"queuedFirstDocNum,omitempty"`
	QueuedFirstDocDate string `json:"queuedFirstDocDate,omitempty"`
	FFDVersion         string `json:"ffdVersion,omitempty"`
	InternetSign       bool   `json:"internetSign,omitempty"`
}

type ReceiptOFDParams struct {
	OFDReceiptStatus ReceiptOFDStatus `json:"ofdReceiptStatus,omitempty"`
	VATAmount        int64            `json:"vatAmount,omitempty"`
}

type ReceiptItemAttribute struct {
	PaymentMethod string `json:"paymentMethod,omitempty"`
	PaymentObject string `json:"paymentObject,omitempty"`
}

type OFDOrderBundleItem struct {
	Name           string                 `json:"name,omitempty"`
	ItemAmount     int64                  `json:"itemAmount,omitempty"`
	ItemAttributes []ReceiptItemAttribute `json:"itemAttributes,omitempty"`
	ItemPrice      int64                  `json:"itemPrice,omitempty"`
	TaxType        int64                  `json:"taxType,omitempty"`
	Quantity       Quantity               `json:"quantity,omitempty"`
}

type ReceiptInfo struct {
	ReceiptStatus           FlexibleInt64        `json:"receiptStatus,omitempty"`
	ReceiptType             string               `json:"receiptType,omitempty"`
	UUID                    string               `json:"uuid,omitempty"`
	OriginalOFDUUID         string               `json:"original_ofd_uuid,omitempty"`
	ShiftNumber             FlexibleInt64        `json:"shift_number,omitempty"`
	FiscalReceiptNumber     FlexibleInt64        `json:"fiscal_receipt_number,omitempty"`
	ReceiptDateTime         FlexibleInt64        `json:"receipt_datetime,omitempty"`
	ReceiptDateTimeLegacy   FlexibleInt64        `json:"receipt_date_time,omitempty"`
	FNNumber                string               `json:"fn_number,omitempty"`
	ECRRegistrationNumber   string               `json:"ecr_registration_number,omitempty"`
	FiscalDocumentNumber    FlexibleInt64        `json:"fiscal_document_number,omitempty"`
	FiscalDocumentAttribute string               `json:"fiscal_document_attribute,omitempty"`
	AmountTotal             string               `json:"amount_total,omitempty"`
	SerialNumber            string               `json:"serial_number,omitempty"`
	OFDReceiptURL           string               `json:"ofd_receipt_url,omitempty"`
	OFD                     ReceiptOFD           `json:"OFD,omitempty"`
	OFDReceiptParams        ReceiptOFDParams     `json:"ofdReceiptParams,omitempty"`
	OFDOrderBundle          []OFDOrderBundleItem `json:"ofdOrderBundle,omitempty"`
}

type Refund struct {
	Date             string `json:"date,omitempty"`
	ExternalRefundID string `json:"externalRefundId,omitempty"`
	ApprovalCode     string `json:"approvalCode,omitempty"`
	ActionCode       string `json:"actionCode,omitempty"`
	ReferenceNumber  string `json:"referenceNumber,omitempty"`
	Amount           int64  `json:"amount,omitempty"`
}

type SecureAuthInfo struct {
	ECI               int64  `json:"eci,omitempty"`
	AuthTypeIndicator string `json:"authTypeIndicator,omitempty"`
	CAVV              string `json:"cavv,omitempty"`
	XID               string `json:"xid,omitempty"`
	ThreeDSProtocol   string `json:"threeDSProtocolVersion,omitempty"`
	RReqTransStatus   string `json:"rreqTransStatus,omitempty"`
	AResTransStatus   string `json:"aresTransStatus,omitempty"`
	PaResStatus       string `json:"paResStatus,omitempty"`
	VeResStatus       string `json:"veResStatus,omitempty"`
	PaResCheckStatus  string `json:"paResCheckStatus,omitempty"`
}

type CardAuthInfo struct {
	MaskedPan       string         `json:"maskedPan,omitempty"`
	Expiration      string         `json:"expiration,omitempty"`
	CardholderName  string         `json:"cardholderName,omitempty"`
	ApprovalCode    string         `json:"approvalCode,omitempty"`
	Pan             string         `json:"pan,omitempty"`
	PaymentSystem   string         `json:"paymentSystem,omitempty"`
	Product         string         `json:"product,omitempty"`
	ProductCategory string         `json:"productCategory,omitempty"`
	CorporateCard   string         `json:"corporateCard,omitempty"`
	SecureAuthInfo  SecureAuthInfo `json:"secureAuthInfo,omitempty"`
}

type BindingInfo struct {
	ClientID        string `json:"clientId,omitempty"`
	BindingID       string `json:"bindingId,omitempty"`
	AuthDateTime    int64  `json:"authDateTime,omitempty"`
	AuthRefNum      string `json:"authRefNum,omitempty"`
	TerminalID      string `json:"terminalId,omitempty"`
	ExternalCreated bool   `json:"externalCreated,omitempty"`
}

type PaymentAmountInfo struct {
	ApprovedAmount  int64  `json:"approvedAmount,omitempty"`
	DepositedAmount int64  `json:"depositedAmount,omitempty"`
	RefundedAmount  int64  `json:"refundedAmount,omitempty"`
	PaymentState    string `json:"paymentState,omitempty"`
	TotalAmount     int64  `json:"totalAmount,omitempty"`
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

type BonusInfo struct {
	ApprovedAmount   int64  `json:"approvedAmountAward,omitempty"`
	DepositedAmount  int64  `json:"depositedAmountAward,omitempty"`
	RefundedAmount   int64  `json:"refundedAmountAward,omitempty"`
	ApprovedBonus    int64  `json:"approvedAmountBonus,omitempty"`
	DepositedBonus   int64  `json:"depositedAmountBonus,omitempty"`
	RefundedBonus    int64  `json:"refundedAmountBonus,omitempty"`
	PCID             string `json:"pcId,omitempty"`
	Successful       string `json:"successful,omitempty"`
	PaymentOperation string `json:"paymentOperation,omitempty"`
}

type LoyaltyInfo struct {
	AwardBonus   BonusInfo `json:"awardBonus,omitempty"`
	PaymentBonus BonusInfo `json:"paymentBonus,omitempty"`
	LoyaltyName  string    `json:"loyaltyName,omitempty"`
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
	Name   string     `json:"name,omitempty"`
	Params JSONParams `json:"params,omitempty"`
}

type RegisterRequest struct {
	UserName                                string              `form:"userName" json:"userName,omitempty"`
	Password                                string              `form:"password" json:"password,omitempty"`
	Token                                   string              `form:"token" json:"token,omitempty"`
	OrderNumber                             string              `form:"orderNumber" json:"orderNumber,omitempty"`
	Amount                                  int64               `form:"amount" json:"amount,omitempty"`
	Currency                                string              `form:"currency" json:"currency,omitempty"`
	ReturnURL                               string              `form:"returnUrl" json:"returnUrl,omitempty"`
	FailURL                                 string              `form:"failUrl" json:"failUrl,omitempty"`
	DynamicCallbackURL                      string              `form:"dynamicCallbackUrl" json:"dynamicCallbackUrl,omitempty"`
	Description                             string              `form:"description" json:"description,omitempty"`
	Language                                string              `form:"language" json:"language,omitempty"`
	IP                                      string              `form:"ip" json:"ip,omitempty"`
	ClientID                                string              `form:"clientId" json:"clientId,omitempty"`
	MerchantLogin                           string              `form:"merchantLogin" json:"merchantLogin,omitempty"`
	CardholderName                          string              `form:"cardholderName" json:"cardholderName,omitempty"`
	JSONParams                              JSONParams          `form:"jsonParams" json:"jsonParams,omitempty"`
	SessionTimeoutSecs                      int64               `form:"sessionTimeoutSecs" json:"sessionTimeoutSecs,omitempty"`
	ExpirationDate                          string              `form:"expirationDate" json:"expirationDate,omitempty"`
	BindingID                               string              `form:"bindingId" json:"bindingId,omitempty"`
	Features                                string              `form:"features" json:"features,omitempty"`
	PostAddress                             string              `form:"postAddress" json:"postAddress,omitempty"`
	OrderBundle                             OrderBundle         `form:"orderBundle" json:"orderBundle,omitempty"`
	AdditionalOFDParams                     AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	TaxSystem                               int64               `form:"taxSystem" json:"taxSystem,omitempty"`
	FeeInput                                int64               `form:"feeInput" json:"feeInput,omitempty"`
	Email                                   string              `form:"email" json:"email,omitempty"`
	MerchantINN                             string              `form:"merchantInn" json:"merchantInn,omitempty"`
	BillingPayerData                        BillingPayerData    `form:"billingPayerData" json:"billingPayerData,omitempty"`
	ShippingPayerData                       ShippingPayerData   `form:"shippingPayerData" json:"shippingPayerData,omitempty"`
	PreOrderPayerData                       PreOrderPayerData   `form:"preOrderPayerData" json:"preOrderPayerData,omitempty"`
	OrderPayerData                          OrderPayerData      `form:"orderPayerData" json:"orderPayerData,omitempty"`
	BillingAndShippingAddressMatchIndicator string              `form:"billingAndShippingAddressMatchIndicator" json:"billingAndShippingAddressMatchIndicator,omitempty"`
	Values                                  url.Values          `form:"-" json:"-"`
}

type RegisterResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	FormURL               string `json:"formUrl,omitempty"`
	OrderID               string `json:"orderId,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type RegisterResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *RegisterResponse
}

type RegisterPreAuthRequest struct {
	UserName                                string              `form:"userName" json:"userName,omitempty"`
	Password                                string              `form:"password" json:"password,omitempty"`
	Token                                   string              `form:"token" json:"token,omitempty"`
	OrderNumber                             string              `form:"orderNumber" json:"orderNumber,omitempty"`
	Amount                                  int64               `form:"amount" json:"amount,omitempty"`
	Currency                                string              `form:"currency" json:"currency,omitempty"`
	ReturnURL                               string              `form:"returnUrl" json:"returnUrl,omitempty"`
	FailURL                                 string              `form:"failUrl" json:"failUrl,omitempty"`
	DynamicCallbackURL                      string              `form:"dynamicCallbackUrl" json:"dynamicCallbackUrl,omitempty"`
	Description                             string              `form:"description" json:"description,omitempty"`
	IP                                      string              `form:"ip" json:"ip,omitempty"`
	Language                                string              `form:"language" json:"language,omitempty"`
	ClientID                                string              `form:"clientId" json:"clientId,omitempty"`
	MerchantLogin                           string              `form:"merchantLogin" json:"merchantLogin,omitempty"`
	CardholderName                          string              `form:"cardholderName" json:"cardholderName,omitempty"`
	JSONParams                              JSONParams          `form:"jsonParams" json:"jsonParams,omitempty"`
	SessionTimeoutSecs                      int64               `form:"sessionTimeoutSecs" json:"sessionTimeoutSecs,omitempty"`
	ExpirationDate                          string              `form:"expirationDate" json:"expirationDate,omitempty"`
	BindingID                               string              `form:"bindingId" json:"bindingId,omitempty"`
	Features                                string              `form:"features" json:"features,omitempty"`
	AutocompletionDate                      string              `form:"autocompletionDate" json:"autocompletionDate,omitempty"`
	AutoReverseDate                         string              `form:"autoReverseDate" json:"autoReverseDate,omitempty"`
	PostAddress                             string              `form:"postAddress" json:"postAddress,omitempty"`
	OrderBundle                             OrderBundle         `form:"orderBundle" json:"orderBundle,omitempty"`
	AdditionalOFDParams                     AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	TaxSystem                               int64               `form:"taxSystem" json:"taxSystem,omitempty"`
	FeeInput                                int64               `form:"feeInput" json:"feeInput,omitempty"`
	Email                                   string              `form:"email" json:"email,omitempty"`
	MerchantINN                             string              `form:"merchantInn" json:"merchantInn,omitempty"`
	BillingPayerData                        BillingPayerData    `form:"billingPayerData" json:"billingPayerData,omitempty"`
	ShippingPayerData                       ShippingPayerData   `form:"shippingPayerData" json:"shippingPayerData,omitempty"`
	PreOrderPayerData                       PreOrderPayerData   `form:"preOrderPayerData" json:"preOrderPayerData,omitempty"`
	OrderPayerData                          OrderPayerData      `form:"orderPayerData" json:"orderPayerData,omitempty"`
	BillingAndShippingAddressMatchIndicator string              `form:"billingAndShippingAddressMatchIndicator" json:"billingAndShippingAddressMatchIndicator,omitempty"`
	Values                                  url.Values          `form:"-" json:"-"`
}

type RegisterPreAuthResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	OrderID               string `json:"orderId,omitempty"`
	FormURL               string `json:"formUrl,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type RegisterPreAuthResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *RegisterPreAuthResponse
}

type PaymentOrderRequest struct {
	UserName                      string            `form:"userName" json:"userName,omitempty"`
	Password                      string            `form:"password" json:"password,omitempty"`
	MDORDER                       string            `form:"MDORDER" json:"MDORDER,omitempty"`
	PAN                           int64             `form:"$PAN" json:"$PAN,omitempty"`
	CVC                           string            `form:"$CVC" json:"$CVC,omitempty"`
	YYYY                          int64             `form:"YYYY" json:"YYYY,omitempty"`
	MM                            int64             `form:"MM" json:"MM,omitempty"`
	EXPIRY                        int64             `form:"$EXPIRY" json:"$EXPIRY,omitempty"`
	SeToken                       string            `form:"seToken" json:"seToken,omitempty"`
	TEXT                          string            `form:"TEXT" json:"TEXT,omitempty"`
	Language                      string            `form:"language" json:"language,omitempty"`
	IP                            string            `form:"ip" json:"ip,omitempty"`
	BindingNotNeeded              bool              `form:"bindingNotNeeded" json:"bindingNotNeeded,omitempty"`
	JSONParams                    JSONParams        `form:"jsonParams" json:"jsonParams,omitempty"`
	ThreeDSSDK                    bool              `form:"threeDSSDK" json:"threeDSSDK,omitempty"`
	Email                         string            `form:"email" json:"email,omitempty"`
	BillingPayerData              BillingPayerData  `form:"billingPayerData" json:"billingPayerData,omitempty"`
	ShippingPayerData             ShippingPayerData `form:"shippingPayerData" json:"shippingPayerData,omitempty"`
	PreOrderPayerData             PreOrderPayerData `form:"preOrderPayerData" json:"preOrderPayerData,omitempty"`
	OrderPayerData                OrderPayerData    `form:"orderPayerData" json:"orderPayerData,omitempty"`
	Tii                           string            `form:"tii" json:"tii,omitempty"`
	ExternalScaExemptionIndicator string            `form:"externalScaExemptionIndicator" json:"externalScaExemptionIndicator,omitempty"`
	ClientBrowserInfo             ClientBrowserInfo `form:"clientBrowserInfo" json:"clientBrowserInfo,omitempty"`
	OriginalPaymentNetRefNum      string            `form:"originalPaymentNetRefNum" json:"originalPaymentNetRefNum,omitempty"`
	OriginalPaymentDate           string            `form:"originalPaymentDate" json:"originalPaymentDate,omitempty"`
	ACSInIFrame                   bool              `form:"acsInIFrame" json:"acsInIFrame,omitempty"`
	SBPSubscriptionToken          string            `form:"sbpSubscriptionToken" json:"sbpSubscriptionToken,omitempty"`
	SBPMemberID                   string            `form:"sbpMemberId" json:"sbpMemberId,omitempty"`
	Values                        url.Values        `form:"-" json:"-"`
}

type PaymentOrderResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Info                  string `json:"info,omitempty"`
	Redirect              string `json:"redirect,omitempty"`
	TermURL               string `json:"termUrl,omitempty"`
	ACSURL                string `json:"acsUrl,omitempty"`
	PaReq                 string `json:"paReq,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type PaymentOrderResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *PaymentOrderResponse
}

type InstantPaymentRequest struct {
	UserName                                string              `form:"userName" json:"userName,omitempty"`
	Password                                string              `form:"password" json:"password,omitempty"`
	Token                                   string              `form:"token" json:"token,omitempty"`
	Amount                                  int64               `form:"amount" json:"amount,omitempty"`
	Currency                                string              `form:"currency" json:"currency,omitempty"`
	ClientID                                string              `form:"clientId" json:"clientId,omitempty"`
	IP                                      string              `form:"ip" json:"ip,omitempty"`
	BindingNotNeeded                        bool                `form:"bindingNotNeeded" json:"bindingNotNeeded,omitempty"`
	OrderNumber                             string              `form:"orderNumber" json:"orderNumber,omitempty"`
	Description                             string              `form:"description" json:"description,omitempty"`
	Language                                string              `form:"language" json:"language,omitempty"`
	BindingID                               string              `form:"bindingId" json:"bindingId,omitempty"`
	PreAuth                                 bool                `form:"preAuth" json:"preAuth,omitempty"`
	Pan                                     string              `form:"pan" json:"pan,omitempty"`
	Cvc                                     string              `form:"cvc" json:"cvc,omitempty"`
	CardHolderName                          string              `form:"cardHolderName" json:"cardHolderName,omitempty"`
	MerchantLogin                           string              `form:"merchantLogin" json:"merchantLogin,omitempty"`
	SessionTimeoutSecs                      int64               `form:"sessionTimeoutSecs" json:"sessionTimeoutSecs,omitempty"`
	AutocompletionDate                      string              `form:"autocompletionDate" json:"autocompletionDate,omitempty"`
	AutoReverseDate                         string              `form:"autoReverseDate" json:"autoReverseDate,omitempty"`
	ExpirationDate                          string              `form:"expirationDate" json:"expirationDate,omitempty"`
	SeToken                                 string              `form:"seToken" json:"seToken,omitempty"`
	BackURL                                 string              `form:"backUrl" json:"backUrl,omitempty"`
	FailURL                                 string              `form:"failUrl" json:"failUrl,omitempty"`
	JSONParams                              JSONParams          `form:"jsonParams" json:"jsonParams,omitempty"`
	Features                                string              `form:"features" json:"features,omitempty"`
	OrderBundle                             OrderBundle         `form:"orderBundle" json:"orderBundle,omitempty"`
	AdditionalOFDParams                     AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	DynamicCallbackURL                      string              `form:"dynamicCallbackUrl" json:"dynamicCallbackUrl,omitempty"`
	ThreeDSServerTransID                    string              `form:"threeDSServerTransId" json:"threeDSServerTransId,omitempty"`
	ThreeDSVer2FinishURL                    string              `form:"threeDSVer2FinishUrl" json:"threeDSVer2FinishUrl,omitempty"`
	ThreeDSMethodNotificationURL            string              `form:"threeDSMethodNotificationUrl" json:"threeDSMethodNotificationUrl,omitempty"`
	ThreeDSVer2MdOrder                      string              `form:"threeDSVer2MdOrder" json:"threeDSVer2MdOrder,omitempty"`
	ThreeDSSDK                              bool                `form:"threeDSSDK" json:"threeDSSDK,omitempty"`
	ThreeDSProtocolVersion                  string              `form:"threeDSProtocolVersion" json:"threeDSProtocolVersion,omitempty"`
	Expiry                                  int64               `form:"expiry" json:"expiry,omitempty"`
	Email                                   string              `form:"email" json:"email,omitempty"`
	Tii                                     string              `form:"tii" json:"tii,omitempty"`
	OriginalPaymentNetRefNum                string              `form:"originalPaymentNetRefNum" json:"originalPaymentNetRefNum,omitempty"`
	OriginalPaymentDate                     string              `form:"originalPaymentDate" json:"originalPaymentDate,omitempty"`
	ExternalScaExemptionIndicator           string              `form:"externalScaExemptionIndicator" json:"externalScaExemptionIndicator,omitempty"`
	BillingPayerData                        BillingPayerData    `form:"billingPayerData" json:"billingPayerData,omitempty"`
	ShippingPayerData                       ShippingPayerData   `form:"shippingPayerData" json:"shippingPayerData,omitempty"`
	PreOrderPayerData                       PreOrderPayerData   `form:"preOrderPayerData" json:"preOrderPayerData,omitempty"`
	OrderPayerData                          OrderPayerData      `form:"orderPayerData" json:"orderPayerData,omitempty"`
	BillingAndShippingAddressMatchIndicator string              `form:"billingAndShippingAddressMatchIndicator" json:"billingAndShippingAddressMatchIndicator,omitempty"`
	SBPTemplateID                           string              `form:"sbpTemplateId" json:"sbpTemplateId,omitempty"`
	Values                                  url.Values          `form:"-" json:"-"`
}

type InstantPaymentResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	Error                 string `json:"error,omitempty"`
	OrderID               string `json:"orderId,omitempty"`
	Info                  string `json:"info,omitempty"`
	Redirect              string `json:"redirect,omitempty"`
	TermURL               string `json:"termUrl,omitempty"`
	ACSURL                string `json:"acsUrl,omitempty"`
	PaReq                 string `json:"paReq,omitempty"`
	OrderStatus           int64  `json:"orderStatus,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type InstantPaymentResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *InstantPaymentResponse
}

type GetOrderStatusExtendedRequest struct {
	UserName      string     `form:"userName" json:"userName,omitempty"`
	Password      string     `form:"password" json:"password,omitempty"`
	Token         string     `form:"token" json:"token,omitempty"`
	OrderID       string     `form:"orderId" json:"orderId,omitempty"`
	OrderNumber   string     `form:"orderNumber" json:"orderNumber,omitempty"`
	Language      string     `form:"language" json:"language,omitempty"`
	MerchantLogin string     `form:"merchantLogin" json:"merchantLogin,omitempty"`
	Values        url.Values `form:"-" json:"-"`
}

type GetOrderStatusExtendedResponse struct {
	ErrorCode             string                 `json:"errorCode,omitempty"`
	ErrorMessage          string                 `json:"errorMessage,omitempty"`
	OrderNumber           string                 `json:"orderNumber,omitempty"`
	OrderStatus           int64                  `json:"orderStatus,omitempty"`
	ActionCode            string                 `json:"actionCode,omitempty"`
	ActionCodeDescription string                 `json:"actionCodeDescription,omitempty"`
	Amount                int64                  `json:"amount,omitempty"`
	Currency              string                 `json:"currency,omitempty"`
	Date                  int64                  `json:"date,omitempty"`
	DepositedDate         int64                  `json:"depositedDate,omitempty"`
	OrderDescription      string                 `json:"orderDescription,omitempty"`
	IP                    string                 `json:"ip,omitempty"`
	AuthRefNum            string                 `json:"authRefNum,omitempty"`
	RefundedDate          int64                  `json:"refundedDate,omitempty"`
	ReversedDate          int64                  `json:"reversedDate,omitempty"`
	PaymentWay            string                 `json:"paymentWay,omitempty"`
	AVSCode               string                 `json:"avsCode,omitempty"`
	Chargeback            bool                   `json:"chargeback,omitempty"`
	AuthDateTime          int64                  `json:"authDateTime,omitempty"`
	TerminalID            string                 `json:"terminalId,omitempty"`
	OrderBundle           OrderBundle            `json:"orderBundle,omitempty"`
	PaymentAmountInfo     PaymentAmountInfo      `json:"paymentAmountInfo,omitempty"`
	Refunds               []Refund               `json:"refunds,omitempty"`
	CardAuthInfo          CardAuthInfo           `json:"cardAuthInfo,omitempty"`
	TransactionAttributes []TransactionAttribute `json:"transactionAttributes,omitempty"`
	PrepaymentMDOrder     string                 `json:"prepaymentMdOrder,omitempty"`
	PartpaymentMDOrders   []string               `json:"partpaymentMdOrders,omitempty"`
	FeUtrnno              int64                  `json:"feUtrnno,omitempty"`
	BindingInfo           BindingInfo            `json:"bindingInfo,omitempty"`
	EfectyOrderInfo       EfectyOrderInfo        `json:"efectyOrderInfo,omitempty"`
	PluginInfo            PluginInfo             `json:"pluginInfo,omitempty"`
	DisplayErrorMessage   string                 `json:"displayErrorMessage,omitempty"`
	TII                   string                 `json:"tii,omitempty"`
	UsedPSDIndicatorValue string                 `json:"usedPsdIndicatorValue,omitempty"`
	LoyaltyInfo           LoyaltyInfo            `json:"loyaltyInfo,omitempty"`
	OFDOrderBundle        []OFDOrderBundleItem   `json:"ofdOrderBundle,omitempty"`
	PayerData             PayerData              `json:"payerData,omitempty"`
	MerchantOrderParams   []NameValuePair        `json:"merchantOrderParams,omitempty"`
	Attributes            []NameValuePair        `json:"attributes,omitempty"`
	BankInfo              BankInfo               `json:"bankInfo,omitempty"`
}

type GetOrderStatusExtendedResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *GetOrderStatusExtendedResponse
}

type DepositRequest struct {
	UserName     string       `form:"userName" json:"userName,omitempty"`
	Password     string       `form:"password" json:"password,omitempty"`
	OrderID      string       `form:"orderId" json:"orderId,omitempty"`
	Amount       string       `form:"amount" json:"amount,omitempty"`
	DepositItems DepositItems `form:"depositItems" json:"depositItems,omitempty"`
	Language     string       `form:"language" json:"language,omitempty"`
	Currency     string       `form:"currency" json:"currency,omitempty"`
	JSONParams   JSONParams   `form:"jsonParams" json:"jsonParams,omitempty"`
	Values       url.Values   `form:"-" json:"-"`
}

type DepositResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type DepositResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *DepositResponse
}

type ReverseRequest struct {
	UserName      string     `form:"userName" json:"userName,omitempty"`
	Password      string     `form:"password" json:"password,omitempty"`
	OrderID       string     `form:"orderId" json:"orderId,omitempty"`
	OrderNumber   string     `form:"orderNumber" json:"orderNumber,omitempty"`
	MerchantLogin string     `form:"merchantLogin" json:"merchantLogin,omitempty"`
	Language      string     `form:"language" json:"language,omitempty"`
	JSONParams    string     `form:"jsonParams" json:"jsonParams,omitempty"`
	Amount        string     `form:"amount" json:"amount,omitempty"`
	Currency      string     `form:"currency" json:"currency,omitempty"`
	Values        url.Values `form:"-" json:"-"`
}

type ReverseResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type ReverseResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *ReverseResponse
}

type RefundRequest struct {
	UserName                string              `form:"userName" json:"userName,omitempty"`
	Password                string              `form:"password" json:"password,omitempty"`
	OrderID                 string              `form:"orderId" json:"orderId,omitempty"`
	Amount                  string              `form:"amount" json:"amount,omitempty"`
	Language                string              `form:"language" json:"language,omitempty"`
	JSONParams              string              `form:"jsonParams" json:"jsonParams,omitempty"`
	ExpectedDepositedAmount int64               `form:"expectedDepositedAmount" json:"expectedDepositedAmount,omitempty"`
	ExternalRefundID        string              `form:"externalRefundId" json:"externalRefundId,omitempty"`
	Currency                string              `form:"currency" json:"currency,omitempty"`
	RefundItems             RefundItems         `form:"refundItems" json:"refundItems,omitempty"`
	AdditionalOFDParams     AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	Values                  url.Values          `form:"-" json:"-"`
}

type RefundResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type RefundResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *RefundResponse
}

type DeclineRequest struct {
	UserName      string     `form:"userName" json:"userName,omitempty"`
	Password      string     `form:"password" json:"password,omitempty"`
	MerchantLogin string     `form:"merchantLogin" json:"merchantLogin,omitempty"`
	Language      string     `form:"language" json:"language,omitempty"`
	OrderID       string     `form:"orderId" json:"orderId,omitempty"`
	OrderNumber   string     `form:"orderNumber" json:"orderNumber,omitempty"`
	Values        url.Values `form:"-" json:"-"`
}

type DeclineResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type DeclineResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *DeclineResponse
}

type ProcessRawSumRefundRequest struct {
	UserName            string              `form:"userName" json:"userName,omitempty"`
	Password            string              `form:"password" json:"password,omitempty"`
	OrderID             string              `form:"orderId" json:"orderId,omitempty"`
	Language            string              `form:"language" json:"language,omitempty"`
	Amount              string              `form:"amount" json:"amount,omitempty"`
	JSONParams          string              `form:"jsonParams" json:"jsonParams,omitempty"`
	AdditionalOFDParams AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	OFDParams           OFDParams           `form:"ofdParams" json:"ofdParams,omitempty"`
	Values              url.Values          `form:"-" json:"-"`
}

type ProcessRawSumRefundResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type ProcessRawSumRefundResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *ProcessRawSumRefundResponse
}

type ProcessRawPositionRefundRequest struct {
	UserName            string              `form:"userName" json:"userName,omitempty"`
	Password            string              `form:"password" json:"password,omitempty"`
	OrderID             string              `form:"orderId" json:"orderId,omitempty"`
	Language            string              `form:"language" json:"language,omitempty"`
	Amount              string              `form:"amount" json:"amount,omitempty"`
	PositionID          int64               `form:"positionId" json:"positionId,omitempty"`
	AdditionalOFDParams AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	Values              url.Values          `form:"-" json:"-"`
}

type ProcessRawPositionRefundResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type ProcessRawPositionRefundResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *ProcessRawPositionRefundResponse
}

type ProcessRawPositionOrderRefundRequest struct {
	UserName            string              `form:"userName" json:"userName,omitempty"`
	Password            string              `form:"password" json:"password,omitempty"`
	OrderID             string              `form:"orderId" json:"orderId,omitempty"`
	Language            string              `form:"language" json:"language,omitempty"`
	Amount              string              `form:"amount" json:"amount,omitempty"`
	PositionID          int64               `form:"positionId" json:"positionId,omitempty"`
	AdditionalOFDParams AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	Values              url.Values          `form:"-" json:"-"`
}

type ProcessRawPositionOrderRefundResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type ProcessRawPositionOrderRefundResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *ProcessRawPositionOrderRefundResponse
}

type PaymentOrderBindingRequest struct {
	UserName                      string     `form:"userName" json:"userName,omitempty"`
	Password                      string     `form:"password" json:"password,omitempty"`
	MdOrder                       string     `form:"mdOrder" json:"mdOrder,omitempty"`
	BindingID                     string     `form:"bindingId" json:"bindingId,omitempty"`
	Language                      string     `form:"language" json:"language,omitempty"`
	IP                            string     `form:"ip" json:"ip,omitempty"`
	Cvc                           string     `form:"cvc" json:"cvc,omitempty"`
	ThreeDSSDK                    bool       `form:"threeDSSDK" json:"threeDSSDK,omitempty"`
	Tii                           string     `form:"tii" json:"tii,omitempty"`
	Email                         string     `form:"email" json:"email,omitempty"`
	ThreeDSProtocolVersion        string     `form:"threeDSProtocolVersion" json:"threeDSProtocolVersion,omitempty"`
	ExternalScaExemptionIndicator string     `form:"externalScaExemptionIndicator" json:"externalScaExemptionIndicator,omitempty"`
	SeToken                       string     `form:"seToken" json:"seToken,omitempty"`
	Values                        url.Values `form:"-" json:"-"`
}

type PaymentOrderBindingResponse struct {
	ErrorCode             string     `json:"errorCode,omitempty"`
	ErrorMessage          string     `json:"errorMessage,omitempty"`
	Redirect              string     `json:"redirect,omitempty"`
	Info                  string     `json:"info,omitempty"`
	Error                 string     `json:"error,omitempty"`
	ProcessingErrorType   string     `json:"processingErrorType,omitempty"`
	DisplayErrorMessage   string     `json:"displayErrorMessage,omitempty"`
	ErrorTypeName         string     `json:"errorTypeName,omitempty"`
	ACSURL                string     `json:"acsUrl,omitempty"`
	PaReq                 string     `json:"paReq,omitempty"`
	TermURL               string     `json:"termUrl,omitempty"`
	BindingID             string     `json:"bindingId,omitempty"`
	SBPC2bInfo            SBPC2BInfo `json:"sbpC2bInfo,omitempty"`
	ActionCode            string     `json:"actionCode,omitempty"`
	ActionCodeDescription string     `json:"actionCodeDescription,omitempty"`
}

type PaymentOrderBindingResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *PaymentOrderBindingResponse
}

type GetBindingsRequest struct {
	ClientID      string     `form:"clientId" json:"clientId,omitempty"`
	Language      string     `form:"language" json:"language,omitempty"`
	UserName      string     `form:"userName" json:"userName,omitempty"`
	Password      string     `form:"password" json:"password,omitempty"`
	BindingID     string     `form:"bindingId" json:"bindingId,omitempty"`
	BindingType   string     `form:"bindingType" json:"bindingType,omitempty"`
	ShowExpired   bool       `form:"showExpired" json:"showExpired,omitempty"`
	MerchantLogin string     `form:"merchantLogin" json:"merchantLogin,omitempty"`
	Values        url.Values `form:"-" json:"-"`
}

type GetBindingsResponse struct {
	ErrorCode             string    `json:"errorCode,omitempty"`
	ErrorMessage          string    `json:"errorMessage,omitempty"`
	Bindings              []Binding `json:"bindings,omitempty"`
	ActionCode            string    `json:"actionCode,omitempty"`
	ActionCodeDescription string    `json:"actionCodeDescription,omitempty"`
}

type GetBindingsResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *GetBindingsResponse
}

type GetBindingsByCardOrIDRequest struct {
	UserName    string     `form:"userName" json:"userName,omitempty"`
	Password    string     `form:"password" json:"password,omitempty"`
	Pan         string     `form:"pan" json:"pan,omitempty"`
	BindingID   string     `form:"bindingId" json:"bindingId,omitempty"`
	ShowExpired bool       `form:"showExpired" json:"showExpired,omitempty"`
	Values      url.Values `form:"-" json:"-"`
}

type GetBindingsByCardOrIDResponse struct {
	ErrorCode             string    `json:"errorCode,omitempty"`
	ErrorMessage          string    `json:"errorMessage,omitempty"`
	Bindings              []Binding `json:"bindings,omitempty"`
	BindingID             string    `json:"bindingId,omitempty"`
	MaskedPan             string    `json:"maskedPan,omitempty"`
	ExpiryDate            string    `json:"expiryDate,omitempty"`
	ClientID              string    `json:"clientId,omitempty"`
	ActionCode            string    `json:"actionCode,omitempty"`
	ActionCodeDescription string    `json:"actionCodeDescription,omitempty"`
}

type GetBindingsByCardOrIDResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *GetBindingsByCardOrIDResponse
}

type UnBindCardRequest struct {
	UserName  string     `form:"userName" json:"userName,omitempty"`
	Password  string     `form:"password" json:"password,omitempty"`
	BindingID string     `form:"bindingId" json:"bindingId,omitempty"`
	Values    url.Values `form:"-" json:"-"`
}

type UnBindCardResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type UnBindCardResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *UnBindCardResponse
}

type BindCardRequest struct {
	UserName  string     `form:"userName" json:"userName,omitempty"`
	Password  string     `form:"password" json:"password,omitempty"`
	BindingID string     `form:"bindingId" json:"bindingId,omitempty"`
	Values    url.Values `form:"-" json:"-"`
}

type BindCardResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type BindCardResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *BindCardResponse
}

type ExtendBindingRequest struct {
	UserName  string     `form:"userName" json:"userName,omitempty"`
	Password  string     `form:"password" json:"password,omitempty"`
	BindingID string     `form:"bindingId" json:"bindingId,omitempty"`
	NewExpiry int64      `form:"newExpiry" json:"newExpiry,omitempty"`
	Language  string     `form:"language" json:"language,omitempty"`
	Values    url.Values `form:"-" json:"-"`
}

type ExtendBindingResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type ExtendBindingResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *ExtendBindingResponse
}

type CreateBindingNoPaymentRequest struct {
	UserName             string               `form:"userName" json:"userName,omitempty"`
	Password             string               `form:"password" json:"password,omitempty"`
	ClientID             string               `form:"clientId" json:"clientId,omitempty"`
	CardholderName       string               `form:"cardholderName" json:"cardholderName,omitempty"`
	ExpiryDate           string               `form:"expiryDate" json:"expiryDate,omitempty"`
	Pan                  string               `form:"pan" json:"pan,omitempty"`
	AdditionalParameters AdditionalParameters `form:"additionalParameters" json:"additionalParameters,omitempty"`
	MerchantLogin        string               `form:"merchantLogin" json:"merchantLogin,omitempty"`
	Email                string               `form:"email" json:"email,omitempty"`
	Phone                string               `form:"phone" json:"phone,omitempty"`
	Values               url.Values           `form:"-" json:"-"`
}

type CreateBindingNoPaymentResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Error                 bool   `json:"error,omitempty"`
	BindingID             string `json:"bindingId,omitempty"`
	ClientID              string `json:"clientId,omitempty"`
	CardholderName        string `json:"cardholderName,omitempty"`
	ExpiryDate            string `json:"expiryDate,omitempty"`
	MaskedPan             string `json:"maskedPan,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type CreateBindingNoPaymentResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *CreateBindingNoPaymentResponse
}

type SBPC2BQRDynamicGetRequest struct {
	UserName                string     `form:"userName" json:"userName,omitempty"`
	Password                string     `form:"password" json:"password,omitempty"`
	MdOrder                 string     `form:"mdOrder" json:"mdOrder,omitempty"`
	Account                 string     `form:"account" json:"account,omitempty"`
	MemberID                string     `form:"memberId" json:"memberId,omitempty"`
	TspMerchantID           string     `form:"tspMerchantId" json:"tspMerchantId,omitempty"`
	PaymentPurpose          string     `form:"paymentPurpose" json:"paymentPurpose,omitempty"`
	RedirectURL             string     `form:"redirectUrl" json:"redirectUrl,omitempty"`
	QRHeight                string     `form:"qrHeight" json:"qrHeight,omitempty"`
	QRWidth                 string     `form:"qrWidth" json:"qrWidth,omitempty"`
	QRFormat                string     `form:"qrFormat" json:"qrFormat,omitempty"`
	CreateSubscription      bool       `form:"createSubscription" json:"createSubscription,omitempty"`
	SubscriptionServiceName string     `form:"subscriptionServiceName" json:"subscriptionServiceName,omitempty"`
	SubscriptionServiceID   string     `form:"subscriptionServiceId" json:"subscriptionServiceId,omitempty"`
	Values                  url.Values `form:"-" json:"-"`
}

type SBPC2BQRDynamicGetResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Payload               string `json:"payload,omitempty"`
	QRID                  string `json:"qrId,omitempty"`
	Status                string `json:"status,omitempty"`
	RenderedQR            string `json:"renderedQr,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type SBPC2BQRDynamicGetResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *SBPC2BQRDynamicGetResponse
}

type SBPC2BQRStatusRequest struct {
	UserName string     `form:"userName" json:"userName,omitempty"`
	Password string     `form:"password" json:"password,omitempty"`
	MdOrder  string     `form:"mdOrder" json:"mdOrder,omitempty"`
	QRID     string     `form:"qrId" json:"qrId,omitempty"`
	Values   url.Values `form:"-" json:"-"`
}

type SBPC2BQRStatusResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Status                string `json:"status,omitempty"`
	QRType                string `json:"qrType,omitempty"`
	TransactionState      string `json:"transactionState,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type SBPC2BQRStatusResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *SBPC2BQRStatusResponse
}

type SBPC2BQRDynamicRejectRequest struct {
	UserName string     `form:"userName" json:"userName,omitempty"`
	Password string     `form:"password" json:"password,omitempty"`
	MdOrder  string     `form:"mdOrder" json:"mdOrder,omitempty"`
	QRID     string     `form:"qrId" json:"qrId,omitempty"`
	Values   url.Values `form:"-" json:"-"`
}

type SBPC2BQRDynamicRejectResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Rejected              bool   `json:"rejected,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type SBPC2BQRDynamicRejectResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *SBPC2BQRDynamicRejectResponse
}

type TemplatesCreateTemplateRequest struct {
	UserName            string     `form:"userName" json:"userName,omitempty"`
	Password            string     `form:"password" json:"password,omitempty"`
	Name                string     `form:"name" json:"name,omitempty"`
	StartDate           string     `form:"startDate" json:"startDate,omitempty"`
	EndDate             string     `form:"endDate" json:"endDate,omitempty"`
	Type                string     `form:"type" json:"type,omitempty"`
	Amount              int64      `form:"amount" json:"amount,omitempty"`
	Currency            string     `form:"currency" json:"currency,omitempty"`
	DistributionChannel string     `form:"distributionChannel" json:"distributionChannel,omitempty"`
	QRTemplate          QRTemplate `form:"qrTemplate" json:"qrTemplate,omitempty"`
	Values              url.Values `form:"-" json:"-"`
}

type TemplatesCreateTemplateResponse struct {
	Name                  string     `json:"name,omitempty"`
	Type                  string     `json:"type,omitempty"`
	TemplateID            string     `json:"templateId,omitempty"`
	Status                string     `json:"status,omitempty"`
	Amount                int64      `json:"amount,omitempty"`
	Currency              string     `json:"currency,omitempty"`
	DistributionChannel   string     `json:"distributionChannel,omitempty"`
	StartDate             string     `json:"startDate,omitempty"`
	EndDate               string     `json:"endDate,omitempty"`
	QRTemplate            QRTemplate `json:"qrTemplate,omitempty"`
	ErrorCode             string     `json:"errorCode,omitempty"`
	ErrorMessage          string     `json:"errorMessage,omitempty"`
	ActionCode            string     `json:"actionCode,omitempty"`
	ActionCodeDescription string     `json:"actionCodeDescription,omitempty"`
}

type TemplatesCreateTemplateResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *TemplatesCreateTemplateResponse
}

type TemplatesGetTemplateDetailsRequest struct {
	UserName   string     `form:"userName" json:"userName,omitempty"`
	Password   string     `form:"password" json:"password,omitempty"`
	TemplateID string     `form:"templateId" json:"templateId,omitempty"`
	Values     url.Values `form:"-" json:"-"`
}

type TemplatesGetTemplateDetailsResponse struct {
	ErrorCode             string     `json:"errorCode,omitempty"`
	ErrorMessage          string     `json:"errorMessage,omitempty"`
	Amount                int64      `json:"amount,omitempty"`
	Currency              string     `json:"currency,omitempty"`
	DistributionChannel   string     `json:"distributionChannel,omitempty"`
	EndDate               string     `json:"endDate,omitempty"`
	Name                  string     `json:"name,omitempty"`
	StartDate             string     `json:"startDate,omitempty"`
	Status                string     `json:"status,omitempty"`
	TemplateID            string     `json:"templateId,omitempty"`
	Type                  string     `json:"type,omitempty"`
	UseDate               string     `json:"useDate,omitempty"`
	QRTemplate            QRTemplate `json:"qrTemplate,omitempty"`
	ActionCode            string     `json:"actionCode,omitempty"`
	ActionCodeDescription string     `json:"actionCodeDescription,omitempty"`
}

type TemplatesGetTemplateDetailsResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *TemplatesGetTemplateDetailsResponse
}

type TemplatesUpdateTemplateRequest struct {
	UserName   string     `form:"userName" json:"userName,omitempty"`
	Password   string     `form:"password" json:"password,omitempty"`
	Name       string     `form:"name" json:"name,omitempty"`
	StartDate  string     `form:"startDate" json:"startDate,omitempty"`
	EndDate    string     `form:"endDate" json:"endDate,omitempty"`
	TemplateID string     `form:"templateId" json:"templateId,omitempty"`
	Status     string     `form:"status" json:"status,omitempty"`
	Values     url.Values `form:"-" json:"-"`
}

type TemplatesUpdateTemplateResponse struct {
	ErrorCode             string     `json:"errorCode,omitempty"`
	ErrorMessage          string     `json:"errorMessage,omitempty"`
	Name                  string     `json:"name,omitempty"`
	Amount                int64      `json:"amount,omitempty"`
	Currency              string     `json:"currency,omitempty"`
	DistributionChannel   string     `json:"distributionChannel,omitempty"`
	EndDate               string     `json:"endDate,omitempty"`
	StartDate             string     `json:"startDate,omitempty"`
	Status                string     `json:"status,omitempty"`
	TemplateID            string     `json:"templateId,omitempty"`
	Type                  string     `json:"type,omitempty"`
	UseDate               string     `json:"useDate,omitempty"`
	QRTemplate            QRTemplate `json:"qrTemplate,omitempty"`
	TakeTax               bool       `json:"takeTax,omitempty"`
	TotalTaxAmount        int64      `json:"totalTaxAmount,omitempty"`
	ActionCode            string     `json:"actionCode,omitempty"`
	ActionCodeDescription string     `json:"actionCodeDescription,omitempty"`
}

type TemplatesUpdateTemplateResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *TemplatesUpdateTemplateResponse
}

type SBPC2BGetBindingsRequest struct {
	ClientID     string     `form:"clientId" json:"clientId,omitempty"`
	UserName     string     `form:"userName" json:"userName,omitempty"`
	Password     string     `form:"password" json:"password,omitempty"`
	BindingID    string     `form:"bindingId" json:"bindingId,omitempty"`
	ShowDisabled bool       `form:"showDisabled" json:"showDisabled,omitempty"`
	Values       url.Values `form:"-" json:"-"`
}

type SBPC2BGetBindingsResponse struct {
	ErrorCode             string    `json:"errorCode,omitempty"`
	ErrorMessage          string    `json:"errorMessage,omitempty"`
	Bindings              []Binding `json:"bindings,omitempty"`
	ActionCode            string    `json:"actionCode,omitempty"`
	ActionCodeDescription string    `json:"actionCodeDescription,omitempty"`
}

type SBPC2BGetBindingsResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *SBPC2BGetBindingsResponse
}

type SBPC2BUnBindRequest struct {
	UserName  string     `form:"userName" json:"userName,omitempty"`
	Password  string     `form:"password" json:"password,omitempty"`
	BindingID string     `form:"bindingId" json:"bindingId,omitempty"`
	Values    url.Values `form:"-" json:"-"`
}

type SBPC2BUnBindResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	UserMessage           string `json:"userMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type SBPC2BUnBindResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *SBPC2BUnBindResponse
}

type Finish3DSPaymentRequest struct {
	UserName string     `form:"userName" json:"userName,omitempty"`
	Password string     `form:"password" json:"password,omitempty"`
	MdOrder  string     `form:"mdOrder" json:"mdOrder,omitempty"`
	PaRes    string     `form:"paRes" json:"paRes,omitempty"`
	Values   url.Values `form:"-" json:"-"`
}

type Finish3DSPaymentResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	Error                 string `json:"error,omitempty"`
	Redirect              string `json:"redirect,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type Finish3DSPaymentResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *Finish3DSPaymentResponse
}

type Finish3DSVer2PaymentRequest struct {
	UserName             string     `form:"userName" json:"userName,omitempty"`
	Password             string     `form:"password" json:"password,omitempty"`
	ThreeDSServerTransID string     `form:"threeDSServerTransId" json:"threeDSServerTransId,omitempty"`
	ThreeDSVer2MdOrder   string     `form:"threeDSVer2MdOrder" json:"threeDSVer2MdOrder,omitempty"`
	Values               url.Values `form:"-" json:"-"`
}

type Finish3DSVer2PaymentResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Redirect              string `json:"redirect,omitempty"`
	Is3DSVer2             bool   `json:"is3DSVer2,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type Finish3DSVer2PaymentResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *Finish3DSVer2PaymentResponse
}

type ThreeDSContinueRequest struct {
	UserName string     `form:"userName" json:"userName,omitempty"`
	Password string     `form:"password" json:"password,omitempty"`
	Token    string     `form:"token" json:"token,omitempty"`
	MdOrder  string     `form:"mdOrder" json:"mdOrder,omitempty"`
	Values   url.Values `form:"-" json:"-"`
}

type ThreeDSContinueResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Info                  string `json:"info,omitempty"`
	Redirect              string `json:"redirect,omitempty"`
	ACSURL                string `json:"acsUrl,omitempty"`
	PackedCReq            string `json:"packedCReq,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type ThreeDSContinueResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *ThreeDSContinueResponse
}

type VerifyCardRequest struct {
	UserName                                string            `form:"userName" json:"userName,omitempty"`
	Password                                string            `form:"password" json:"password,omitempty"`
	Token                                   string            `form:"token" json:"token,omitempty"`
	Amount                                  int64             `form:"amount" json:"amount,omitempty"`
	Currency                                string            `form:"currency" json:"currency,omitempty"`
	Pan                                     string            `form:"pan" json:"pan,omitempty"`
	Cvc                                     string            `form:"cvc" json:"cvc,omitempty"`
	Expiry                                  int64             `form:"expiry" json:"expiry,omitempty"`
	CardholderName                          string            `form:"cardholderName" json:"cardholderName,omitempty"`
	BackURL                                 string            `form:"backUrl" json:"backUrl,omitempty"`
	FailURL                                 string            `form:"failUrl" json:"failUrl,omitempty"`
	Description                             string            `form:"description" json:"description,omitempty"`
	Language                                string            `form:"language" json:"language,omitempty"`
	ReturnURL                               string            `form:"returnUrl" json:"returnUrl,omitempty"`
	ThreeDSServerTransID                    string            `form:"threeDSServerTransId" json:"threeDSServerTransId,omitempty"`
	ThreeDSVer2FinishURL                    string            `form:"threeDSVer2FinishUrl" json:"threeDSVer2FinishUrl,omitempty"`
	ThreeDSVer2MdOrder                      string            `form:"threeDSVer2MdOrder" json:"threeDSVer2MdOrder,omitempty"`
	ThreeDSSDK                              bool              `form:"threeDSSDK" json:"threeDSSDK,omitempty"`
	BillingPayerData                        BillingPayerData  `form:"billingPayerData" json:"billingPayerData,omitempty"`
	ShippingPayerData                       ShippingPayerData `form:"shippingPayerData" json:"shippingPayerData,omitempty"`
	PreOrderPayerData                       PreOrderPayerData `form:"preOrderPayerData" json:"preOrderPayerData,omitempty"`
	OrderPayerData                          OrderPayerData    `form:"orderPayerData" json:"orderPayerData,omitempty"`
	BillingAndShippingAddressMatchIndicator string            `form:"billingAndShippingAddressMatchIndicator" json:"billingAndShippingAddressMatchIndicator,omitempty"`
	Values                                  url.Values        `form:"-" json:"-"`
}

type VerifyCardResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	OrderID               string `json:"orderId,omitempty"`
	OrderNumber           string `json:"orderNumber,omitempty"`
	AuthCode              int64  `json:"authCode,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
	Time                  int64  `json:"time,omitempty"`
	Eci                   int64  `json:"eci,omitempty"`
	Amount                int64  `json:"amount,omitempty"`
	Currency              string `json:"currency,omitempty"`
	Rrn                   int64  `json:"rrn,omitempty"`
	ACSURL                string `json:"acsUrl,omitempty"`
	TermURL               string `json:"termUrl,omitempty"`
	PaReq                 string `json:"paReq,omitempty"`
}

type VerifyCardResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *VerifyCardResponse
}

type CloseOFDReceiptRequest struct {
	MdOrder             string              `form:"mdOrder" json:"mdOrder,omitempty"`
	OrderNumber         string              `form:"orderNumber" json:"orderNumber,omitempty"`
	UserName            string              `form:"userName" json:"userName,omitempty"`
	Password            string              `form:"password" json:"password,omitempty"`
	Amount              int64               `form:"amount" json:"amount,omitempty"`
	AdditionalOFDParams AdditionalOFDParams `form:"additionalOfdParams" json:"additionalOfdParams,omitempty"`
	MerchantLogin       string              `form:"merchantLogin" json:"merchantLogin,omitempty"`
	OrderBundle         OrderBundle         `form:"orderBundle" json:"orderBundle,omitempty"`
	Values              url.Values          `form:"-" json:"-"`
}

type CloseOFDReceiptResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Success               bool   `json:"success,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type CloseOFDReceiptResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *CloseOFDReceiptResponse
}

type GetReceiptStatusRequest struct {
	UserName    string     `form:"userName" json:"userName,omitempty"`
	Password    string     `form:"password" json:"password,omitempty"`
	OrderID     string     `form:"orderId" json:"orderId,omitempty"`
	OrderNumber string     `form:"orderNumber" json:"orderNumber,omitempty"`
	UUID        string     `form:"uuid" json:"uuid,omitempty"`
	Language    string     `form:"language" json:"language,omitempty"`
	Values      url.Values `form:"-" json:"-"`
}

type GetReceiptStatusResponse struct {
	ErrorCode             string        `json:"errorCode,omitempty"`
	ErrorMessage          string        `json:"errorMessage,omitempty"`
	OrderNumber           string        `json:"orderNumber,omitempty"`
	OrderID               string        `json:"orderId,omitempty"`
	Receipt               []ReceiptInfo `json:"receipt,omitempty"`
	ActionCode            string        `json:"actionCode,omitempty"`
	ActionCodeDescription string        `json:"actionCodeDescription,omitempty"`
}

type GetReceiptStatusResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *GetReceiptStatusResponse
}

type TemplatesCreateRequest struct {
	UserName string     `form:"userName" json:"userName,omitempty"`
	Password string     `form:"password" json:"password,omitempty"`
	Template Template   `form:"template" json:"template,omitempty"`
	Values   url.Values `form:"-" json:"-"`
}

type TemplatesCreateResponse struct {
	Status                string        `json:"status,omitempty"`
	Error                 TemplateError `json:"error,omitempty"`
	Template              Template      `json:"template,omitempty"`
	ErrorCode             string        `json:"errorCode,omitempty"`
	ErrorMessage          string        `json:"errorMessage,omitempty"`
	ActionCode            string        `json:"actionCode,omitempty"`
	ActionCodeDescription string        `json:"actionCodeDescription,omitempty"`
}

type TemplatesCreateResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *TemplatesCreateResponse
}

type TemplatesGetRequest struct {
	UserName string     `form:"userName" json:"userName,omitempty"`
	Password string     `form:"password" json:"password,omitempty"`
	Template Template   `form:"template" json:"template,omitempty"`
	Values   url.Values `form:"-" json:"-"`
}

type TemplatesGetResponse struct {
	Status                string        `json:"status,omitempty"`
	Error                 TemplateError `json:"error,omitempty"`
	Template              Template      `json:"template,omitempty"`
	ErrorCode             string        `json:"errorCode,omitempty"`
	ErrorMessage          string        `json:"errorMessage,omitempty"`
	ActionCode            string        `json:"actionCode,omitempty"`
	ActionCodeDescription string        `json:"actionCodeDescription,omitempty"`
}

type TemplatesGetResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *TemplatesGetResponse
}

type TemplatesUpdateRequest struct {
	UserName string     `form:"userName" json:"userName,omitempty"`
	Password string     `form:"password" json:"password,omitempty"`
	Template Template   `form:"template" json:"template,omitempty"`
	Values   url.Values `form:"-" json:"-"`
}

type TemplatesUpdateResponse struct {
	Status                string        `json:"status,omitempty"`
	Error                 TemplateError `json:"error,omitempty"`
	ErrorCode             string        `json:"errorCode,omitempty"`
	ErrorMessage          string        `json:"errorMessage,omitempty"`
	ActionCode            string        `json:"actionCode,omitempty"`
	ActionCodeDescription string        `json:"actionCodeDescription,omitempty"`
}

type TemplatesUpdateResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *TemplatesUpdateResponse
}

type TemplatesGetListRequest struct {
	UserName      string     `form:"userName" json:"userName,omitempty"`
	Password      string     `form:"password" json:"password,omitempty"`
	MerchantLogin string     `form:"merchant_login" json:"merchant_login,omitempty"`
	Status        string     `form:"status" json:"status,omitempty"`
	QRTemplate    QRTemplate `form:"qrTemplate" json:"qrTemplate,omitempty"`
	Values        url.Values `form:"-" json:"-"`
}

type TemplatesGetListResponse struct {
	Status                string        `json:"status,omitempty"`
	Error                 TemplateError `json:"error,omitempty"`
	Templates             []Template    `json:"templates,omitempty"`
	ErrorCode             string        `json:"errorCode,omitempty"`
	ErrorMessage          string        `json:"errorMessage,omitempty"`
	ActionCode            string        `json:"actionCode,omitempty"`
	ActionCodeDescription string        `json:"actionCodeDescription,omitempty"`
}

type TemplatesGetListResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *TemplatesGetListResponse
}

type Finish3DSRequest struct {
	Values url.Values `form:"-" json:"-"`
}

type Finish3DSResponse struct {
	ErrorCode             string `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ActionCode            string `json:"actionCode,omitempty"`
	ActionCodeDescription string `json:"actionCodeDescription,omitempty"`
}

type Finish3DSResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Response   *Finish3DSResponse
}
