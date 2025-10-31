package template

import (
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/transip"
	"github.com/libdns/transip/client"
	"github.com/pbergman/provider"
	"go.uber.org/zap"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *transip.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.transip",
		New: func() caddy.Module { return &Provider{new(transip.Provider)} },
	}
}

func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()

	p.Provider.AuthLogin = repl.ReplaceAll(p.Provider.AuthLogin, "")
	p.Provider.PrivateKey = repl.ReplaceAll(p.Provider.PrivateKey, "")

	if p.Provider.DebugLevel > 0 {
		p.Provider.DebugOut = zap.NewStdLog(ctx.Logger()).Writer()
	}

	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	providername [<login> <private_key>] {
//	    login 				<login username>
//		private_key 		<private key>
//		full_zone_control 	<update whole zone file at once>
//		debug_level 		<debug level for client 0, 1, 2 or 3>
//		expiration_time     <specifies the time-to-live for an authentication token>
//		not_global_key      <set to true to generate keys that are restricted to clients with IP addresses>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {

		if d.NextArg() {
			p.Provider.AuthLogin = d.Val()
		}

		if d.NextArg() {
			p.Provider.PrivateKey = d.Val()
		}

		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "login":

				if p.Provider.AuthLogin != "" {
					return d.Err("Account Name already set")
				}

				if d.NextArg() {
					p.Provider.AuthLogin = d.Val()
				}

				if d.NextArg() {
					return d.ArgErr()
				}

			case "private_key":

				if p.Provider.PrivateKey != "" {
					return d.Err("Private Key Path already set")
				}

				if d.NextArg() {
					p.Provider.PrivateKey = d.Val()
				}

				if d.NextArg() {
					return d.ArgErr()
				}

			case "debug_level":

				if d.NextArg() {

					if val, err := strconv.Atoi(d.Val()); err == nil && (val > 0 && val < 4) {
						p.Provider.DebugLevel = provider.OutputLevel(val)
					}
				}

				if d.NextArg() {
					return d.ArgErr()
				}

			case "not_global_key":

				if d.NextArg() {
					if val, err := strconv.ParseBool(d.Val()); err == nil {
						p.Provider.AuthNotGlobalKey = val
					}
				}

				if d.NextArg() {
					return d.ArgErr()
				}

			case "expiration_time":

				if p.Provider.ExpirationTime() != "" {
					return d.Err("Expiration Time already set")
				}

				if d.NextArg() {
					p.Provider.AuthExpirationTime = client.ExpirationTime(d.Val())
				}

				if d.NextArg() {
					return d.ArgErr()
				}

			case "full_zone_control":

				if d.NextArg() {
					if val, err := strconv.ParseBool(d.Val()); err == nil && val {
						p.Provider.ClientControl = client.FullZoneControl
					}
				}

				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.AuthLogin == "" {
		return d.Err("Missing account login")
	}
	if p.Provider.PrivateKey == "" {
		return d.Err("Missing private key")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
