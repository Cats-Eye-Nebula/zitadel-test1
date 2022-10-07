---
title: zitadel/system.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)


## SystemService {#zitadelsystemv1systemservice}


### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)

Indicates if ZITADEL is running.
It respondes as soon as ZITADEL started



    GET: /healthz


### ListInstances

> **rpc** ListInstances([ListInstancesRequest](#listinstancesrequest))
[ListInstancesResponse](#listinstancesresponse)

Returns a list of ZITADEL instances



    POST: /instances/_search


### GetInstance

> **rpc** GetInstance([GetInstanceRequest](#getinstancerequest))
[GetInstanceResponse](#getinstanceresponse)

Returns the detail of an instance



    GET: /instances/{instance_id}


### AddInstance

> **rpc** AddInstance([AddInstanceRequest](#addinstancerequest))
[AddInstanceResponse](#addinstanceresponse)

deprecated:  Creates a new instance with all needed setup data
This might take some time



    POST: /instances


### UpdateInstance

> **rpc** UpdateInstance([UpdateInstanceRequest](#updateinstancerequest))
[UpdateInstanceResponse](#updateinstanceresponse)

Updates name of an existing instance



    PUT: /instances/{instance_id}


### CreateInstance

> **rpc** CreateInstance([CreateInstanceRequest](#createinstancerequest))
[CreateInstanceResponse](#createinstanceresponse)

Creates a new instance with all needed setup data
This might take some time



    POST: /instances/_create


### RemoveInstance

> **rpc** RemoveInstance([RemoveInstanceRequest](#removeinstancerequest))
[RemoveInstanceResponse](#removeinstanceresponse)

Removes a instances
This might take some time



    DELETE: /instances/{instance_id}


### ExistsDomain

> **rpc** ExistsDomain([ExistsDomainRequest](#existsdomainrequest))
[ExistsDomainResponse](#existsdomainresponse)

Checks if a domain exists



    POST: /domains/{domain}/_exists


### ListDomains

> **rpc** ListDomains([ListDomainsRequest](#listdomainsrequest))
[ListDomainsResponse](#listdomainsresponse)

Returns the custom domains of an instance



    POST: /instances/{instance_id}/domains/_search


### AddDomain

> **rpc** AddDomain([AddDomainRequest](#adddomainrequest))
[AddDomainResponse](#adddomainresponse)

Returns the domain of an instance



    POST: /instances/{instance_id}/domains


### RemoveDomain

> **rpc** RemoveDomain([RemoveDomainRequest](#removedomainrequest))
[RemoveDomainResponse](#removedomainresponse)

Returns the domain of an instance



    DELETE: /instances/{instance_id}/domains/{domain}


### SetPrimaryDomain

> **rpc** SetPrimaryDomain([SetPrimaryDomainRequest](#setprimarydomainrequest))
[SetPrimaryDomainResponse](#setprimarydomainresponse)

Returns the domain of an instance



    POST: /instances/{instance_id}/domains/_set_primary


### ListViews

> **rpc** ListViews([ListViewsRequest](#listviewsrequest))
[ListViewsResponse](#listviewsresponse)

Returns all stored read models of ZITADEL
views are used for search optimisation and optimise request latencies
they represent the delta of the event happend on the objects



    POST: /views/_search


### ClearView

> **rpc** ClearView([ClearViewRequest](#clearviewrequest))
[ClearViewResponse](#clearviewresponse)

Truncates the delta of the change stream
be carefull with this function because ZITADEL has to
recompute the deltas after they got cleared.
Search requests will return wrong results until all deltas are recomputed



    POST: /views/{database}/{view_name}


### ListFailedEvents

> **rpc** ListFailedEvents([ListFailedEventsRequest](#listfailedeventsrequest))
[ListFailedEventsResponse](#listfailedeventsresponse)

Returns event descriptions which cannot be processed.
It's possible that some events need some retries.
For example if the SMTP-API wasn't able to send an email at the first time



    POST: /failedevents/_search


### RemoveFailedEvent

> **rpc** RemoveFailedEvent([RemoveFailedEventRequest](#removefailedeventrequest))
[RemoveFailedEventResponse](#removefailedeventresponse)

Deletes the event from failed events view.
the event is not removed from the change stream
This call is usefull if the system was able to process the event later.
e.g. if the second try of sending an email was successful. the first try produced a
failed event. You can find out if it worked on the `failure_count`



    DELETE: /failedevents/{database}/{view_name}/{failed_sequence}







## Messages


### AddDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### AddDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddInstanceRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| first_org_name |  string | - | string.max_len: 200<br />  |
| custom_domain |  string | - | string.max_len: 200<br />  |
| owner_user_name |  string | - | string.max_len: 200<br />  |
| owner_email |  AddInstanceRequest.Email | - | message.required: true<br />  |
| owner_profile |  AddInstanceRequest.Profile | - | message.required: false<br />  |
| owner_password |  AddInstanceRequest.Password | - | message.required: false<br />  |
| default_language |  string | - | string.max_len: 10<br />  |




### AddInstanceRequest.Email



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| is_email_verified |  bool | - |  |




### AddInstanceRequest.Password



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| password |  string | - | string.max_len: 200<br />  |
| password_change_required |  bool | - |  |




### AddInstanceRequest.Profile



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.max_len: 200<br />  |
| last_name |  string | - | string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |




### AddInstanceResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### ChangeSubscriptionRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| subscription_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| request_limit |  uint64 | - |  |
| action_mins_limit |  uint64 | - |  |




### ChangeSubscriptionResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ClearViewRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| view_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ClearViewResponse
This is an empty response




### CreateInstanceRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| first_org_name |  string | - | string.max_len: 200<br />  |
| custom_domain |  string | - | string.max_len: 200<br />  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) user.human |  CreateInstanceRequest.Human | oneof field for the user managing the instance |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) user.machine |  CreateInstanceRequest.Machine | - |  |
| default_language |  string | - | string.max_len: 10<br />  |




### CreateInstanceRequest.Email



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| is_email_verified |  bool | - |  |




### CreateInstanceRequest.Human



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  CreateInstanceRequest.Email | - | message.required: true<br />  |
| profile |  CreateInstanceRequest.Profile | - | message.required: false<br />  |
| password |  CreateInstanceRequest.Password | - | message.required: false<br />  |
| user_name |  string | - | string.max_len: 200<br />  |




### CreateInstanceRequest.Machine



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.max_len: 200<br />  |
| name |  string | - | string.max_len: 200<br />  |
| personal_access_token |  CreateInstanceRequest.PersonalAccessToken | - |  |
| machine_key |  CreateInstanceRequest.MachineKey | - |  |




### CreateInstanceRequest.MachineKey



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.authn.v1.KeyType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |
| expiration_date |  google.protobuf.Timestamp | - |  |




### CreateInstanceRequest.Password



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| password |  string | - | string.max_len: 200<br />  |
| password_change_required |  bool | - |  |




### CreateInstanceRequest.PersonalAccessToken



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| expiration_date |  google.protobuf.Timestamp | - |  |




### CreateInstanceRequest.Profile



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.max_len: 200<br />  |
| last_name |  string | - | string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |




### CreateInstanceResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| pat |  string | - |  |
| machine_key |  bytes | - |  |




### ExistsDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ExistsDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| exists |  bool | - |  |




### FailedEvent



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - |  |
| view_name |  string | - |  |
| failed_sequence |  uint64 | - |  |
| failure_count |  uint64 | - |  |
| error_message |  string | - |  |




### GetInstanceRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetInstanceResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance |  zitadel.instance.v1.InstanceDetail | - |  |




### GetUsageRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetUsageResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| executed_requests |  uint64 | - |  |
| executed_action_mins |  uint64 | - |  |




### HealthzRequest
This is an empty request




### HealthzResponse
This is an empty response




### ListDomainsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | list limitations and ordering | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | - |  |
| sorting_column |  zitadel.instance.v1.DomainFieldName | the field the result is sorted |  |
| queries | repeated zitadel.instance.v1.DomainSearchQuery | criterias the client is looking for |  |




### ListDomainsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.instance.v1.DomainFieldName | - |  |
| result | repeated zitadel.instance.v1.Domain | - |  |




### ListFailedEventsRequest
This is an empty request




### ListFailedEventsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated FailedEvent | TODO: list details |  |




### ListInstancesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| sorting_column |  zitadel.instance.v1.FieldName | the field the result is sorted |  |
| queries | repeated zitadel.instance.v1.Query | criterias the client is looking for |  |




### ListInstancesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.instance.v1.FieldName | - |  |
| result | repeated zitadel.instance.v1.Instance | - |  |




### ListViewsRequest
This is an empty request




### ListViewsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated View | TODO: list details |  |




### RemoveDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveFailedEventRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| view_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| failed_sequence |  uint64 | - |  |




### RemoveFailedEventResponse
This is an empty response




### RemoveInstanceRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveInstanceResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetPrimaryDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### SetPrimaryDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateInstanceRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance_id |  string | - |  |
| instance_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateInstanceResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### View



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - |  |
| view_name |  string | - |  |
| processed_sequence |  uint64 | - |  |
| event_timestamp |  google.protobuf.Timestamp | The timestamp the event occured |  |
| last_successful_spooler_run |  google.protobuf.Timestamp | - |  |
| instance |  string | - |  |






