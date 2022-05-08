// https://docs.cypress.io/api/introduction/api.html

describe('My First Test', () => {
  it('visits the app root url', () => {
    cy.intercept('GET', '/api/v1/info', {
      statusCode: 200,
      body: {
        version: 'integration-test',
      },
    })
    cy.visit('/')
    cy.contains('div', 'Cert Manager Selfservice')
  })
})
