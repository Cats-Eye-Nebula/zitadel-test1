function CheckInitPwPolicy() {
    let policyElement = document.getElementById("password");
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("passwordconfirm").value;
    let button = document.getElementById("init-button");

    if (pwNew == "" || pwNewConfirmation == "" || pwNew != pwNewConfirmation) {
        button.classList.add("disabled");
    } else {
        button.classList.remove("disabled");
    }

    ComplexityPolicyCheck(policyElement, button, pwNew);
}

let newPW = document.getElementById("password");
newPW.addEventListener('input', CheckInitPwPolicy);

let newPWConfirmation= document.getElementById("passwordconfirm");
newPWConfirmation.addEventListener('input', CheckInitPwPolicy);
