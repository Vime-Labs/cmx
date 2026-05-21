package validate_test

import (
	"testing"

	"github.com/Vime-Labs/cmx/internal/validate"
)

func TestPorts(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"3000", false},
		{"3000,3001", false},
		{"80", false},
		{"65535", false},
		{"1", false},
		{"0", true},
		{"65536", true},
		{"abc", true},
		{"3000,abc", true},
		{"3000, 3001", false}, // espaço tolerado no split
		{"", true},
	}
	for _, c := range cases {
		err := validate.Ports(c.input)
		if (err != nil) != c.wantErr {
			t.Errorf("Ports(%q): got err=%v, wantErr=%v", c.input, err, c.wantErr)
		}
	}
}

func TestRepoFormat(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"Vime-Labs/cmx", false},
		{"org/repo.git", false},
		{"org/repo-name", false},
		{"org/repo_name", false},
		{"invalid", true},
		{"org/", true},
		{"/repo", true},
		{"org/repo/extra", true},
		{"", true},
	}
	for _, c := range cases {
		err := validate.RepoFormat(c.input)
		if (err != nil) != c.wantErr {
			t.Errorf("RepoFormat(%q): got err=%v, wantErr=%v", c.input, err, c.wantErr)
		}
	}
}

func TestBranchName(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"main", false},
		{"feature/my-feature", false},
		{"fix-123", false},
		{"branch with space", true},
		{"branch\twith\ttab", true},
		{"", true},
	}
	for _, c := range cases {
		err := validate.BranchName(c.input)
		if (err != nil) != c.wantErr {
			t.Errorf("BranchName(%q): got err=%v, wantErr=%v", c.input, err, c.wantErr)
		}
	}
}

func TestResourceName(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"my-app", false},
		{"mydb", false},
		{"my app", true},
		{"my/app", true},
		{"my\\app", true},
		{"", true},
	}
	for _, c := range cases {
		err := validate.ResourceName(c.input)
		if (err != nil) != c.wantErr {
			t.Errorf("ResourceName(%q): got err=%v, wantErr=%v", c.input, err, c.wantErr)
		}
	}
}

func TestKeyValue(t *testing.T) {
	cases := []struct {
		args      []string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{[]string{"KEY=VALUE"}, "KEY", "VALUE", false},
		{[]string{"KEY", "VALUE"}, "KEY", "VALUE", false},
		{[]string{"KEY=VALUE=WITH=EQUALS"}, "KEY", "VALUE=WITH=EQUALS", false},
		{[]string{"=VALUE"}, "", "", true},
		{[]string{"NOEQUALS"}, "", "", true},
		{[]string{"KEY", "VALUE", "EXTRA"}, "", "", true},
		{[]string{}, "", "", true},
	}
	for _, c := range cases {
		k, v, err := validate.KeyValue(c.args)
		if (err != nil) != c.wantErr {
			t.Errorf("KeyValue(%v): got err=%v, wantErr=%v", c.args, err, c.wantErr)
			continue
		}
		if !c.wantErr && (k != c.wantKey || v != c.wantValue) {
			t.Errorf("KeyValue(%v): got (%q, %q), want (%q, %q)", c.args, k, v, c.wantKey, c.wantValue)
		}
	}
}
