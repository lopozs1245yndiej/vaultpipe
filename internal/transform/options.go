package transform

// Options holds configuration for building a Transformer from named presets.
type Options struct {
	Uppercase   bool
	Prefix      string
	TrimSpace   bool
	ReplaceFrom string
	ReplaceTo   string
}

// NewFromOptions constructs a Transformer from an Options struct.
func NewFromOptions(opts Options) *Transformer {
	var fns []TransformFunc

	if opts.ReplaceFrom != "" {
		fns = append(fns, ReplaceKeyChars(opts.ReplaceFrom, opts.ReplaceTo))
	}
	if opts.Uppercase {
		fns = append(fns, UppercaseKeys())
	}
	if opts.Prefix != "" {
		fns = append(fns, PrefixKeys(opts.Prefix))
	}
	if opts.TrimSpace {
		fns = append(fns, TrimValueSpace())
	}

	return New(fns...)
}
