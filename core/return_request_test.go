package core

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetReturnRequestByReceiptId(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		wantURL := Schema + "://" + Host + "/v2/providers/openapi/apis/api/v6/vendors/A00012345/returnRequests/365937"
		if req.Method != http.MethodGet || req.URL.String() != wantURL {
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL)
		}
		if req.Header.Get("X-MARKET") != "TW" || req.Header.Get("X-Requested-By") != "A00012345" || !strings.Contains(req.Header.Get("Authorization"), "access-key=test-key") {
			t.Fatal("missing market, vendor or authorization header")
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Request: req, Body: io.NopCloser(strings.NewReader(`{
			"code":"200","message":"OK","data":[{
				"receiptId":365937,"orderId":500004398,"receiptType":"RETURN",
				"receiptStatus":"RELEASE_STOP_UNCHECKED","requesterRealPhoneNumber":null,
				"completeConfirmDate":"","returnShippingCharge":{"currencyCode":"TWD","units":-3000,"nanos":0},
				"returnItems":[{"vendorItemId":3000001893,"shipmentBoxId":123456789012345678,"cancelCount":1,"purchaseCount":2}],
				"returnDeliveryDtos":[{"deliveryCompanyCode":"DIRECT","deliveryInvoiceNo":"201807261200"},{"deliveryInvoiceNo":null}]
			}]
		}`))}, nil
	})
	client := NewCoupangClient("test-key", "test-secret", "A00012345")
	resp, err := client.GetReturnRequestByReceiptId(context.Background(), 365937)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Code != "200" || resp.Message != "OK" || len(resp.Data) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	item := resp.Data[0]
	if item.ReceiptId != 365937 || item.OrderId != 500004398 || item.ReceiptStatus != "RELEASE_STOP_UNCHECKED" || item.RequesterRealPhoneNumber != nil || item.CompleteConfirmDate != "" || item.ReturnShippingCharge.Units != -3000 {
		t.Fatalf("unexpected return: %+v", item)
	}
	if len(item.ReturnItems) != 1 || item.ReturnItems[0].ShipmentBoxId != 123456789012345678 || item.ReturnItems[0].CancelCount != 1 || item.ReturnItems[0].PurchaseCount != 2 {
		t.Fatalf("unexpected return items: %+v", item.ReturnItems)
	}
	if len(item.ReturnDeliveryDtos) != 2 || item.ReturnDeliveryDtos[0].DeliveryInvoiceNo == nil || *item.ReturnDeliveryDtos[0].DeliveryInvoiceNo != "201807261200" || item.ReturnDeliveryDtos[1].DeliveryInvoiceNo != nil {
		t.Fatalf("unexpected delivery invoices: %+v", item.ReturnDeliveryDtos)
	}
}

func TestGetReturnRequestByReceiptIdValidation(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	http.DefaultTransport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("invalid input must not send an HTTP request")
		return nil, nil
	})
	for _, tt := range []struct {
		vendorID  string
		receiptId int64
		want      string
	}{
		{"A00012345", 0, "receiptId"},
		{"A00012345", -1, "receiptId"},
		{"", 365937, "vendorId"},
		{" \t", 365937, "vendorId"},
	} {
		client := NewCoupangClient("", "", tt.vendorID)
		if _, err := client.GetReturnRequestByReceiptId(context.Background(), tt.receiptId); err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Fatalf("error = %v, want %s validation error", err, tt.want)
		}
	}
}

func TestGetReturnRequestByReceiptIdErrors(t *testing.T) {
	transportErr := errors.New("connection failed")
	for _, tt := range []struct {
		name string
		body string
		err  error
		want string
	}{
		{"transport", "", transportErr, "connection failed"},
		{"invalid JSON", "<html>error</html>", nil, "decode single return request"},
		{"API error", `{"code":"400","message":"ReceiptId doesn't belong to the vendorId, or cannot find the Receipt by the given id"}`, nil, "ReceiptId doesn't belong to the vendorId"},
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
			resp, err := client.GetReturnRequestByReceiptId(context.Background(), 365937)
			if resp != nil || err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("response = %+v, error = %v, want %s", resp, err, tt.want)
			}
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Fatalf("transport error was not preserved: %v", err)
			}
		})
	}
}
