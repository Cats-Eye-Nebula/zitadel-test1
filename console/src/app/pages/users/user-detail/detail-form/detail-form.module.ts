import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

import { DetailFormComponent } from './detail-form.component';

@NgModule({
    declarations: [
        DetailFormComponent,
    ],
    imports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        TranslateModule,
        MatSelectModule,
        MatButtonModule,
        MatIconModule,
        TranslateModule,
        FormFieldModule,
    ],
    exports: [
        DetailFormComponent,
    ],
})
export class DetailFormModule { }
