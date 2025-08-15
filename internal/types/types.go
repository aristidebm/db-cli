package types

type Format string

const (
	JSON     Format = "json"
	CSV             = "csv"
	HTML            = "html"
	MARKDOWN        = "md"
	LATEX           = "latex"
	ASCIIDOC        = "asciidoc"
	DEFAULT         = "default"
)

type Source struct {
	URL     string `toml:"url"`
	Ping    string `toml:"ping"`
	Connect string `toml:"connect"`
	Query   string `toml:"query"`
	Format  Format `toml:"format"`
}

type Driver struct {
	Ping    string `toml:"ping"`
	Connect string `toml:"connect"`
	Query   string `toml:"query"`
	Format  Format `toml:"format"`
}
