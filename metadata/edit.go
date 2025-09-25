package metadata

type EditMetadataFunc func(Metadata) error

func EditChain(fns ...EditMetadataFunc) EditMetadataFunc {
	return func(meta Metadata) error {
		for _, fn := range fns {
			if err := fn(meta); err != nil {
				return err
			}
		}
		return nil
	}
}

func Replace(key string, value string) EditMetadataFunc {
	return func(meta Metadata) error {
		meta[key] = value
		return nil
	}
}

func ReplaceRegistrationEndpoint(url string) EditMetadataFunc {
	return Replace(RegistrationEndpointKey, url)
}

func ReplaceAuthorizationEndpoint(url string) EditMetadataFunc {
	return Replace(AuthorizationEndpointKey, url)
}

func ReplaceTokenEndpoint(url string) EditMetadataFunc {
	return Replace(TokenEndpointKey, url)
}
