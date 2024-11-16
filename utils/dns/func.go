package dns

import "strings"

func ConvertIPToDNSProvider(ip string) string {
	if provider, exists := ipToDnsProvider[ip]; exists {
		return provider
	}
	for prefix, provider := range ipToDnsProvider {
		if strings.HasPrefix(ip, prefix) {
			return provider
		}
	}
	return ProviderUnknown
}
