package util

import "os"

func IsRoot() bool {
	return os.Geteuid() == 0
}

const runAsRootConfirm = "Warning: Running doth as root is not recommended. Please run doth as a regular user.\nDo you want to continue anyway?"

func ConfirmRunIfRoot() bool {
	if IsRoot() {
		return ConfirmAction(runAsRootConfirm) == nil
	}

	return true
}
