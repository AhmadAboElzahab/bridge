package seeder

import (
	"io"
	"log"
	"net/http"
	"os"
)

func EnsureDataFilesExist() {
	downloadIfMissing("data/countryInfo.txt", "https://download.geonames.org/export/dump/countryInfo.txt")
}
func downloadIfMissing(path, url string) {
	// Ensure the data folder exists
	dir := "data"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatalf("❌ Failed to create directory: %v", err)
		}
	}

	// Check if the file already exists
	if _, err := os.Stat(path); err == nil {
		log.Printf("✅ File already exists: %s\n", path)
		return
	}

	// Download the file
	log.Printf("⬇️  Downloading %s...\n", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("❌ Failed to download %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("❌ Bad status while downloading %s: %s", url, resp.Status)
	}

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		log.Fatalf("❌ Failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("❌ Failed to save file: %v", err)
	}

	log.Printf("✅ Downloaded and saved %s\n", path)
}
