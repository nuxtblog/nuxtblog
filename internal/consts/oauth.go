package consts

// OAuthProviderDef describes a built-in OAuth provider implemented in code.
// Credentials (client_id / client_secret) live in config.yaml — never here.
// To add a new provider: implement oauth.Provider, register via init(), then add an entry below.
type OAuthProviderDef struct {
	Slug    string // must match the oauth registry key and config key (auth.oauth.<slug>)
	LabelZh string
	LabelEn string
	Icon    string // Tabler icon name used by the frontend (i-tabler-brand-xxx)
}

// OAuth provider slug constants — use these instead of magic strings.
const (
	OAuthProviderGitHub = "github"
	OAuthProviderGoogle = "google"
	OAuthProviderQQ     = "qq"
)

// BuiltinOAuthProviders lists every OAuth provider implemented in this codebase.
// The admin settings page reads this list to show which providers can be configured.
// Adding a provider requires: (1) implement oauth.Provider, (2) register via init(),
// (3) add an entry here.
var BuiltinOAuthProviders = []OAuthProviderDef{
	{Slug: OAuthProviderGitHub, LabelZh: "GitHub", LabelEn: "GitHub", Icon: "i-tabler-brand-github"},
	{Slug: OAuthProviderGoogle, LabelZh: "Google", LabelEn: "Google", Icon: "i-tabler-brand-google"},
	{Slug: OAuthProviderQQ,     LabelZh: "QQ",     LabelEn: "QQ",     Icon: "i-tabler-brand-qq"},
}
