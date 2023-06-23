package sgr_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	. "github.com/pamburus/go-tst/tst"

	"github.com/pamburus/go-ansi-esc/sgr"
)

func TestWriter(tt *testing.T) {
	t := New(tt)

	t.Run("Record", func(t Test) {
		rec := record{
			ts:      time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC),
			level:   0,
			logger:  []byte("tst"),
			message: []byte("hello"),
			fields: []field{
				{[]byte("f1"), "v1"},
				{[]byte("f2"), 10},
				{[]byte("f3"), true},
				{[]byte("f4"), []field{
					{[]byte("f5"), "v5"},
					{[]byte("f6"), 20},
					{[]byte("f7"), false},
				}},
			},
		}

		buf := bytes.NewBuffer(nil)
		err := rec.Write(buf)
		t.Expect(err).To(HaveNotOccurred())
		t.Expect(buf.String()).ToEqual(
			"\x1b[90m2023-06-01T12:00:00Z\x1b[0m |\x1b[36mINF\x1b[0m| \x1b[90mtst:\x1b[0m \x1b[97mhello\x1b[0m \x1b[48;2;10;10;30;92;4mf1\x1b[90;24m:\x1b[96mv1\x1b[0m \x1b[48;2;10;10;30;92;4mf2\x1b[90;24m:\x1b[94m10\x1b[0m \x1b[48;2;10;10;30;92;4mf3\x1b[90;24m:\x1b[91mtrue\x1b[0m \x1b[48;2;10;10;30;92;4mf4\x1b[90;24m:{ \x1b[92;4mf5\x1b[90;24m:\x1b[96mv5\x1b[39m \x1b[92;4mf6\x1b[90;24m:\x1b[94m20\x1b[39m \x1b[92;4mf7\x1b[90;24m:\x1b[91mfalse\x1b[90m }\x1b[0m\n",
		)
	})

	t.Run("Reset", func(t Test) {
		buf := bytes.NewBuffer(nil)
		writer := sgr.NewWriter(buf)
		writer.SetForegroundColor(sgr.Blue)
		writer.SetModes(sgr.Italic.ModeSet(), sgr.ModeAdd)
		writer.PushUnderlineColor(sgr.Red)
		t.Expect(writer.Write([]byte("a"))).ToSucceed().AndResult().To(Equal(1))
		writer.PopUnderlineColor()
		writer.Reset()
		t.Expect(writer.Flush()).ToSucceed()
		t.Expect(buf.String()).ToEqual("\x1b[34;58;5;1;3ma\x1b[0m")
	})

	t.Run("Error", func(t Test) {
		writer := sgr.NewWriter(failingWriter{})
		writer.SetForegroundColor(sgr.Blue)
		t.Expect(writer.Write([]byte("a"))).ToFailWith(errFailingWriterError)
		t.Expect(writer.Flush()).ToFailWith(errFailingWriterError)
	})
}

