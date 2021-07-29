package util

import "testing"

func TestGitClone(t *testing.T) {
	repoUrl, hash := "https://github.com/org/repo.git", "123123412345235432"
	expected := "https://github.com/org/repo/archive/123123412345235432.zip"
	url := ArchiveUrl(repoUrl, hash)
	if url != expected {
		t.Errorf("Incorrect URL created for the archive\n\n should have been: %s\n\nbut was: %s", expected, url)
	}
}
