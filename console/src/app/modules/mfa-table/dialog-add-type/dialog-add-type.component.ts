import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { MultiFactorType as AdminMultiFactorType } from 'src/app/proto/generated/admin_pb';
import { MultiFactorType as MgmtMultiFactorType } from 'src/app/proto/generated/management_pb';

enum LoginMethodComponentType {
    MultiFactor = 1,
    SecondFactor = 2,
}

@Component({
    selector: 'app-dialog-add-type',
    templateUrl: './dialog-add-type.component.html',
    styleUrls: ['./dialog-add-type.component.scss'],
})
export class DialogAddTypeComponent {
    public LoginMethodComponentType: any = LoginMethodComponentType;
    public newMfaType!: AdminMultiFactorType | MgmtMultiFactorType;
    public availableMfaTypes: Array<AdminMultiFactorType | MgmtMultiFactorType> = [];
    constructor(public dialogRef: MatDialogRef<DialogAddTypeComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) {
        this.availableMfaTypes = data.types;
    }

    public closeDialog(): void {
        this.dialogRef.close();
    }

    public closeDialogWithCode(): void {
        this.dialogRef.close(this.newMfaType);
    }
}
