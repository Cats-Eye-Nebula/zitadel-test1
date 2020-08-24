import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { ProjectGrantView } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-granted-project-list',
    templateUrl: './granted-project-list.component.html',
    styleUrls: ['./granted-project-list.component.scss'],
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
export class GrantedProjectListComponent implements OnInit, OnDestroy {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;

    public dataSource: MatTableDataSource<ProjectGrantView.AsObject> =
        new MatTableDataSource<ProjectGrantView.AsObject>();
    @ViewChild(MatPaginator) public paginator!: MatPaginator;

    public grantedProjectList: ProjectGrantView.AsObject[] = [];
    public displayedColumns: string[] = ['select', 'name', 'resourceOwnerName', 'state', 'creationDate', 'changeDate'];
    public selection: SelectionModel<ProjectGrantView.AsObject> = new SelectionModel<ProjectGrantView.AsObject>(true, []);

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public grid: boolean = true;
    private subscription?: Subscription;

    constructor(private router: Router,
        public translate: TranslateService,
        private projectService: ManagementService,
        private toast: ToastService,
    ) { }

    public ngOnInit(): void {
        this.getData(10, 0);
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
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
        this.getData(event.pageSize, event.pageIndex);
    }

    public addProject(): void {
        this.router.navigate(['/projects', 'create']);
    }

    private async getData(limit: number, offset: number): Promise<void> {
        this.loadingSubject.next(true);
        this.projectService.SearchGrantedProjects(limit, offset).then(res => {
            const response = res.toObject();
            this.grantedProjectList = response.resultList;
            this.totalResult = response.totalResult;
            if (response.viewTimestamp) {
                this.viewTimestamp = response.viewTimestamp;
            }
            if (this.totalResult > 5) {
                this.grid = false;
            }
            this.dataSource.data = this.grantedProjectList;
            this.loadingSubject.next(false);
        }).catch(error => {
            console.error(error);
            this.toast.showError(error);
            this.loadingSubject.next(false);
        });
    }

    public refreshPage(): void {
        this.selection.clear();
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
    }
}
