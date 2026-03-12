package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCoupang(t *testing.T) {
	accessKey := os.Getenv("COUPANG_ACCESS_KEY")
	secretKey := os.Getenv("COUPANG_SECRET_KEY")
	vendorID := os.Getenv("COUPANG_VENDOR_ID")
	if accessKey == "" || secretKey == "" || vendorID == "" {
		t.Skip("set COUPANG_ACCESS_KEY/COUPANG_SECRET_KEY/COUPANG_VENDOR_ID to run integration test")
	}

	client := NewCoupangClient(accessKey, secretKey, vendorID)

	//orderId := "111100001543008"
	//resp, err := client.GetOrderByOrderId(orderId)
	//if err != nil {
	//	t.Logf("GetOrderByOrderId failed: %v", err)
	//} else {
	//	fmt.Printf("Order Details: %+v\n", resp)
	//}

	// Example: Query last hour
	//now := time.Now()
	//// Ref 2: yyyy-mm-ddT01:00%2B09:00
	//from := now.Add(-1*time.Hour).Format("2006-01-02T15:04:05") + "+09:00"
	//to := now.Format("2006-01-02T15:04:05") + "+09:00"
	//
	//resp, err := client.GetOrderListMinutely(from, to, "ACCEPT")
	//if err != nil {
	//	t.Logf("GetOrderListMinutely failed: %v", err)
	//} else {
	//	fmt.Printf("Minutely Orders: %+v\n", resp)
	//}

	now := time.Now()

	from := now.AddDate(0, 0, -7).Format("2006-01-02") + " 09:00"
	to := now.Format("2006-01-02") + " 09:00"

	req := &GetOrderListRequest{
		CreatedAtFrom: from,
		CreatedAtTo:   to,
		Status:        "",
		MaxPerPage:    50,
		NextToken:     "",
	}
	resp, err := client.GetOrderListDaily(context.Background(), req)
	if err != nil {
		t.Logf("GetOrderListDaily failed (might be due to invalid keys/vendorID): %v", err)
	} else {
		fmt.Printf("Daily Orders: %+v\n", resp)
	}

	//shipmentBoxId := "626887864795136"
	//resp, err := client.GetOrderByShipmentBoxId(shipmentBoxId)
	//if err != nil {
	//	t.Logf("GetOrderByShipmentBoxId failed: %v", err)
	//} else {
	//	fmt.Printf("ShipmentBox Order: %+v\n", resp)
	//}
}

func TestReturnShippingCentersResponseUnmarshal(t *testing.T) {
	raw := `{
  "code": 200,
  "message": "SUCCESS",
  "data": {
    "content": [
      {
        "vendorId": "A00012345",
        "returnCenterCode": "1000557004",
        "shippingPlaceName": "32777 R",
        "deliverCode": "HANJIN",
        "deliverName": "한진택배",
        "goodsflowStatus": "승인",
        "errorMessage": "",
        "createdAt": 1581195908000,
        "vendorCreditFee02kg": null,
        "vendorCreditFee05kg": 3000,
        "vendorCreditFee10kg": 3000,
        "vendorCreditFee20kg": 3000,
        "vendorCashFee02kg": null,
        "vendorCashFee05kg": 3000,
        "vendorCashFee10kg": 3000,
        "vendorCashFee20kg": 3000,
        "consumerCashFee02kg": null,
        "consumerCashFee05kg": 3000,
        "consumerCashFee10kg": 3000,
        "consumerCashFee20kg": 3000,
        "returnFee02kg": null,
        "returnFee05kg": 3000,
        "returnFee10kg": 3000,
        "returnFee20kg": 3000,
        "usable": true,
        "placeAddresses": [
          {
            "addressType": "JIBUN",
            "countryCode": "KR",
            "companyContactNumber": "02-111-1111",
            "phoneNumber2": "000-0000-0000",
            "returnZipCode": "42701",
            "returnAddress": "대구시 달서구 성서공단로",
            "returnAddressDetail": "235"
          }
        ]
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 749,
      "totalElements": 7488,
      "countPerPage": 10
    }
  }
}`

	var resp ReturnShippingCentersResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("unexpected code: %d", resp.Code)
	}
	if resp.Data.Pagination.CurrentPage != 1 || resp.Data.Pagination.CountPerPage != 10 {
		t.Fatalf("unexpected pagination: %+v", resp.Data.Pagination)
	}
	if len(resp.Data.Content) != 1 {
		t.Fatalf("unexpected content length: %d", len(resp.Data.Content))
	}
	if resp.Data.Content[0].ReturnCenterCode != "1000557004" {
		t.Fatalf("unexpected returnCenterCode: %s", resp.Data.Content[0].ReturnCenterCode)
	}
	if resp.Data.Content[0].VendorCreditFee05kg == nil || *resp.Data.Content[0].VendorCreditFee05kg != 3000 {
		t.Fatalf("unexpected vendorCreditFee05kg: %v", resp.Data.Content[0].VendorCreditFee05kg)
	}
}

