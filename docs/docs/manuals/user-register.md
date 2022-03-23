---
title: User Register
---

## Organization and user registration

ZITADEL allows users to register a organization and/or user with just a few steps.

A. Register an organization

 1. Create an organization
 2. Verify your email
 3. Login to ZITADEL and manage the organization

B. Create User
 1. An administrator can create and manage users within console.

C. Enable selfregistration for User

 1. Create an organization as above
 2. Create custom policy
 3. Enable the "Register allowed" flag in the Login Policy
 4. Connect your application and add the applications [scope](https://docs.zitadel.ch/architecture/#Custom_Scopes) to the redirect URL.

This will enable the register option in the login dialog and will register the user within your organization if he does not already have an account.

Register Organization
![Register Organization](/img/register.gif)


Create User
![Create User](/img/create-user.gif)


Enable Selfregister
![Enable Selfregister](/img/enable-selfregister.gif)



##  Self Register

When self registration is enabled, users can register themselfes in the organanization without any administrative effort.


Self Register
![Self Register](/img/self-register.gif)

