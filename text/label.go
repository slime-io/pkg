package text

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	hexTable     = "0123456789abcdef"
	escapeSep    = '-'
	escapeSepStr = string('-')
	valueSep     = '.'
	valueSepStr  = string('.')
)

func escapeLabelValue(v string, buf *bytes.Buffer) {
	isValidChar := func(b byte) bool {
		if 'a' <= b && b <= 'z' {
			return true
		}
		if 'A' <= b && b <= 'Z' {
			return true
		}
		if '0' <= b && b <= '9' {
			return true
		}
		if '_' == b {
			return true
		}
		return false
	}

	for idx := 0; idx < len(v); idx++ {
		c := v[idx]
		if isValidChar(c) {
			buf.WriteByte(c)
		} else {
			buf.WriteByte(hexTable[c>>4])
			buf.WriteByte(escapeSep)
			buf.WriteByte(hexTable[c&0x0f])
		}
	}
}

func EscapeLabelValue(v string) string {
	buf := &bytes.Buffer{}
	escapeLabelValue(v, buf)
	return buf.String()
}

// EscapeLabelValues escapes all values and joins them to a single string
// example: printData2,printData1$ -> printData12-4.printData2 while hex('$') == "24"
func EscapeLabelValues(values []string) string {
	buf := &bytes.Buffer{}

	for _, v := range values {
		if v == "" {
			continue
		}
		escapeLabelValue(v, buf)
		buf.WriteByte(valueSep)
	}

	ret := buf.String()
	if l := len(ret); l > 0 && ret[l-1] == valueSep {
		// remove trailing sep
		ret = ret[:l-1]
	}
	return ret
}

func UnescapeLabelValue(part string) (string, error) {
	// fromHexChar converts a hex character into its value and a success flag.
	fromHexChar := func(c byte) (byte, bool) {
		switch {
		case '0' <= c && c <= '9':
			return c - '0', true
		case 'a' <= c && c <= 'f':
			return c - 'a' + 10, true
		case 'A' <= c && c <= 'F':
			return c - 'A' + 10, true
		}

		return 0, false
	}

	methodParts := strings.Split(part, escapeSepStr)
	if len(methodParts) == 1 {
		return part, nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(part)))
	for idx := 1; idx < len(methodParts); idx++ {
		mpPrev, mpNext := methodParts[idx-1], methodParts[idx]
		if mpPrev == "" || mpNext == "" {
			return "", fmt.Errorf("skip invalid method %s", part)
		}

		buf.WriteString(mpPrev[:len(mpPrev)-1])

		a, ok := fromHexChar(mpPrev[len(mpPrev)-1])
		if !ok {
			return "", fmt.Errorf("skip invalid method %s", part)
		}
		b, ok := fromHexChar(mpNext[0])
		if !ok {
			return "", fmt.Errorf("skip invalid method %s", part)
		}
		buf.WriteByte((a << 4) | b)

		methodParts[idx] = mpNext[1:]
	}

	buf.WriteString(methodParts[len(methodParts)-1])

	return buf.String(), nil
}

func UnescapeLabelValues(s string) ([]string, error) {
	var ret []string
	if s == "" {
		return ret, nil
	}

	parts := strings.Split(s, valueSepStr)
	ret = make([]string, 0, len(parts))

	// fromHexChar converts a hex character into its value and a success flag.
	fromHexChar := func(c byte) (byte, bool) {
		switch {
		case '0' <= c && c <= '9':
			return c - '0', true
		case 'a' <= c && c <= 'f':
			return c - 'a' + 10, true
		case 'A' <= c && c <= 'F':
			return c - 'A' + 10, true
		}

		return 0, false
	}

	for _, part := range parts {
		methodParts := strings.Split(part, escapeSepStr)
		if len(methodParts) == 1 {
			ret = append(ret, part)
			continue
		}

		buf := bytes.NewBuffer(make([]byte, 0, len(part)))
		for idx := 1; idx < len(methodParts); idx++ {
			mpPrev, mpNext := methodParts[idx-1], methodParts[idx]
			if mpPrev == "" || mpNext == "" {
				return nil, fmt.Errorf("skip invalid method %s", part)
			}

			buf.WriteString(mpPrev[:len(mpPrev)-1])

			a, ok := fromHexChar(mpPrev[len(mpPrev)-1])
			if !ok {
				return nil, fmt.Errorf("skip invalid method %s", part)
			}
			b, ok := fromHexChar(mpNext[0])
			if !ok {
				return nil, fmt.Errorf("skip invalid method %s", part)
			}
			buf.WriteByte((a << 4) | b)

			methodParts[idx] = mpNext[1:]
		}
		buf.WriteString(methodParts[len(methodParts)-1])

		ret = append(ret, buf.String())
	}

	return ret, nil
}
