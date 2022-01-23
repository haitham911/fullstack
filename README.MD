Tasks
● REST API should be implemented consuming and producing “application/json”
● Implement user model with username, password, deposit and role (seller or buyer
) fields
● Implement product model with amountAvailable, cost, productName and sellerId
fields
● Implement an authentication method (basic, oAuth, JWT or something else, the
choice is yours)
● All endpoints should be authenticated except the user endpoints
● Implement CRUD for user and product models
● Implement /buy endpoint (accepts productId, amount of products) so users with
a “buyer” role can buy products with the money they’ve deposited.
Evaluation criteria:
● Language/Framework of choice best practises
● Edge cases covered
● Write tests for buy and one CRUD endpoint of your choice
● Add Readme file