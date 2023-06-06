# cloudflare-aname
A simple CLI tool to achieve similar functionality of ANAME, ALIAS and CNAME flattening. Turns any CNAME into an A and AAAA records.

## Usage
Run `cloudflare-aname` in a terminal with `-conf <path-to-config>` arguments.

### Configuration file
The configuration file uses YAML syntax.
```yaml
cloudflare:
  api-token: <cloudflate-api-token>
  zone-id: <cloudflare-zone-id>
record:
  name: something.somecloudflaredomain.com
  target: something.elb.eu-central-1.amazonaws.com
  ttl: 60 # Default value: 60 seconds
```

## Copyright and license
Code released under the [Apache License, Version 2.0](https://github.com/surfshark/cloudflare-aname/blob/master/LICENSE).

## Notice of Non-Affiliation and Disclaimer
We are not affiliated, associated, authorized, endorsed by, or in any way officially connected with Cloudflare, or any of its subsidiaries or its affiliates. The official Cloudflare website can be found at https://www.cloudflare.com/.

The name Cloudflare as well as related names, marks, emblems and images are registered trademarks of their respective owners.
