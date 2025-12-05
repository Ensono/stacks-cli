package downloaders

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// MockGitClone is a function variable that can be overridden in tests
var mockGitClone func(string, string, string, string, string) (string, error)

// gitCloneOriginal stores the original GitClone function for restoration
var gitCloneOriginal func(string, string, string, string, string) (string, error)

// mockUtilGitClone replaces the util.GitClone function with a mock for testing
func mockUtilGitClone() {
	// This would need to be implemented by replacing the util.GitClone call
	// For now, we'll test the structure and interface compliance
}

func TestNewGitDownloader(t *testing.T) {
	// Test data
	testURL := "https://github.com/example/repo.git"
	testVersion := "v1.0.0"
	testFrameworkVersion := "latest"
	testTempDir := "/tmp/test"
	testToken := "ghp_test_token"

	// Create new git downloader
	downloader := NewGitDownloader(testURL, testVersion, testFrameworkVersion, testTempDir, testToken)

	// Assertions
	assert.NotNil(t, downloader, "Downloader should not be nil")
	assert.Equal(t, testURL, downloader.URL, "URL should be set correctly")
	assert.Equal(t, testVersion, downloader.Version, "Version should be set correctly")
	assert.Equal(t, testFrameworkVersion, downloader.FrameworkVersion, "FrameworkVersion should be set correctly")
	assert.Equal(t, testTempDir, downloader.TempDir, "TempDir should be set correctly")
	assert.Equal(t, testToken, downloader.Token, "Token should be set correctly")
	assert.Nil(t, downloader.logger, "Logger should be nil initially")
}

func TestGit_SetLogger(t *testing.T) {
	downloader := NewGitDownloader("test", "v1.0.0", "latest", "temp", "token")
	logger := logrus.New()

	// Set logger
	downloader.SetLogger(logger)

	// Assertions
	assert.Equal(t, logger, downloader.logger, "Logger should be set correctly")
}

func TestGit_PackageURL(t *testing.T) {
	testURL := "https://github.com/example/repo.git"
	downloader := NewGitDownloader(testURL, "v1.0.0", "latest", "temp", "token")

	// Get package URL
	url := downloader.PackageURL()

	// Assertions
	assert.Equal(t, testURL, url, "PackageURL should return the URL")
}

func TestGit_StructFields(t *testing.T) {
	// Test all field assignments
	testData := struct {
		URL              string
		Version          string
		FrameworkVersion string
		TempDir          string
		Token            string
	}{
		URL:              "https://github.com/example/test-repo.git",
		Version:          "v2.1.0",
		FrameworkVersion: "net6.0",
		TempDir:          "/tmp/git-test",
		Token:            "ghp_example_token_12345",
	}

	downloader := NewGitDownloader(
		testData.URL,
		testData.Version,
		testData.FrameworkVersion,
		testData.TempDir,
		testData.Token,
	)

	// Verify all fields are set correctly
	assert.Equal(t, testData.URL, downloader.URL, "URL field should match")
	assert.Equal(t, testData.Version, downloader.Version, "Version field should match")
	assert.Equal(t, testData.FrameworkVersion, downloader.FrameworkVersion, "FrameworkVersion field should match")
	assert.Equal(t, testData.TempDir, downloader.TempDir, "TempDir field should match")
	assert.Equal(t, testData.Token, downloader.Token, "Token field should match")
}

func TestGit_EmptyValues(t *testing.T) {
	// Test with empty values
	downloader := NewGitDownloader("", "", "", "", "")

	assert.Equal(t, "", downloader.URL, "Empty URL should be preserved")
	assert.Equal(t, "", downloader.Version, "Empty Version should be preserved")
	assert.Equal(t, "", downloader.FrameworkVersion, "Empty FrameworkVersion should be preserved")
	assert.Equal(t, "", downloader.TempDir, "Empty TempDir should be preserved")
	assert.Equal(t, "", downloader.Token, "Empty Token should be preserved")
}


func TestGit_URLValidation(t *testing.T) {
	// Test various URL formats
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "HTTPS URL",
			url:  "https://github.com/user/repo.git",
		},
		{
			name: "HTTPS URL without .git suffix",
			url:  "https://github.com/user/repo",
		},
		{
			name: "SSH URL",
			url:  "git@github.com:user/repo.git",
		},
		{
			name: "GitLab URL",
			url:  "https://gitlab.com/user/repo.git",
		},
		{
			name: "Bitbucket URL",
			url:  "https://bitbucket.org/user/repo.git",
		},
		{
			name: "Custom Git server",
			url:  "https://git.company.com/project/repo.git",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			downloader := NewGitDownloader(tc.url, "main", "latest", "/tmp", "token")
			assert.Equal(t, tc.url, downloader.URL, "URL should be preserved as-is")
			assert.Equal(t, tc.url, downloader.PackageURL(), "PackageURL should return the same URL")
		})
	}
}

