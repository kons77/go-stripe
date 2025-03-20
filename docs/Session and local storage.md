

### **Preventing Form Resubmission** (vid 46).

The receipt page is actually a direct result of the post and it's good practice to redirect people somewhere else so they can't accidentally post that data twice. 

**Problem:** Reloading the "Payment Succeeded" page prompts the user to resubmit the form, which could accidentally reprocess the payment.

**Solution:** Use a session to store necessary data and redirect the user after submission.



### Save sessions in the DB 

Use this https://github.com/alexedwards/scs/tree/master/mysqlstore 

```
go get github.com/alexedwards/scs/mysqlstore
```



### Local Storage and safety

Question: **Is it "safe" to save the token in the localStorage?** How is it compared to use cookie as the storage?

Trevor: That's a good question. Here is a well written article that deals with that exact issue:

Is putting JWTs in local storage "bad"? https://www.ducktypelabs.com/is-localstorage-bad/

**When to Use `localStorage`?**

- Storing authentication tokens (**but it's safer to use `sessionStorage` or `cookies`** if higher security is needed).
- Saving user settings (e.g., theme, language).
- Caching data to avoid reloading it from the server on every request.

⚠️ **Do not use `localStorage` for sensitive data** (such as passwords), as attackers can access it through XSS attacks.

