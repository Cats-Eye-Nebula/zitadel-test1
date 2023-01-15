import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable, Subject } from 'rxjs';
import { take } from 'rxjs/operators';
import {
  UpdateIDPJWTConfigRequest,
  UpdateIDPOIDCConfigRequest,
  UpdateIDPRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { IDP, IDPState, IDPStylingType, OIDCMappingField } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  UpdateOrgIDPJWTConfigRequest,
  UpdateOrgIDPOIDCConfigRequest,
  UpdateOrgIDPRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';

@Component({
  selector: 'cnsl-idp',
  templateUrl: './idp.component.html',
  styleUrls: ['./idp.component.scss'],
})
export class IdpComponent implements OnDestroy {
  public mappingFields: OIDCMappingField[] = [];
  public styleFields: IDPStylingType[] = [];

  public showIdSecretSection: boolean = false;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public idp?: IDP.AsObject;
  private destroy$: Subject<void> = new Subject();
  public projectId: string = '';

  public idpForm!: UntypedFormGroup;
  public oidcConfigForm!: UntypedFormGroup;
  public jwtConfigForm!: UntypedFormGroup;

  IDPState: any = IDPState;

  public canWrite: Observable<boolean> = this.authService.isAllowed([
    this.serviceType === PolicyComponentServiceType.ADMIN
      ? 'iam.idp.write'
      : this.serviceType === PolicyComponentServiceType.MGMT
      ? 'org.idp.write'
      : '',
  ]);

  constructor(
    private toast: ToastService,
    private injector: Injector,
    private route: ActivatedRoute,
    private router: Router,
    private _location: Location,
    private authService: GrpcAuthService,
    private dialog: MatDialog,
    breadcrumbService: BreadcrumbService,
  ) {
    this.idpForm = new UntypedFormGroup({
      id: new UntypedFormControl({ disabled: true, value: '' }, [Validators.required]),
      name: new UntypedFormControl('', [Validators.required]),
      stylingType: new UntypedFormControl('', [Validators.required]),
      autoRegister: new UntypedFormControl(false, [Validators.required]),
    });

    this.oidcConfigForm = new UntypedFormGroup({
      clientId: new UntypedFormControl('', [Validators.required]),
      clientSecret: new UntypedFormControl(''),
      issuer: new UntypedFormControl('', [Validators.required]),
      scopesList: new UntypedFormControl([], []),
      displayNameMapping: new UntypedFormControl(0),
      usernameMapping: new UntypedFormControl(0),
    });

    this.jwtConfigForm = new UntypedFormGroup({
      jwtEndpoint: new UntypedFormControl('', [Validators.required]),
      issuer: new UntypedFormControl('', [Validators.required]),
      keysEndpoint: new UntypedFormControl('', [Validators.required]),
      headerName: new UntypedFormControl('', [Validators.required]),
    });

    const serviceType = this.route.snapshot.data.serviceType;
    if (serviceType !== undefined) {
      this.serviceType = serviceType;

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
            type: BreadcrumbType.INSTANCE,
            name: 'Instance',
            routerLink: ['/instance'],
          });
          breadcrumbService.setBreadcrumb([iamBread]);
          break;
      }

      this.mappingFields = [
        OIDCMappingField.OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
        OIDCMappingField.OIDC_MAPPING_FIELD_EMAIL,
      ];
      this.styleFields = [IDPStylingType.STYLING_TYPE_UNSPECIFIED, IDPStylingType.STYLING_TYPE_GOOGLE];
    }

    const idpId = this.route.snapshot.paramMap.get('id');
    if (idpId) {
      this.checkWrite();

      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        (this.service as ManagementService).getOrgIDPByID(idpId).then((resp) => {
          if (resp.idp) {
            this.idp = resp.idp;
            this.idpForm.patchValue(this.idp);
            if (this.idp.oidcConfig) {
              this.oidcConfigForm.patchValue(this.idp.oidcConfig);
            } else if (this.idp.jwtConfig) {
              this.jwtConfigForm.patchValue(this.idp.jwtConfig);
              this.jwtIssuer?.setValue(this.idp.jwtConfig.issuer);
            }
          }
        });
      } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
        (this.service as AdminService).getIDPByID(idpId).then((resp) => {
          if (resp.idp) {
            this.idp = resp.idp;
            this.idpForm.patchValue(this.idp);
            if (this.idp.oidcConfig) {
              this.oidcConfigForm.patchValue(this.idp.oidcConfig);
            } else if (this.idp.jwtConfig) {
              this.jwtConfigForm.patchValue(this.idp.jwtConfig);
              this.jwtIssuer?.setValue(this.idp.jwtConfig.issuer);
            }
          }
        });
      }
    }
  }

  public checkWrite(): void {
    this.canWrite.pipe(take(1)).subscribe((canWrite) => {
      if (canWrite) {
        this.idpForm.enable();
        this.oidcConfigForm.enable();
      } else {
        this.idpForm.disable();
        this.oidcConfigForm.disable();
      }
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public deleteIdp(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'IDP.DELETE_TITLE',
        descriptionKey: 'IDP.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (this.idp && resp) {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
          (this.service as ManagementService)
            .removeOrgIDP(this.idp.id)
            .then(() => {
              this.toast.showInfo('IDP.TOAST.DELETED', true);
              this.router.navigate(this.backroutes);
            })
            .catch((error: any) => {
              this.toast.showError(error);
            });
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
          (this.service as AdminService)
            .removeIDP(this.idp.id)
            .then(() => {
              this.toast.showInfo('IDP.TOAST.DELETED', true);
              this.router.navigate(this.backroutes);
            })
            .catch((error: any) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public changeState(state: IDPState): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      if (state === IDPState.IDP_STATE_ACTIVE && this.idp) {
        (this.service as ManagementService)
          .reactivateOrgIDP(this.idp.id)
          .then(() => {
            this.idp!.state = state;
            this.toast.showInfo('IDP.TOAST.REACTIVATED', true);
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      } else if (state === IDPState.IDP_STATE_INACTIVE && this.idp) {
        (this.service as ManagementService)
          .deactivateOrgIDP(this.idp.id)
          .then(() => {
            this.idp!.state = state;
            this.toast.showInfo('IDP.TOAST.DEACTIVATED', true);
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      if (state === IDPState.IDP_STATE_ACTIVE && this.idp) {
        (this.service as AdminService)
          .reactivateIDP(this.idp.id)
          .then(() => {
            this.idp!.state = state;
            this.toast.showInfo('IDP.TOAST.REACTIVATED', true);
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      } else if (state === IDPState.IDP_STATE_INACTIVE && this.idp) {
        (this.service as AdminService)
          .deactivateIDP(this.idp.id)
          .then(() => {
            this.idp!.state = state;
            this.toast.showInfo('IDP.TOAST.DEACTIVATED', true);
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    }
  }

  public updateIdp(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT && this.idp) {
      const req = new UpdateOrgIDPRequest();

      req.setIdpId(this.idp.id);
      req.setName(this.name?.value);
      req.setStylingType(this.stylingType?.value);
      req.setAutoRegister(this.autoRegister?.value);

      (this.service as ManagementService)
        .updateOrgIDP(req)
        .then(() => {
          this.toast.showInfo('IDP.TOAST.SAVED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN && this.idp) {
      const req = new UpdateIDPRequest();

      req.setIdpId(this.idp.id);
      req.setName(this.name?.value);
      req.setStylingType(this.stylingType?.value);
      req.setAutoRegister(this.autoRegister?.value);

      (this.service as AdminService)
        .updateIDP(req)
        .then(() => {
          this.toast.showInfo('IDP.TOAST.SAVED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public updateOidcConfig(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT && this.idp) {
      const req = new UpdateOrgIDPOIDCConfigRequest();

      req.setIdpId(this.idp.id);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);
      req.setUsernameMapping(this.usernameMapping?.value);
      req.setDisplayNameMapping(this.displayNameMapping?.value);

      (this.service as ManagementService)
        .updateOrgIDPOIDCConfig(req)
        .then((oidcConfig) => {
          this.toast.showInfo('IDP.TOAST.SAVED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN && this.idp) {
      const req = new UpdateIDPOIDCConfigRequest();

      req.setIdpId(this.idp.id);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);
      req.setUsernameMapping(this.usernameMapping?.value);
      req.setDisplayNameMapping(this.displayNameMapping?.value);

      (this.service as AdminService)
        .updateIDPOIDCConfig(req)
        .then((oidcConfig) => {
          this.toast.showInfo('IDP.TOAST.SAVED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public updateJwtConfig(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT && this.idp) {
      const req = new UpdateOrgIDPJWTConfigRequest();

      req.setIdpId(this.idp.id);
      req.setIssuer(this.jwtIssuer?.value);
      req.setHeaderName(this.headerName?.value);
      req.setJwtEndpoint(this.jwtEndpoint?.value);
      req.setKeysEndpoint(this.keyEndpoint?.value);

      (this.service as ManagementService)
        .updateOrgIDPJWTConfig(req)
        .then((jwtConfig) => {
          this.toast.showInfo('IDP.TOAST.SAVED', true);
          // this.router.navigate(['idp', ]);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN && this.idp) {
      const req = new UpdateIDPJWTConfigRequest();

      req.setIdpId(this.idp.id);
      req.setIssuer(this.jwtIssuer?.value);
      req.setHeaderName(this.headerName?.value);
      req.setJwtEndpoint(this.jwtEndpoint?.value);
      req.setKeysEndpoint(this.keyEndpoint?.value);

      (this.service as AdminService)
        .updateIDPJWTConfig(req)
        .then((jwtConfig) => {
          this.toast.showInfo('IDP.TOAST.SAVED', true);
          // this.router.navigate(['idp', ]);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public close(): void {
    this._location.back();
  }

  public addScope(event: MatChipInputEvent): void {
    const input = event.chipInput?.inputElement;
    const value = event.value.trim();

    if (value !== '') {
      if (this.scopesList?.value) {
        this.scopesList.value.push(value);
        if (input) {
          input.value = '';
        }
      }
    }
  }

  public removeScope(uri: string): void {
    if (this.scopesList?.value) {
      const index = this.scopesList?.value.indexOf(uri);

      if (index !== undefined && index >= 0) {
        this.scopesList?.value.splice(index, 1);
      }
    }
  }

  public get backroutes(): string[] {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return ['/org', 'policy', 'login'];
      case PolicyComponentServiceType.ADMIN:
        return ['/instance', 'policy', 'login'];
    }
  }

  public get name(): AbstractControl | null {
    return this.idpForm.get('name');
  }

  public get stylingType(): AbstractControl | null {
    return this.idpForm.get('stylingType');
  }

  public get autoRegister(): AbstractControl | null {
    return this.idpForm.get('autoRegister');
  }

  public get clientId(): AbstractControl | null {
    return this.oidcConfigForm.get('clientId');
  }

  public get clientSecret(): AbstractControl | null {
    return this.oidcConfigForm.get('clientSecret');
  }

  public get issuer(): AbstractControl | null {
    return this.oidcConfigForm.get('issuer');
  }

  public get scopesList(): AbstractControl | null {
    return this.oidcConfigForm.get('scopesList');
  }

  public get displayNameMapping(): AbstractControl | null {
    return this.oidcConfigForm.get('displayNameMapping');
  }

  public get usernameMapping(): AbstractControl | null {
    return this.oidcConfigForm.get('usernameMapping');
  }

  public get jwtIssuer(): AbstractControl | null {
    return this.jwtConfigForm.get('issuer');
  }

  public get jwtEndpoint(): AbstractControl | null {
    return this.jwtConfigForm.get('jwtEndpoint');
  }

  public get keyEndpoint(): AbstractControl | null {
    return this.jwtConfigForm.get('keysEndpoint');
  }

  public get headerName(): AbstractControl | null {
    return this.jwtConfigForm.get('headerName');
  }
}
