import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { ProjectRoleAdd } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-project-role-create',
    templateUrl: './project-role-create.component.html',
    styleUrls: ['./project-role-create.component.scss'],
    animations: [
        trigger('list', [
            transition(':enter', [
                query('@animate',
                    stagger(80, animateChild()),
                ),
            ]),
        ]),
        trigger('animate', [
            transition(':enter', [
                style({ opacity: 0, transform: 'translateY(-100%)' }),
                animate('100ms', style({ opacity: 1, transform: 'translateY(0)' })),
            ]),
            transition(':leave', [
                style({ opacity: 1, transform: 'translateY(0)' }),
                animate('100ms', style({ opacity: 0, transform: 'translateY(100%)' })),
            ]),
        ]),
    ],
})
export class ProjectRoleCreateComponent implements OnInit, OnDestroy {
    private subscription?: Subscription;
    public projectId: string = '';

    public formArray!: FormArray;
    public formGroup!: FormGroup;
    public createSteps: number = 1;
    public currentCreateStep: number = 1;

    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private toast: ToastService,
        private projectService: ProjectService,
        private _location: Location,
    ) {
        this.formGroup = new FormGroup({
            key: new FormControl('', [Validators.required]),
            displayName: new FormControl(''),
            group: new FormControl(''),
        });

        this.formArray = new FormArray([this.formGroup]);
    }

    public addEntry(): void {
        const newGroup = new FormGroup({
            key: new FormControl(''),
            displayName: new FormControl(''),
            group: new FormControl('', [Validators.required]),
        });

        this.formArray.push(newGroup);
    }

    public removeEntry(index: number): void {
        this.formArray.removeAt(index);
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    private getData({ projectid }: Params): void {
        this.projectId = projectid;
    }

    public addRole(): void {
        const promises = this.formArray.value.map((role: ProjectRoleAdd.AsObject) => {
            role.id = this.projectId;
            return this.projectService.AddProjectRole(role);
        });

        Promise.all(promises).then(() => {
            this.router.navigate(['projects', this.projectId]);
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }


    public close(): void {
        this._location.back();
    }
}
