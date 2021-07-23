// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
// ###############################


it('LOGIN: Fill in credentials and login', () => {
    //console login
    cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    //wait for console to load
    cy.wait(5000)
})


describe('PROJECT: show Projects ', () => {
    it('PROJECT: show Projects ', () => {
        cy.get('a[href*="projects"').eq(0).click()
        cy.url().should('contain', '/projects')
    })
})

describe('PROJECT: add Project ', () => {

    it('PROJECT: add Project ', () => {
        cy.get('a[href*="projects"').eq(0).click()
        cy.url().should('contain', '/projects')
        cy.get('.add-project-button').click()
        cy.get('input').type("newProjectToTest")
        cy.get('[type^=submit]').click()
    })
})

describe('PROJECT: create app in Project ', () => {

    it('PROJECT: create app ', () => {
        //click on org to clear screen
        cy.get('a[href*="org"').eq(0).click()
        cy.wait(1000)
        cy.get('a[href*="projects"').eq(0).click()
        cy.url().should('contain', '/projects')
        cy.wait(1000)
        cy.get('.card').contains("newProjectToTest").click()
        cy.get('.cnsl-app-card').filter(':contains("add")').click()
        cy.get('[formcontrolname^=name]').type("newAppToTest")
        // select webapp
        cy.get('[for^=WEB]').click()
        cy.get('[type^=submit]').filter(':contains("Continue")').should('be.visible').eq(0).click()
        //select authentication
        cy.get('[for^=PKCE]').click()
        cy.get('[type^=submit]').filter(':contains("Continue")').should('be.visible').eq(1).click()
        //enter URL
        cy.get('cnsl-redirect-uris').eq(0).type("https://testurl.org")
        cy.get('cnsl-redirect-uris').eq(1).type("https://testlogouturl.org")
        cy.get('[type^=submit]').filter(':contains("Continue")').should('be.visible').eq(2).click()
        cy.get('button').filter(':contains("Create")').should('be.visible').click()
        //wait for application to be created
        cy.wait(3000)
        //TODO: check client ID/Secret
        cy.get('button').filter(':contains("Close")').should('be.visible').click()
    })
})

describe('PROJECT: delete Project ', () => {
    it('PROJECT: delete Project ', () => {
        cy.log(`PROJECT: delete project`);
        //click on org to clear screen
        cy.get('a[href*="org"').eq(0).click()
        //click on Projects 
        cy.get('a[href*="projects"').eq(0).click()
        cy.url().should('contain', '/projects')
        cy.wait(3000)
        //TODO variable for regex
        cy.get('tr').filter(':contains("demo")').find('button').click()
    })
})

