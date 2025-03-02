# Go Stripe

This is the repository for a secure ecommerce application built while following Trevor Sawler's course:
["Building Web Applications with Go - Intermediate Level"](https://www.udemy.com/course/building-web-applications-with-go-intermediate-level/) 

Learn how to build a secure ecommerce application with Go (go-stripe) 


### Technical Stack
- Built with Go version 1.23.3

### Dependencies
- [Chi Router](https://github.com/go-chi/chi/v5) - routing and middleware
- [Chi CORS](https://github.com/go-chi/cors) - CORS net/http middleware
- [Alex Edwards SCS](https://github.com/alexedwards/scs/v2) - session management 
- [Alex Edwards mysqlstore](https://github.com/alexedwards/scs/mysqlstore) - to store session in the db
- [GoDotEnv](https://github.com/joho/godotenv) - loads environment variables from .env files
- [Go Stripe](https://github.com/stripe/stripe-go) - Go library for the Stripe API
- [MySQL Driver](https://github.com/go-sql-driver/mysql) - Go MySQL Driver
- [Soda CLI](https://gobuffalo.io/documentation/database/soda/) - database migrations
- [Go Simple Mail](https://github.com/xhit/go-simple-mail) - sending emails
- [BW Marin Go alone](https://github.com/bwmarrin/go-alone) - Go MAC signer

### UI Components
- [Stripe Virtual Terminal](https://dashboard.stripe.com/) - integrates the Stripe Virtual Terminal for processing payments 

### Development Tools
- [Air](https://github.com/air-verse/air) - live reload for Go apps
- [Make](https://www.gnu.org/software/make) - controls the generation of executables and other non-source files 