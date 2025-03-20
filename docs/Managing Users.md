# Managing Users

(section 13, video 149, 157)

### Prevent to delete yourself

When loading the form:

- For a new user, fields remain empty, and the delete button is hidden.
- For an existing user, we fetch their details using the user ID from the URL, populate the fields, and leave the password fields empty.

A user shouldn’t be able to delete themselves. That's generally considered bad practice.  Even if there's only one user in the database, the delete button should be disabled to prevent accidental self-deletion. Users will attempt anything if given the chance, so we assume the worst and prevent this action.

To hide the delete button for the logged-in user, check the front-end session and compare the user being viewed with the logged-in user. Update both the front end and back end accordingly.

### WebSockets 

In this course, we won’t receive messages from the client-only send notifications when a user is deleted (“User ID X has been deleted”).

We call `upgradeConnection` because WebSockets upgrade the connection, allowing two-way communication. While we don’t need client-to-server messaging now, setting it up could be useful for future features.
