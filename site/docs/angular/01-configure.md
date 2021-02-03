---
title: Configure Zitadel
---

### Setup Application and get Keys

Before we can start building our application we have do do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your [Project](https://console.zitadel.ch/projects) and add a new application at the top of the page.
Select Web Application and continue.
We recommend that you use [Authorization Code](architecture#Authorization_Code) in combination with [Proof Key for Code Exchange](architecture#Proof_Key_for_Code_Exchange) for all web applications.

> Make sure Authentication Method is set to `NONE` and the Application Type is set to `SPA` or `NATIVE`.

#### Redirect URLs

A redirect URL is a URL in your application where ZITADEL redirects the user after they have authenticated. Set your url to the domain the web app will be deployed to or use `localhost:4200` for development as Angular will be running on port 4200.

> If you are following along with the sample project you downloaded from our Templates, you should set the Allowed Callback URL to http://localhost:4200/auth/callback You will have to set dev mode to `true` as this will enable unsecure http for the moment.

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the post redirectURI field.

Continue and Create the application.

#### Client ID and Secret

After successful app creation a popup will appear showing you your clientID as well as a secret.
Copy your client ID as it will be needed in the next step.

> Note: You will be able to regenerate the secret at a later time, though it is not needed for SPAs with Authorization Code Flow.
