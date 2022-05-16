import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { Component } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Action, ActionState, Flow, FlowType, TriggerType } from 'src/app/proto/generated/zitadel/action_pb';
import { SetTriggerActionsRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddFlowDialogComponent } from './add-flow-dialog/add-flow-dialog.component';

@Component({
  selector: 'cnsl-actions',
  templateUrl: './actions.component.html',
  styleUrls: ['./actions.component.scss'],
})
export class ActionsComponent {
  public flow!: Flow.AsObject;
  public flowType: FlowType = FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION;

  public typeControl: FormControl = new FormControl(FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION);

  public typesForSelection: FlowType[] = [FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION];

  public selection: Action.AsObject[] = [];
  public InfoSectionType: any = InfoSectionType;
  public ActionKeysType: any = ActionKeysType;

  public maxActions: number | null = null;
  public ActionState: any = ActionState;

  constructor(
    private mgmtService: ManagementService,
    breadcrumbService: BreadcrumbService,
    private dialog: MatDialog,
    private toast: ToastService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);

    this.loadFlow();
  }

  private loadFlow() {
    this.mgmtService.getFlow(this.flowType).then((flowResponse) => {
      if (flowResponse.flow) this.flow = flowResponse.flow;
    });
  }

  public clearFlow(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CLEAR',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'FLOWS.DIALOG.CLEAR.TITLE',
        descriptionKey: 'FLOWS.DIALOG.CLEAR.DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .clearFlow(this.flowType)
          .then((resp) => {
            this.loadFlow();
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public openAddTrigger(): void {
    const dialogRef = this.dialog.open(AddFlowDialogComponent, {
      data: {
        flowType: this.flowType,
        triggerType: TriggerType.TRIGGER_TYPE_POST_AUTHENTICATION,
        actions: this.selection && this.selection.length ? this.selection : [],
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: SetTriggerActionsRequest) => {
      if (req) {
        this.mgmtService
          .setTriggerActions(req.getActionIdsList(), req.getFlowType(), req.getTriggerType())
          .then((resp) => {
            this.loadFlow();
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });
  }

  drop(triggerActionsListIndex: number, array: any[], event: CdkDragDrop<Action.AsObject[]>) {
    moveItemInArray(array, event.previousIndex, event.currentIndex);
    this.saveFlow(triggerActionsListIndex);
  }

  saveFlow(index: number) {
    this.mgmtService
      .setTriggerActions(
        this.flow.triggerActionsList[index].actionsList.map((action) => action.id),
        this.flowType,
        this.flow.triggerActionsList[index].triggerType,
      )
      .then((updateResponse) => {
        this.toast.showInfo('FLOWS.TOAST.ACTIONSSET', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }
}
