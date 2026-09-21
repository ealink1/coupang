package core

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestGetReturnRequestList(t *testing.T) {
	tests := []struct {
		name string
		req  GetReturnRequestListRequest
		want url.Values
	}{
		{
			name: "daily pagination",
			req:  GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24", CreatedAtTo: "2025-07-25", Status: "UC", CancelType: "RETURN", MaxPerPage: 50, NextToken: "next+/=token"},
			want: url.Values{"createdAtFrom": {"2025-07-24"}, "createdAtTo": {"2025-07-25"}, "status": {"UC"}, "cancelType": {"RETURN"}, "maxPerPage": {"50"}, "nextToken": {"next+/=token"}},
		},
		{
			name: "minute query",
			req:  GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24T00:01", CreatedAtTo: "2025-07-24T23:59", SearchType: "timeFrame", Status: "RU"},
			want: url.Values{"createdAtFrom": {"2025-07-24T00:01"}, "createdAtTo": {"2025-07-24T23:59"}, "searchType": {"timeFrame"}, "status": {"RU"}},
		},
		{
			name: "cancel by order",
			req:  GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24", CreatedAtTo: "2025-07-24", CancelType: "CANCEL", OrderId: 28000008707838},
			want: url.Values{"createdAtFrom": {"2025-07-24"}, "createdAtTo": {"2025-07-24"}, "cancelType": {"CANCEL"}, "orderId": {"28000008707838"}},
		},
		{
			name: "cancel list defaults",
			req:  GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24", CreatedAtTo: "2025-07-24", CancelType: "CANCEL"},
			want: url.Values{"createdAtFrom": {"2025-07-24"}, "createdAtTo": {"2025-07-24"}, "cancelType": {"CANCEL"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := http.DefaultTransport
			t.Cleanup(func() { http.DefaultTransport = original })
			http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				wantURL, _ := url.Parse(Schema + "://" + Host + "/v2/providers/openapi/apis/api/v6/vendors/A00012345/returnRequests")
				if req.Method != http.MethodGet || req.URL.Host != wantURL.Host || req.URL.Path != wantURL.Path {
					t.Fatalf("unexpected request: %s %s", req.Method, req.URL)
				}
				if !reflect.DeepEqual(req.URL.Query(), tt.want) {
					t.Fatalf("query = %v, want %v", req.URL.Query(), tt.want)
				}
				if req.Header.Get("X-MARKET") != "TW" || req.Header.Get("X-Requested-By") != "A00012345" || !strings.Contains(req.Header.Get("Authorization"), "access-key=test-key") {
					t.Fatalf("missing market, vendor or authorization header")
				}
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{
					"code":200,"message":"OK","nextToken":"next-page",
					"data":[{"receiptId":50229613,"orderId":28000008707838,"paymentId":28000009486604,
					"receiptType":"RETURN","receiptStatus":"RETURNS_UNCHECKED",
					"createdAt":"2025-01-15T14:17:13.973885-08:00","requesterRealPhoneNumber":null,
					"completeConfirmDate":"","enclosePrice":{"currencyCode":"TWD","units":0,"nanos":0},
					"returnShippingCharge":{"currencyCode":"TWD","units":-3000,"nanos":-500000000},
					"returnItems":[{"shipmentBoxId":123456789012345678,"vendorItemId":3187044096,"purchaseCount":2,"cancelCount":1}],
					"returnDeliveryDtos":[{"deliveryCompanyCode":"TWL_711","deliveryInvoiceNo":null},{"deliveryInvoiceNo":""},{"deliveryInvoiceNo":"12345"}]}]
				}`)), Request: req}, nil
			})
			client := NewCoupangClient("test-key", "test-secret", "A00012345")
			resp, err := client.GetReturnRequestList(context.Background(), &tt.req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.NextToken != "next-page" || len(resp.Data) != 1 {
				t.Fatalf("unexpected response: %+v", resp)
			}
			item := resp.Data[0]
			if item.OrderId != 28000008707838 || len(item.ReturnItems) != 1 || item.ReturnItems[0].ShipmentBoxId != 123456789012345678 || item.ReturnItems[0].CancelCount != 1 {
				t.Fatalf("unexpected order or item: %+v", item)
			}
			if item.ReturnShippingCharge.Units != -3000 || item.ReturnShippingCharge.Nanos != -500000000 || item.RequesterRealPhoneNumber != nil || item.CompleteConfirmDate != "" {
				t.Fatalf("unexpected nullable fields or money: %+v", item)
			}
			if len(item.ReturnDeliveryDtos) != 3 || item.ReturnDeliveryDtos[0].DeliveryInvoiceNo != nil || item.ReturnDeliveryDtos[1].DeliveryInvoiceNo == nil || *item.ReturnDeliveryDtos[1].DeliveryInvoiceNo != "" || item.ReturnDeliveryDtos[2].DeliveryInvoiceNo == nil || *item.ReturnDeliveryDtos[2].DeliveryInvoiceNo != "12345" {
				t.Fatalf("unexpected delivery invoices: %+v", item.ReturnDeliveryDtos)
			}
		})
	}
}

func TestGetReturnRequestListValidation(t *testing.T) {
	tests := []struct {
		name   string
		change func(*GetReturnRequestListRequest)
	}{
		{"missing from", func(r *GetReturnRequestListRequest) { r.CreatedAtFrom = "" }},
		{"invalid date", func(r *GetReturnRequestListRequest) { r.CreatedAtTo = "2025-02-30" }},
		{"reversed dates", func(r *GetReturnRequestListRequest) { r.CreatedAtTo = "2025-07-23" }},
		{"over 31 days", func(r *GetReturnRequestListRequest) { r.CreatedAtTo = "2025-08-25" }},
		{"invalid search type", func(r *GetReturnRequestListRequest) { r.SearchType = "typo" }},
		{"cancel with status", func(r *GetReturnRequestListRequest) { r.CancelType = "CANCEL" }},
		{"invalid cancel type", func(r *GetReturnRequestListRequest) { r.CancelType = "typo" }},
		{"invalid status", func(r *GetReturnRequestListRequest) { r.Status = "typo" }},
		{"missing status and order", func(r *GetReturnRequestListRequest) { r.Status = "" }},
		{"negative page size", func(r *GetReturnRequestListRequest) { r.MaxPerPage = -1 }},
		{"negative order", func(r *GetReturnRequestListRequest) { r.OrderId = -1 }},
	}
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	http.DefaultTransport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("invalid input must not send an HTTP request")
		return nil, nil
	})
	client := NewCoupangClient("", "", "A00012345")
	if _, err := client.GetReturnRequestList(context.Background(), nil); err == nil {
		t.Fatal("expected nil request error")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24", CreatedAtTo: "2025-07-24", Status: "UC"}
			tt.change(&req)
			if _, err := client.GetReturnRequestList(context.Background(), &req); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
	for _, change := range []func(*GetReturnRequestListRequest){
		func(r *GetReturnRequestListRequest) { r.NextToken = "token" },
		func(r *GetReturnRequestListRequest) { r.MaxPerPage = 50 },
		func(r *GetReturnRequestListRequest) { r.OrderId = 123 },
	} {
		req := GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24T00:01", CreatedAtTo: "2025-07-24T23:59", SearchType: "timeFrame", Status: "UC"}
		change(&req)
		if _, err := client.GetReturnRequestList(context.Background(), &req); err == nil {
			t.Fatal("expected incompatible minute query error")
		}
	}
	// Exactly 31 days and return queries by order ID are supported.
	req := GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24", CreatedAtTo: "2025-08-24", OrderId: 28000008707838}
	if err := validateGetReturnRequestListRequest(&req); err != nil {
		t.Fatal(err)
	}
	client.VendorID = " "
	if _, err := client.GetReturnRequestList(context.Background(), &req); err == nil {
		t.Fatal("expected vendor ID error")
	}
}

func TestGetReturnRequestListErrors(t *testing.T) {
	transportErr := errors.New("connection failed")
	for _, tt := range []struct {
		name string
		body string
		err  error
		want string
	}{
		{"transport", "", transportErr, "connection failed"},
		{"malformed JSON", "<html>error</html>", nil, "decode return request list"},
		{"API error", `{"code":400,"message":"Invalid vendor ID"}`, nil, "Invalid vendor ID"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			original := http.DefaultTransport
			t.Cleanup(func() { http.DefaultTransport = original })
			http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if tt.err != nil {
					return nil, tt.err
				}
				return &http.Response{StatusCode: 400, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(tt.body)), Request: req}, nil
			})
			client := NewCoupangClient("test-key", "test-secret", "A00012345")
			_, err := client.GetReturnRequestList(context.Background(), &GetReturnRequestListRequest{CreatedAtFrom: "2025-07-24", CreatedAtTo: "2025-07-24", CancelType: "CANCEL"})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %s", err, tt.want)
			}
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Fatalf("transport error was not preserved: %v", err)
			}
		})
	}
}
