package core

// GetReturnRequestListRequest 查询退货/取消订单记录。
type GetReturnRequestListRequest struct {
	// SearchType 留空按天分页；timeFrame 按分钟查询，不支持分页和 OrderId。
	SearchType    string `json:"searchType,omitempty"`
	CreatedAtFrom string `json:"createdAtFrom"`        // yyyy-MM-dd；分钟查询为 yyyy-MM-ddTHH:mm
	CreatedAtTo   string `json:"createdAtTo"`          // 查询范围最长 31 天
	Status        string `json:"status,omitempty"`     // RU、UC、CC、PR；CANCEL 查询不可用
	CancelType    string `json:"cancelType,omitempty"` // RETURN（默认）或 CANCEL
	NextToken     string `json:"nextToken,omitempty"`
	MaxPerPage    int    `json:"maxPerPage,omitempty"` // 留空使用服务端默认值 50
	OrderId       int64  `json:"orderId,omitempty"`
}

type ReturnRequestListResponse struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	Data      []ReturnRequest `json:"data"`
	NextToken string          `json:"nextToken,omitempty"`
}

type ReturnRequest struct {
	ReceiptId                int64                `json:"receiptId"`
	OrderId                  int64                `json:"orderId"`
	PaymentId                int64                `json:"paymentId"`
	ReceiptType              string               `json:"receiptType"`
	ReceiptStatus            string               `json:"receiptStatus"`
	CreatedAt                string               `json:"createdAt"`
	ModifiedAt               string               `json:"modifiedAt"`
	RequesterName            string               `json:"requesterName"`
	RequesterPhoneNumber     string               `json:"requesterPhoneNumber"`
	RequesterRealPhoneNumber *string              `json:"requesterRealPhoneNumber"`
	RequesterAddress         string               `json:"requesterAddress"`
	RequesterAddressDetail   string               `json:"requesterAddressDetail"`
	RequesterZipCode         string               `json:"requesterZipCode"`
	CancelReasonCategory1    string               `json:"cancelReasonCategory1"`
	CancelReasonCategory2    string               `json:"cancelReasonCategory2"`
	CancelReason             string               `json:"cancelReason"`
	CancelCountSum           int                  `json:"cancelCountSum"`
	ReturnDeliveryId         int64                `json:"returnDeliveryId"`
	ReturnDeliveryType       string               `json:"returnDeliveryType"`
	ReleaseStopStatus        string               `json:"releaseStopStatus"`
	EnclosePrice             ReturnRequestMoney   `json:"enclosePrice"`
	FaultByType              string               `json:"faultByType"`
	PreRefund                bool                 `json:"preRefund"`
	CompleteConfirmType      string               `json:"completeConfirmType"`
	CompleteConfirmDate      string               `json:"completeConfirmDate"` // 未确认时可能为空
	ReturnItems              []ReturnRequestItem  `json:"returnItems"`
	ReturnDeliveryDtos       []ReturnDeliveryInfo `json:"returnDeliveryDtos"`
	ReasonCode               string               `json:"reasonCode"`
	ReasonCodeText           string               `json:"reasonCodeText"`
	ReturnShippingCharge     ReturnRequestMoney   `json:"returnShippingCharge"`
}

type ReturnRequestMoney struct {
	CurrencyCode string `json:"currencyCode"`
	Units        int64  `json:"units"`
	Nanos        int32  `json:"nanos"`
}

type ReturnRequestItem struct {
	VendorItemPackageId   int64  `json:"vendorItemPackageId"`
	VendorItemPackageName string `json:"vendorItemPackageName"`
	VendorItemId          int64  `json:"vendorItemId"`
	VendorItemName        string `json:"vendorItemName"`
	CancelCount           int    `json:"cancelCount"`
	PurchaseCount         int    `json:"purchaseCount"`
	ShipmentBoxId         int64  `json:"shipmentBoxId"`
	SellerProductId       int64  `json:"sellerProductId"`
	SellerProductName     string `json:"sellerProductName"`
	ReleaseStatus         string `json:"releaseStatus"`
	CancelCompleteUser    string `json:"cancelCompleteUser"`
}

type ReturnDeliveryInfo struct {
	DeliveryCompanyCode string  `json:"deliveryCompanyCode"`
	DeliveryInvoiceNo   *string `json:"deliveryInvoiceNo"` // 可能为 null 或空字符串
}
