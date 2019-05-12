package main

import "testing"

func Test_parseGitURL(t *testing.T) {
	u := "git@bitbucket.org:someuser/my-cool-repo.git"
	expected := "https://bitbucket.org/someuser/my-cool-repo"

	res, err := parseGitURL(u)
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	if res != expected {
		t.Errorf("expected '%s' got '%s'", expected, res)
	}
}

func Test_getURL(t *testing.T) {
	f, err := readConfigFile("testdata/config")
	if err != nil {
		t.Fatalf("could not read test config: %v", err)
	}
	u, err := getURL(f)
	if err != nil {
		t.Fatalf("failed to get URL: %v", err)
	}
	expected := "git@bitbucket.org:someuser/my-cool-repo.git"
	if u != expected {
		t.Errorf("expected %s got %s", u, expected)
	}
}
