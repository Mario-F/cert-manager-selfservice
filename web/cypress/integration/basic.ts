// https://docs.cypress.io/api/introduction/api.html

const mockInfo = () => {
  cy.intercept('GET', '/api/v1/info', {
    statusCode: 200,
    body: {
      version: 'integration-test',
    },
  })
}

describe('checking absolute basics', () => {
  it('visit root and test for brand', () => {
    mockInfo()
    cy.visit('/')
    cy.contains('div', 'Cert Manager Selfservice')
  })
})
