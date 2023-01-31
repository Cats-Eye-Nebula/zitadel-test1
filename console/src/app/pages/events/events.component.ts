import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatSort, Sort } from '@angular/material/sort';
import { MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { BehaviorSubject, Observable, Subject, takeUntil } from 'rxjs';
import { ListEventsRequest, ListEventsResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { LiveAnnouncer } from '@angular/cdk/a11y';
import { ToastService } from 'src/app/services/toast.service';
import { DisplayJsonDialogComponent } from 'src/app/modules/display-json-dialog/display-json-dialog.component';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

enum EventFieldName {
  EDITOR = 'editor',
  AGGREGATE = 'aggregate',
  RESOURCEOWNER = 'resourceOwner',
  SEQUENCE = 'sequence',
  CREATIONDATE = 'creationDate',
  TYPE = 'type',
  PAYLOAD = 'payload',
}

type LoadRequest = {
  req: ListEventsRequest;
  override: boolean;
};

@Component({
  selector: 'cnsl-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.scss'],
})
export class EventsComponent implements OnDestroy {
  public INITPAGESIZE = 20;
  public sortAsc = false;
  private destroy$: Subject<void> = new Subject();

  public displayedColumns: string[] = [
    EventFieldName.TYPE,
    EventFieldName.AGGREGATE,
    EventFieldName.RESOURCEOWNER,
    EventFieldName.EDITOR,
    EventFieldName.SEQUENCE,
    EventFieldName.CREATIONDATE,
    EventFieldName.PAYLOAD,
  ];

  public currentRequest$: BehaviorSubject<LoadRequest> = new BehaviorSubject<LoadRequest>({
    req: new ListEventsRequest().setLimit(this.INITPAGESIZE),
    override: true,
  });

  @ViewChild(MatSort) public sort!: MatSort;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<Event> = new MatTableDataSource<Event>([]);

  public _done: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public done: Observable<boolean> = this._done.asObservable();

  public _loading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private _data: BehaviorSubject<Event[]> = new BehaviorSubject<Event[]>([]);

  constructor(
    private adminService: AdminService,
    private breadcrumbService: BreadcrumbService,
    private _liveAnnouncer: LiveAnnouncer,
    private toast: ToastService,
    private dialog: MatDialog,
  ) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);

    this.currentRequest$.pipe(takeUntil(this.destroy$)).subscribe(({ req, override }) => {
      this._loading.next(true);
      this.adminService
        .listEvents(req)
        .then((res: ListEventsResponse) => {
          if (override) {
            this._data = new BehaviorSubject<Event[]>([]);
            this.dataSource = new MatTableDataSource<Event>([]);
          }

          const eventList = res.getEventsList();
          this._data.next(eventList);

          const concat = this.dataSource.data.concat(eventList);
          this.dataSource = new MatTableDataSource<Event>(concat);

          this._loading.next(false);

          if (eventList.length === 0) {
            this._done.next(true);
          } else {
            this._done.next(false);
          }
        })
        .catch((error) => {
          this.toast.showError(error);
          this._loading.next(false);
          this._data.next([]);
        });
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public loadEvents(filteredRequest: ListEventsRequest, override: boolean = false): void {
    this.currentRequest$.next({ req: filteredRequest, override });
    const { req } = this.currentRequest$.value;
  }

  public refresh(): void {
    const req = new ListEventsRequest();
    req.setLimit(this.paginator.pageSize);
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
      this.sortAsc = sortState.direction === 'asc';

      const { req } = this.currentRequest$.value;

      req.setLimit(this.INITPAGESIZE);
      req.setAsc(this.sortAsc ? true : false);

      this.loadEvents(req, true);
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  public openDialog(event: Event): void {
    this.dialog.open(DisplayJsonDialogComponent, {
      data: {
        event: event,
      },
      width: '450px',
    });
  }

  public more(): void {
    const sequence = this.getCursor();
    const { req } = this.currentRequest$.value;
    req.setSequence(sequence);
    this.loadEvents(req);
  }

  public filterChanged(filterRequest: ListEventsRequest) {
    console.log('filterchanged');

    const req = new ListEventsRequest();
    req.setLimit(this.INITPAGESIZE);
    req.setAsc(this.sortAsc ? true : false);

    req.setAggregateTypesList(filterRequest.getAggregateTypesList());
    req.setAggregateId(filterRequest.getAggregateId());
    req.setEventTypesList(filterRequest.getEventTypesList());
    req.setEditorUserId(filterRequest.getEditorUserId());
    req.setResourceOwner(filterRequest.getResourceOwner());
    req.setSequence(filterRequest.getSequence());
    req.setCreationDate(filterRequest.getCreationDate());

    console.log(req.toObject());
    this.loadEvents(req, true);
  }

  private getCursor(): number {
    const current = this._data.value;

    if (current.length) {
      const sequence = current[current.length - 1].toObject().sequence;
      return sequence;
    }
    return 0;
  }
}
