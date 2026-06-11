# Get the status of a payment request
 Get the status of a payment request      

> Documentation Index
> -------------------
> 
> Fetch the complete documentation index at: [https://docs.xendit.co/llms.txt](https://docs.xendit.co/llms.txt)
> 
> Use this file to discover all available pages before exploring further.

What information are you looking for?

For archived content, access the previous documentation [here](https://archive.docs.xendit.co) or the previous API reference [here](https://archive.developers.xendit.co/).

Get the status of a payment request
===================================

*   Updated on Apr 27, 2026
*   Published on Jan 9, 2025

[Prev](https://docs.xendit.co/apidocs/payment-webhook-notification "Payment webhook notification") [Next](https://docs.xendit.co/apidocs/cancel-payment-request "Cancel a payment request")

Get

/v3/payment\_requests/{payment\_request\_id}

Get payment request status

Security

HTTP

Type basic

Header parameters

api-version

string

Valid values\[ "2024-11-11" \]

Path parameters

payment\_request\_id

stringRequired

Min length39

Max length39

Examplepr-8877c08a-740d-4153-9816-3d744ed197a5

Responses

200

Fetch Payment Request Status

application/json

getPaymentResponseCards getPaymentResponseRedirect getPaymentResponsePresentToCustomer

getPaymentResponseCards

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "90392f42-d98a-49ef-a7f3-abcezas123",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas123",
  "customer_id": "cust-90392f42-d98a-49ef-a7f3-abcezas123",
  "type": "PAY_AND_SAVE",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 10000.01,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "masked_card_number": "2222XXXXXXXX8888",
      "expiry_year": "2027",
      "expiry_month": "12",
      "cardholder_first_name": "John",
      "cardholder_last_name": "Doe",
      "cardholder_email": "john.doe@example.com",
      "cardholder_phone_number": "+6212345678902"
    },
    "skip_three_ds": false,
    "card_on_file_type": "CUSTOMER_UNSCHEDULED",
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success",
    "billing_information": {
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "phone_number": "+6212345678904",
      "city": "Singapore",
      "country": "SG",
      "postal_code": "644228",
      "street_line1": "Merlion Bay Sands Suites",
      "street_line2": "21-37",
      "province_state": "Singapore"
    },
    "statement_descriptor": "Goods & Services",
    "recurring_configuration": {
      "recurring_expiry": "2025-12-31",
      "recurring_frequency": 30
    }
  },
  "actions": [
    {
      "type": "REDIRECT_CUSTOMER",
      "value": "https://xendit.co/success",
      "descriptor": "WEB_URL"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Payment for invoice #INV-2025-001",
  "metadata": {
    "invoice_id": "INV-2025-001",
    "customer_type": "business"
  },
  "shipping_information": {
    "city": "Singapore",
    "country": "SG",
    "postal_code": "644228",
    "street_line1": "Merlion Bay Sands Suites",
    "street_line2": "21-37",
    "province_state": "Singapore"
  },
  "items": [
    {
      "reference_id": "item-123",
      "type": "PHYSICAL_PRODUCT",
      "name": "Vyson Dacuum Cleaner",
      "net_unit_amount": 10000.01,
      "quantity": 1,
      "category": "HOME_APPLIANCES"
    }
  ],
  "created": "2021-12-31T23:59:59Z",
  "updated": "2021-12-31T23:59:59Z"
}
```


getPaymentResponseRedirect

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "90392f42-d98a-49ef-a7f3-abcezas123",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas123",
  "customer_id": "cust-90392f42-d98a-49ef-a7f3-abcezas123",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 10000.01,
  "capture_method": "AUTOMATIC",
  "channel_code": "DANA",
  "channel_properties": {
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [
    {
      "type": "REDIRECT_CUSTOMER",
      "value": "https://xendit.co/success",
      "descriptor": "WEB_URL"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Payment for invoice #INV-2025-001",
  "metadata": {
    "invoice_id": "INV-2025-001",
    "customer_type": "business"
  },
  "shipping_information": {
    "city": "Singapore",
    "country": "SG",
    "postal_code": "644228",
    "street_line1": "Merlion Bay Sands Suites",
    "street_line2": "21-37",
    "province_state": "Singapore"
  },
  "items": [
    {
      "reference_id": "item-123",
      "type": "PHYSICAL_PRODUCT",
      "name": "Vyson Dacuum Cleaner",
      "net_unit_amount": 10000.01,
      "quantity": 1,
      "category": "HOME_APPLIANCES"
    }
  ],
  "created": "2021-12-31T23:59:59Z",
  "updated": "2021-12-31T23:59:59Z"
}
```


getPaymentResponsePresentToCustomer

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "90392f42-d98a-49ef-a7f3-abcezas123",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas123",
  "customer_id": "cust-90392f42-d98a-49ef-a7f3-abcezas123",
  "type": "REUSABLE_PAYMENT_CODE",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 10000.01,
  "capture_method": "AUTOMATIC",
  "channel_code": "BRI_VIRTUAL_ACCOUNT",
  "channel_properties": {
    "expires_at": "2024-12-31T23:59:59Z"
  },
  "actions": [
    {
      "type": "PRESENT_TO_CUSTOMER",
      "descriptor": "VIRTUAL_ACCOUNT_NUMBER",
      "value": "1251255"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Payment for invoice #INV-2025-001",
  "metadata": {
    "invoice_id": "INV-2025-001",
    "customer_type": "business"
  },
  "shipping_information": {
    "city": "Singapore",
    "country": "SG",
    "postal_code": "644228",
    "street_line1": "Merlion Bay Sands Suites",
    "street_line2": "21-37",
    "province_state": "Singapore"
  },
  "items": [
    {
      "reference_id": "item-123",
      "type": "PHYSICAL_PRODUCT",
      "name": "Vyson Dacuum Cleaner",
      "net_unit_amount": 10000.01,
      "quantity": 1,
      "category": "HOME_APPLIANCES"
    }
  ],
  "created": "2021-12-31T23:59:59Z",
  "updated": "2021-12-31T23:59:59Z"
}
```


Expand All

object

Payment request object

business\_id

string

Xendit-generated identifier for the business that owns the transaction

Example5f27a14a9bf05c73dd040bc8

reference\_id

string

A reference ID from merchants to identify their request. For "CARDS" channel code, reference ID must be unique.

Min length1

Max length255

payment\_request\_id

string

Xendit unique Payment Request ID generated as reference after creation of payment request.

Examplepr-1102feb0-bb79-47ae-9d1e-e69394d3949c

payment\_token\_id

string

Xendit unique Payment Token ID generated as reference for reusable payment details of the end user.

Examplept-cc3938dc-c2a5-43c4-89d7-7570793348c2

customer\_id

string

Xendit unique Capture ID generated as reference for the end user

Max length41

Examplecust-b98d6f63-d240-44ec-9bd5-aa42954c4f48

latest\_payment\_id

string

Latest Payment ID linked to the payment request.

Examplepy-1402feb0-bb79-47ae-9d1e-e69394d3949c

type

string

The payment collection intent type for the payment request.

PAY: Create a payment request that is able to receive one payment.

PAY\_AND\_SAVE: Create a payment request that is able to receive one payment. If the payment is successful, a reusable payment token will be returned for subsequent payment requests.

REUSABLE\_PAYMENT\_CODE: Create a payment request that is able to receive multiple payments. This is only used for repeat use payment method like a static QR, a predefined OTC payment code or a predefined Virtual Account number.

Valid values\[ "PAY", "PAY\_AND\_SAVE", "REUSABLE\_PAYMENT\_CODE" \]

country

string

ISO 3166-1 alpha-2 two-letter country code for the country of transaction.

Valid values\[ "ID", "PH", "VN", "TH", "SG", "MY", "HK", "MX" \]

ExampleID

currency

string

ISO 4217 three-letter currency code for the payment.

Valid values\[ "IDR", "PHP", "VND", "THB", "SGD", "MYR", "USD", "HKD", "AUD", "GBP", "EUR", "JPY", "MXN" \]

ExampleIDR

request\_amount

number

The intended payment amount to be collected from the end user.

Minimum0.0

Example10000.0

capture\_method

string

AUTOMATIC: payment capture will be processed immediately after payment request is created. MANUAL: payment capture requires merchant's trigger via payment capture endpoint before being processed

Valid values\[ "AUTOMATIC", "MANUAL" \]

Default"AUTOMATIC"

ExampleAUTOMATIC

channel\_code

string

Channel code used to select the payment method provider.

channel\_properties

object (Payments\_API\_ChannelProperties)

Data required to initiate transaction with payment method provider. Refer to the Channel Data Finder widget in the channel\_code field above for the full list of required properties for each channel.

actions

Array of object (Payments\_API\_Actions)

object

Actions object contains possible next steps merchants can take to proceed with payment collection from end user

type

The type of action that merchant system will need to handle to complete payment.

Valid values\[ "PRESENT\_TO\_CUSTOMER", "REDIRECT\_CUSTOMER", "API\_POST\_REQUEST" \]

descriptor

The type of action that merchant system will need to handle to complete payment.

Valid values\[ "CAPTURE\_PAYMENT", "PAYMENT\_CODE", "QR\_STRING", "VIRTUAL\_ACCOUNT\_NUMBER", "WEB\_URL", "DEEPLINK\_URL", "VALIDATE\_OTP", "RESEND\_OTP" \]

value

string

The specific value that will be used by merchant to complete the action

status

string

Status of the payment request.

Valid values\[ "ACCEPTING\_PAYMENTS", "REQUIRES\_ACTION", "AUTHORIZED", "CANCELED", "EXPIRED", "SUCCEEDED", "FAILED" \]

ExampleSUCCEEDED

failure\_code

string

Failure codes for payments.

Valid values\[ "ACCOUNT\_ACCESS\_BLOCKED", "INVALID\_MERCHANT\_SETTINGS", "INVALID\_ACCOUNT\_DETAILS", "PAYMENT\_ATTEMPT\_COUNTS\_EXCEEDED", "USER\_DEVICE\_UNREACHABLE", "CHANNEL\_UNAVAILABLE", "INSUFFICIENT\_BALANCE", "ACCOUNT\_NOT\_ACTIVATED", "INVALID\_TOKEN", "SERVER\_ERROR", "PARTNER\_TIMEOUT\_ERROR", "TIMEOUT\_ERROR", "USER\_DECLINED\_PAYMENT", "USER\_DID\_NOT\_AUTHORIZE", "PAYMENT\_REQUEST\_EXPIRED", "FAILURE\_DETAILS\_UNAVAILABLE", "EXPIRED\_OTP", "INVALID\_OTP", "PAYMENT\_AMOUNT\_LIMITS\_EXCEEDED", "OTP\_ATTEMPT\_COUNTS\_EXCEEDED", "CARD\_DECLINED", "DECLINED\_BY\_ISSUER", "ISSUER\_UNAVAILABLE", "INVALID\_CVV", "DECLINED\_BY\_PROCESSOR", "CAPTURE\_AMOUNT\_EXCEEDED ", "AUTHENTICATION\_FAILED", "PROCESSOR\_ERROR", "EXPIRED\_CARD", "STOLEN\_CARD", "INACTIVE\_OR\_UNAUTHORIZED\_CARD", "INVALID\_MERCHANT\_CREDENTIALS", "SUSPECTED\_FRAUDULENT" \]

ExampleCARD\_DECLINED

description

string

A custom description for the Payment Request.

Min length1

Max length1000

ExamplePayment for your order #123

metadata

object (Payments\_API\_MerchantMetadata)

Key-value entries for your custom data. You can specify up to 50 keys, with key names up to 40 characters and values up to 500 characters. This is for your convenience. Xendit will not use this data for any processing.

Example{ "my\_custom\_id": "merchant-123", "my\_custom\_order\_id": "order-123" }

items

Array of object (Payments\_API\_XenditStandardItem)

Array of objects describing the item/s attached to the payment.

object

reference\_id

string

Merchant provided identifier for the item

Min length1

Max length255

type

Type of item

Valid values\[ "DIGITAL\_PRODUCT", "PHYSICAL\_PRODUCT", "DIGITAL\_SERVICE", "PHYSICAL\_SERVICE", "FEE" \]

name

string

Name of item

Min length1

Max length255

net\_unit\_amount

number

Net amount to be charged per unit

quantity

integer

Number of units of this item in the basket

Minimum1.0

url

string

URL of the item. Must be HTTPS or HTTP

image\_url

string

URL of the image of the item. Must be HTTPS or HTTP

category

string

Category for item

Max length255

subcategory

string

Sub-category for item

Max length255

description

string

Description of item

Max length255

metadata

object (Payments\_API\_MerchantMetadata)

Key-value entries for your custom data. You can specify up to 50 keys, with key names up to 40 characters and values up to 500 characters. This is for your convenience. Xendit will not use this data for any processing.

Example{ "my\_custom\_id": "merchant-123", "my\_custom\_order\_id": "order-123" }

shipping\_information

object (Payments\_API\_XenditStandardShippingInformation)

country

2-letter ISO 3166-2 country code for the customer’s shipping country

Valid values\[ "ID", "PH", "VN", "TH", "SG", "MY", "MX" \]

street\_line1

string

Building name and apartment unit number

Min length1

Max length255

street\_line2

string

Building street address

Min length1

Max length255

city

string

City, village or town as appropriate

Min length1

Max length255

province\_state

string

Either one of (whichever is applicable): Geographic area, province, or region / Formal state designation within country

Min length1

Max length255

postal\_code

string

Postal, zip or rural delivery code, if applicable

Min length1

Max length255

created

string (date-time)

ISO 8601 date-time format.

Example2021-12-31T23:59:59Z

updated

string (date-time)

ISO 8601 date-time format.

Example2021-12-31T23:59:59Z

400

Bad request

application/json

OneOf

Payments\_API\_Http400ApiValidationError

object (Payments\_API\_Http400ApiValidationError)

error\_code

string

Valid values\[ "API\_VALIDATION\_ERROR" \]

message

string

Fields or values in the payment request does not comply with our API specification. Check the specific error message for debugging.

404

Not found

application/json

OneOf

Payments\_API\_Http404DataNotFound

object (Payments\_API\_Http404DataNotFound)

error\_code

string

Valid values\[ "DATA\_NOT\_FOUND" \]

message

string

ID specified in request cannot be found.

500

Internal server error

application/json

OneOf

Payments\_API\_Http500ServerError

object (Payments\_API\_Http500ServerError)

error\_code

string

Valid values\[ "SERVER\_ERROR" \]

message

string

An unexpected error occured, our team has been notified and will troubleshoot the issue

Was this article helpful?

Yes No