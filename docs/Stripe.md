# Stripe



### Stripe Keys

You'll have to use your own private and publishable test keys there. Just go to your stripe dashboard, copy your publishable key, your test key, not your live key and put it here 

```js
<script>
    let card;
    let stripe;
    
    //  initialize with my publishable key
    stripe = Stripe('')  <--- publishable key here
    ..... 
<script>
```

**Stripe Publishable key** is used on public facing web pages to identify what Stripe account is associated with the transactions that take place.

### Payment Intent

A **PaymentIntent** transitions through multiple statuses throughout its lifetime as it interfaces with Stripe.js to perform authentication flows and ultimately creates at most one successful charge

> The payment intent will be getting doesn't actually charge the credit card. It is the first instance of **payment intent** which will change its status throughout the life cycle of a transaction. So initially we're just making sure that everything is valid, that we have a credit card, that we can charge that credit card and we get an ID back our payment intent ID. So now we want to actually try charging a credit card and also take care of the situation where we're charging a card, but something goes wrong.

**Meaningful function names like `charge` wrapper func with `paymentIntent` inside** 

Now I use create `paymentIntent` here because I'm using Stripe, but at some point down the road I might want to modify my code base to work with, say, Stripe and PayPal and maybe a Canadian bank like Moneris. 

And at that point I could use the repository pattern and have `charge refund`, `partial refund`, whatever those may be, and those would be **meaningful names** regardless as to the payment gateway we're using. Of course, I'm not going to do that in this course, but it's the sort of thing you may want to keep in mind as you build any payments solution in go.

### Stripe library version mis-match

> John: *Trevor, I seem adversely affected by using the 'v75' code after two years or more.*
>
> *The params struct built using the 'v75' code is incompatible with the paymentintent structure in the new library code imported by VSCode.*
>
> *Do you recommend digging into historical docs to find a payment intent compatible with 'v75' or abandoning the use of 'v75' entirely?*
>
> **Trevor:** Either way would work, but I would probably go with the updating to a more recent version. Doing so on this course is on my list of things to get to.

### **Logging transactions** - store ONLY the last 4 digits

Storing cards in most countries is now in **violation of the merchant agreement**; you just need to store the last four digits, which you get back from Stripe (as we do in the course; if you're not there yet, then you will be soon), as well as the paymentIntent. Everything you need related to that transaction is available through just the paymentIntent object.

