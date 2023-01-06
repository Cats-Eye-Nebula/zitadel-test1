import { login, User } from 'support/login/users';
import { API } from './types';

const authHeaderKey = 'Authorization',
  orgIdHeaderKey = 'x-zitadel-orgid';

export function apiAuth(): Cypress.Chainable<API> {
  return login(User.IAMAdminUser, 'Password1!', false).then((token) => {
    return <API>{
      token: token,
      mgmtBaseURL: `${Cypress.env('BACKEND_URL')}/management/v1`,
      adminBaseURL: `${Cypress.env('BACKEND_URL')}/admin/v1`,
    };
  });
}

export function requestHeaders(api: API, orgId?: number): object {
  const headers = { [authHeaderKey]: `Bearer ${api.token}` };
  if (orgId) {
    headers[orgIdHeaderKey] = orgId;
  }
  return headers;
}
