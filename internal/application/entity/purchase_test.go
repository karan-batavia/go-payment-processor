package entity

import (
	"errors"
	"reflect"
	"testing"
)

func TestCreatePurchase(t *testing.T) {
	items := []string{"Item 1", "Item 2"}
	purchase := NewPurchase(9.99, items, 5)

	if purchase == nil {
		t.Error("purchase should have been created")
		return
	}

	if purchase.Value != 9.99 {
		t.Error("purchase value should be 9.99")
	}

	if !reflect.DeepEqual(purchase.Items, items) {
		t.Errorf("purchase items should be %v", items)
	}

	if purchase.Installments != 5 {
		t.Error("purchase installments shoudl be 5")
	}
}

func TestPurchaseFields(t *testing.T) {
	testCases := []struct {
		Test         string
		Value        float64
		Items        []string
		Installments int
		err          error
	}{
		{
			"value is negative",
			-1,
			[]string{"Item 1", "Item 2"},
			1,
			PurchaseValueIsInvalidErr,
		},
		{
			"value is zero",
			0,
			[]string{"Item 1", "Item 2"},
			1,
			PurchaseValueIsInvalidErr,
		},
		{
			"items is nil",
			1,
			nil,
			1,
			PurchaseItemsIsRequiredErr,
		},
		{
			"items is empty",
			1,
			[]string{},
			1,
			PurchaseItemsIsRequiredErr,
		},
		{
			"items has empty elements",
			1,
			[]string{"Item 1", ""},
			1,
			PurchaseItemsAreInvalidErr,
		},
		{
			"installments is negative",
			1,
			[]string{"Item 1", "Item 2"},
			-1,
			PurchaseInstallmentsIsInvalidErr,
		},
		{
			"installments is zero",
			1,
			[]string{"Item 1", "Item 2"},
			0,
			PurchaseInstallmentsIsInvalidErr,
		},
		{
			"all fields are valid",
			1,
			[]string{"Item 1", "Item 2"},
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			err := NewPurchase(tc.Value, tc.Items, tc.Installments).Validate()
			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}
