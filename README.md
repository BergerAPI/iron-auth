# OAuth2 iron.sh
Very simple implementation of OAuth2 ([RFC6749](https://www.rfc-editor.org/rfc/rfc6749)) for managing the login-process on [iron.sh](https://iron.sh/) services.

## The flow
OAuth2 provides detailed information about how things need to be done. I tried implementing everything one to one.
The client-application creates an authentication link, which after authentication by the user is turned into a grant and transmitted to the client via query-parameters and redirects. This grant can now be converted to an access key, which allows the client to access the data of the user. For this conversion, the client requires a client-secret, which further proves the identity of the client.

[OAuth2 Flow](/docs/oauth2-flow.png)

### The authorize flow
The `/oauth/authorize` endpoint has a side effect: in case of unauthorized access to the endpoint, it will redirect the user to the login page. The query parameters, that should be provided to `/oauth/authorize` will be passed to the `/login` page, where they will be passed into the form. When the user clicks the submit button, the form will be processed, where another `/oauth/authorize` link will be generated (only if these parameters did exist in the first place), and the user will be redirected to it. Now the flow continues where it stopped due to the user being not authorization.

## Endpoints
- `POST /login`: authenticates the user, using the form-encoded body
- `GET /login`: provides an interface for the user to sign in
- `GET /oauth/authorize`: generates a grant for a certain client
- `POST /oauth/token`: converts a grant into an access token

## Goal
In the first place I wanted to learn how OAuth2 works, and what is a better learning experience, than programming it myself. On top of gaining this valuable insight, I also tried to keep the binary and dependency count as low as possible. The result is a very lightweight application, that responds fast and can handle many requests per second.