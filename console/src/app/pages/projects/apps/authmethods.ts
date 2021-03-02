import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import { APIAuthMethodType, APIConfig, OIDCAuthMethodType, OIDCConfig, OIDCGrantType, OIDCResponseType } from 'src/app/proto/generated/management_pb';

export const CODE_METHOD: RadioItemAuthType = {
    key: 'CODE',
    titleI18nKey: 'APP.AUTHMETHODS.CODE.TITLE',
    descI18nKey: 'APP.AUTHMETHODS.CODE.DESCRIPTION',
    disabled: false,
    prefix: 'CODE',
    background: 'rgb(89 115 128)',
    responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
    grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
    authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC,
    recommended: false,
};
export const PKCE_METHOD: RadioItemAuthType = {
    key: 'PKCE',
    titleI18nKey: 'APP.AUTHMETHODS.PKCE.TITLE',
    descI18nKey: 'APP.AUTHMETHODS.PKCE.DESCRIPTION',
    disabled: false,
    prefix: 'PKCE',
    background: 'rgb(80 110 92)',
    responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
    grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
    authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
    recommended: true,
};
export const POST_METHOD: RadioItemAuthType = {
    key: 'POST',
    titleI18nKey: 'APP.AUTHMETHODS.POST.TITLE',
    descI18nKey: 'APP.AUTHMETHODS.POST.DESCRIPTION',
    disabled: false,
    prefix: 'POST',
    background: 'rgb(144 75 75)',
    responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
    grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
    authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST,
    notRecommended: true,
};
export const PK_JWT_METHOD: RadioItemAuthType = {
    key: 'PK_JWT',
    titleI18nKey: 'APP.AUTHMETHODS.PK_JWT.TITLE',
    descI18nKey: 'APP.AUTHMETHODS.PK_JWT.DESCRIPTION',
    disabled: false,
    prefix: 'JWT',
    background: 'rgb(89, 93, 128)',
    responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
    grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
    authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_PRIVATE_KEY_JWT,
    apiAuthMethod: APIAuthMethodType.APIAUTHMETHODTYPE_PRIVATE_KEY_JWT,
    // recommended: true,
};
export const BASIC_AUTH_METHOD: RadioItemAuthType = {
    key: 'BASIC',
    titleI18nKey: 'APP.AUTHMETHODS.BASIC.TITLE',
    descI18nKey: 'APP.AUTHMETHODS.BASIC.DESCRIPTION',
    disabled: false,
    prefix: 'BASIC',
    background: 'rgb(144 75 75)',
    responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
    grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
    authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST,
    apiAuthMethod: APIAuthMethodType.APIAUTHMETHODTYPE_BASIC,
};

export const IMPLICIT_METHOD: RadioItemAuthType = {
    key: 'IMPLICIT',
    titleI18nKey: 'APP.AUTHMETHODS.IMPLICIT.TITLE',
    descI18nKey: 'APP.AUTHMETHODS.IMPLICIT.DESCRIPTION',
    disabled: false,
    prefix: 'IMP',
    background: 'rgb(144 75 75)',
    responseType: OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN,
    grantType: OIDCGrantType.OIDCGRANTTYPE_IMPLICIT,
    authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
    notRecommended: true,
};

export const CUSTOM_METHOD: RadioItemAuthType = {
    key: 'CUSTOM',
    titleI18nKey: 'APP.AUTHMETHODS.CUSTOM.TITLE',
    descI18nKey: 'APP.AUTHMETHODS.CUSTOM.DESCRIPTION',
    disabled: false,
    prefix: 'CUSTOM',
    background: '#333',
};

export function getPartialConfigFromAuthMethod(authMethod: string): {
    oidc?: Partial<OIDCConfig.AsObject>;
    api?: Partial<APIConfig.AsObject>;
} | undefined {
    let config: {
        oidc?: Partial<OIDCConfig.AsObject>,
        api?: Partial<APIConfig.AsObject>,
    };
    switch (authMethod) {
        case CODE_METHOD.key:
            config = {
                oidc: {
                    responseTypesList: [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                    grantTypesList: [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                    authMethodType: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC,
                },
            };
            return config;
        case PKCE_METHOD.key:
            config = {
                oidc: {
                    responseTypesList: [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                    grantTypesList: [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                    authMethodType: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
                }
            };
            return config;
        case POST_METHOD.key:
            config = {
                oidc: {
                    responseTypesList: [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                    grantTypesList: [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                    authMethodType: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST,
                }
            };
            return config;
        case PK_JWT_METHOD.key:
            config = {
                oidc: {
                    responseTypesList: [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                    grantTypesList: [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                    authMethodType: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_PRIVATE_KEY_JWT,
                },
                api: {
                    authMethodType: APIAuthMethodType.APIAUTHMETHODTYPE_PRIVATE_KEY_JWT,
                }
            };
            return config;
        case BASIC_AUTH_METHOD.key:
            config = {
                oidc: {
                    authMethodType: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC,
                },
                api: {
                    authMethodType: APIAuthMethodType.APIAUTHMETHODTYPE_BASIC,
                }
            };
            return config;
        case IMPLICIT_METHOD.key:
            config = {
                oidc: {
                    responseTypesList: [OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN_TOKEN],
                    grantTypesList: [OIDCGrantType.OIDCGRANTTYPE_IMPLICIT],
                    authMethodType: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
                },
                api: {
                    authMethodType: APIAuthMethodType.APIAUTHMETHODTYPE_PRIVATE_KEY_JWT,
                }
            };
            return config;
        default:
            return undefined;
    }
}

export function getAuthMethodFromPartialConfig(config: { oidc?: Partial<OIDCConfig.AsObject>, api?: Partial<APIConfig.AsObject>; }): string {
    if (config?.oidc) {
        const toCheck = [config.oidc.responseTypesList, config.oidc.grantTypesList, config.oidc.authMethodType];
        const code = JSON.stringify(
            [
                [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC,
            ]
        );

        const pkce = JSON.stringify(
            [
                [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
            ]
        );

        const post = JSON.stringify(
            [
                [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST,
            ]
        );

        const pk_jwt = JSON.stringify(
            [
                [OIDCResponseType.OIDCRESPONSETYPE_CODE],
                [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE],
                OIDCAuthMethodType.OIDCAUTHMETHODTYPE_PRIVATE_KEY_JWT,
            ]
        );

        const implicit = JSON.stringify(
            [
                [OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN_TOKEN],
                [OIDCGrantType.OIDCGRANTTYPE_IMPLICIT],
                OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
            ]
        );

        switch (JSON.stringify(toCheck)) {
            case code: return CODE_METHOD.key;
            case pkce: return PKCE_METHOD.key;
            case post: return POST_METHOD.key;
            case pk_jwt: return PK_JWT_METHOD.key;
            case implicit: return IMPLICIT_METHOD.key;
            default:
                return CUSTOM_METHOD.key;
        }
    } else if (config.api && config.api.authMethodType !== undefined) {
        switch (config.api.authMethodType.toString()) {
            case APIAuthMethodType.APIAUTHMETHODTYPE_PRIVATE_KEY_JWT.toString(): return PK_JWT_METHOD.key;
            case APIAuthMethodType.APIAUTHMETHODTYPE_BASIC.toString(): return BASIC_AUTH_METHOD.key;
            default:
                return CUSTOM_METHOD.key;
        }
    } else {
        return CUSTOM_METHOD.key;
    }
}