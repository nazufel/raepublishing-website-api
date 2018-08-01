# raepublishing-website-api
Backend api for Rae Publishing's website and blog.


## Possible Endpoint Layout

```
var host string = http://raepublishing.com/api/v1/

// GET all users. Only one method applied here and that's GET
host/users/ - returns all users, will apply filters later

// Where the real CRUD happens
host/users/:id
```

## List of Must Haves (in no particular order)

* Create session package to control dialing to db.
* CRUD of Blog article resources - Sort of implimented.
* Serving of static pages: ```About```, ```Contact Us```, ```index```, ```site map```.
* Need to use cookies, properly by warning users and have them opt in, but off by default.
* Wrap page loads in DB transactions - not sure if this is totally necessary, but the [Go Buffalo](https://gobuffalo.io/en) project does this.
* String validation of all type-able fields to avoid scripts and DB comments from being submitted.
* Use oAuth for users to login and make comments
* Filtering of results: limit - ```?limit=10```, sort - ```?sortby=name&order=asc```, etc. More filtering will be addressed as needed.
* Limit endpoints to a specific resource or collection of resources. Good API first design uses [nouns](https://www.youtube.com/watch?v=sMKsmZbpyjE), not verbs.
* API versions
* Responding with status codes and custom errors where needed
* URL redirection and returning where needed
* Containerize the app
* Change dependency from Niemeyer's MGO to [Globalsign](https://github.com/globalsign/mgo).
* Change Controllers to Handlers
* Separate routes from ```main.go```.
* Unit Tests: Handler and routes
* Documentation such as [Swagger](https://github.com/go-swagger/go-swagger)
* Handle time. Probably write to DB in UTC, but allow the front end to look up user timezone and adjust UTC to local time when presenting. Will have to think about how to handle global times.
