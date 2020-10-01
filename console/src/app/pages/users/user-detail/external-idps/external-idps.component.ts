import { SelectionModel } from '@angular/cdk/collections';
import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';

import { ExternalIDPView as AuthExternalIDPView } from '../../../../proto/generated/auth_pb';
import {
    ExternalIDPSearchResponse,
    ExternalIDPView as MgmtExternalIDPView,
} from '../../../../proto/generated/management_pb';
import { GrpcAuthService } from '../../../../services/grpc-auth.service';
import { ManagementService } from '../../../../services/mgmt.service';
import { ToastService } from '../../../../services/toast.service';

@Component({
    selector: 'app-external-idps',
    templateUrl: './external-idps.component.html',
    styleUrls: ['./external-idps.component.scss'],
})
export class ExternalIdpsComponent implements OnInit {
    @Input() service!: GrpcAuthService | ManagementService;
    @Input() userId!: string;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public externalIdpResult!: ExternalIDPSearchResponse.AsObject;
    public dataSource: MatTableDataSource<MgmtExternalIDPView.AsObject | AuthExternalIDPView.AsObject>
        = new MatTableDataSource<MgmtExternalIDPView.AsObject | AuthExternalIDPView.AsObject>();
    public selection: SelectionModel<MgmtExternalIDPView.AsObject | AuthExternalIDPView.AsObject>
        = new SelectionModel<MgmtExternalIDPView.AsObject | AuthExternalIDPView.AsObject>(true, []);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    @Input() public displayedColumns: string[] = ['idpConfigId', 'idpName', 'externalUserId', 'externalUserDisplayName'];

    constructor(private toast: ToastService) { }

    ngOnInit(): void {
        this.getData(10, 0);
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.data.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.data.forEach(row => this.selection.select(row));
    }

    public changePage(event: PageEvent): void {
        this.getData(event.pageSize, event.pageIndex * event.pageSize);
    }

    private async getData(limit: number, offset: number): Promise<void> {
        this.loadingSubject.next(true);

        let promise;
        if (this.service instanceof ManagementService) {
            promise = (this.service as ManagementService).SearchUserExternalIDPs(limit, offset, this.userId);
        } else if (this.service instanceof GrpcAuthService) {
            promise = (this.service as GrpcAuthService).SearchMyExternalIdps(limit, offset);
        }

        if (promise) {
            promise.then(resp => {
                this.externalIdpResult = resp.toObject();
                this.dataSource.data = this.externalIdpResult.resultList;
                console.log(this.externalIdpResult.resultList);
                this.loadingSubject.next(false);
            }).catch((error: any) => {
                this.toast.showError(error);
                this.loadingSubject.next(false);
            });
        }
    }

    public refreshPage(): void {
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
    }
}
