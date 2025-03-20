### Two approaches (monolithic vs distributed) 

(vid 53 section 4)

If you have separate applications for the front end and back end, a purist approach would avoid using posts directly on the front end. Instead, you'd rely on the back-end API for handling everything.

We'll cover that as well, but it's useful to learn both methods. Some of you might prefer a monolithic application, where the front end and back end run together in a single application. This approach works well for that scenario. So, don't worry — we'll create some posts on the back end before we're done.

> The original text talks about two different approaches to building applications:
>
> **Distributed (or Service-Oriented) Architecture:**
> In this setup, you have separate applications for the front end (what users see) and the back end (where data is processed). A more "purist" approach here means using the back-end API for all interactions, like fetching or submitting posts. The front end simply makes requests to the back end, which handles everything.
>
> **Monolithic Architecture:**
> This is when both the front end and back end are combined into one single application. In this case, the front end might directly handle posts without relying on an API. This approach is simpler and easier to manage for smaller applications.

**Copying `main.go` content from `cmd/web` to `cmd/api`**

You might wonder why the types weren’t placed in a shared package. That's because over time the configuration for the API and the configuration for the front end might be markedly different. And I just like to keep things as clean as I possibly can.



