package template

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	transip "github.com/libdns/transip"
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

// TODO: This is just an example. Useful to allow env variable placeholders; update accordingly.
// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.AccountName = repl.ReplaceAll(p.Provider.AccountName, "")
	p.Provider.PrivateKeyPath = repl.ReplaceAll(p.Provider.PrivateKeyPath, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	providername [<account_name> <private_key_path>] {
//	    api_token <api_token>
//		private_key_path <private_key_path>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.AccountName = d.Val()
		}
		if d.NextArg() {
			p.Provider.PrivateKeyPath = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "account_name":
				if p.Provider.AccountName != "" {
					return d.Err("Account Name already set")
				}
				if d.NextArg() {
					p.Provider.AccountName = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "private_key_path":
				if p.Provider.PrivateKeyPath != "" {
					return d.Err("Private Key Path already set")
				}
				if d.NextArg() {
					p.Provider.PrivateKeyPath = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.AccountName == "" {
		return d.Err("Missing AccountName")
	}
	if p.Provider.PrivateKeyPath == "" {
		return d.Err("Missing PrivateKeyPath")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
