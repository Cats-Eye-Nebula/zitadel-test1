---
title: Scopes
---

ZITADEL supports the usage of scopes as way of requesting information from the IAM and also instruct ZITADEL to do certain operations.

## Standard Scopes

| Scopes         | Example          | Description                                                                    |
|:---------------|:-----------------|--------------------------------------------------------------------------------|
| openid         | `openid`         | When using openid connect this is a mandatory scope                            |
| profile        | `profile`        | Optional scope to request the profile of the subject                           |
| email          | `email`          | Optional scope to request the email of the subject                             |
| address        | `address`        | Optional scope to request the address of the subject                           |
| offline_access | `offline_access` | Optional scope to request a refresh_token (only possible when using code flow) |

## Custom Scopes

> This feature is not yet released

## Reserved Scopes

In addition to the standard compliant scopes we utilize the following scopes.

| Scopes                                           | Example                                                                        | Description                                                                                                                                                                                                                                                                                                                                                             |
|:-------------------------------------------------|:-------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| urn:zitadel:iam:org:project:role:{rolename}      | `urn:zitadel:iam:org:project:role:user`                                        | By using this scope a client can request the claim urn:zitadel:iam:roles:rolename} to be asserted when possible. As an alternative approach you can enable all roles to be asserted from the [project](../../guides/usage/projects) a client belongs to. |
| urn:zitadel:iam:org:domain:primary:{domainname}  | `urn:zitadel:iam:org:domain:primary:acme.ch`                                   | When requesting this scope **ZITADEL** will enforce that the user is a member of the selected organization. If the organization does not exist a failure is displayed                                                                                                                                                                                                   |
| urn:zitadel:iam:role:{rolename}                  |                                                                                |                                                                                                                                                                                                                                                                                                                                                                         |
| `urn:zitadel:iam:org:project:id:{projectid}:aud` | ZITADEL's Project id is `urn:zitadel:iam:org:project:id:69234237810729019:aud` | By adding this scope, the requested projectid will be added to the audience of the access and id token                                                                                                                                                                                                                                                                  |
| urn:zitadel:iam:user:metadata | `urn:zitadel:iam:user:metadata` | By adding this scope, the metadata of the user will be included in the token. The values are base64 encoded.                                                                                                                                                                                                                                                                |
| urn:zitadel:iam:user:resourceowner | `urn:zitadel:iam:user:resourceowner` | By adding this scope, the resourceowner (id, name, primary_domain) of the user will be included in the token.                                                                                                                                                                                                                                                                |
| urn:zitadel:iam:org:idp:id:{idp_id} | `urn:zitadel:iam:org:idp:id:76625965177954913` | By adding this scope the user will directly be redirected to the identity provider to authenticate. Make sure you also send the primary domain scope if a custom login policy is configured. Otherwise the system will not be able to identify the identity provider.                                                                                                                                                                                                                                                             |

> If access to ZITADEL's API's is needed with a service user the scope `urn:zitadel:iam:org:project:id:69234237810729019:aud` needs to be used with the JWT Profile request
