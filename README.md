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
				"login": "YOUR_TRANSIP_ACCOUNT_NAME",
				"private_key": "PATH_TO_YOUR_TRANSIP_PRIVATE_KEY"
			}
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns transip <login> <private_kKey>
}
```

```
# one site
tls {
	dns transip <login> <private_key>
}
```

or alternatively:


```
tls {
	dns transip {
	    login 				<login username>
		private_key 		<private key>
		full_zone_control 	<update whole zone file at once>
    	debug_level 		<debug level for client 0, 1, 2 or 3>
		expiration_time     <specifies the time-to-live for an authentication token>
		not_global_key      <set to true to generate keys that are restricted to clients with IP addresses>	
	}
}
```
