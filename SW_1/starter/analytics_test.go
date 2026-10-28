package analytics_test

import (
	"reflect"
	"testing"

	analytics "example.com/sales-analytics"
)

func TestAnalyzeSalesMainScenario(t *testing.T) {
	sales := []analytics.Sale{
		{Product: "Кофе", Category: "Напитки", Quantity: 2, UnitPriceKopecks: 15_000},
		{Product: "Чай", Category: "Напитки", Quantity: 3, UnitPriceKopecks: 10_000},
		{Product: "Кружка", Category: "Посуда", Quantity: 1, UnitPriceKopecks: 25_000},
	}

	got := analytics.AnalyzeSales(sales)
	want := analytics.Report{
		RevenueKopecks:    85_000,
		UnitsSold:         6,
		RevenueByCategory: map[string]int64{"Напитки": 60_000, "Посуда": 25_000},
		TopProduct:        "Чай",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AnalyzeSales() = %#v, want %#v", got, want)
	}
}

func TestAnalyzeSalesSkipsCancelledAndCountsInvalid(t *testing.T) {
	sales := []analytics.Sale{
		{Product: "Кофе", Category: "Напитки", Quantity: 1, UnitPriceKopecks: 15_000},
		{Product: "", Category: "Напитки", Quantity: 1, UnitPriceKopecks: 10_000},
		{Product: "Отмена", Category: "", Quantity: 0, UnitPriceKopecks: 0, Cancelled: true},
	}

	got := analytics.AnalyzeSales(sales)
	if got.RevenueKopecks != 15_000 || got.UnitsSold != 1 || got.InvalidSales != 1 {
		t.Fatalf("unexpected report: %#v", got)
	}
}

// Check empty input
func TestAnalyzeSalesEmptyInputReturnsInitializedReport(t *testing.T) {
	got := analytics.AnalyzeSales(nil)
	want := analytics.Report{RevenueByCategory: map[string]int64{}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AnalyzeSales(nil) = %#v, want %#v", got, want)
	}
	if got.RevenueByCategory == nil {
		t.Fatal("RevenueByCategory must be initialized for empty input")
	}
}

// Check each version of invalid input
func TestAnalyzeSalesCountsEachInvalidField(t *testing.T) {
	sales := []analytics.Sale{
		{Category: "Напитки", Quantity: 1, UnitPriceKopecks: 100},
		{Product: "Кофе", Quantity: 1, UnitPriceKopecks: 100},
		{Product: "Кофе", Category: "Напитки", Quantity: 0, UnitPriceKopecks: 100},
		{Product: "Кофе", Category: "Напитки", Quantity: -1, UnitPriceKopecks: 100},
		{Product: "Кофе", Category: "Напитки", Quantity: 1, UnitPriceKopecks: 0},
		{Product: "Кофе", Category: "Напитки", Quantity: 1, UnitPriceKopecks: -1},
	}

	got := analytics.AnalyzeSales(sales)
	want := analytics.Report{
		RevenueByCategory: map[string]int64{},
		InvalidSales:      6,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AnalyzeSales() = %#v, want %#v", got, want)
	}
}

// Function should not change input
func TestAnalyzeSalesDoesNotModifyInput(t *testing.T) {
	sales := []analytics.Sale{
		{Product: "Кофе", Category: "Напитки", Quantity: 2, UnitPriceKopecks: 15_000},
		{Product: "", Category: "Посуда", Quantity: 1, UnitPriceKopecks: 25_000},
	}
	original := append([]analytics.Sale(nil), sales...)

	analytics.AnalyzeSales(sales)

	if !reflect.DeepEqual(sales, original) {
		t.Fatalf("AnalyzeSales() modified input: got %#v, want %#v", sales, original)
	}
}
