package main

import (
	"context"
	"coupang/core"
	"fmt"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

	//10	2026-02-04 14:32:09.89015+08	A1580015Y3	cb2ee9a6-c4dc-4dfb-9a70-7d479addaf78	b732e10b84c9b1cd86e9f933ebb20b8bab9cf163	1

	//client := core.NewCoupangClient("4b66358f-efbc-490c-a7ca-63e7d100f901", "14a02ef7123c8da073672a3e955200eeb10c81fb", "A158000L2X")
	client := core.NewCoupangClient("cb2ee9a6-c4dc-4dfb-9a70-7d479addaf78", "b732e10b84c9b1cd86e9f933ebb20b8bab9cf163", "A1580015Y3")
	// now := time.Now()

	// from := now.AddDate(0, 0, 0).Format("2006-01-02") + "+08:00"
	// to := now.AddDate(0, 0, 0).Format("2006-01-02") + "+08:00"

	// //from := ""
	// //to := ""

	// fmt.Println("-----------------")
	// fmt.Println("-----------------")
	// fmt.Println("-----------------")
	// fmt.Println("-----------------")
	// fmt.Println("-------ACCEPT-----------")
	// req := &core.GetOrderListRequest{
	// 	CreatedAtFrom: from,
	// 	CreatedAtTo:   to,
	// 	Status:        "ACCEPT",
	// 	MaxPerPage:    0,
	// 	NextToken:     "",
	// }
	// resp, err := client.GetOrderListDaily(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Daily Orders: %+v\n", resp)
	// }

	// fmt.Println("-------INSTRUCT-----------")
	// req.Status = "INSTRUCT"
	// resp, err = client.GetOrderListDaily(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Daily Orders: %+v\n", resp)
	// }

	// fmt.Println("-------DEPARTURE-----------")
	// req.Status = "DELIVERING"
	// resp, err = client.GetOrderListDaily(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Daily Orders: %+v\n", resp)
	// }

	// fmt.Println("-------DELIVERING-----------")
	// req.Status = "DELIVERING"
	// resp, err = client.GetOrderListDaily(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Daily Orders: %+v\n", resp)
	// }
	// fmt.Println("-------FINAL_DELIVERY-----------")
	// req.Status = "FINAL_DELIVERY"
	// resp, err = client.GetOrderListDaily(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Daily Orders: %+v\n", resp)
	// }

	// fmt.Println("-------FINAL_DELIVERY-----------")
	// req.Status = "FINAL_DELIVERY"
	// resp, err = client.GetOrderListDaily(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Daily Orders: %+v\n", resp)
	// }

	// fmt.Println("-------NONE_TRACKING-----------")
	// req.Status = "NONE_TRACKING"
	// resp, err = client.GetOrderListDaily(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Daily Orders: %+v\n", resp)
	// }

	//fmt.Println("------------------")
	//
	//shipmentBoxId := "626887864795136"
	//_resp, err := client.GetOrderByShipmentBoxId(context.Background(), shipmentBoxId)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("ShipmentBox Order: %+v\n", _resp)
	//}
	//
	//fmt.Println("------------------")
	//orderId := "127100001830279"
	//orderresp, err := client.GetOrderByOrderId(context.Background(), orderId)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("Order Details: %+v\n", orderresp)
	//}

	// invResp, err := client.GetSellerProductInventories(context.Background(), []int64{1134696112653240})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Inventory Details: %+v\n", invResp)
	// }

	productResp, err := client.GetSellerProduct(context.Background(), 1134696112653240)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//fmt.Printf("Product Details: %+v\n", productResp)
	}
	for _, item := range productResp.Data.Items {
		for _, image := range item.Images {
			if image.ImageType == "REPRESENTATION" {
				fmt.Println(image.CdnPath)
			}
		}
	}
}
