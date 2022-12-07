import { Component, Input, OnInit } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { SetDefaultLanguageResponse, SetSecurityPolicyRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-security-policy',
  templateUrl: './security-policy.component.html',
  styleUrls: ['./security-policy.component.scss'],
})
export class SecurityPolicyComponent implements OnInit {
  public originsList: string[] = ['asdf', 'asdf'];
  public enabled: boolean = false;

  public loading: boolean = false;

  @Input() public originsControl: UntypedFormControl = new UntypedFormControl({ value: [], disabled: false });

  constructor(private service: AdminService, private toast: ToastService) {}

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.service.getSecurityPolicy().then((securityPolicy) => {
      if (securityPolicy.policy) {
        this.enabled = securityPolicy.policy?.enableIframeEmbedding;
        this.originsList = securityPolicy.policy?.allowedOriginsList;
      }
    });
  }

  private updateData(): Promise<SetDefaultLanguageResponse.AsObject> {
    const req = new SetSecurityPolicyRequest();
    req.setAllowedOriginsList(this.originsList);
    req.setEnableIframeEmbedding(this.enabled);
    return (this.service as AdminService).setSecurityPolicy(req);
  }

  public savePolicy(): void {
    const prom = this.updateData();
    this.loading = true;
    if (prom) {
      prom
        .then(() => {
          this.toast.showInfo('POLICY.SECURITY_POLICY.SAVED', true);
          this.loading = false;
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.loading = false;
          this.toast.showError(error);
        });
    }
  }

  public add(input: any): void {
    if (this.originsControl.valid) {
      if (input.value !== '' && input.value !== ' ' && input.value !== '/') {
        this.originsList.push(input.value);
      }
      if (input) {
        input.value = '';
      }
    }
  }

  public remove(redirect: any): void {
    const index = this.originsList.indexOf(redirect);

    if (index >= 0) {
      this.originsList.splice(index, 1);
    }
  }
}
