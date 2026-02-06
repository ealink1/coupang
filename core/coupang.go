package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

const (
	Schema = "https"
	Host   = "api-gateway.coupang.com"
)

type CoupangClient struct {
	AccessKey string
	SecretKey string
	VendorID  string
}

func NewCoupangClient(accessKey, secretKey, vendorID string) *CoupangClient {
	return &CoupangClient{
		AccessKey: accessKey,
		SecretKey: secretKey,
		VendorID:  vendorID,
	}
}

// GetOrderListDaily queries order list by day
func (c *CoupangClient) GetOrderListDaily(ctx context.Context, req *GetOrderListRequest) (*OrderListResponse, error) {
	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v5/vendors/%s/ordersheets", c.VendorID)
	params := url.Values{}
	params.Set("createdAtFrom", req.CreatedAtFrom)
	params.Set("createdAtTo", req.CreatedAtTo)
	params.Set("status", req.Status)
	if req.MaxPerPage > 0 {
		params.Set("maxPerPage", fmt.Sprintf("%d", req.MaxPerPage))
	}
	params.Set("searchType", req.SearchType)
	params.Set("nextToken", req.NextToken)

	var resp OrderListResponse
	respStr, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}

	log.Println(respStr)
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

// GetOrderListMinutely queries order list by minute
func (c *CoupangClient) GetOrderListMinutely(ctx context.Context, createdAtFrom, createdAtTo string, status string) (*OrderListResponse, error) {
	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v5/vendors/%s/ordersheets", c.VendorID)
	params := url.Values{}
	params.Set("createdAtFrom", createdAtFrom)
	params.Set("createdAtTo", createdAtTo)
	params.Set("status", status)
	params.Set("searchType", "timeFrame")

	var resp OrderListResponse
	respStr, err := c.doRequest(ctx, "GET", path, params)
	log.Println(respStr)
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

// GetOrderByShipmentBoxId queries single order by shipmentBoxId
func (c *CoupangClient) GetOrderByShipmentBoxId(ctx context.Context, shipmentBoxId string) (*SingleOrderResponse, error) {
	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v5/vendors/%s/ordersheets/%s", c.VendorID, shipmentBoxId)
	params := url.Values{}

	var resp SingleOrderResponse
	respStr, err := c.doRequest(ctx, "GET", path, params)
	log.Println(respStr)
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

// GetOrderByOrderId queries single order by orderId
func (c *CoupangClient) GetOrderByOrderId(ctx context.Context, orderId string) (*SingleOrderListResponse, error) {
	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v5/vendors/%s/%s/ordersheets", c.VendorID, orderId)
	params := url.Values{}

	var resp SingleOrderListResponse
	respStr, err := c.doRequest(ctx, "GET", path, params)
	log.Println(respStr)
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (c *CoupangClient) GetSellerProductInventories(ctx context.Context, sellerProductIds []int64) (*InventoryResponse, error) {
	if len(sellerProductIds) == 0 {
		return nil, fmt.Errorf("sellerProductIds is required")
	}
	if len(sellerProductIds) > 50 {
		return nil, fmt.Errorf("maximum 50 sellerProductIds allowed")
	}

	path := "/v2/providers/seller_api/apis/api/v1/marketplace/seller-products/inventories"
	body := BatchInventoryRequest{
		SellerProductIds: sellerProductIds,
	}

	var resp InventoryResponse
	respStr, err := c.doPostJSON(ctx, path, body)
	if err != nil {
		return nil, err
	}

	log.Println(respStr)
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CoupangClient) GetSellerProduct(ctx context.Context, sellerProductId int64) (*ProductDetailResponse, error) {
	if sellerProductId == 0 {
		return nil, fmt.Errorf("sellerProductId is required")
	}

	path := fmt.Sprintf("/v2/providers/seller_api/apis/api/v1/marketplace/seller-products/%d", sellerProductId)
	params := url.Values{}

	var resp ProductDetailResponse
	respStr, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}

	log.Println(respStr)
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
