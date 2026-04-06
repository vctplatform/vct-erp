package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Fetches all Vietnam administrative divisions from provinces.open-api.vn v2
// and saves to backend/data/vietnam_divisions.json

type Province struct {
	Name         string `json:"name"`
	Code         int    `json:"code"`
	DivisionType string `json:"division_type"`
	Codename     string `json:"codename"`
	PhoneCode    int    `json:"phone_code"`
	Wards        []Ward `json:"wards"`
}

type Ward struct {
	Name         string `json:"name"`
	Code         int    `json:"code"`
	DivisionType string `json:"division_type"`
	Codename     string `json:"codename"`
	ProvinceCode int    `json:"province_code"`
}

func main() {
	// Step 1: Get all provinces
	fmt.Println("📡 Fetching provinces list...")
	resp, err := http.Get("https://provinces.open-api.vn/api/v2/p/")
	if err != nil {
		fmt.Printf("❌ Failed to fetch provinces: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var provinces []Province
	if err := json.Unmarshal(body, &provinces); err != nil {
		fmt.Printf("❌ Failed to parse provinces: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Found %d provinces\n", len(provinces))

	// Step 2: Fetch wards for each province
	totalWards := 0
	for i := range provinces {
		p := &provinces[i]
		url := fmt.Sprintf("https://provinces.open-api.vn/api/v2/p/%d?depth=2", p.Code)
		fmt.Printf("  [%d/%d] Fetching wards for %s (code=%d)...", i+1, len(provinces), p.Name, p.Code)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf(" ❌ Error: %v\n", err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var full Province
		if err := json.Unmarshal(body, &full); err != nil {
			fmt.Printf(" ❌ Parse error: %v\n", err)
			continue
		}
		p.Wards = full.Wards
		totalWards += len(full.Wards)
		fmt.Printf(" %d wards\n", len(full.Wards))

		// Rate limit: 100ms between requests
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n📊 Total: %d provinces, %d wards\n", len(provinces), totalWards)

	// Step 3: Write to JSON file
	outDir := "data"
	os.MkdirAll(outDir, 0755)
	outPath := outDir + "/vietnam_divisions.json"

	data, err := json.MarshalIndent(provinces, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outPath, data, 0644); err != nil {
		fmt.Printf("❌ Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Saved to %s (%d bytes)\n", outPath, len(data))
}
