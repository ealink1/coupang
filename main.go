package main

import (
	"context"
	"coupang/config"
	"coupang/core"
	"flag"
	"fmt"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

	flag.Parse()

	cfg := config.GetCfg()

	client := core.NewCoupangClient(cfg.Coupang.ApiKey, cfg.Coupang.SecretKey, cfg.Coupang.VendorId)
	//now := time.Now()

	//from := fmt.Sprintf("%s+08:00", now.AddDate(0, 0, -4).Format("2006-01-02")) // 不能有空格
	//to := fmt.Sprintf("%s+08:00", now.AddDate(0, 0, 0).Format("2006-01-02"))

	// //from := ""
	// //to := ""

	//fmt.Println("-----------------")
	//fmt.Println("-----------------")
	//fmt.Println("-----------------")
	//fmt.Println("-----------------")
	//fmt.Println("-------ACCEPT-----------")
	//req := &core.GetOrderListRequest{
	//	CreatedAtFrom: from,
	//	CreatedAtTo:   to,
	//	Status:        "ACCEPT",
	//	MaxPerPage:    50,
	//	NextToken:     "",
	//}
	//resp, err := client.GetOrderListDaily(context.Background(), req)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("Daily Orders: %+v\n", resp)
	//}

	//fmt.Println("-------INSTRUCT-----------")
	//req.Status = "INSTRUCT"
	//resp, err = client.GetOrderListDaily(context.Background(), req)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("Daily Orders: %+v\n", resp)
	//}

	//fmt.Println("-------DEPARTURE-----------")
	//req.Status = "DELIVERING"
	//resp, err = client.GetOrderListDaily(context.Background(), req)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("Daily Orders: %+v\n", resp)
	//}

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

	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// orderId := "8089769166"
	// orderresp, err := client.GetOrderByOrderId(context.Background(), orderId)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Order Details: %+v\n", orderresp)
	// }

	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//orderresp, err := client.GetReturnShippingCenters(context.Background(), &core.GetReturnShippingCentersRequest{
	//	PageNum:  1,
	//	PageSize: 50,
	//})
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("Order Details: %+v\n", orderresp)
	//}

	fmt.Println("------------------")
	fmt.Println("------------------")
	fmt.Println("------------------")
	fmt.Println("------------------")
	orderresp, err := client.ArrangeShipmentDirectIntegration(context.Background(), []core.ArrangeShipmentDirectIntegrationRequestItem{
		//{
		//	ShipmentBoxId:       661725267329026,
		//	DeliveryCompanyCode: "TWL_711",
		//	ReturnCenterCode:    "1002359671",
		//},
		{
			ShipmentBoxId:       663859534741504,
			DeliveryCompanyCode: "TWL_FM",
			ReturnCenterCode:    "1002431063",
		},
		{
			ShipmentBoxId:       662556762865664,
			DeliveryCompanyCode: "TWL_FM",
			ReturnCenterCode:    "1002431063",
		},
		//{
		//	ShipmentBoxId:             663836294021120,
		//	DeliveryCompanyCode:       "TWL_KERRY",
		//	OutboundShippingPlaceCode: 23556341,
		//},
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Order Details: %+v\n", orderresp)
	}

	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// orderresp, err = client.ArrangeShipmentDirectIntegration(context.Background(), []core.ArrangeShipmentDirectIntegrationRequestItem{
	// 	{
	// 		ShipmentBoxId:       659053428523008,
	// 		DeliveryCompanyCode: "TWL_FM",
	// 		// ReturnCenterCode:    "10000000000000000000000000000000",
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Order Details: %+v\n", orderresp)
	// }

	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// fmt.Println("------------------")
	// orderresp, err = client.ArrangeShipmentDirectIntegration(context.Background(), []core.ArrangeShipmentDirectIntegrationRequestItem{
	// 	{
	// 		ShipmentBoxId:       659160945377281,
	// 		DeliveryCompanyCode: "TWL_711",
	// 		// ReturnCenterCode:    "10000000000000000000000000000000",
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Order Details: %+v\n", orderresp)
	// }

	// invResp, err := client.GetSellerProductInventories(context.Background(), []int64{1134696112653240})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Printf("Inventory Details: %+v\n", invResp)
	// }
	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//for _, datum := range orderresp.Data {
	//	for _, item := range datum.OrderItems {
	//
	//		productResp, _err := client.GetSellerProduct(context.Background(), item.SellerProductId)
	//		if _err != nil {
	//			fmt.Println(_err.Error())
	//		} else {
	//			fmt.Printf("Product Details: %+v\n", productResp)
	//		}
	//	}
	//	//for _, item := range productResp.Data.Items {
	//	//	for _, image := range item.Images {
	//	//		if image.ImageType == "REPRESENTATION" {
	//	//			fmt.Println(image.CdnPath)
	//	//		}
	//	//	}
	//	//}
	//}

	//invoices, err := client.DownloadDirectIntegrationInvoices(context.Background(), &core.DownloadDirectIntegrationInvoicesRequest{
	//	DeliveryCompanyCode: "TWL_KERRY",
	//	InvoicePrintDtoList: []core.DirectIntegrationInvoice{
	//		{
	//			InvoiceNumber: "40529136745",
	//			ShipmentBoxId: 663883689672705,
	//		},
	//	},
	//})
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//marshal, _ := json.Marshal(invoices)
	//fmt.Println(string(marshal))

	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//fmt.Println("------------------")
	//// orderId := "8089769166"
	//orderresp, err := client.GetOutboundShippingPlaces(context.Background(), &core.GetOutboundShippingPlacesRequest{
	//	PageNum:  1,
	//	PageSize: 50,
	//})
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("Order Details: %+v\n", orderresp)
	//}

}
