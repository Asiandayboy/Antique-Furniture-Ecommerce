1/15
- implement session management [DONE]

1/16
- send session id in response in login handler [DONE]
- create ListFurniture/GetFurnitures/GetFurniture endpoints [DONE]
  - create tests for them [DONE] (1/22)


1/22
- update to include authentication to necessary endpoint handlers [DONE]
- update those tests as well ^ [DONE]

1/23 
- Need to link userID with each furniture listing [DONE]
- /account endpoint (GET)
- change Session struct type to use string for SessionID and not UUID, so this way [DONE]
  it can integrate nicely with the template version of CreateSession
  - tempalte version already started: Finish it [DONE]




1/25
- instead of sending the character-based represention of images
in the JSON object of the body, use multipart/form-data so that you can
upload the file and then include the JSON object, sending two data types
in one HTTP request. Access the image file through form and the JSON data
through the request body. 

When parsing the image data, you'll get the binary data of the images; Save 
that into the DB
- file upload that the browser handles
- HTTP POST with content 
- choose to send a separate request to retrieve the binary data for the image file


1/26
- update DeleteSession to delete nested maps????? or no?
- there is no expiry on sessionID
- implement /logout with test [DONE?]
- 1/26 updating types (users.go, auth.go, furniture.go) [DONE?]

1/28
- AuthMiddleware [WIP]
 - refactor handlers to include AuthMiddleware [DONE]
- !implement /checkout and test [WIP]

1/29
- commited request methods middleware changes [DONE]

1/30 
- (NOTES): Disregard caring about scaling; you're barely sending or storing 
any large amounts of data besides the images. Ur fine. Just do it 
  - so just ignore the note in furniture.go for GetFurnitures
- [!] return all the errors at once for ValidateListForm
- NOTES /checkout: [
  - receipts are saved the receipts collection [DONE]
  - when a user requests their order history, they will query the receipts
  collection using their userID
  - so, the receipts collection acts also as every user's history as well
  - delete furnitureListing from listings  [DONE]

  - when a purchased is successfully processed, that money needs to go to the 
  seller, and a percentage goes to the website. 98 seller/2 website split [DONE]
  - [
    This means, each user account must have a balance
    Balance must be Decimal128 (FIGURE THIS OUT!!!!!!!!!!!!!!!!! 1/31 1:34 AM)
  ]
]

2/1
- wrote tests and functons for float64 and dec128 conversions [DONE]


2/3
- incorporate images into checkout session [WIP]
- 2/4: Finished changing to multipart -> work on getting it to work with stripe now in checkout session
  --[bugged: images won't upload -> link won't work :(]
- finish fulfillment of order through webhook [DONE]
  - still need to email receipt and/or send the receipt in a request or something 



2/5
- Finished PUT and GET /account
- work on addresses endpoint

2/7
- implemented POST /account/address and test

2/8
- wrote test for PUT /account/address

2/13
As of go version 1.22, there is now enhanced routing patterns and wildcards
- don't need the custom method middleware anymore

2/14
- committed test/impl address GET
- committed test/impl address DELETE
- committed test/impl address PUT

2/20
- need to add confirmPass check in auth.go





# ENDPOINTS
===========
/login                            **DONE**
/signup                           **DONE**
/logout                           **DONE** 
/list_furniture                   **DONE**
/get_furniture?listingid={id}     **DONE**
/get_furnitures                   **DONE[?]**
/checkout                         **DONE**
/account?userid={id}              **DONE**[i think?]



FRONTEND



2/21
- session management could use some work cuz it's kinda stupid rn [DONE]
- delete session id cookie when /logout is hit [DONE]

2/22
- conditionally render navbar links [DONE]
- use some sort of react context to pass auth prop down component tree [DONE]
- build fundamenetal design of dashboard



# Planned
=========
- CSS styles for login and signup [DONE]
- CSS styles for navbar [DONE]
- CSS styles for dashboard [DONE]
- populate dashboard with user data [DONE]
- populate market with user data [DONE]
- populate MyAddresses with user data [DONE]
- create component for listing furniture [DONE]


plans established on 3/5
- create form for listing furniture [DONE]
- refine rendering for market listings [DONE]
  - render images for each listing, etc
- add styles to MyAddresses [LATER]


3/12 - 3/13
- add checkout page [DONE]
  - could use some sort of loading????
- render listings conditionally based on bought or not (SOLD tag) [DONE]


3/19
- styles for List a Furniture [DONE]
- styles for dashboard [DONE]
- styles for login and signup [DONE]
- styles for shopping cart page [DONE] - *3/25*


3/25
- styles for Your Furniture Listings [DONE] - *3/25*
- styles for detailed view of Furniture listing in Market [DONE] - *3/25*
- edit and delete functionality for shipping addresses and styles [DONE] - *3/26*
- edit functionality for LoginSecurity and styles [DONE] *3/29/24*

3/26
- styles for detailed purchase history [DONE] *3/29/24*



# Things left to do
- LATEST FURNITURE LISTED [DONE] *3/30/24*
- SUBSCRIBE [DONE] *4/4/24* 
  - email subscription notification for new listings
    - backend: write logic 
      - logic for subscription state [DONE] *4/2/24*
      - logic for sending email [DONE] *4/4/24*
    - frontend: create button with CSS styles on market page [DONE] *4/2/24*
- input validation for POST? [MAYBE]
- delete button for shopping cart items [DONE] *3/30/24*
- documentation



# Final touches 4/18
- input validation for login page [DONE] *4/18*
- better error visuals for input validations (login and signup) [DONE] *4/18*
- use selenium to test signup and login page


# notes
- focus on UX/documentation/input validation