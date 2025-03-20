# Certain explanations

(Section 3)

I'm structuring a single codebase to generate two binaries: one for the front end and one for the back end. While not necessary for small projects, this setup supports scalability and shows how to build multiple binaries from one codebase.

```
/cmd
----/api
----/web
```

For caching, I'll append a version number to CSS and JavaScript files. Incrementing it forces browsers to load the latest version, avoiding manual cache clearing.

```
const cssVersion = "1"
```

For security, Stripe keys will be read from environment variables instead of command-line flags to prevent exposure in process listings.

```
ini


CopyEdit
```

```
cfg.stripe.key = os.Getenv("STRIPE_KEY")
cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

```

### Access to API 

[Chi CORS net/http middleware](https://github.com/go-chi/cors) will give me a really easy middleware package that I can use to determine who has the rights to access my API.







