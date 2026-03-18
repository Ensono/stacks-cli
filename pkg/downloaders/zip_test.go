package downloaders

import (
	"archive/zip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewZipDownloader(t *testing.T) {
	testURL := "https://github.com/example/repo/releases/download/v1.0.0/module-1.0.0.zip"
	testVersion := "v1.0.0"
	testCacheDir := "/tmp/cache"
	testTempDir := "/tmp/test"
	testToken := "ghp_test_token"

	downloader := NewZipDownloader(testURL, testVersion, testCacheDir, testTempDir, testToken)

	assert.NotNil(t, downloader, "Downloader should not be nil")
	assert.Equal(t, testURL, downloader.URL, "URL should be set correctly")
	assert.Equal(t, testVersion, downloader.Version, "Version should be set correctly")
	assert.Equal(t, testCacheDir, downloader.CacheDir, "CacheDir should be set correctly")
	assert.Equal(t, testTempDir, downloader.TempDir, "TempDir should be set correctly")
	assert.Equal(t, testToken, downloader.Token, "Token should be set correctly")
	assert.Nil(t, downloader.logger, "Logger should be nil initially")
}

func TestZip_SetLogger(t *testing.T) {
	downloader := NewZipDownloader("test", "v1.0.0", "cache", "temp", "")
	logger := logrus.New()

	downloader.SetLogger(logger)

	assert.Equal(t, logger, downloader.logger, "Logger should be set correctly")
}

func TestZip_PackageURL(t *testing.T) {
	testURL := "https://github.com/example/repo/releases/download/v1.0.0/module-1.0.0.zip"
	downloader := NewZipDownloader(testURL, "v1.0.0", "cache", "temp", "")

	url := downloader.PackageURL()

	assert.Equal(t, testURL, url, "PackageURL should return the URL")
}

func TestZip_EmptyValues(t *testing.T) {
	downloader := NewZipDownloader("", "", "", "", "")

	assert.Empty(t, downloader.URL, "URL should be empty")
	assert.Empty(t, downloader.Version, "Version should be empty")
	assert.Empty(t, downloader.CacheDir, "CacheDir should be empty")
	assert.Empty(t, downloader.TempDir, "TempDir should be empty")
	assert.Empty(t, downloader.Token, "Token should be empty")
}

func TestZip_IndependentInstances(t *testing.T) {
	d1 := NewZipDownloader("url1", "v1", "cache1", "temp1", "token1")
	d2 := NewZipDownloader("url2", "v2", "cache2", "temp2", "token2")

	assert.Equal(t, "url1", d1.URL, "First instance URL should remain unchanged")
	assert.Equal(t, "url2", d2.URL, "Second instance URL should be set correctly")
	assert.NotEqual(t, d1.URL, d2.URL, "Instances should have independent URLs")
}

// createTestZip creates a zip file containing the given files at the specified path.
// files is a map of filename -> content.
func createTestZip(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()

	f, err := os.Create(zipPath)
	require.NoError(t, err, "Should create zip file")
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	for name, content := range files {
		fw, err := w.Create(name)
		require.NoError(t, err, "Should create file in zip: %s", name)
		_, err = fw.Write([]byte(content))
		require.NoError(t, err, "Should write content to file in zip: %s", name)
	}
}

func TestZip_Get_Success(t *testing.T) {
	// Create a test zip file to serve
	serverDir := t.TempDir()
	zipPath := filepath.Join(serverDir, "module-1.0.0.zip")
	createTestZip(t, zipPath, map[string]string{
		"my-module/main.tf":      "resource \"azurerm_resource_group\" \"example\" {}",
		"my-module/variables.tf": "variable \"name\" { type = string }",
		"my-module/outputs.tf":   "output \"id\" { value = azurerm_resource_group.example.id }",
	})

	// Serve the zip file via httptest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, zipPath)
	}))
	defer server.Close()

	cacheDir := filepath.Join(t.TempDir(), "cache")
	tempDir := filepath.Join(t.TempDir(), "temp")

	downloader := NewZipDownloader(
		fmt.Sprintf("%s/module-1.0.0.zip", server.URL),
		"v1.0.0",
		cacheDir,
		tempDir,
		"",
	)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	downloader.SetLogger(logger)

	dir, err := downloader.Get()

	assert.NoError(t, err, "Get should not return an error")
	assert.NotEmpty(t, dir, "Get should return a directory path")

	// Verify the cached zip file exists
	assert.FileExists(t, filepath.Join(cacheDir, "module-1.0.0.zip"), "Zip file should be cached")

	// Verify extracted files exist (util.Unzip returns the first subdirectory)
	assert.FileExists(t, filepath.Join(dir, "main.tf"), "main.tf should be extracted")
	assert.FileExists(t, filepath.Join(dir, "variables.tf"), "variables.tf should be extracted")
	assert.FileExists(t, filepath.Join(dir, "outputs.tf"), "outputs.tf should be extracted")

	// Verify content
	content, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	assert.NoError(t, err, "Should be able to read extracted file")
	assert.Contains(t, string(content), "azurerm_resource_group", "File content should be preserved")
}

