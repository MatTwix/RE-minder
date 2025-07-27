package oauth

var providers = make(map[string]OAuthProvider)

func RegisterProvider(p OAuthProvider) {
	providers[p.Platform()] = p
}

func GetProvider(platform string) (OAuthProvider, bool) {
	p, ok := providers[platform]
	return p, ok
}
