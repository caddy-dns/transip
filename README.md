TransIP module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with [TransIP](https://www.transip.eu/).
It makes use of [libdns/transip](https://github.com/libdns/transip)

## Caddy module name

```
dns.providers.transip
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "transip",
				"account_name": "YOUR_TRANSIP_ACCOUNT_NAME",
				"private_key_path": "PATH_TO_YOUR_TRANSIP_PRIVATE_KEY"
			}
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns transip <accountName> <privateKeyPath>
}
```

```
# one site
tls {
	dns transip <accountName> <privateKeyPath>
}
```

or alternatively:


```
tls {
	dns transip {
		account_name <accountName> 
		private_key_path <privateKeyPath>
	}
}
```
