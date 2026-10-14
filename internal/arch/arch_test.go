package arch

import (
	"reflect"
	"testing"
)

func TestParseModified(t *testing.T) {
	out := []byte(`Name            : filesystem
Version         : 2025.10.12-1
Description     : Base Arch Linux files
Backup Files    : /etc/crypttab [unreadable]
                  /etc/fstab [modified]
                  /etc/group [modified]
                  /etc/host.conf [unmodified]

Name            : openssh
Backup Files    : /etc/ssh/sshd_config [modified]

Name            : zsh
Backup Files    : /etc/zsh/zprofile [unmodified]
`)

	want := []ModifiedFile{
		{"filesystem", "/etc/fstab"},
		{"filesystem", "/etc/group"},
		{"openssh", "/etc/ssh/sshd_config"},
	}
	if got := parseModified(out); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
