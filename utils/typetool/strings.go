// description:
// @author renshiwei
// Date: 2022/11/16 10:44

package util

import "strings"

// TrimLeft0x Remove 0x to the left of the string.
func TrimLeft0x(s string) string {
	return strings.TrimLeft(s, "0x")
}
