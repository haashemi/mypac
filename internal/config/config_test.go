package config

import "testing"

func TestParse(t *testing.T) {
	c, err := parse([]byte(`
packages:
  - name: Base
    packages: [base]
lostfiles:
  /opt/foo: manually installed foo
modified:
  ignore_files: [/etc/locale.gen]
`))
	if err != nil {
		t.Fatal(err)
	}
	if c.Packages[0].Name != "Base" || c.LostFiles["/opt/foo"] == "" || c.Modified.IgnoreFiles[0] != "/etc/locale.gen" {
		t.Fatalf("unexpected config: %+v", c)
	}

	legacy, err := parse([]byte("- name: Base\n  packages: [base]\n"))
	if err != nil || legacy.Packages[0].Packages[0] != "base" {
		t.Fatalf("legacy format: %v %+v", err, legacy)
	}
}

func TestLostFileReason(t *testing.T) {
	c := &Config{LostFiles: map[string]string{
		"/opt/foo":        "foo",
		"/opt/foo/cache/": "foo cache",
		"/etc/x.conf":     "x",
	}}

	for file, want := range map[string]string{
		"/opt/foo":            "foo",
		"/opt/foo/bin/foo":    "foo",
		"/opt/foo/cache/a":    "foo cache",
		"/opt/foobar":         "",
		"/etc/x.conf":         "x",
		"/etc/x.conf.pacsave": "",
	} {
		if _, got, _ := c.LostFileReason(file); got != want {
			t.Errorf("%s: got %q, want %q", file, got, want)
		}
	}
}
