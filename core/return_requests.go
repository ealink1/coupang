package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func validateGetReturnRequestListRequest(req *GetReturnRequestListRequest) error {
	if req == nil {
		return fmt.Errorf("request is required")
	}
	layout := "2006-01-02"
	switch req.SearchType {
	case "":
	case "timeFrame":
		layout = "2006-01-02T15:04"
		if req.NextToken != "" || req.MaxPerPage != 0 || req.OrderId != 0 {
			return fmt.Errorf("nextToken, maxPerPage and orderId are not supported with searchType=timeFrame")
		}
	default:
		return fmt.Errorf("searchType must be empty or timeFrame")
	}
	from, err := time.Parse(layout, req.CreatedAtFrom)
	if err != nil {
		return fmt.Errorf("invalid createdAtFrom: %w", err)
	}
	to, err := time.Parse(layout, req.CreatedAtTo)
	if err != nil {
		return fmt.Errorf("invalid createdAtTo: %w", err)
	}
	if to.Before(from) || to.Sub(from) > 31*24*time.Hour {
		return fmt.Errorf("createdAtTo must be on or after createdAtFrom and within 31 days")
	}
	if req.MaxPerPage < 0 || req.OrderId < 0 {
		return fmt.Errorf("maxPerPage and orderId must not be negative")
	}
	switch req.CancelType {
	case "", "RETURN":
		if req.Status == "" && req.OrderId == 0 {
			return fmt.Errorf("status or orderId is required for return queries")
		}
	case "CANCEL":
		if req.Status != "" {
			return fmt.Errorf("status is not supported with cancelType=CANCEL")
		}
	default:
		return fmt.Errorf("cancelType must be RETURN or CANCEL")
	}
	switch req.Status {
	case "", "RU", "UC", "CC", "PR":
	default:
		return fmt.Errorf("status must be RU, UC, CC or PR")
	}
	return nil
}

// GetReturnRequestList 查询一页退货/取消订单记录，NextToken 由调用方传入以获取下一页。
// 台湾市场，使用现有客户端的签名和 X-MARKET: TW 请求头。
// https://developers.coupang.com/zh-TW/api/returns/return-cancellation-request-list-query
func (c *CoupangClient) GetReturnRequestList(ctx context.Context, req *GetReturnRequestListRequest) (*ReturnRequestListResponse, error) {
	//if err := validateGetReturnRequestListRequest(req); err != nil {
	//	return nil, err
	//}
	if strings.TrimSpace(c.VendorID) == "" {
		return nil, fmt.Errorf("vendorId is required")
	}
	path := fmt.Sprintf("/v2/providers/openapi/apis/api/v6/vendors/%s/returnRequests", c.VendorID)
	params := url.Values{}
	params.Set("createdAtFrom", req.CreatedAtFrom)
	params.Set("createdAtTo", req.CreatedAtTo)
	for key, value := range map[string]string{
		"searchType": req.SearchType,
		"status":     req.Status,
		"cancelType": req.CancelType,
		"nextToken":  req.NextToken,
	} {
		if value != "" {
			params.Set(key, value)
		}
	}
	if req.MaxPerPage > 0 {
		params.Set("maxPerPage", strconv.Itoa(req.MaxPerPage))
	}
	if req.OrderId > 0 {
		params.Set("orderId", strconv.FormatInt(req.OrderId, 10))
	}
	body, err := c.doRequest(ctx, http.MethodGet, path, params)
	if err != nil {
		return nil, err
	}
	var resp ReturnRequestListResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, fmt.Errorf("decode return request list: %w", err)
	}
	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("return request list failed with code %d: %s", resp.Code, resp.Message)
	}
	return &resp, nil
}
