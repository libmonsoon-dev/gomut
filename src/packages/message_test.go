package packages

import (
	"fmt"
	"strings"
)

type Message struct {
	PkgName  string
	FileName string
	Code     string
}

func (m Message) AssertMatch(node Node) error {
	if pkgName := node.Package.String(); pkgName != m.PkgName {
		return fmt.Errorf("packages not equal:\ngot:      %#v\nexpected: %#v", pkgName, m.PkgName)
	}

	if !strings.HasSuffix(node.FileName, m.FileName) {
		return fmt.Errorf("FileName not match:\ngot:      %#v\nexpected: %#v", node.FileName, m.FileName)
	}

	if actualCode := node.String(); actualCode != m.Code {
		return fmt.Errorf("code not equal:\n\ngot:      %#v\nexpected: %#v", actualCode, m.Code)
	}

	return nil
}