func TestZip_Get_CachedFile(t *testing.T) {
	// Create a test zip file pre-cached
	cacheDir := filepath.Join(t.TempDir(), "cache")
	tempDir := filepath.Join(t.TempDir(), "temp")

	require.NoError(t, os.MkdirAll(cacheDir, 0755))

	zipPath := filepath.Join(cacheDir, "module-1.0.0.zip")
	createTestZip(t, zipPath, map[string]string{
		"my-module/main.tf": "# cached module",
	})

	// Track whether the server was called
	serverCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		http.ServeFile(w, r, zipPath)
	}))
	defer server.Close()

	downloader := NewZipDownloader(
		fmt.Sprintf("%s/module-1.0.0.zip", server.URL),
		"v1.0.0",
		cacheDir,
		tempDir,
		"",
	)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	downloader.SetLogger(logger)

	dir, err := downloader.Get()

	assert.NoError(t, err, "Get should not return an error")
	assert.NotEmpty(t, dir, "Get should return a directory path")
	assert.False(t, serverCalled, "Server should not be called when zip is already cached")
}

func TestZip_Get_InvalidURL(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	tempDir := filepath.Join(t.TempDir(), "temp")

	downloader := NewZipDownloader(
		"http://localhost:1/nonexistent.zip",
		"v1.0.0",
		cacheDir,
		tempDir,
		"",
	)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	downloader.SetLogger(logger)

	_, err := downloader.Get()

	assert.Error(t, err, "Get should return an error for unreachable URL")
}

func TestZip_Get_CorruptedZip(t *testing.T) {
	// Serve non-zip content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is not a zip file"))
	}))
	defer server.Close()

	cacheDir := filepath.Join(t.TempDir(), "cache")
	tempDir := filepath.Join(t.TempDir(), "temp")

	downloader := NewZipDownloader(
		fmt.Sprintf("%s/bad.zip", server.URL),
		"v1.0.0",
		cacheDir,
		tempDir,
		"",
	)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	downloader.SetLogger(logger)

	_, err := downloader.Get()

	assert.Error(t, err, "Get should return an error for corrupted zip file")
}

func TestZip_Get_WithToken(t *testing.T) {
	// Verify the token is forwarded (server checks for Authorization header)
	expectedToken := "ghp_test_token_12345"

	serverDir := t.TempDir()
	zipPath := filepath.Join(serverDir, "module-1.0.0.zip")
	createTestZip(t, zipPath, map[string]string{
		"my-module/main.tf": "# private module",
	})

	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		http.ServeFile(w, r, zipPath)
	}))
	defer server.Close()

	cacheDir := filepath.Join(t.TempDir(), "cache")
	tempDir := filepath.Join(t.TempDir(), "temp")

	downloader := NewZipDownloader(
		fmt.Sprintf("%s/module-1.0.0.zip", server.URL),
		"v1.0.0",
		cacheDir,
		tempDir,
		expectedToken,
	)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	downloader.SetLogger(logger)

	_, err := downloader.Get()

	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, fmt.Sprintf("token %s", expectedToken), receivedAuth, "Authorization header should contain the token")
}

func TestZip_Get_MultipleModulesInZip(t *testing.T) {
	// Simulate a mono-repo release artifact with multiple modules
	serverDir := t.TempDir()
	zipPath := filepath.Join(serverDir, "modules-1.0.0.zip")
	createTestZip(t, zipPath, map[string]string{
		"infra-modules/networking/main.tf":      "# networking module",
		"infra-modules/networking/variables.tf": "# networking vars",
		"infra-modules/compute/main.tf":         "# compute module",
		"infra-modules/compute/variables.tf":    "# compute vars",
		"infra-modules/storage/main.tf":         "# storage module",
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, zipPath)
	}))
	defer server.Close()

	cacheDir := filepath.Join(t.TempDir(), "cache")
	tempDir := filepath.Join(t.TempDir(), "temp")

	downloader := NewZipDownloader(
		fmt.Sprintf("%s/modules-1.0.0.zip", server.URL),
		"v1.0.0",
		cacheDir,
		tempDir,
		"",
	)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	downloader.SetLogger(logger)

	dir, err := downloader.Get()

	assert.NoError(t, err, "Get should not return an error")
	assert.DirExists(t, dir, "Returned directory should exist")

	// Verify all modules from the mono-repo are extracted
	assert.FileExists(t, filepath.Join(dir, "networking/main.tf"), "networking module should be extracted")
	assert.FileExists(t, filepath.Join(dir, "compute/main.tf"), "compute module should be extracted")
	assert.FileExists(t, filepath.Join(dir, "storage/main.tf"), "storage module should be extracted")
}
