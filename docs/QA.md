# Questions and answers 



###  Just curious why you went about with this project structure

Zachary (lecture 13)

Not really a complaint or argument and I know there is no right or wrong way to structure a project but isn't it kind of standard in most go projects to do a structure similar to the following:

```
cmd/
  _appname_/
    main.go
internal/
  internal app logic here
pkg/
  anything consumed by outside go projects
web/
  templates/
  other static assets
```

Trevor - Instructor

Not in my experience. For some time, a lot of people thought that [this GitHub repository](https://github.com/golang-standards/project-layout) was the standard, accepted way to organize a Go project, simply because the author decided to call it a standard. Then they realized that just because you call something a standard, or, as the author puts it, "a set of common historical and emerging project layout patterns in the Go ecosystem," does not make it so.

And just for clarification, the `cmd` folder rarely holds `__appname__`, at least in my experience; instead, it holds one or more folders that describes what the main package inside that folder (and it is very nearly always the main package) does, so `cmd/api `is the api build, `cmd/cli` is the command line app, `cmd/web `is the web app, etc. Of course, as you say, there is no right or wrong way to structure a project.

I use this layout simply because it works for me, and is remarkably similar to the layouts used in many of the Go shops that I have done work with over the years.

***



#### app *application Receiver

*Sindri:* *Not related to this section, but I'm wondering why do we keep creating functions with a receiver to a pointer of the application struct. I get that we want a pointer so that we can modify the contents of it/prevent copying application type but it sort feels like these functions are basically just global functions since the application is a global type.*

Trevor: The app variable does not need to be global. I can be set in the main function, and then is limited in scope to the functions that use it as a receiver.

*Sindri: That makes sense. What I was trying to say, maybe less clearly, was if the application type is the receiver for all these functions isn't that basically too big of scope? It's basically global because all the functions have access to a complete shared type. Just wondering from a design perspective how else this could be implemented.*

Trevor: Any handler should, in the vast majority of cases, have access to the receiver if you want to share the config, database pool, etc. Anything that does not need this sort of thing does not need to use the receiver, obviously. I wish I could give you a "this approach covers all situations" answer, but the reality is, as is the case with so many things in computer science, that it is completely dependent on what problem you are trying to solve.

***



### Querying SQL in Web Templates instead of API

**Question** 

I am getting confused about the structure of this application. Usually, I would think that querying the database happens in the API instead of the web. However, in this application, we are querying the database from the web/handlers.go itself.

Is this the idea? Correct me if I am wrong.
API is for calling Stripe.
Web is for calling our services like querying our own database.

In both routes-api.go and routes.go, we are calling GetWidget. Why do we need to call them twice?

**Answer (Trevor)**

This is the idea, and you have it mostly right.

When you say "usually, I would think that querying the database happens in the API," you are referring to sites that do not have server side rendering, and all rendering is done on the client site (think React or Vue). This is not a React/Vue application; everything is rendered server side. When you say "in this application, we are querying the database from the `web/handlers.go`" you are not quite right; querying is not done at the handler level; querying is done in the models package, which is called by the handlers.

We call GetWidget in both routes-api.go and routes.go because the same codebase is being used to generate **two different applications** -- one that uses one set of routes, and one that uses another.

***



### Why int instead of decimal

**Question**: You are using int and last 2 digits as a floating points like. You actually storing prices in cent and mention floating errors. What could be errors you mention I will face in production. What could happen if I store them like decimal(15,2)? 

**Trevor's Answer**: This what is known as a rounding error. Computers are not perfect when it comes to floating point errors --- here is a reasonably good overview of the problem:

https://blog.penjee.com/what-are-floating-point-errors-answered/

**Question to asnwer**: The blog __ talk about the float type. is very different to decimal(15,2). The correct way to store in any database any money type is using decimal(XX,2). I usually use decimal(16,2)