import { Component, EventEmitter, Input, Output } from '@angular/core';
import { OIDCAuthMethodType, OIDCGrantType, OIDCResponseType } from 'src/app/proto/generated/management_pb';

export interface RadioItemAuthType {
    key: string;
    titleI18nKey: string;
    descI18nKey: string;
    checked: boolean,
    disabled: boolean,
    prefix: string;
    background: string;
    responseType: OIDCResponseType;
    grantType: OIDCGrantType;
    authMethod: OIDCAuthMethodType;
    recommended: boolean;
    notRecommended: boolean;
}

@Component({
    selector: 'app-auth-method-radio',
    templateUrl: './app-auth-method-radio.component.html',
    styleUrls: ['./app-auth-method-radio.component.scss'],
})
export class AppAuthMethodRadioComponent {
    selected: string = '';
    @Input() authMethods!: RadioItemAuthType[];
    @Output() selectedType: EventEmitter<string> = new EventEmitter();

    public emitChange(): void {
        console.log('ch');
        this.selectedType.emit(this.selected);
    }
}