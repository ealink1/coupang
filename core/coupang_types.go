package core

import "time"

// Order struct represents the order information from Coupang API
type Order struct {
	ShipmentBoxId int64     `json:"shipmentBoxId"`
	OrderId       int64     `json:"orderId"`
	OrderedAt     time.Time `json:"orderedAt"`
	Orderer       struct {
		Name          string      `json:"name"`
		Email         string      `json:"email"`
		SafeNumber    string      `json:"safeNumber"`
		OrdererNumber interface{} `json:"ordererNumber"`
	} `json:"orderer"`
	PaidAt        time.Time `json:"paidAt"`
	Status        string    `json:"status"`
	ShippingPrice struct {
		CurrencyCode string `json:"currencyCode"`
		Units        int    `json:"units"`
		Nanos        int    `json:"nanos"`
	} `json:"shippingPrice"`
	CodDeliveryFee struct {
		CurrencyCode string `json:"currencyCode"`
		Units        int    `json:"units"`
		Nanos        int    `json:"nanos"`
	} `json:"codDeliveryFee"`
	RemotePrice struct {
		CurrencyCode string `json:"currencyCode"`
		Units        int    `json:"units"`
		Nanos        int    `json:"nanos"`
	} `json:"remotePrice"`
	RemoteArea         bool   `json:"remoteArea"`
	ParcelPrintMessage string `json:"parcelPrintMessage"`
	SplitShipping      bool   `json:"splitShipping"`
	AbleSplitShipping  bool   `json:"ableSplitShipping"`
	Receiver           struct {
		Name           string      `json:"name"`
		SafeNumber     string      `json:"safeNumber"`
		ReceiverNumber interface{} `json:"receiverNumber"`
		Addr1          string      `json:"addr1"`
		Addr2          string      `json:"addr2"`
		PostCode       string      `json:"postCode"`
	} `json:"receiver"`
	OrderItems []struct {
		VendorItemPackageId   int    `json:"vendorItemPackageId"`
		VendorItemPackageName string `json:"vendorItemPackageName"`
		ProductId             int64  `json:"productId"`
		VendorItemId          int64  `json:"vendorItemId"`
		VendorItemName        string `json:"vendorItemName"`
		ShippingCount         int    `json:"shippingCount"`
		SalesPrice            struct {
			CurrencyCode string `json:"currencyCode"`
			Units        int    `json:"units"`
			Nanos        int    `json:"nanos"`
		} `json:"salesPrice"`
		OrderPrice struct {
			CurrencyCode string `json:"currencyCode"`
			Units        int    `json:"units"`
			Nanos        int    `json:"nanos"`
		} `json:"orderPrice"`
		DiscountPrice struct {
			CurrencyCode string `json:"currencyCode"`
			Units        int    `json:"units"`
			Nanos        int    `json:"nanos"`
		} `json:"discountPrice"`
		InstantCouponDiscount struct {
			CurrencyCode string `json:"currencyCode"`
			Units        int    `json:"units"`
			Nanos        int    `json:"nanos"`
		} `json:"instantCouponDiscount"`
		DownloadableCouponDiscount struct {
			CurrencyCode string `json:"currencyCode"`
			Units        int    `json:"units"`
			Nanos        int    `json:"nanos"`
		} `json:"downloadableCouponDiscount"`
		CoupangDiscount struct {
			CurrencyCode string `json:"currencyCode"`
			Units        int    `json:"units"`
			Nanos        int    `json:"nanos"`
		} `json:"coupangDiscount"`
		ExternalVendorSkuCode      string      `json:"externalVendorSkuCode"`
		EtcInfoHeader              interface{} `json:"etcInfoHeader"`
		EtcInfoValue               interface{} `json:"etcInfoValue"`
		EtcInfoValues              interface{} `json:"etcInfoValues"`
		SellerProductId            int64       `json:"sellerProductId"`
		SellerProductName          string      `json:"sellerProductName"`
		SellerProductItemName      string      `json:"sellerProductItemName"`
		FirstSellerProductItemName string      `json:"firstSellerProductItemName"`
		CancelCount                int         `json:"cancelCount"`
		HoldCountForCancel         int         `json:"holdCountForCancel"`
		EstimatedShippingDate      string      `json:"estimatedShippingDate"`
		PlannedShippingDate        string      `json:"plannedShippingDate"`
		InvoiceNumberUploadDate    interface{} `json:"invoiceNumberUploadDate"`
		ExtraProperties            struct {
			M1VARIATIONID   string `json:"M1_VARIATION_ID,omitempty"`
			M1ORIGINALPRICE string `json:"M1_ORIGINAL_PRICE,omitempty"`
		} `json:"extraProperties"`
		PricingBadge           bool        `json:"pricingBadge"`
		UsedProduct            bool        `json:"usedProduct"`
		ConfirmDate            interface{} `json:"confirmDate"`
		DeliveryChargeTypeName string      `json:"deliveryChargeTypeName"`
		UpBundleVendorItemId   interface{} `json:"upBundleVendorItemId"`
		UpBundleVendorItemName interface{} `json:"upBundleVendorItemName"`
		UpBundleSize           interface{} `json:"upBundleSize"`
		TaxReceiptInfo         interface{} `json:"taxReceiptInfo"`
		Canceled               bool        `json:"canceled"`
		UpBundleItem           bool        `json:"upBundleItem"`
	} `json:"orderItems"`
	OverseaShippingInfoDto struct {
		PersonalCustomsClearanceCode string `json:"personalCustomsClearanceCode"`
		OrdererSsn                   string `json:"ordererSsn"`
		OrdererPhoneNumber           string `json:"ordererPhoneNumber"`
	} `json:"overseaShippingInfoDto"`
	DeliveryCompanyName string      `json:"deliveryCompanyName"`
	InvoiceNumber       string      `json:"invoiceNumber"`
	InTrasitDateTime    interface{} `json:"inTrasitDateTime"`
	DeliveredDate       interface{} `json:"deliveredDate"`
	Refer               string      `json:"refer"`
	ShipmentType        string      `json:"shipmentType"`
	IsCod               bool        `json:"isCod"`
	ExtraProperties     struct {
		TaxReceiptInfo string `json:"taxReceiptInfo"`
	} `json:"extraProperties"`
}

