import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { LoginMethodComponentType } from 'src/app/modules/mfa-table/mfa-table.component';
import {
    IdpProviderView as AdminIdpProviderView,
    IdpStylingType,
    PasswordlessType as AdminPasswordlessType,
} from 'src/app/proto/generated/admin_pb';
import {
    IdpProviderView as MgmtIdpProviderView,
    PasswordlessType as MgmtPasswordlessType,
} from 'src/app/proto/generated/management_pb';
import {
    GetLoginPolicyResponse as AdminGetLoginPolicyResponse,
    UpdateLoginPolicyRequest,
    UpdateLoginPolicyResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { IDP, IDPLoginPolicyLink } from 'src/app/proto/generated/zitadel/idp_pb';
import {
    AddCustomLoginPolicyRequest,
    GetLoginPolicyResponse as MgmtGetLoginPolicyResponse,
} from 'src/app/proto/generated/zitadel/management_pb';
import { LoginPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { CnslLinks } from '../../links/links.component';
import {
    IAM_COMPLEXITY_LINK,
    IAM_LABEL_LINK,
    IAM_POLICY_LINK,
    ORG_COMPLEXITY_LINK,
    ORG_IAM_POLICY_LINK,
} from '../../policy-grid/policy-links';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { AddIdpDialogComponent } from './add-idp-dialog/add-idp-dialog.component';

@Component({
    selector: 'app-login-policy',
    templateUrl: './login-policy.component.html',
    styleUrls: ['./login-policy.component.scss'],
})
export class LoginPolicyComponent implements OnDestroy {
    public LoginMethodComponentType: any = LoginMethodComponentType;
    public passwordlessTypes: Array<AdminPasswordlessType | MgmtPasswordlessType> = [];
    public loginData!: LoginPolicy.AsObject;

    private sub: Subscription = new Subscription();
    public service!: ManagementService | AdminService;
    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    public idps: MgmtIdpProviderView.AsObject[] | AdminIdpProviderView.AsObject[] = [];

    public loading: boolean = false;
    public disabled: boolean = true;

    public IdpStylingType: any = IdpStylingType;
    public nextLinks: CnslLinks[] = [];
    constructor(
        private route: ActivatedRoute,
        private toast: ToastService,
        private dialog: MatDialog,
        private injector: Injector,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.serviceType = data.serviceType;
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    this.service = this.injector.get(ManagementService as Type<ManagementService>);
                    this.passwordlessTypes = [MgmtPasswordlessType.PASSWORDLESSTYPE_ALLOWED,
                    MgmtPasswordlessType.PASSWORDLESSTYPE_NOT_ALLOWED];
                    this.nextLinks = [
                        ORG_COMPLEXITY_LINK,
                        ORG_IAM_POLICY_LINK,
                    ];
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.service = this.injector.get(AdminService as Type<AdminService>);
                    this.passwordlessTypes = [AdminPasswordlessType.PASSWORDLESSTYPE_ALLOWED,
                    AdminPasswordlessType.PASSWORDLESSTYPE_NOT_ALLOWED];
                    this.nextLinks = [
                        IAM_COMPLEXITY_LINK,
                        IAM_POLICY_LINK,
                        IAM_LABEL_LINK,
                    ];
                    break;
            }

            return this.route.params;
        })).subscribe(() => {
            this.fetchData();
        });
    }

    private fetchData(): void {
        this.getData().then(resp => {
            if (resp.policy) {
                this.loginData = resp.policy;
                this.loading = false;
                this.disabled = ((this.loginData as LoginPolicy.AsObject)?.pb_default) ?? false;
            }
        });
        this.getIdps().then(idps => {
            this.idps = idps;
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData():
        Promise<AdminGetLoginPolicyResponse.AsObject | MgmtGetLoginPolicyResponse.AsObject> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).getLoginPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).getLoginPolicy();
        }
    }

    private async getIdps(): Promise<IDPLoginPolicyLink.AsObject[]> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).listLoginPolicyIDPs()
                    .then((providers) => {
                        return providers.resultList;
                    });
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).listLoginPolicyIDPs()
                    .then((providers) => {
                        return providers.resultList;
                    });
        }
    }

    private async updateData():
        Promise<UpdateLoginPolicyResponse.AsObject> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                const mgmtreq = new AddCustomLoginPolicyRequest();
                mgmtreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
                mgmtreq.setAllowRegister(this.loginData.allowRegister);
                mgmtreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
                mgmtreq.setForceMfa(this.loginData.forceMfa);
                mgmtreq.setPasswordlessType(this.loginData.passwordlessType);
                // console.log(mgmtreq.toObject());
                if ((this.loginData as LoginPolicy.AsObject).isDefault) {
                    return (this.service as ManagementService).addCustomLoginPolicy(mgmtreq);
                } else {
                    return (this.service as ManagementService).updateCustomLoginPolicy(mgmtreq);
                }
            case PolicyComponentServiceType.ADMIN:
                const adminreq = new UpdateLoginPolicyRequest();
                adminreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
                adminreq.setAllowRegister(this.loginData.allowRegister);
                adminreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
                adminreq.setForceMfa(this.loginData.forceMfa);
                adminreq.setPasswordlessType(this.loginData.passwordlessType);

                // console.log(adminreq.toObject());

                return (this.service as AdminService).updateLoginPolicy(adminreq);
        }
    }

    public savePolicy(): void {
        this.updateData().then(() => {
            this.toast.showInfo('POLICY.LOGIN_POLICY.SAVED', true);
            this.loading = true;
            setTimeout(() => {
                this.fetchData();
            }, 2000);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public removePolicy(): void {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            (this.service as ManagementService).resetLoginPolicyToDefault().then(() => {
                this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
                this.loading = true;
                setTimeout(() => {
                    this.fetchData();
                }, 2000);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public openDialog(): void {
        const dialogRef = this.dialog.open(AddIdpDialogComponent, {
            data: {
                serviceType: this.serviceType,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp && resp.idp) {
                this.addIdp(resp.idp).then(() => {
                    this.loading = true;
                    setTimeout(() => {
                        this.fetchData();
                    }, 2000);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    private addIdp(idp: IDP.AsObject | IDP.AsObject): Promise<any> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).addIDPToLoginPolicy(idp.id);
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).addIDPToLoginPolicy(idp.id);
        }
    }

    public removeIdp(idp: AdminIdpProviderView.AsObject | MgmtIdpProviderView.AsObject): void {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                (this.service as ManagementService).removeIDPFromLoginPolicy(idp.idpConfigId).then(() => {
                    const index = (this.idps as MgmtIdpProviderView.AsObject[]).findIndex(temp => temp === idp);
                    if (index > -1) {
                        this.idps.splice(index, 1);
                    }
                });
                break;
            case PolicyComponentServiceType.ADMIN:
                (this.service as AdminService).removeIDPFromLoginPolicy(idp.idpConfigId).then(() => {
                    const index = (this.idps as AdminIdpProviderView.AsObject[]).findIndex(temp => temp === idp);
                    if (index > -1) {
                        this.idps.splice(index, 1);
                    }
                });
                break;
        }
    }

    public get isDefault(): boolean {
        if (this.loginData && this.serviceType === PolicyComponentServiceType.MGMT) {
            return (this.loginData as LoginPolicy.AsObject).pb_default;
        } else {
            return false;
        }
    }
}