func BenchmarkWriter(b *testing.B) {
	rec := record{
		ts:      time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC),
		level:   0,
		logger:  []byte("tst"),
		message: []byte("hello"),
		fields: []field{
			{[]byte("f1"), "v1"},
			{[]byte("f2"), 10},
			{[]byte("f3"), true},
			{[]byte("f4"), []field{
				{[]byte("f5"), "v5"},
				{[]byte("f6"), 20},
				{[]byte("f7"), false},
			}},
		},
	}

	buf := bytes.NewBuffer(make([]byte, 0, 2048))
	r := newRenderer(&rec, buf)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i != b.N; i++ {
		buf.Reset()
		err := r.run()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ---

type record struct {
	ts      time.Time
	level   int
	logger  []byte
	message []byte
	fields  []field
}

func (r *record) Write(target io.Writer) error {
	return newRenderer(r, target).run()
}

// ---

func newRenderer(r *record, target io.Writer) *renderer {
	return &renderer{
		r,
		sgr.NewWriter(target),
		make([]byte, 0, 36),
	}
}

type renderer struct {
	r     *record
	w     *sgr.Writer
	tsBuf []byte
}

func (r renderer) run() error {
	r.tsBuf = r.tsBuf[0:0]

	return r.render(
		r.ts,
		r.space,
		r.levelDecorated,
		r.space,
		r.logger,
		r.space,
		r.message,
		r.space,
		r.fields,
		r.end,
	)
}

func (r *renderer) render(parts ...func() error) error {
	for _, part := range parts {
		err := part()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *renderer) ts() error {
	r.w.PushForegroundColor(sgr.BrightBlack)

	err := r.write(r.r.ts.UTC().AppendFormat(r.tsBuf, time.RFC3339Nano))

	r.w.PopForegroundColor()

	return err
}

func (r *renderer) levelDecorated() error {
	r.w.PushForegroundColor(sgr.Default)

	err := r.write(pipeSign)
	if err != nil {
		return err
	}

	err = r.level()
	if err != nil {
		return err
	}

	err = r.write(pipeSign)

	r.w.PopForegroundColor()

	return err
}

func (r *renderer) level() error {
	settings, ok := levelsSettings[r.r.level]
	if !ok {
		settings.name = []byte("???")
	}

	r.w.PushForegroundColor(settings.color)

	err := r.write(settings.name)

	r.w.PopForegroundColor()

	return err
}

func (r *renderer) logger() error {
	if r.r.logger == nil {
		return nil
	}

	r.w.PushForegroundColor(sgr.BrightBlack)

	err := r.write(r.r.logger)
	if err != nil {
		return err
	}

	err = r.punctuation(colon)

	r.w.PopForegroundColor()

	return err
}

func (r *renderer) message() error {
	r.w.PushForegroundColor(sgr.BrightWhite)

	err := r.write(r.r.message)

	r.w.PopForegroundColor()

	return err
}

func (r *renderer) fields() error {
	return r.customFields(r.r.fields...)
}

func (r *renderer) field(f field) error {
	r.w.PushBackgroundColor(sgr.RGB(10, 10, 30))

	err := r.key(f.key)
	if err != nil {
		return err
	}

	err = r.punctuation(colon)
	if err != nil {
		return err
	}

	err = r.value(f.value)

	r.w.PopBackgroundColor()

	return err
}

func (r *renderer) object(fields ...field) error {
	err := r.punctuation(openingBrace)
	if err != nil {
		return err
	}

	err = r.customFields(fields...)
	if err != nil {
		return err
	}

	return r.punctuation(closingBrace)
}

func (r *renderer) customFields(fields ...field) error {
	for i, field := range fields {
		if i != 0 {
			err := r.space()
			if err != nil {
				return err
			}
		}

		err := r.field(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *renderer) key(value []byte) error {
	r.w.PushForegroundColor(sgr.BrightGreen)
	r.w.PushModes(sgr.Underlined.ModeSet(), sgr.ModeAdd)

	err := r.write(value)

	r.w.PopModes()
	r.w.PopForegroundColor()

	return err
}

func (r *renderer) value(value any) error {
	switch v := value.(type) {
	case []field:
		return r.object(v...)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return r.valueWithColor(value, sgr.BrightBlue.Color())
	case bool, nil:
		return r.valueWithColor(value, sgr.BrightRed.Color())
	case string:
		return r.valueWithColor(value, sgr.BrightCyan.Color())
	default:
		return r.valueWithColor(value, sgr.BrightWhite.Color())
	}
}

func (r *renderer) valueWithColor(value any, color sgr.Color) error {
	r.w.PushForegroundColor(color)

	err := r.format("%v", value)

	r.w.PopForegroundColor()

	return err
}

func (r *renderer) space() error {
	return r.write(space)
}

func (r *renderer) end() error {
	return r.write(eol)
}

func (r *renderer) punctuation(data []byte) error {
	r.w.PushForegroundColor(sgr.BrightBlack)

	err := r.write(data)

	r.w.PopForegroundColor()

	return err
}

func (r *renderer) write(data []byte) error {
	_, err := r.w.Write(data)

	return err
}

func (r *renderer) format(format string, a ...any) error {
	_, err := fmt.Fprintf(r.w, format, a...)

	return err
}

type field struct {
	key   []byte
	value any
}

// ---

var levelsSettings = map[int]struct {
	color sgr.Color
	name  []byte
}{
	-4: {
		sgr.Magenta.Color(),
		[]byte("DBG"),
	},
	0: {
		sgr.Cyan.Color(),
		[]byte("INF"),
	},
	4: {
		sgr.BrightYellow.Color(),
		[]byte("WRN"),
	},
	8: {
		sgr.BrightRed.Color(),
		[]byte("ERR"),
	},
}

// ---

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errFailingWriterError
}

// ---

var errFailingWriterError = errors.New("test writer error")

// ---

var space = []byte{' '}
var pipeSign = []byte("|")
var colon = []byte(":")
var openingBrace = []byte("{ ")
var closingBrace = []byte(" }")
var eol = []byte{'\n'}
