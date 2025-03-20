### Sending email 

Lections 99, 100 Section 8 

I'm adding a function for sending emails in the backend. While emails might also be sent from the frontend, I'll assume they could be different, so I'll duplicate the functionality.

To send emails, we need an SMTP server.

If you already have one, you can use it, making this step optional. Otherwise, **Mailtrap** provides a free service for developers. Sign up at [mailtrap.io](https://mailtrap.io), then open the **demo inbox** and click **"Show Credentials"** to find the SMTP settings (host, port, username, and password).

Since we don’t want to send real emails during testing, Mailtrap captures them in a safe environment. Alternatives like **MailHog** exist, but for this course, we’ll use Mailtrap.