import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { MultiFactorType } from 'src/app/proto/generated/zitadel/policy_pb';
import { MultiFactor, MultiFactorState } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { _base64ToArrayBuffer } from '../../u2f-util';
import { DialogOtpComponent } from '../dialog-otp/dialog-otp.component';
import { DialogU2FComponent, U2FComponentDestination } from '../dialog-u2f/dialog-u2f.component';

export interface WebAuthNOptions {
    challenge: string;
    rp: { name: string, id: string; };
    user: { name: string, id: string, displayName: string; };
    pubKeyCredParams: any;
    authenticatorSelection: { userVerification: string; };
    timeout: number;
    attestation: string;
}

@Component({
    selector: 'app-auth-user-mfa',
    templateUrl: './auth-user-mfa.component.html',
    styleUrls: ['./auth-user-mfa.component.scss'],
})
export class AuthUserMfaComponent implements OnInit, OnDestroy {
    public displayedColumns: string[] = ['type', 'attr', 'state', 'actions'];
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    @ViewChild(MatTable) public table!: MatTable<MultiFactor.AsObject>;
    @ViewChild(MatSort) public sort!: MatSort;
    public dataSource!: MatTableDataSource<MultiFactor.AsObject>;

    public MultiFactorState: any = MultiFactorState;

    public error: string = '';
    public otpAvailable: boolean = false;

    constructor(
        private service: GrpcAuthService,
        private toast: ToastService,
        private dialog: MatDialog
    ) { }

    public ngOnInit(): void {
        this.getMFAs();
    }

    public ngOnDestroy(): void {
        this.loadingSubject.complete();
    }

    public addOTP(): void {
        this.service.addMyMultiFactorOTP().then((otpresp) => {
            const otp = otpresp;
            const dialogRef = this.dialog.open(DialogOtpComponent, {
                data: otp.url,
                width: '400px',
            });

            dialogRef.afterClosed().subscribe((code) => {
                if (code) {
                    this.service.verifyMyMultiFactorOTP(code).then(() => {
                        this.getMFAs();
                    });
                }
            });
        }, error => {
            this.toast.showError(error);
        });
    }

    public addU2F(): void {
        this.service.addMyMultiFactorU2F().then((u2fresp) => {
            const credOptions: CredentialCreationOptions = JSON.parse(atob(u2fresp.key?.publicKey as string));

            if (credOptions.publicKey?.challenge) {
                credOptions.publicKey.challenge = _base64ToArrayBuffer(credOptions.publicKey.challenge as any);
                credOptions.publicKey.user.id = _base64ToArrayBuffer(credOptions.publicKey.user.id as any);
                if (credOptions.publicKey.excludeCredentials) {
                    credOptions.publicKey.excludeCredentials.map(cred => {
                        cred.id = _base64ToArrayBuffer(cred.id as any);
                        return cred;
                    });
                }
                console.log(credOptions);
                const dialogRef = this.dialog.open(DialogU2FComponent, {
                    width: '400px',
                    data: {
                        credOptions,
                        type: U2FComponentDestination.MFA,
                    },
                });

                dialogRef.afterClosed().subscribe(done => {
                    if (done) {
                        this.getMFAs();
                    } else {
                        this.getMFAs();
                    }
                });
            }

        }, error => {
            this.toast.showError(error);
        });
    }

    public getMFAs(): void {
        this.service.listMyMultiFactors().then(mfas => {
            const list = mfas.resultList;
            this.dataSource = new MatTableDataSource(list);
            this.dataSource.sort = this.sort;

            const index = list.findIndex(mfa => mfa.otp);
            if (index === -1) {
                this.otpAvailable = true;
            }
        }).catch(error => {
            this.error = error.message;
        });
    }

    public deleteMFA(type: MultiFactorType, id?: string): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.MFA.DIALOG.MFA_DELETE_TITLE',
                descriptionKey: 'USER.MFA.DIALOG.MFA_DELETE_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                if (type === MultiFactorType.OTP) {
                    this.service.removeMyMultiFactorOTP().then(() => {
                        this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

                        const index = this.dataSource.data.findIndex(mfa => mfa.type === type);
                        if (index > -1) {
                            this.dataSource.data.splice(index, 1);
                        }
                        this.getMFAs();
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                } else if (type === MultiFactorType.MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION && id) {
                    this.service.removeMyMultiFactorU2F(id).then(() => {
                        this.toast.showInfo('USER.TOAST.U2FREMOVED', true);

                        const index = this.dataSource.data.findIndex(mfa => mfa.type === type);
                        if (index > -1) {
                            this.dataSource.data.splice(index, 1);
                        }
                        this.getMFAs();
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }


}
