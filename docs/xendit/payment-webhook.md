> ## Documentation Index
>
> Fetch the complete documentation index at: [https://docs.xendit.co/llms.txt](https://docs.xendit.co/llms.txt)
>
> Use this file to discover all available pages before exploring further.

For archived content, access the previous documentation [here](https://archive.docs.xendit.co/) or the previous API reference [here](https://archive.developers.xendit.co/).

# Payment webhook notification

- Updated on Apr 26, 2026
- Published on Jan 22, 2025

[Prev](https://docs.xendit.co/apidocs/create-payment-request "Create a payment request")[Next](https://docs.xendit.co/apidocs/get-payment-request "Get the status of a payment request")

Post

/your\_payment\_webhook\_url

Webhook notification that will be sent to your defined webhook url for updates to payment status.

Header parameters

x-callback-token

string

Your Xendit unique webhook token to verify the origin of the webhook. It is highly recommended for your integration to verify this value.

Body parameters

Payment capture status callback

Show example

Show example

application/json

Code snippet

```json
{
  "paymentCapture": {
    "value": {
      "event": "payment.capture",
      "business_id": "6094fa76c2fd53701b8e079c",
      "created": "2021-12-02T14:52:21.566Z",
      "data": {
        "payment_id": "py-1fdaf346-dd2e-4b6c-b938-124c7167a822",
        "business_id": "6094fa76c2fd53701b8e079c",
        "status": "SUCCEEDED",
        "payment_request_id": "pr-1fdaf346-dd2e-4b6c-b938-124c7167a822",
        "request_amount": 10000,
        "customer_id": "cust-5ed61c4e-499f-49bd-9d90-f3f45028a7a3",
        "channel_code": "BRI_VIRTUAL_ACCOUNT",
        "country": "ID",
        "currency": "IDR",
        "reference_id": "example_reference_id",
        "description": "Payment description",
        "channel_properties": {
          "failure_return_url": "https://xendit.co/failure",
          "success_return_url": "https://xendit.co/success"
        },
        "type": "SINGLE_PAYMENT",
        "created": "2021-12-02T14:52:21.566Z",
        "updated": "2021-12-02T14:52:21.566Z"
      }
    }
  },
  "paymentAuthorization": {
    "value": {
      "event": "payment.authorization",
      "business_id": "6094fa76c2fd53701b8e079c",
      "created": "2021-12-02T14:52:21.566Z",
      "data": {
        "payment_id": "py-1fdaf346-dd2e-4b6c-b938-124c7167a822",
        "business_id": "6094fa76c2fd53701b8e079c",
        "status": "AUTHORIZED",
        "payment_request_id": "pr-1fdaf346-dd2e-4b6c-b938-124c7167a822",
        "request_amount": 10000,
        "customer_id": "cust-5ed61c4e-499f-49bd-9d90-f3f45028a7a3",
        "channel_code": "CARDS",
        "country": "PH",
        "currency": "PHP",
        "reference_id": "example_reference_id",
        "description": "Payment description",
        "channel_properties": {
          "failure_return_url": "https://xendit.co/failure",
          "success_return_url": "https://xendit.co/success"
        },
        "type": "SINGLE_PAYMENT",
        "created": "2021-12-02T14:52:21.566Z",
        "updated": "2021-12-02T14:52:21.566Z"
      }
    }
  },
  "paymentFailure": {
    "value": {
      "event": "payment.failure",
      "business_id": "6094fa76c2fd53701b8e079c",
      "created": "2021-12-02T14:52:21.566Z",
      "data": {
        "payment_id": "py-1fdaf346-dd2e-4b6c-b938-124c7167a822",
        "business_id": "6094fa76c2fd53701b8e079c",
        "status": "FAILED",
        "payment_request_id": "pr-1fdaf346-dd2e-4b6c-b938-124c7167a822",
        "request_amount": 10000,
        "customer_id": "cust-5ed61c4e-499f-49bd-9d90-f3f45028a7a3",
        "channel_code": "CARDS",
        "country": "TH",
        "currency": "THB",
        "reference_id": "example_reference_id",
        "description": "Payment description",
        "failure_code": "INSUFFICIENT_BALANCE",
        "channel_properties": {
          "failure_return_url": "https://xendit.co/failure",
          "success_return_url": "https://xendit.co/success"
        },
        "type": "SINGLE_PAYMENT",
        "created": "2021-12-02T14:52:21.566Z",
        "updated": "2021-12-02T14:52:21.566Z"
      }
    }
  }
}
```

JSON

Copy

Collapse all

object

Payment capture status callback for payment

event

string

Webhook event names for payment capture status updates.

Valid values\[\
"payment.capture",\
"payment.authorization",\
"payment.failure"\
\]

business\_id

string

Xendit-generated identifier for the business that owns the transaction

Example5f27a14a9bf05c73dd040bc8

created

string (date-time)

Timestamp of webhook delivery attempt in ISO 8601 date-time format.

Example2021-12-31T23:59:59Z

data

object (Payments\_API\_PaymentSchema)

Payment object

payment\_id

string

Xendit unique Payment ID generated as reference for a payment.

Examplepy-1402feb0-bb79-47ae-9d1e-e69394d3949c

business\_id

string

Xendit-generated identifier for the business that owns the transaction

Example5f27a14a9bf05c73dd040bc8

reference\_id

string

A Reference ID from merchants to identify their request.

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

type

string

The payment collection intent type for the payment request.

PAY: Create a payment request that is able to receive one payment.

PAY\_AND\_SAVE: Create a payment request that is able to receive one payment. If the payment is successful, a reusable payment token will be returned for subsequent payment requests.

REUSABLE\_PAYMENT\_CODE: Create a payment request that is able to receive multiple payments. This is only used for repeat use payment method like a static QR, a predefined OTC payment code or a predefined Virtual Account number.

Valid values\[\
"PAY",\
"PAY\_AND\_SAVE",\
"REUSABLE\_PAYMENT\_CODE"\
\]

country

string

ISO 3166-1 alpha-2 two-letter country code for the country of transaction.

Valid values\[\
"ID",\
"PH",\
"VN",\
"TH",\
"SG",\
"MY",\
"HK",\
"MX"\
\]

ExampleID

currency

string

ISO 4217 three-letter currency code for the payment.

Valid values\[\
"IDR",\
"PHP",\
"VND",\
"THB",\
"SGD",\
"MYR",\
"USD",\
"HKD",\
"AUD",\
"GBP",\
"EUR",\
"JPY",\
"MXN"\
\]

ExampleIDR

request\_amount

number

The intended payment amount to be collected from the end user.

Minimum0.0

Example10000.0

capture\_method

string

AUTOMATIC: payment capture will be processed immediately after payment request is created.
MANUAL: payment capture requires merchant's trigger via payment capture endpoint before being processed

Valid values\[\
"AUTOMATIC",\
"MANUAL"\
\]

Default"AUTOMATIC"

ExampleAUTOMATIC

channel\_code

string

Channel code used to select the payment method provider.

Xendit Documentation Widget

Channel Data FinderData required to initiate transaction with payment method provider. Use routing payment channels mapping for full list of data required.

captures

Array of object (Payments\_API\_Capture)

object

Capture object contains information about the capture that was performed

capture\_timestamp

string (date-time)

ISO 8601 date-time format.

Example2021-12-31T23:59:59Z

capture\_id

string

Xendit unique Capture ID generated as reference for a single capture.

Examplecap-1502feb0-bb79-47ae-9d1e-e69394d3949c

capture\_amount

number

The payment amount captured for this payment. Maximum capture amount can only be equal or lesser than the authorized amount value.

Minimum0.0

Example10000.0

status

string

Status of the payment.

Valid values\[\
"AUTHORIZED",\
"CANCELED",\
"SUCCEEDED",\
"FAILED",\
"EXPIRED",\
"PENDING"\
\]

ExampleSUCCEEDED

payment\_details

object (Payments\_API\_PaymentDetails)

Payment information provided by the payment method provider. Fields returned are dependent on what is made available by the provider.

Xendit Documentation Widget

Available Payment DetailsFind channel codes and their payment details for integration.

authorization\_data

object (Payments\_API\_AuthorizationData)

Specific to cards transaction only. Details about the card authorization processing.

authorization\_code

string

Authorization approval code from the scheme. 6 alphanumeric characters.

cvn\_verification\_result

string

Whether CVN input matches with the issuer's data.

Valid values\[\
"M",\
"N"\
\]

address\_verification\_result

string

Whether the end user's address input matches with the issuer's data.

Valid values\[\
"M",\
"N"\
\]

retrieval\_reference\_number

string

Receipt reference number communicated to the end user by their card issuer for this specific payment. This a commonly used reference number for the end users to raise tickets.

network\_response\_code

string

The response code returned by the scheme (Visa, Mastercard, JCB, China Unionpay or Amex).

network\_response\_code\_descriptor

string

Description of the response code.

network\_transaction\_id

string

Transaction ID received from the card scheme. Only available for merchants on switcher model.

acquirer\_merchant\_id

string

Acquirer's record of the MID that was used to process this transaction. Only available for merchants on switcher model.

reconciliation\_id

string

Acquirer's transaction record of the payment on their settlement statement. Only available for merchants on switcher model.

authentication\_data

object (Payments\_API\_AuthenticationData)

Specific to cards transaction only. Details about the card authentication.

flow

string

Indicates the flow that was used for the 3DS authentication.

Valid values\[\
"FULL\_AUTH",\
"FRICTIONLESS"\
\]

a\_res

object

Details about the card authentication response from the 3DS server.

eci

string

Payment system-specific value provided by the ACS or DS to indicate the results of the attempt to authenticate the Cardholder.

message\_version

string

The 3DS protocol version which has been used to perform 3DS.

authentication\_value

string

The result value from the 3DS transaction received from the ACS. This value is no longer present on responses after 45 days have passed after the authentication. Note that Mastercard and Visa use a different underlying format.

ds\_trans\_id

string

Universally unique transaction identifier assigned by the DS to identify a single transaction.

issuer\_name

string

Name of the payment method provider used by the end user.

payer\_account\_number

string

Account number of the end user making the payment from the payment method provider's records.

payer\_name

string

Name of the end user making the payment from the payment method provider's records.

receipt\_id

string

Receipt reference number communicated to the end user by their payment method provider for this specific payment. This a commonly used reference number for the end users to raise tickets.

remark

string

Remarks about this specific payment from the payment method provider's records.

network

string

Payment network which the payment was processed over.

fund\_source

string

Information about what was used by the end user to complete the payment. e.g. balance, installment, credit.

failure\_code

string

Failure codes for payments.

Valid values\[\
"ACCOUNT\_ACCESS\_BLOCKED",\
"INVALID\_MERCHANT\_SETTINGS",\
"INVALID\_ACCOUNT\_DETAILS",\
"PAYMENT\_ATTEMPT\_COUNTS\_EXCEEDED",\
"USER\_DEVICE\_UNREACHABLE",\
"CHANNEL\_UNAVAILABLE",\
"INSUFFICIENT\_BALANCE",\
"ACCOUNT\_NOT\_ACTIVATED",\
"INVALID\_TOKEN",\
"SERVER\_ERROR",\
"PARTNER\_TIMEOUT\_ERROR",\
"TIMEOUT\_ERROR",\
"USER\_DECLINED\_PAYMENT",\
"USER\_DID\_NOT\_AUTHORIZE",\
"PAYMENT\_REQUEST\_EXPIRED",\
"FAILURE\_DETAILS\_UNAVAILABLE",\
"EXPIRED\_OTP",\
"INVALID\_OTP",\
"PAYMENT\_AMOUNT\_LIMITS\_EXCEEDED",\
"OTP\_ATTEMPT\_COUNTS\_EXCEEDED",\
"CARD\_DECLINED",\
"DECLINED\_BY\_ISSUER",\
"ISSUER\_UNAVAILABLE",\
"INVALID\_CVV",\
"DECLINED\_BY\_PROCESSOR",\
"CAPTURE\_AMOUNT\_EXCEEDED ",\
"AUTHENTICATION\_FAILED",\
"PROCESSOR\_ERROR",\
"EXPIRED\_CARD",\
"STOLEN\_CARD",\
"INACTIVE\_OR\_UNAUTHORIZED\_CARD",\
"INVALID\_MERCHANT\_CREDENTIALS",\
"SUSPECTED\_FRAUDULENT"\
\]

ExampleCARD\_DECLINED

metadata

object (Payments\_API\_MerchantMetadata)

Key-value entries for your custom data.
You can specify up to 50 keys, with key names up to 40 characters and values up to 500 characters.
This is for your convenience. Xendit will not use this data for any processing.

Example{
"my\_custom\_id": "merchant-123",
"my\_custom\_order\_id": "order-123"
}

created

string (date-time)

ISO 8601 date-time format.

Example2021-12-31T23:59:59Z

updated

string (date-time)

ISO 8601 date-time format.

Example2021-12-31T23:59:59Z

Responses

200

OK

Was this article helpful?

Yes  No

Previous article

Create a payment request

Next article

Get the status of a payment request