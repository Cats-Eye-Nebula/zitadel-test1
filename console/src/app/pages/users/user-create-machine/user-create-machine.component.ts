import { Component, OnDestroy } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { CreateMachineRequest } from 'src/app/proto/generated/admin_pb';
import { UserResponse } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

function noEmailValidator(c: AbstractControl): any {
    const EMAIL_REGEXP: RegExp = /^((?!@).)*$/gm;
    if (!c.parent || !c) {
        return;
    }
    const username = c.parent.get('userName');

    if (!username) {
        return;
    }

    return EMAIL_REGEXP.test(username.value) ? null : {
        noEmailValidator: {
            valid: false,
        },
    };
}

@Component({
    selector: 'app-user-create-machine',
    templateUrl: './user-create-machine.component.html',
    styleUrls: ['./user-create-machine.component.scss'],
})
export class UserCreateMachineComponent implements OnDestroy {
    public user: CreateMachineRequest.AsObject = new CreateMachineRequest().toObject();
    public userForm!: FormGroup;

    private sub: Subscription = new Subscription();
    public loading: boolean = false;

    constructor(
        private router: Router,
        private toast: ToastService,
        public userService: ManagementService,
        private fb: FormBuilder,
    ) {
        this.initForm();
    }

    private initForm(): void {
        this.userForm = this.fb.group({
            userName: ['',
                [
                    Validators.required,
                    Validators.minLength(2),
                ],
            ],
            name: ['', [Validators.required]],
            description: ['', []],
        });
    }

    public createUser(): void {
        this.user = this.userForm.value;

        this.loading = true;

        const machineReq = new CreateMachineRequest();
        machineReq.setDescription(this.description?.value);
        machineReq.setName(this.name?.value);

        this.userService
            .CreateUserMachine(this.userName?.value, machineReq)
            .then((data: UserResponse) => {
                this.loading = false;
                this.toast.showInfo('USER.TOAST.CREATED', true);
                const id = data.getId();
                if (id) {
                    this.router.navigate(['users', id]);
                }
            })
            .catch((error: any) => {
                this.loading = false;
                this.toast.showError(error);
            });
    }

    ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    public get name(): AbstractControl | null {
        return this.userForm.get('name');
    }
    public get description(): AbstractControl | null {
        return this.userForm.get('description');
    }
    public get userName(): AbstractControl | null {
        return this.userForm.get('userName');
    }
}
