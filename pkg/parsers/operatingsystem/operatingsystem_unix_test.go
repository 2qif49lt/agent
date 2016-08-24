// +build linux freebsd

package operatingsystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGetOperatingSystem(t *testing.T) {
	var backup = etcOsRelease

	invalids := []struct {
		content       string
		errorExpected string
	}{
		{
			`PRETTY_NAME=Source Mage GNU/Linux
PRETTY_NAME=Ubuntu 14.04.LTS`,
			"PRETTY_NAME needs to be enclosed by quotes if they have spaces: Source Mage GNU/Linux",
		},
		{
			`PRETTY_NAME="Ubuntu Linux
PRETTY_NAME=Ubuntu 14.04.LTS`,
			"PRETTY_NAME is invalid: invalid command line string",
		},
		{
			`PRETTY_NAME=Ubuntu'
PRETTY_NAME=Ubuntu 14.04.LTS`,
			"PRETTY_NAME is invalid: invalid command line string",
		},
		{
			`PRETTY_NAME'
PRETTY_NAME=Ubuntu 14.04.LTS`,
			"PRETTY_NAME needs to be enclosed by quotes if they have spaces: Ubuntu 14.04.LTS",
		},
	}

	valids := []struct {
		content  string
		expected string
	}{
		{
			`NAME="Ubuntu"
PRETTY_NAME_AGAIN="Ubuntu 14.04.LTS"
VERSION="14.04, Trusty Tahr"
ID=ubuntu
ID_LIKE=debian
VERSION_ID="14.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"`,
			"Linux",
		},
		{
			`NAME="Ubuntu"
VERSION="14.04, Trusty Tahr"
ID=ubuntu
ID_LIKE=debian
VERSION_ID="14.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"`,
			"Linux",
		},
		{
			`NAME=Gentoo
ID=gentoo
PRETTY_NAME="Gentoo/Linux"
ANSI_COLOR="1;32"
HOME_URL="http://www.gentoo.org/"
SUPPORT_URL="http://www.gentoo.org/main/en/support.xml"
BUG_REPORT_URL="https://bugs.gentoo.org/"
`,
			"Gentoo/Linux",
		},
		{
			`NAME="Ubuntu"
VERSION="14.04, Trusty Tahr"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 14.04 LTS"
VERSION_ID="14.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"`,
			"Ubuntu 14.04 LTS",
		},
		{
			`NAME="Ubuntu"
VERSION="14.04, Trusty Tahr"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME='Ubuntu 14.04 LTS'`,
			"Ubuntu 14.04 LTS",
		},
		{
			`PRETTY_NAME=Source
NAME="Source Mage"`,
			"Source",
		},
		{
			`PRETTY_NAME=Source
PRETTY_NAME="Source Mage"`,
			"Source Mage",
		},
	}

	dir := os.TempDir()
	etcOsRelease = filepath.Join(dir, "etcOsRelease")

	defer func() {
		os.Remove(etcOsRelease)
		etcOsRelease = backup
	}()

	for _, elt := range invalids {
		if err := ioutil.WriteFile(etcOsRelease, []byte(elt.content), 0600); err != nil {
			t.Fatalf("failed to write to %s: %v", etcOsRelease, err)
		}
		s, err := GetOperatingSystem()
		if err == nil || err.Error() != elt.errorExpected {
			t.Fatalf("Expected an error %q, got %q (err: %v)", elt.errorExpected, s, err)
		}
	}

	for _, elt := range valids {
		if err := ioutil.WriteFile(etcOsRelease, []byte(elt.content), 0600); err != nil {
			t.Fatalf("failed to write to %s: %v", etcOsRelease, err)
		}
		s, err := GetOperatingSystem()
		if err != nil || s != elt.expected {
			t.Fatalf("Expected %q, got %q (err: %v)", elt.expected, s, err)
		}
	}
}

func TestOsReleaseFallback(t *testing.T) {
	var backup = etcOsRelease
	var altBackup = altOsRelease
	dir := os.TempDir()
	etcOsRelease = filepath.Join(dir, "etcOsRelease")
	altOsRelease = filepath.Join(dir, "altOsRelease")

	defer func() {
		os.Remove(dir)
		etcOsRelease = backup
		altOsRelease = altBackup
	}()
	content := `NAME=Gentoo
ID=gentoo
PRETTY_NAME="Gentoo/Linux"
ANSI_COLOR="1;32"
HOME_URL="http://www.gentoo.org/"
SUPPORT_URL="http://www.gentoo.org/main/en/support.xml"
BUG_REPORT_URL="https://bugs.gentoo.org/"
`
	if err := ioutil.WriteFile(altOsRelease, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write to %s: %v", etcOsRelease, err)
	}
	s, err := GetOperatingSystem()
	if err != nil || s != "Gentoo/Linux" {
		t.Fatalf("Expected %q, got %q (err: %v)", "Gentoo/Linux", s, err)
	}
}
