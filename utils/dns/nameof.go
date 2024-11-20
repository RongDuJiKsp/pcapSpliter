package dns

const (
	ProviderCloudflare = "Cloudflare"
	ProviderGoogle     = "Google"
	ProviderQuad9      = "Quad9"
	ProviderServeroid  = "Serveroid"
	ProviderAdGuard    = "AdGuard (DNS-unfiltered)"
	ProviderTencent    = "Tencent"
	ProviderTWNIC      = "TWNIC"
	ProviderProxy      = "Proxy"
	ProviderUnknown    = "Unknown Provider"
)

var ipToDnsProvider = map[string]string{
	"1.0.0.1":         ProviderCloudflare,
	"1.1.1.1":         ProviderCloudflare,
	"104.16.248.249":  ProviderCloudflare,
	"104.16.249.249":  ProviderCloudflare,
	"8.8.8.8":         ProviderGoogle,
	"8.8.4.4":         ProviderGoogle,
	"9.9.9.9":         ProviderQuad9,
	"9.9.9.10":        ProviderQuad9,
	"9.9.9.11":        ProviderQuad9,
	"149.112.112.10":  ProviderQuad9,
	"149.112.112.11":  ProviderQuad9,
	"149.112.112.112": ProviderQuad9,
	"176.103.130.131": ProviderServeroid,
	"176.103.130.130": ProviderServeroid,
	"94.140.14.140":   ProviderAdGuard,
	"94.140.14.141":   ProviderAdGuard,
	"119.29.29.29":    ProviderTencent,
	"1.12.12.12":      ProviderTencent,
	"120.53.53.53":    ProviderTencent,
	"162.14.21.56":    ProviderTencent,
	"162.14.21.178":   ProviderTencent,
	"101.101.101.101": ProviderTWNIC,
	"101.102.103.104": ProviderTWNIC,
}

var prefixProviders = map[string]string{
	"198.18": ProviderProxy,
}
