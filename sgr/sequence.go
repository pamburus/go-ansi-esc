// Package sgr provides facilities for dealing with CSI/SGR ANSI Escape Sequences in a convenient way.
//
// ANSI Escape Sequences are defined in ECMA-48 standard and then extended in ISO 8613-6 [CCITT Recommendation T.416]).
package sgr

// ---

// Sequence is a sequence of commands.
type Sequence []Command

// Render produces binary string corresponding to the contained commands starting with "\x1b[" (ESC/CSI) and ending with 'm'.
func (s Sequence) Render(buf []byte) []byte {
	buf = append(buf, seqBegin...)
	l := len(buf)
	for i := range s {
		if len(buf) != l {
			buf = append(buf, seqNext)
		}
		l = len(buf)
		buf = s[i].render(buf)
	}
	buf = append(buf, seqEnd)

	return buf
}

// Bytes returns binary string with the serialized sequence, the same as produces by Render but collected into a new byte slice.
func (s Sequence) Bytes() []byte {
	return s.Render(make([]byte, 0, 16))
}

// ---

const (
	seqBegin = "\x1b["
	seqNext  = ';'
	seqEnd   = 'm'
	seqReset = seqBegin + string(seqEnd)
)