type OrderListResponse struct {
	Code      int     `json:"code"`
	Message   string  `json:"message"`
	Data      []Order `json:"data"`
	NextToken string  `json:"nextToken,omitempty"`
}

type SingleOrderResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Order  `json:"data"`
}

type SingleOrderListResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    []Order `json:"data"`
}

type GetOrderListRequest struct {
	CreatedAtFrom string `json:"createdAtFrom"`
	CreatedAtTo   string `json:"createdAtTo"`
	Status        string `json:"status"`
	NextToken     string `json:"nextToken"`
	MaxPerPage    int    `json:"maxPerPage"`
	SearchType    string `json:"searchType"`
}

type InventoryItem struct {
	SellerProductItemId int64  `json:"sellerProductItemId"`
	VendorItemId        int64  `json:"vendorItemId"`
	ItemName            string `json:"itemName"`
	ExternalVendorSku   string `json:"externalVendorSku"`
	AmountInStock       int64  `json:"amountInStock"`
	SalePrice           int64  `json:"salePrice"`
	OnSale              bool   `json:"onSale"`
}

type ProductInventory struct {
	SellerProductId    int64           `json:"sellerProductId"`
	SellerProductName  string          `json:"sellerProductName"`
	DisplayProductName string          `json:"displayProductName"`
	GeneralProductName string          `json:"generalProductName"`
	VendorId           string          `json:"vendorId"`
	Items              []InventoryItem `json:"items"`
}

type InventoryResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    []ProductInventory `json:"data"`
}

type BatchInventoryRequest struct {
	SellerProductIds []int64 `json:"sellerProductIds"`
}

type ProductDetailResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TrackingId                interface{} `json:"trackingId"`
		SellerProductId           int64       `json:"sellerProductId"`
		SellerProductName         string      `json:"sellerProductName"`
		DisplayCategoryCode       int         `json:"displayCategoryCode"`
		CategoryId                int         `json:"categoryId"`
		ProductId                 int64       `json:"productId"`
		VendorId                  string      `json:"vendorId"`
		MdId                      string      `json:"mdId"`
		MdName                    string      `json:"mdName"`
		SaleStartedAt             string      `json:"saleStartedAt"`
		SaleEndedAt               string      `json:"saleEndedAt"`
		DisplayProductName        string      `json:"displayProductName"`
		Brand                     string      `json:"brand"`
		GeneralProductName        string      `json:"generalProductName"`
		ProductGroup              string      `json:"productGroup"`
		StatusName                string      `json:"statusName"`
		DeliveryMethod            string      `json:"deliveryMethod"`
		DeliveryCompanyCode       string      `json:"deliveryCompanyCode"`
		DeliveryChargeType        string      `json:"deliveryChargeType"`
		DeliveryCharge            int         `json:"deliveryCharge"`
		FreeShipOverAmount        int         `json:"freeShipOverAmount"`
		DeliveryChargeOnReturn    int         `json:"deliveryChargeOnReturn"`
		DeliverySurcharge         int         `json:"deliverySurcharge"`
		RemoteAreaDeliverable     string      `json:"remoteAreaDeliverable"`
		BundlePackingDelivery     int         `json:"bundlePackingDelivery"`
		UnionDeliveryType         string      `json:"unionDeliveryType"`
		ReturnCenterCode          string      `json:"returnCenterCode"`
		ReturnChargeName          string      `json:"returnChargeName"`
		CompanyContactNumber      string      `json:"companyContactNumber"`
		ReturnZipCode             string      `json:"returnZipCode"`
		ReturnAddress             string      `json:"returnAddress"`
		ReturnAddressDetail       string      `json:"returnAddressDetail"`
		ReturnCharge              int         `json:"returnCharge"`
		ExchangeType              string      `json:"exchangeType"`
		AfterServiceInformation   string      `json:"afterServiceInformation"`
		AfterServiceContactNumber string      `json:"afterServiceContactNumber"`
		OutboundShippingPlaceCode int         `json:"outboundShippingPlaceCode"`
		ContributorType           string      `json:"contributorType"`
		VendorUserId              string      `json:"vendorUserId"`
		Requested                 bool        `json:"requested"`
		Items                     []struct {
			OfferCondition            string      `json:"offerCondition"`
			OfferDescription          string      `json:"offerDescription"`
			SellerProductItemId       int64       `json:"sellerProductItemId"`
			VendorItemId              int64       `json:"vendorItemId"`
			ItemId                    int64       `json:"itemId"`
			ItemName                  string      `json:"itemName"`
			OriginalPrice             int         `json:"originalPrice"`
			SalePrice                 int         `json:"salePrice"`
			SupplyPrice               int         `json:"supplyPrice"`
			MaximumBuyCount           int         `json:"maximumBuyCount"`
			MaximumBuyForPerson       int         `json:"maximumBuyForPerson"`
			OutboundShippingTimeDay   int         `json:"outboundShippingTimeDay"`
			MaximumBuyForPersonPeriod int         `json:"maximumBuyForPersonPeriod"`
			UnitCount                 int         `json:"unitCount"`
			AdultOnly                 string      `json:"adultOnly"`
			FreePriceType             interface{} `json:"freePriceType"`
			TaxType                   string      `json:"taxType"`
			ParallelImported          string      `json:"parallelImported"`
			OverseasPurchased         string      `json:"overseasPurchased"`
			ExternalVendorSku         string      `json:"externalVendorSku"`
			PccNeeded                 bool        `json:"pccNeeded"`
			BestPriceGuaranteed3P     bool        `json:"bestPriceGuaranteed3P"`
			EmptyBarcode              bool        `json:"emptyBarcode"`
			EmptyBarcodeReason        string      `json:"emptyBarcodeReason"`
			Barcode                   string      `json:"barcode"`
			SaleAgentCommission       float64     `json:"saleAgentCommission"`
			ModelNo                   string      `json:"modelNo"`
			Images                    []struct {
				ImageOrder int    `json:"imageOrder"`
				ImageType  string `json:"imageType"`
				CdnPath    string `json:"cdnPath"`
				VendorPath string `json:"vendorPath"`
			} `json:"images"`
			Notices []struct {
				NoticeCategoryName       string `json:"noticeCategoryName"`
				NoticeCategoryDetailName string `json:"noticeCategoryDetailName"`
				Content                  string `json:"content"`
			} `json:"notices"`
			Attributes []struct {
				AttributeTypeName  string `json:"attributeTypeName"`
				AttributeValueName string `json:"attributeValueName"`
				Exposed            string `json:"exposed"`
				Editable           bool   `json:"editable"`
			} `json:"attributes"`
			Contents []struct {
				ContentsType   string `json:"contentsType"`
				ContentDetails []struct {
					Content    string `json:"content"`
					DetailType string `json:"detailType"`
				} `json:"contentDetails"`
			} `json:"contents"`
			Certifications []struct {
				CertificationType        string        `json:"certificationType"`
				CertificationCode        string        `json:"certificationCode"`
				CertificationAttachments []interface{} `json:"certificationAttachments"`
			} `json:"certifications"`
			ExtraProperties struct {
				M1VARIATIONID   string `json:"M1_VARIATION_ID"`
				M1ORIGINALPRICE string `json:"M1_ORIGINAL_PRICE"`
			} `json:"extraProperties"`
			SearchTags      []interface{} `json:"searchTags"`
			IsAutoGenerated bool          `json:"isAutoGenerated"`
		} `json:"items"`
		RequiredDocuments []interface{} `json:"requiredDocuments"`
		ExtraInfoMessage  string        `json:"extraInfoMessage"`
		Manufacture       string        `json:"manufacture"`
		ProductOrigin     string        `json:"productOrigin"`
		RoleCode          int           `json:"roleCode"`
		Status            string        `json:"status"`
		BundleInfo        struct {
			BundleType string `json:"bundleType"`
		} `json:"bundleInfo"`
		MultiShippingInfos []struct {
			VendorInventoryId             int64       `json:"vendorInventoryId"`
			DeliveryChargeType            string      `json:"deliveryChargeType"`
			DeliveryChargeTypeDescription string      `json:"deliveryChargeTypeDescription"`
			DeliveryCompanyCode           string      `json:"deliveryCompanyCode"`
			DeliveryCompanyDescription    string      `json:"deliveryCompanyDescription"`
			DeliveryCharge                int         `json:"deliveryCharge"`
			DeliveryChargeMoney           interface{} `json:"deliveryChargeMoney"`
			FreeShipOverAmount            int         `json:"freeShipOverAmount"`
			FreeShipOverAmountMoney       interface{} `json:"freeShipOverAmountMoney"`
			DeliveryChargeOnReturn        int         `json:"deliveryChargeOnReturn"`
			DeliveryChargeOnReturnMoney   interface{} `json:"deliveryChargeOnReturnMoney"`
			DeliverySurcharge             int         `json:"deliverySurcharge"`
			DeliverySurchargeMoney        interface{} `json:"deliverySurchargeMoney"`
			BundlePackingDelivery         int         `json:"bundlePackingDelivery"`
			DeliveryMethod                string      `json:"deliveryMethod"`
			DeliveryMethodDescription     string      `json:"deliveryMethodDescription"`
			DeliveryType                  string      `json:"deliveryType"`
			DeliveryTypeDescription       string      `json:"deliveryTypeDescription"`
			OutboundShippingPlaceId       int         `json:"outboundShippingPlaceId"`
			RemoteAreaDeliverable         interface{} `json:"remoteAreaDeliverable"`
			DeliveryCompanyType           string      `json:"deliveryCompanyType"`
			OutboundShippingTime          int         `json:"outboundShippingTime"`
			DeliveryCompanyPersistCode    string      `json:"deliveryCompanyPersistCode"`
			ShippingInfoMetaData          struct {
			} `json:"shippingInfoMetaData"`
		} `json:"multiShippingInfos"`
		MultiReturnInfos []struct {
			PickUpBranchId            *int        `json:"pickUpBranchId"`
			PickUpBranchGroupCode     string      `json:"pickUpBranchGroupCode"`
			PickUpBranchType          string      `json:"pickUpBranchType"`
			SellerProductId           int64       `json:"sellerProductId"`
			ReturnCenterCode          string      `json:"returnCenterCode"`
			ReturnChargeName          string      `json:"returnChargeName"`
			CompanyContactNumber      string      `json:"companyContactNumber"`
			ReturnZipCode             *string     `json:"returnZipCode"`
			ReturnAddress             string      `json:"returnAddress"`
			ReturnAddressDetail       string      `json:"returnAddressDetail"`
			ReturnCharge              int         `json:"returnCharge"`
			ReturnChargeMoney         interface{} `json:"returnChargeMoney"`
			ExchangeType              interface{} `json:"exchangeType"`
			ReturnChargeVendor        string      `json:"returnChargeVendor"`
			AfterServiceInformation   interface{} `json:"afterServiceInformation"`
			AfterServiceContactNumber interface{} `json:"afterServiceContactNumber"`
		} `json:"multiReturnInfos"`
	} `json:"data"`
}
