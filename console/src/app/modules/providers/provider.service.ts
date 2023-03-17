import { Location } from '@angular/common';
import { Injectable, Injector } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { BehaviorSubject, take } from 'rxjs';
import { Options, Provider, ProviderType } from 'src/app/proto/generated/zitadel/idp_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { AbstractProvider } from './abstract-provider';

import { ActivatedRoute } from '@angular/router';
import {
  AddGenericOAuthProviderRequest as AdminAddGenericOAuthProviderRequest,
  AddGenericOIDCProviderRequest as AdminAddGenericOIDCProviderRequest,
  AddGitHubEnterpriseServerProviderRequest as AdminGitHubEnterpriseServerProviderRequest,
  AddJWTProviderRequest as AdminAddJWTProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  AddGenericOAuthProviderRequest as MgmtAddGenericOAuthProviderRequest,
  AddGenericOIDCProviderRequest as MgmtAddGenericOIDCProviderRequest,
  AddGitHubEnterpriseServerProviderRequest as MgmtGitHubEnterpriseServerProviderRequest,
  AddJWTProviderRequest as MgmtAddJWTProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
@Injectable({
  providedIn: 'root',
})
export class ProviderService implements AbstractProvider {
  public id: string | null = '';
  public showOptional: boolean = false;
  public options: Options = new Options();

  private loading$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public form!: FormGroup;
  private service!: ManagementService | AdminService;
  public serviceType!: PolicyComponentServiceType;

  constructor(
    private injector: Injector,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
    private route: ActivatedRoute,
  ) {
    this.route.data.pipe(take(1)).subscribe((data) => {
      this.serviceType = data.serviceType;

      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.service = this.injector.get(ManagementService as Type<ManagementService>);

          const bread: Breadcrumb = {
            type: BreadcrumbType.ORG,
            routerLink: ['/org'],
          };

          breadcrumbService.setBreadcrumb([bread]);
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);

          const iamBread = new Breadcrumb({
            type: BreadcrumbType.ORG,
            name: 'Instance',
            routerLink: ['/instance'],
          });
          breadcrumbService.setBreadcrumb([iamBread]);
          break;
      }

      this.id = this.route.snapshot.paramMap.get('id');

      if (this.id) {
        this.clientSecret?.setValidators([]);
        this.getData(this.id);
      }
    });
  }

  public getData(id: string, providerType: string): Promise<Provider.AsObject | undefined> {
    this.loading$.next(true);
    const req =
      this.serviceType === PolicyComponentServiceType.ADMIN
        ? new AdminGetProviderByIDRequest()
        : new MgmtGetProviderByIDRequest();

    req.setId(id);

    return this.service.getProviderByID(req).then((resp) => {
      const provider = resp.idp;
      this.loading$.next(false);
      if (provider?.config?.oidc) {
        this.form.patchValue(provider.config.oidc);
      }
      return provider;
    });
  }

  public addProvider(form: FormGroup<any>, providerType: ProviderType): Promise<any> {
    switch (providerType) {
      case ProviderType.PROVIDER_TYPE_GITHUB:
        const req =
          this.serviceType === PolicyComponentServiceType.MGMT
            ? new MgmtAddGenericOIDCProviderRequest()
            : new AdminAddGenericOIDCProviderRequest();
        req.setName(this.form.get('name')?.value);
        req.setClientId(this.form.get('clientId')?.value);
        req.setClientSecret(this.form.get('clientSecret')?.value);
        req.setIssuer(this.form.get('issuer')?.value);
        req.setScopesList(this.form.get('scopesList')?.value);
        return this.service.addGitHubProvider(req);
      case ProviderType.PROVIDER_TYPE_GITHUB_ES:
        const reqGES =
          this.serviceType === PolicyComponentServiceType.MGMT
            ? new MgmtGitHubEnterpriseServerProviderRequest()
            : new AdminGitHubEnterpriseServerProviderRequest();
        reqGES.setName(this.form.get('name')?.value);
        reqGES.setAuthorizationEndpoint(this.form.get('authorizationEndpoint')?.value);
        reqGES.setTokenEndpoint(this.form.get('tokenEndpoint')?.value);
        reqGES.setUserEndpoint(this.form.get('userEndpoint')?.value);
        reqGES.setClientId(this.form.get('clientId')?.value);
        reqGES.setClientSecret(this.form.get('clientSecret')?.value);
        reqGES.setScopesList(this.form.get('scopesList')?.value);
        return this.service.addGitHubEnterpriseServerProvider(reqGES);
      case ProviderType.PROVIDER_TYPE_GITLAB:
        const reqG =
          this.serviceType === PolicyComponentServiceType.MGMT
            ? new MgmtAddGenericOIDCProviderRequest()
            : new AdminAddGenericOIDCProviderRequest();
        reqG.setName(this.form.get('name')?.value);
        reqG.setClientId(this.form.get('clientId')?.value);
        reqG.setClientSecret(this.form.get('clientSecret')?.value);
        reqG.setIssuer(this.form.get('issuer')?.value);
        reqG.setScopesList(this.form.get('scopesList')?.value);
        return this.service.addGitLabProvider(reqG);
      case ProviderType.PROVIDER_TYPE_GITLAB_SELF_HOSTED:
        const reqGSH =
          this.serviceType === PolicyComponentServiceType.MGMT
            ? new MgmtAddGenericOIDCProviderRequest()
            : new AdminAddGenericOIDCProviderRequest();
        reqGSH.setName(this.form.get('name')?.value);
        reqGSH.setClientId(this.form.get('clientId')?.value);
        reqGSH.setClientSecret(this.form.get('clientSecret')?.value);
        reqGSH.setIssuer(this.form.get('issuer')?.value);
        reqGSH.setScopesList(this.form.get('scopesList')?.value);
        return this.service.addGitLabSelfHostedProvider(reqGSH);
      case ProviderType.PROVIDER_TYPE_AZURE_AD:
        break;
      case ProviderType.PROVIDER_TYPE_GOOGLE:
        break;
      case ProviderType.PROVIDER_TYPE_JWT:
        const reqJWT =
          this.serviceType === PolicyComponentServiceType.MGMT
            ? new MgmtAddJWTProviderRequest()
            : new AdminAddJWTProviderRequest();
        reqJWT.setName(this.form.get('name')?.value);
        reqJWT.setHeaderName(this.form.get('headerName')?.value);
        reqJWT.setIssuer(this.form.get('issuer')?.value);
        reqJWT.setJwtEndpoint(this.form.get('jwtEndpoint')?.value);
        reqJWT.setKeysEndpoint(this.form.get('keysEndpoint')?.value);
        reqJWT.setProviderOptions(this.options);
        break;
      case ProviderType.PROVIDER_TYPE_OAUTH:
        const reqOAUTH =
          this.serviceType === PolicyComponentServiceType.MGMT
            ? new MgmtAddGenericOAuthProviderRequest()
            : new AdminAddGenericOAuthProviderRequest();
        reqOAUTH.setName(this.form.get('name')?.value);
        reqOAUTH.setAuthorizationEndpoint(this.form.get('authorizationEndpoint')?.value);
        reqOAUTH.setIdAttribute(this.form.get('idAttribute')?.value);
        reqOAUTH.setTokenEndpoint(this.form.get('tokenEndpoint')?.value);
        reqOAUTH.setUserEndpoint(this.form.get('userEndpoint')?.value);
        reqOAUTH.setClientId(this.form.get('clientId')?.value);
        reqOAUTH.setClientSecret(this.form.get('clientSecret')?.value);
        reqOAUTH.setScopesList(this.form.get('scopesList')?.value);
        break;
      case ProviderType.PROVIDER_TYPE_OIDC:
        break;
    }
  }

  public updateProvider(id: string, form: FormGroup<any>): Promise<Provider.AsObject> {
    throw new Error('Method not implemented.');
  }

  public navigateBack(): void {
    this._location.back();
  }
}
