package core

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCoupang(t *testing.T) {

	vendorID := "A158000L2X"
	client := NewCoupangClient("4b66358f-efbc-490c-a7ca-63e7d100f901", "14a02ef7123c8da073672a3e955200eeb10c81fb", vendorID)

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
