package types

type Format string

const (
	JSON      Format = "json"
	CSV              = "csv"
	HTML             = "html"
	MARKDOWN         = "md"
	LATEX            = "latex"
	ASCIIDOC         = "asciidoc"
	UNALIGNED        = "unaligned"
)

type Source struct {
	URL         string `toml:"url"`
	Interactive string `toml:"interactive"`
	Format      Format `toml:"format"`
}

type Scheme struct {
	Interactive string `toml:"interactive"`
	Format      Format `toml:"format"`
}
