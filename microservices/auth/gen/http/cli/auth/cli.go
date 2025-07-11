// Code generated by goa v3.21.1, DO NOT EDIT.
//
// auth HTTP client CLI support package
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/auth/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	authc "object-t.com/hackz-giganoto/microservices/auth/gen/http/auth/client"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `auth (introspect|auth-url|oauth-callback)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` auth introspect --body '{
      "token": "Recusandae non cum perspiciatis error."
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, any, error) {
	var (
		authFlags = flag.NewFlagSet("auth", flag.ContinueOnError)

		authIntrospectFlags    = flag.NewFlagSet("introspect", flag.ExitOnError)
		authIntrospectBodyFlag = authIntrospectFlags.String("body", "REQUIRED", "")

		authAuthURLFlags = flag.NewFlagSet("auth-url", flag.ExitOnError)

		authOauthCallbackFlags     = flag.NewFlagSet("oauth-callback", flag.ExitOnError)
		authOauthCallbackCodeFlag  = authOauthCallbackFlags.String("code", "REQUIRED", "")
		authOauthCallbackStateFlag = authOauthCallbackFlags.String("state", "REQUIRED", "")
	)
	authFlags.Usage = authUsage
	authIntrospectFlags.Usage = authIntrospectUsage
	authAuthURLFlags.Usage = authAuthURLUsage
	authOauthCallbackFlags.Usage = authOauthCallbackUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "auth":
			svcf = authFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "auth":
			switch epn {
			case "introspect":
				epf = authIntrospectFlags

			case "auth-url":
				epf = authAuthURLFlags

			case "oauth-callback":
				epf = authOauthCallbackFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     any
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "auth":
			c := authc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "introspect":
				endpoint = c.Introspect()
				data, err = authc.BuildIntrospectPayload(*authIntrospectBodyFlag)
			case "auth-url":
				endpoint = c.AuthURL()
			case "oauth-callback":
				endpoint = c.OauthCallback()
				data, err = authc.BuildOauthCallbackPayload(*authOauthCallbackCodeFlag, *authOauthCallbackStateFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// authUsage displays the usage of the auth command and its subcommands.
func authUsage() {
	fmt.Fprintf(os.Stderr, `Authentication service that converts opaque tokens to internal JWT tokens for Kong Gateway
Usage:
    %[1]s [globalflags] auth COMMAND [flags]

COMMAND:
    introspect: Introspect opaque token and return internal JWT token for Kong Gateway
    auth-url: Get GitHub OAuth authorization URL with state parameter
    oauth-callback: Handle GitHub OAuth callback and return opaque token

Additional help:
    %[1]s auth COMMAND --help
`, os.Args[0])
}
func authIntrospectUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] auth introspect -body JSON

Introspect opaque token and return internal JWT token for Kong Gateway
    -body JSON: 

Example:
    %[1]s auth introspect --body '{
      "token": "Recusandae non cum perspiciatis error."
   }'
`, os.Args[0])
}

func authAuthURLUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] auth auth-url

Get GitHub OAuth authorization URL with state parameter

Example:
    %[1]s auth auth-url
`, os.Args[0])
}

func authOauthCallbackUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] auth oauth-callback -code STRING -state STRING

Handle GitHub OAuth callback and return opaque token
    -code STRING: 
    -state STRING: 

Example:
    %[1]s auth oauth-callback --code "Laborum quam." --state "Quia sapiente est sed accusamus temporibus."
`, os.Args[0])
}
