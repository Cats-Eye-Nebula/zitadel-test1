import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';
import { PolicyComponentServiceType, PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';

import { OrgCreateComponent } from './org-create/org-create.component';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgGridComponent } from './org-grid/org-grid.component';

const routes: Routes = [
    {
        path: 'create',
        component: OrgCreateComponent,
        canActivate: [RoleGuard],
        data: {
            roles: ['(org.create)?(iam.write)?'],
        },
        loadChildren: () => import('./org-create/org-create.module').then(m => m.OrgCreateModule),
    },
    {
        path: 'idp/:id',
        loadChildren: () => import('src/app/modules/idp/idp.module').then(m => m.IdpModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['iam.idp.read'],
            serviceType: PolicyComponentServiceType.ADMIN,
        },
    },
    {
        path: 'idp/create',
        loadChildren: () => import('src/app/modules/idp-create/idp-create.module').then(m => m.IdpCreateModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['org.idp.write'],
        },
    },
    {
        path: 'policy',
        children: [
            {
                path: PolicyComponentType.AGE,
                loadChildren: () => import('src/app/modules/policies/password-age-policy/password-age-policy.module')
                    .then(m => m.PasswordAgePolicyModule),
            },
            {
                path: PolicyComponentType.LOCKOUT,
                loadChildren: () => import('src/app/modules/policies/password-lockout-policy/password-lockout-policy.module')
                    .then(m => m.PasswordLockoutPolicyModule),
            },
            {
                path: PolicyComponentType.COMPLEXITY,
                loadChildren: () => import('src/app/modules/policies/password-complexity-policy/password-complexity-policy.module')
                    .then(m => m.PasswordComplexityPolicyModule),
            },
            {
                path: PolicyComponentType.IAM,
                loadChildren: () => import('src/app/modules/policies/password-iam-policy/password-iam-policy.module')
                    .then(m => m.PasswordIamPolicyModule),
            },
            {
                path: PolicyComponentType.LOGIN,
                data: {
                    serviceType: PolicyComponentServiceType.MGMT,
                },
                loadChildren: () => import('src/app/modules/policies/login-policy/login-policy.module')
                    .then(m => m.LoginPolicyModule),
            },
        ],
    },
    {
        path: 'members',
        loadChildren: () => import('./org-members/org-members.module').then(m => m.OrgMembersModule),
    },
    {
        path: '',
        component: OrgDetailComponent,
    },
    {
        path: 'overview',
        component: OrgGridComponent,
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class OrgsRoutingModule { }