func TestArrangeShipmentDirectIntegrationResponseUnmarshal(t *testing.T) {
	raw := `{
  "12345": {
    "success": true,
    "errorKey": "",
    "errorMessage": "",
    "invoiceNumber": "12345556",
    "shipmentBoxId": 12345,
    "validationCode": "",
    "invoiceNumberInvalidTime":"2025-09-02T13:54:39+08:00"
  }
}`

	var resp ArrangeShipmentDirectIntegrationResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
}

func TestArrangeShipmentDirectIntegrationValidate(t *testing.T) {
	client := NewCoupangClient("", "", "A00012345")

	if _, err := client.ArrangeShipmentDirectIntegration(context.Background(), nil); err == nil {
		t.Fatalf("expected error for empty items")
	}

	if _, err := client.ArrangeShipmentDirectIntegration(context.Background(), []ArrangeShipmentDirectIntegrationRequestItem{
		{ShipmentBoxId: 1, DeliveryCompanyCode: "TWL_FM", ReturnCenterCode: "1100047243"},
		{ShipmentBoxId: 2, DeliveryCompanyCode: "TWL_711", ReturnCenterCode: "1100047243"},
	}); err == nil {
		t.Fatalf("expected error for mixed deliveryCompanyCode")
	}

	if _, err := client.ArrangeShipmentDirectIntegration(context.Background(), []ArrangeShipmentDirectIntegrationRequestItem{
		{ShipmentBoxId: 1, DeliveryCompanyCode: "TWL_FM"},
	}); err == nil {
		t.Fatalf("expected error for missing returnCenterCode")
	}

	if _, err := client.ArrangeShipmentDirectIntegration(context.Background(), []ArrangeShipmentDirectIntegrationRequestItem{
		{ShipmentBoxId: 1, DeliveryCompanyCode: "TWL_KERRY"},
	}); err == nil {
		t.Fatalf("expected error for missing outboundShippingPlaceCode")
	}
}

func TestDownloadDirectIntegrationInvoicesValidate(t *testing.T) {
	if err := validateDownloadDirectIntegrationInvoicesRequest(nil); err == nil {
		t.Fatalf("expected error for nil request")
	}

	if err := validateDownloadDirectIntegrationInvoicesRequest(&DownloadDirectIntegrationInvoicesRequest{}); err == nil {
		t.Fatalf("expected error for empty request")
	}

	if err := validateDownloadDirectIntegrationInvoicesRequest(&DownloadDirectIntegrationInvoicesRequest{
		DeliveryCompanyCode: "TWL_KERRY",
		InvoicePrintDtoList: nil,
	}); err == nil {
		t.Fatalf("expected error for empty invoicePrintDtoList")
	}

	if err := validateDownloadDirectIntegrationInvoicesRequest(&DownloadDirectIntegrationInvoicesRequest{
		DeliveryCompanyCode: "TWL_KERRY",
		InvoicePrintDtoList: []DirectIntegrationInvoice{{ShipmentBoxId: 1, InvoiceNumber: "x"}, {ShipmentBoxId: 2, InvoiceNumber: "x"}, {ShipmentBoxId: 3, InvoiceNumber: "x"}, {ShipmentBoxId: 4, InvoiceNumber: "x"}, {ShipmentBoxId: 5, InvoiceNumber: "x"}, {ShipmentBoxId: 6, InvoiceNumber: "x"}},
	}); err == nil {
		t.Fatalf("expected error for kerry size limit")
	}

	if err := validateDownloadDirectIntegrationInvoicesRequest(&DownloadDirectIntegrationInvoicesRequest{
		DeliveryCompanyCode: "TWL_FM",
		InvoicePrintDtoList: []DirectIntegrationInvoice{{ShipmentBoxId: 0, InvoiceNumber: "x"}},
	}); err == nil {
		t.Fatalf("expected error for shipmentBoxId")
	}

	if err := validateDownloadDirectIntegrationInvoicesRequest(&DownloadDirectIntegrationInvoicesRequest{
		DeliveryCompanyCode: "TWL_711",
		InvoicePrintDtoList: []DirectIntegrationInvoice{{ShipmentBoxId: 1, InvoiceNumber: ""}},
	}); err == nil {
		t.Fatalf("expected error for invoiceNumber")
	}
}