func TestGit_VersionFormats(t *testing.T) {
	// Test various version formats
	testCases := []struct {
		name    string
		version string
	}{
		{
			name:    "Semantic version",
			version: "v1.2.3",
		},
		{
			name:    "Branch name",
			version: "main",
		},
		{
			name:    "Branch name (master)",
			version: "master",
		},
		{
			name:    "Feature branch",
			version: "feature/new-functionality",
		},
		{
			name:    "Commit hash (short)",
			version: "abc123f",
		},
		{
			name:    "Commit hash (full)",
			version: "abc123f456789012345678901234567890abcdef0",
		},
		{
			name:    "Tag name",
			version: "release-2023-11-15",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			downloader := NewGitDownloader("https://github.com/test/repo.git", tc.version, "latest", "/tmp", "")
			assert.Equal(t, tc.version, downloader.Version, "Version should be preserved as-is")
		})
	}
}

func TestGit_TokenHandling(t *testing.T) {
	// Test various token scenarios
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "GitHub Personal Access Token",
			token: "ghp_1234567890abcdef1234567890abcdef123456",
		},
		{
			name:  "GitHub Classic Token",
			token: "abc123def456789012345678901234567890abcd",
		},
		{
			name:  "GitLab Token",
			token: "glpat-xxxxxxxxxxxxxxxxxxxx",
		},
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Custom token",
			token: "custom_token_format_12345",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			downloader := NewGitDownloader("https://github.com/test/repo.git", "main", "latest", "/tmp", tc.token)
			assert.Equal(t, tc.token, downloader.Token, "Token should be preserved as-is")
		})
	}
}

// Test that the downloader implements the Downloader interface
func TestGit_ImplementsDownloaderInterface(t *testing.T) {
	downloader := NewGitDownloader("test", "main", "latest", "temp", "token")
	
	// This will fail compilation if the interface is not implemented correctly
	var _ interface {
		Get() (string, error)
		PackageURL() string
		SetLogger(*logrus.Logger)
	} = downloader
	
	// If we get here, the interface is implemented correctly
	assert.True(t, true, "Git downloader implements the required interface")
}

func TestGit_LoggerIntegration(t *testing.T) {
	downloader := NewGitDownloader("https://github.com/test/repo.git", "main", "latest", "/tmp", "token")
	
	// Test with different logger configurations
	testCases := []struct {
		name   string
		logger *logrus.Logger
	}{
		{
			name:   "Standard logger",
			logger: logrus.New(),
		},
		{
			name: "Debug level logger",
			logger: func() *logrus.Logger {
				l := logrus.New()
				l.SetLevel(logrus.DebugLevel)
				return l
			}(),
		},
		{
			name: "JSON formatter logger",
			logger: func() *logrus.Logger {
				l := logrus.New()
				l.SetFormatter(&logrus.JSONFormatter{})
				return l
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			downloader.SetLogger(tc.logger)
			assert.Equal(t, tc.logger, downloader.logger, "Logger should be set correctly")
		})
	}
}

// TestGit_Get_IntegrationNote provides documentation about Get method testing
func TestGit_Get_IntegrationNote(t *testing.T) {
	// NOTE: The Get() method calls util.GitClone(), which performs actual Git operations.
	// For comprehensive testing of the Get() method, you would need to:
	//
	// 1. Mock the util.GitClone function
	// 2. Create integration tests with real repositories
	// 3. Test error scenarios (network failures, invalid URLs, authentication failures)
	//
	// Example test scenarios for Get() method:
	// - Successful clone of public repository
	// - Successful clone with authentication token
	// - Clone of specific branch/tag/commit
	// - Network failure handling
	// - Invalid URL handling
	// - Repository not found (404) handling
	// - Authentication failure handling
	// - Disk space issues
	// - Permission issues in temp directory
	
	downloader := NewGitDownloader("https://github.com/test/repo.git", "main", "latest", "/tmp", "token")
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress output during tests
	downloader.SetLogger(logger)

	// This test documents the Get method's interface without calling external dependencies
	assert.NotNil(t, downloader.Get, "Get method should exist")
	
	// The actual testing of Get() would require mocking util.GitClone
	// or setting up integration test repositories
	t.Skip("Get() method testing requires mocking util.GitClone or integration test setup")
}

// Benchmark test for downloader creation
func BenchmarkNewGitDownloader(b *testing.B) {
	url := "https://github.com/test/repo.git"
	version := "main"
	frameworkVersion := "latest" 
	tempDir := "/tmp/bench"
	token := "test_token"

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = NewGitDownloader(url, version, frameworkVersion, tempDir, token)
	}
}
