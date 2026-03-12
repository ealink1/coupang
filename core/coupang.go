package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/url"
	"strings"
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

func (c *CoupangClient) GetReturnShippingCenters(ctx context.Context, req *GetReturnShippingCentersRequest) (*ReturnShippingCentersResponse, error) {
	pageNum := 1
	pageSize := 10
	if req != nil {
		if req.PageNum > 0 {
			pageNum = req.PageNum
		}
		if req.PageSize > 0 {
			pageSize = req.PageSize
		}
	}
	if pageSize > 50 {
		return nil, fmt.Errorf("pageSize must be <= 50")
	}

	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v5/vendors/%s/returnShippingCenters", c.VendorID)
	params := url.Values{}
	params.Set("pageNum", fmt.Sprintf("%d", pageNum))
	params.Set("pageSize", fmt.Sprintf("%d", pageSize))

	var resp ReturnShippingCentersResponse
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

func (c *CoupangClient) ArrangeShipmentDirectIntegration(ctx context.Context, items []ArrangeShipmentDirectIntegrationRequestItem) (*ArrangeShipmentDirectIntegrationResponse, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items is required")
	}

	deliveryCompanyCode := ""
	for _, item := range items {
		if item.ShipmentBoxId <= 0 {
			return nil, fmt.Errorf("shipmentBoxId must be > 0")
		}
		if item.DeliveryCompanyCode == "" {
			return nil, fmt.Errorf("deliveryCompanyCode is required")
		}
		if deliveryCompanyCode == "" {
			deliveryCompanyCode = item.DeliveryCompanyCode
		} else if deliveryCompanyCode != item.DeliveryCompanyCode {
			return nil, fmt.Errorf("only one type of deliveryCompanyCode is supported for bulk delivery")
		}

		switch item.DeliveryCompanyCode {
		case "TWL_FM", "TWL_711":
			if item.ReturnCenterCode == "" {
				return nil, fmt.Errorf("returnCenterCode is required for deliveryCompanyCode: %s", item.DeliveryCompanyCode)
			}
		case "TWL_KERRY":
			if item.OutboundShippingPlaceCode <= 0 {
				return nil, fmt.Errorf("outboundShippingPlaceCode is required for deliveryCompanyCode: %s", item.DeliveryCompanyCode)
			}
		}
	}

	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v1/vendors/%s/orders/directIntegration/arrangeShipment", c.VendorID)

	var resp ArrangeShipmentDirectIntegrationResponse
	respStr, err := c.doPostJSON(ctx, path, items)
	if err != nil {
		return nil, err
	}

	log.Println(respStr)
	if err := json.Unmarshal([]byte(respStr), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func validateDownloadDirectIntegrationInvoicesRequest(req *DownloadDirectIntegrationInvoicesRequest) error {
	if req == nil {
		return fmt.Errorf("request is required")
	}
	if req.DeliveryCompanyCode == "" {
		return fmt.Errorf("deliveryCompanyCode is required")
	}
	if len(req.InvoicePrintDtoList) == 0 {
		return fmt.Errorf("invoicePrintDtoList is required")
	}

	limit := 0
	switch req.DeliveryCompanyCode {
	case "TWL_KERRY":
		limit = 5
	case "TWL_FM":
		limit = 30
	case "TWL_711":
		limit = 40
	default:
		return fmt.Errorf("unsupported deliveryCompanyCode")
	}
	if len(req.InvoicePrintDtoList) > limit {
		return fmt.Errorf("request size exceeds maximum limit of %d. requestSize: %d", limit, len(req.InvoicePrintDtoList))
	}

	for _, item := range req.InvoicePrintDtoList {
		if item.ShipmentBoxId <= 0 {
			return fmt.Errorf("shipmentBoxId must be > 0")
		}
		if item.InvoiceNumber == "" {
			return fmt.Errorf("invoiceNumber is required")
		}
	}
	return nil
}

func parseContentDispositionFilename(contentDisposition string) string {
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return ""
	}
	return params["filename"]
}

func (c *CoupangClient) DownloadDirectIntegrationInvoices(ctx context.Context, req *DownloadDirectIntegrationInvoicesRequest) (*DownloadDirectIntegrationInvoicesFile, error) {
	if err := validateDownloadDirectIntegrationInvoicesRequest(req); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v1/vendors/%s/orders/directIntegration/downloadInvoices", c.VendorID)
	statusCode, header, body, err := c.doPostJSONWithHeaders(ctx, path, req)
	if err != nil {
		return nil, err
	}
	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	contentDisposition := header.Get("Content-Disposition")
	return &DownloadDirectIntegrationInvoicesFile{
		ContentType:        header.Get("Content-Type"),
		ContentDisposition: contentDisposition,
		Filename:           parseContentDispositionFilename(contentDisposition),
		Content:            body,
	}, nil
}

func validateGetOutboundShippingPlacesRequest(req *GetOutboundShippingPlacesRequest) error {
	if req == nil {
		return fmt.Errorf("request is required")
	}

	hasPage := req.PageNum > 0 || req.PageSize > 0
	hasCodes := len(req.PlaceCodes) > 0
	hasNames := len(req.PlaceNames) > 0
	if !hasPage && !hasCodes && !hasNames {
		return fmt.Errorf("(pageNum & pageSize) or placeCodes or placeNames must be provided to call this API")
	}

	if hasPage {
		if req.PageNum < 1 {
			return fmt.Errorf("pageNum must be >= 1")
		}
		if req.PageSize <= 0 {
			return fmt.Errorf("pageSize must be > 0")
		}
		if req.PageSize > 50 {
			return fmt.Errorf("pageSize must be <= 50")
		}
	}

	if hasCodes {
		for _, v := range req.PlaceCodes {
			if v <= 0 {
				return fmt.Errorf("placeCodes must be > 0")
			}
		}
	}

	if hasNames {
		hasNonEmpty := false
		for _, v := range req.PlaceNames {
			if strings.TrimSpace(v) != "" {
				hasNonEmpty = true
				break
			}
		}
		if !hasNonEmpty {
			return fmt.Errorf("placeNames must contain non-empty values")
		}
	}

	return nil
}

func (c *CoupangClient) GetOutboundShippingPlaces(ctx context.Context, req *GetOutboundShippingPlacesRequest) (*OutboundShippingPlacesResponse, error) {
	if err := validateGetOutboundShippingPlacesRequest(req); err != nil {
		return nil, err
	}

	path := "/v2/providers/marketplace_openapi/apis/api/v2/vendor/shipping-place/outbound"
	params := url.Values{}

	if len(req.PlaceCodes) > 0 {
		parts := make([]string, 0, len(req.PlaceCodes))
		for _, v := range req.PlaceCodes {
			if v > 0 {
				parts = append(parts, fmt.Sprintf("%d", v))
			}
		}
		if len(parts) == 0 {
			return nil, fmt.Errorf("placeCodes must contain positive values")
		}
		params.Set("placeCodes", strings.Join(parts, ","))
	} else if len(req.PlaceNames) > 0 {
		parts := make([]string, 0, len(req.PlaceNames))
		for _, v := range req.PlaceNames {
			v = strings.TrimSpace(v)
			if v != "" {
				parts = append(parts, v)
			}
		}
		if len(parts) == 0 {
			return nil, fmt.Errorf("placeNames must contain non-empty values")
		}
		params.Set("placeNames", strings.Join(parts, ","))
	} else {
		params.Set("pageNum", fmt.Sprintf("%d", req.PageNum))
		params.Set("pageSize", fmt.Sprintf("%d", req.PageSize))
	}

	var resp OutboundShippingPlacesResponse
	respStr, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}

	log.Println(respStr)
	if err := json.Unmarshal([]byte(respStr), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
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