func TestParseContentDispositionFilename(t *testing.T) {
	got := parseContentDispositionFilename(`attachment; filename="7e1beb74-11e9-407c-9775-04696ea3db2f.pdf"`)
	if got != "7e1beb74-11e9-407c-9775-04696ea3db2f.pdf" {
		t.Fatalf("unexpected filename: %s", got)
	}
}

//func TestCoupangAPIs(t *testing.T) {
//	// Parse flags to allow -f override
//	flag.Parse()
//
//	//ctx := loadConfig()
//	//fmt.Printf("Loaded Config Coupang: %+v\n", ctx.Config.Coupang)
//
//	// VendorID should be replaced with actual VendorID
//	// Since it's not in config, we use a placeholder or one from environment
//	vendorID := "A00012345"
//
//	client := NewCoupangClient(ctx.Config.Coupang.AccessKey, ctx.Config.Coupang.SecretKey, vendorID)
//
//	t.Run("GetOrderListDaily", func(t *testing.T) {
//		// Example: Query yesterday's orders
//		now := time.Now()
//		// Coupang API expects yyyy-MM-dd or with time?
//		// Ref says "yyyy-mm-dd%2B09:00" for query param
//		// doRequest encodes params, so we pass "yyyy-mm-dd+09:00"
//		// or just "yyyy-mm-dd" if time is not needed for daily?
//		// Ref 1 example: createdAtFrom=2025-07-15%2B09:00
//		// So we should provide it with timezone.
//
//		from := now.AddDate(0, 0, -7).Format("2006-01-02") + "+09:00"
//		to := now.Format("2006-01-02") + "+09:00"
//
//		resp, err := client.GetOrderListDaily(from, to, "ACCEPT", 10)
//		if err != nil {
//			t.Logf("GetOrderListDaily failed (might be due to invalid keys/vendorID): %v", err)
//		} else {
//			fmt.Printf("Daily Orders: %+v\n", resp)
//		}
//	})
//
//	t.Run("GetOrderListMinutely", func(t *testing.T) {
//		// Example: Query last hour
//		now := time.Now()
//		// Ref 2: yyyy-mm-ddT01:00%2B09:00
//		from := now.Add(-1*time.Hour).Format("2006-01-02T15:04") + "+09:00"
//		to := now.Format("2006-01-02T15:04") + "+09:00"
//
//		resp, err := client.GetOrderListMinutely(from, to, "ACCEPT")
//		if err != nil {
//			t.Logf("GetOrderListMinutely failed: %v", err)
//		} else {
//			fmt.Printf("Minutely Orders: %+v\n", resp)
//		}
//	})
//
//	t.Run("GetOrderByShipmentBoxId", func(t *testing.T) {
//		// Placeholder ID
//		shipmentBoxId := "123456789"
//		resp, err := client.GetOrderByShipmentBoxId(shipmentBoxId)
//		if err != nil {
//			t.Logf("GetOrderByShipmentBoxId failed: %v", err)
//		} else {
//			fmt.Printf("ShipmentBox Order: %+v\n", resp)
//		}
//	})
//
//	t.Run("GetOrderByOrderId", func(t *testing.T) {
//		// Placeholder ID
//		orderId := "500000596"
//		resp, err := client.GetOrderByOrderId(orderId)
//		if err != nil {
//			t.Logf("GetOrderByOrderId failed: %v", err)
//		} else {
//			fmt.Printf("Order by ID: %+v\n", resp)
//		}
//	})
//}
