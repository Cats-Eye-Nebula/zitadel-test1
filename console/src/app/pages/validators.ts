import { UntypedFormControl } from '@angular/forms';

export function symbolValidator(c: UntypedFormControl): any {
    const REGEXP: RegExp = /[^a-z0-9]/gi;

    return REGEXP.test(c.value) ? null : {
        invalid: true,
        symbolValidator: {
            valid: false,
        },
    };
}

export function numberValidator(c: UntypedFormControl): any {
    const REGEXP = /[0-9]/g;

    return REGEXP.test(c.value) ? null : {
        invalid: true,
        numberValidator: {
            valid: false,
        },
    };
}

export function upperCaseValidator(c: UntypedFormControl): any {
    const REGEXP = /[A-Z]/g;

    return REGEXP.test(c.value) ? null : {
        invalid: true,
        upperCaseValidator: {
            valid: false,
        },
    };
}

export function lowerCaseValidator(c: UntypedFormControl): any {
    const REGEXP = /[a-z]/g;

    return REGEXP.test(c.value) ? null : {
        invalid: true,
        lowerCaseValidator: {
            valid: false,
        },
    };
}
