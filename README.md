# raepublishing-website-api
Backend api for Rae Publishing's website and blog.


## List of Must Haves

### Session

#### Cookies

* Need to use cookies, properly
* Warn users and have them opt in, but off by default.

#### DB

Wrap page loads in DB transactions

### Security

#### Cross-site scripting and DB injection attacks

String validation of all type-able fields to avoid scripts and DB comments from being submitted.

#### oAuth

Use oAuth for users to login and make comments
