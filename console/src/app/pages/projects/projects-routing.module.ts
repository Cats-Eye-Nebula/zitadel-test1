import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { GrantedProjectDetailComponent } from './granted-project-detail/granted-project-detail.component';
import { OwnedProjectDetailComponent } from './owned-project-detail/owned-project-detail.component';
import { ProjectsComponent } from './projects.component';

const routes: Routes = [
    {
        path: '',
        component: ProjectsComponent,
        data: { animation: 'HomePage' },
    },
    {
        path: 'create',
        loadChildren: () => import('../project-create/project-create.module').then(m => m.ProjectCreateModule),
    },
    {
        path: ':id/grant/:grantId',
        component: GrantedProjectDetailComponent,
        data: { animation: 'HomePage' },
    },
    {
        path: ':id',
        component: OwnedProjectDetailComponent,
        data: { animation: 'HomePage' },
    },
    {
        path: ':projectid/members',
        loadChildren: () => import('./project-members/project-members.module').then(m => m.ProjectMembersModule),
    },
    {
        path: ':projectid/apps',
        data: { animation: 'AddPage' },
        loadChildren: () => import('../apps/apps.module').then(m => m.AppsModule),
    },
    {
        path: ':projectid/roles/create',
        loadChildren: () => import('../project-role-create/project-role-create.module').then(m => m.ProjectRoleCreateModule),

    },
    {
        path: ':projectid/grants/create',
        loadChildren: () => import('../project-grant-create/project-grant-create.module')
            .then(m => m.ProjectGrantCreateModule),
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class ProjectsRoutingModule { }
