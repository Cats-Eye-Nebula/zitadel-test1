import { SidenavSetting } from '../sidenav/sidenav.component';

export const GENERAL: SidenavSetting = { id: 'general', i18nKey: 'SETTINGS.LIST.GENERAL' };

export const LOGIN: SidenavSetting = {
  id: 'login',
  i18nKey: 'SETTINGS.LIST.LOGIN',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
};
export const LOCKOUT: SidenavSetting = {
  id: 'lockout',
  i18nKey: 'SETTINGS.LIST.LOCKOUT',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
};
export const COMPLEXITY: SidenavSetting = {
  id: 'complexity',
  i18nKey: 'SETTINGS.LIST.COMPLEXITY',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
};

export const IDP: SidenavSetting = { id: 'idp', i18nKey: 'SETTINGS.LIST.IDP', groupI18nKey: 'SETTINGS.GROUPS.LOGIN' };

export const NOTIFICATIONPROVIDERS: SidenavSetting = {
  id: 'notificationproviders',
  i18nKey: 'SETTINGS.LIST.NOTIFICATIONPROVIDERS',
  groupI18nKey: 'SETTINGS.GROUPS.NOTIFICATIONS',
};

export const NOTIFICATIONS: SidenavSetting = {
  id: 'notifications',
  i18nKey: 'SETTINGS.LIST.NOTIFICATIONS',
  groupI18nKey: 'SETTINGS.GROUPS.NOTIFICATIONS',
};

export const MESSAGETEXTS: SidenavSetting = {
  id: 'messagetexts',
  i18nKey: 'SETTINGS.LIST.MESSAGETEXTS',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
};

export const LOGINTEXTS: SidenavSetting = {
  id: 'logintexts',
  i18nKey: 'SETTINGS.LIST.LOGINTEXTS',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
};

export const PRIVACYPOLICY: SidenavSetting = {
  id: 'privacypolicy',
  i18nKey: 'SETTINGS.LIST.PRIVACYPOLICY',
  groupI18nKey: 'SETTINGS.GROUPS.OTHER',
};

export const BRANDING: SidenavSetting = {
  id: 'branding',
  i18nKey: 'SETTINGS.LIST.BRANDING',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
};
