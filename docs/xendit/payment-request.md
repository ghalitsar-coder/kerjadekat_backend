# Create a payment request
 Create a payment request      

> Documentation Index
> -------------------
> 
> Fetch the complete documentation index at: [https://docs.xendit.co/llms.txt](https://docs.xendit.co/llms.txt)
> 
> Use this file to discover all available pages before exploring further.

What information are you looking for?

For archived content, access the previous documentation [here](https://archive.docs.xendit.co) or the previous API reference [here](https://archive.developers.xendit.co/).

Create a payment request
========================

*   Updated on Apr 27, 2026
*   Published on Jan 9, 2025

[Prev](https://docs.xendit.co/apidocs/introduction "Introduction") [Next](https://docs.xendit.co/apidocs/payment-webhook-notification "Payment webhook notification")

Post

/v3/payment\_requests

Create payment request. Initiates payment collection from end user.

Security

HTTP

Type basic

Header parameters

api-version

string

Valid values\[ "2024-11-11" \]

for-user-id

string

The XenPlatform subaccount user id that will perform this transaction.

with-split-rule

string

The XenPlatform split rule id that will be applied to this transaction.

Body parameters

application/json

PAY\_Cards\_3DS\_Auth PAY\_Cards\_No\_3DS PAY\_Cards\_Manual\_Capture PAY\_AND\_SAVE\_Cards PAY\_WithToken\_Cards PAY\_WithToken\_TOUCHNGO PAY\_Ewallet\_Shopeepay PAY\_QR\_PromptPay PAY\_DirectDebit\_BPI REUSABLE\_PAYMENT\_CODE\_VA\_VietCapital

PAY\_Cards\_3DS\_Auth

```
{
  "reference_id": "order_123456_3ds",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 100000,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "cvn": "123",
      "card_number": "4000000000001091",
      "expiry_year": "2025",
      "expiry_month": "12",
      "cardholder_first_name": "John",
      "cardholder_last_name": "Doe",
      "cardholder_email": "john.doe@example.com",
      "cardholder_phone_number": "+628123456789"
    },
    "skip_three_ds": false,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Payment for Order #123456",
  "metadata": {
    "order_id": "123456",
    "customer_type": "premium"
  }
}
```


PAY\_Cards\_No\_3DS

```
{
  "reference_id": "order_123457_no3ds",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 50000,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "cvn": "456",
      "card_number": "5200000000001096",
      "expiry_year": "2026",
      "expiry_month": "06",
      "cardholder_first_name": "Jane",
      "cardholder_last_name": "Doe",
      "cardholder_email": "jane.doe@example.com",
      "cardholder_phone_number": "+6312345678901"
    },
    "skip_three_ds": true,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Quick checkout for Order #123457",
  "metadata": {
    "order_id": "123457",
    "checkout_type": "express"
  }
}
```


PAY\_Cards\_Manual\_Capture

```
{
  "reference_id": "booking_123458",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 250000,
  "capture_method": "MANUAL",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "cvn": "789",
      "card_number": "4000000000000002",
      "expiry_year": "2027",
      "expiry_month": "03",
      "cardholder_first_name": "John",
      "cardholder_last_name": "Doe",
      "cardholder_email": "john.doe@example.com",
      "cardholder_phone_number": "+6212345678902"
    },
    "skip_three_ds": false,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Hotel booking pre-authorization #123458",
  "metadata": {
    "booking_id": "123458",
    "booking_type": "hotel"
  }
}
```


PAY\_AND\_SAVE\_Cards

```
{
  "reference_id": "order_123459_save",
  "customer": {
    "reference_id": "customer_789",
    "type": "INDIVIDUAL",
    "individual_detail": {
      "given_names": "John",
      "surname": "Doe"
    },
    "email": "john.doe@example.com",
    "mobile_number": "+6212345678901"
  },
  "type": "PAY_AND_SAVE",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 500000,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "cvn": "654",
      "card_number": "4000000000001109",
      "expiry_year": "2026",
      "expiry_month": "11",
      "cardholder_first_name": "John",
      "cardholder_last_name": "Doe",
      "cardholder_email": "john.doe@example.com",
      "cardholder_phone_number": "+6212345678901"
    },
    "skip_three_ds": false,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Initial subscription payment - Premium Plan",
  "metadata": {
    "order_id": "123459",
    "subscription_plan": "premium_monthly"
  }
}
```


PAY\_WithToken\_Cards

```
{
  "reference_id": "recurring_payment_123460",
  "payment_token_id": "pt-90392f42-d98a-49ef-a7f3-abcezas123",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 500000,
  "capture_method": "AUTOMATIC",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "skip_three_ds": true,
    "card_on_file_type": "RECURRING",
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Monthly subscription - Month 2",
  "metadata": {
    "subscription_id": "sub_789",
    "billing_cycle": "2"
  }
}
```


PAY\_WithToken\_TOUCHNGO

```
{
  "reference_id": "recurring_touchngo_123461",
  "payment_token_id": "pt-tng-90392f42-d98a-49ef-a7f3",
  "type": "PAY",
  "country": "MY",
  "currency": "MYR",
  "request_amount": 150,
  "capture_method": "AUTOMATIC",
  "channel_properties": {
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Monthly toll payment - Touch n Go",
  "metadata": {
    "account_id": "tng_123461",
    "payment_type": "toll_recurring"
  }
}
```


PAY\_Ewallet\_Shopeepay

```
{
  "reference_id": "order_ph_987654",
  "type": "PAY",
  "country": "PH",
  "currency": "PHP",
  "request_amount": 1250,
  "capture_method": "AUTOMATIC",
  "channel_code": "SHOPEEPAY",
  "channel_properties": {
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Online purchase - Order #987654",
  "metadata": {
    "order_id": "987654",
    "store_location": "manila_mall"
  }
}
```


PAY\_QR\_PromptPay

```
{
  "reference_id": "qr_th_321654",
  "type": "PAY",
  "country": "TH",
  "currency": "THB",
  "request_amount": 1500,
  "capture_method": "AUTOMATIC",
  "channel_code": "QRPROMPTPAY",
  "channel_properties": {
    "expires_at": "2025-01-15T23:59:59Z",
    "qr_string_type": "DYNAMIC"
  },
  "description": "Restaurant bill - Table 15",
  "metadata": {
    "order_id": "321654",
    "table_number": "15",
    "branch": "bangkok_central",
    "service_type": "dine_in"
  }
}
```


PAY\_DirectDebit\_BPI

```
{
  "reference_id": "dd_ph_456789_with_cust",
  "type": "PAY",
  "country": "PH",
  "currency": "PHP",
  "request_amount": 2500,
  "capture_method": "AUTOMATIC",
  "customer": {
    "type": "INDIVIDUAL",
    "reference_id": "customer_bpi_456789",
    "individual_detail": {
      "given_names": "Juan",
      "surname": "Dela Cruz"
    },
    "email": "juan.delacruz@example.com",
    "mobile_number": "+639123456789"
  },
  "channel_code": "BPI_DIRECT_DEBIT",
  "channel_properties": {
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "description": "Bill payment - Account #456789",
  "metadata": {
    "bill_id": "456789",
    "customer_account": "ACC456789"
  }
}
```


REUSABLE\_PAYMENT\_CODE\_VA\_VietCapital

```
{
  "reference_id": "rpc_vn_123456",
  "type": "REUSABLE_PAYMENT_CODE",
  "country": "VN",
  "currency": "VND",
  "request_amount": 10000,
  "channel_code": "VIETCAPITAL_VIRTUAL_ACCOUNT",
  "channel_properties": {
    "display_name": "xenCommerce Shop"
  }
}
```


Expand All

OneOf

Payments\_API\_Pay

object (Payments\_API\_Pay)

reference\_id

string Required

A reference ID from merchants to identify their request. For "CARDS" channel code, reference ID must be unique.

Min length1

Max length255

type

string Required

PAY: Create a payment request that is able to receive one payment.

Valid values\[ "PAY" \]

country

string Required

ISO 3166-1 alpha-2 two-letter country code for the country of transaction.

Valid values\[ "ID", "PH", "VN", "TH", "SG", "MY", "HK", "MX" \]

ExampleID

currency

string Required

ISO 4217 three-letter currency code for the payment.

Valid values\[ "IDR", "PHP", "VND", "THB", "SGD", "MYR", "USD", "HKD", "AUD", "GBP", "EUR", "JPY", "MXN" \]

ExampleIDR

channel\_code

string Required

Channel code used to select the payment method provider.

channel\_properties

object (Payments\_API\_ChannelPropertiesWidgetPay) Required

Data required to initiate transaction with payment method provider. Refer to the Channel Data Finder widget in the channel\_code field above for the full list of required properties for each channel.

request\_amount

number Required

The intended payment amount to be collected from the end user.

Minimum0.0

Example10000.0

capture\_method

string

AUTOMATIC: payment capture will be processed immediately after payment request is created. MANUAL: payment capture requires merchant's trigger via payment capture endpoint before being processed

Valid values\[ "AUTOMATIC", "MANUAL" \]

Default"AUTOMATIC"

ExampleAUTOMATIC

description

string

A custom description for the Payment Request.

Min length1

Max length1000

ExamplePayment for your order #123

customer\_id

string

Xendit unique Capture ID generated as reference for the end user

Max length41

Examplecust-b98d6f63-d240-44ec-9bd5-aa42954c4f48

customer

object (Payments\_API\_XenditStandardCustomer)

type

string Required

Type of customer

Valid values\[ "INDIVIDUAL" \]

reference\_id

string Required

Merchant provided identifier for the customer. Must be unique. Alphanumeric no special characters allowed

Min length1

Max length255

email

string (email)

E-mail address of customer. Maximum length 50 characters

Min length4

Max length50

mobile\_number

string

Mobile number of customer in E.164 format +(country code)(subscriber number)

Min length1

Max length50

individual\_detail

object (Payments\_API\_XenditStandardIndividualDetail) Required

given\_names

string Required

Primary or first name/s of customer. Alphanumeric. No special characters is allowed.

Min length1

Max length50

surname

string

Last or family name of customer. Alphanumeric. No special characters is allowed.

Min length1

Max length50

nationality

string

Country code for customer nationality. ISO 3166-1 alpha-2 Country Code

Min length2

Max length2

place\_of\_birth

string

City or other relevant location for the customer birth place. Alphanumeric. No special characters is allowed.

Min length1

Max length60

date\_of\_birth

string

Date of birth of the customer. Format: YYYY-MM-DD

Min length10

Max length10

gender

Gender of customer

Valid values\[ "MALE", "FEMALE", "OTHER" \]

items

Array of object (Payments\_API\_XenditStandardItem)

Array of objects describing the item/s attached to the payment.

object

reference\_id

string Required

Merchant provided identifier for the item

Min length1

Max length255

type

Type of item

Valid values\[ "DIGITAL\_PRODUCT", "PHYSICAL\_PRODUCT", "DIGITAL\_SERVICE", "PHYSICAL\_SERVICE", "FEE" \]

name

string Required

Name of item

Min length1

Max length255

net\_unit\_amount

number Required

Net amount to be charged per unit

quantity

integer Required

Number of units of this item in the basket

Minimum1.0

url

string

URL of the item. Must be HTTPS or HTTP

image\_url

string

URL of the image of the item. Must be HTTPS or HTTP

category

string Required

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

metadata

object (Payments\_API\_MerchantMetadata)

Key-value entries for your custom data. You can specify up to 50 keys, with key names up to 40 characters and values up to 500 characters. This is for your convenience. Xendit will not use this data for any processing.

Example{ "my\_custom\_id": "merchant-123", "my\_custom\_order\_id": "order-123" }

Payments\_API\_PayAndSave

object (Payments\_API\_PayAndSave)

reference\_id

string Required

A reference ID from merchants to identify their request. For "CARDS" channel code, reference ID must be unique.

Min length1

Max length255

type

string Required

PAY\_AND\_SAVE: Create a payment request that is able to receive one payment. If the payment is successful, a reusable payment token will be returned in the callback as saved payment information for subsequent payment requests.

Valid values\[ "PAY\_AND\_SAVE" \]

country

string Required

ISO 3166-1 alpha-2 two-letter country code for the country of transaction.

Valid values\[ "ID", "PH", "VN", "TH", "SG", "MY", "HK", "MX" \]

ExampleID

currency

string Required

ISO 4217 three-letter currency code for the payment.

Valid values\[ "IDR", "PHP", "VND", "THB", "SGD", "MYR", "USD", "HKD", "AUD", "GBP", "EUR", "JPY", "MXN" \]

ExampleIDR

channel\_code

string Required

Channel code used to select the payment method provider.

channel\_properties

object (Payments\_API\_ChannelPropertiesWidgetPayAndSave) Required

Data required to initiate transaction with payment method provider. Refer to the Channel Data Finder widget in the channel\_code field above for the full list of required properties for each channel.

request\_amount

number Required

The intended payment amount to be collected from the end user.

Minimum0.0

Example10000.0

capture\_method

string

AUTOMATIC: payment capture will be processed immediately after payment request is created. MANUAL: payment capture requires merchant's trigger via payment capture endpoint before being processed

Valid values\[ "AUTOMATIC", "MANUAL" \]

Default"AUTOMATIC"

ExampleAUTOMATIC

customer\_id

string

Xendit unique Capture ID generated as reference for the end user

Max length41

Examplecust-b98d6f63-d240-44ec-9bd5-aa42954c4f48

customer

object (Payments\_API\_XenditStandardCustomer)

type

string Required

Type of customer

Valid values\[ "INDIVIDUAL" \]

reference\_id

string Required

Merchant provided identifier for the customer. Must be unique. Alphanumeric no special characters allowed

Min length1

Max length255

email

string (email)

E-mail address of customer. Maximum length 50 characters

Min length4

Max length50

mobile\_number

string

Mobile number of customer in E.164 format +(country code)(subscriber number)

Min length1

Max length50

individual\_detail

object (Payments\_API\_XenditStandardIndividualDetail) Required

given\_names

string Required

Primary or first name/s of customer. Alphanumeric. No special characters is allowed.

Min length1

Max length50

surname

string

Last or family name of customer. Alphanumeric. No special characters is allowed.

Min length1

Max length50

nationality

string

Country code for customer nationality. ISO 3166-1 alpha-2 Country Code

Min length2

Max length2

place\_of\_birth

string

City or other relevant location for the customer birth place. Alphanumeric. No special characters is allowed.

Min length1

Max length60

date\_of\_birth

string

Date of birth of the customer. Format: YYYY-MM-DD

Min length10

Max length10

gender

Gender of customer

Valid values\[ "MALE", "FEMALE", "OTHER" \]

description

string

A custom description for the Payment Request.

Min length1

Max length1000

ExamplePayment for your order #123

items

Array of object (Payments\_API\_XenditStandardItem)

Array of objects describing the item/s attached to the payment.

object

reference\_id

string Required

Merchant provided identifier for the item

Min length1

Max length255

type

Type of item

Valid values\[ "DIGITAL\_PRODUCT", "PHYSICAL\_PRODUCT", "DIGITAL\_SERVICE", "PHYSICAL\_SERVICE", "FEE" \]

name

string Required

Name of item

Min length1

Max length255

net\_unit\_amount

number Required

Net amount to be charged per unit

quantity

integer Required

Number of units of this item in the basket

Minimum1.0

url

string

URL of the item. Must be HTTPS or HTTP

image\_url

string

URL of the image of the item. Must be HTTPS or HTTP

category

string Required

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

metadata

object (Payments\_API\_MerchantMetadata)

Key-value entries for your custom data. You can specify up to 50 keys, with key names up to 40 characters and values up to 500 characters. This is for your convenience. Xendit will not use this data for any processing.

Example{ "my\_custom\_id": "merchant-123", "my\_custom\_order\_id": "order-123" }

Payments\_API\_PayWithToken

object (Payments\_API\_PayWithToken)

reference\_id

string Required

A reference ID from merchants to identify their request. For "CARDS" channel code, reference ID must be unique.

Min length1

Max length255

type

string Required

PAY: Create a payment request that is able to receive one payment.

Valid values\[ "PAY" \]

payment\_token\_id

string Required

Xendit unique Payment Token ID generated as reference for reusable payment details of the end user.

Examplept-cc3938dc-c2a5-43c4-89d7-7570793348c2

country

string Required

ISO 3166-1 alpha-2 two-letter country code for the country of transaction.

Valid values\[ "ID", "PH", "VN", "TH", "SG", "MY", "HK", "MX" \]

ExampleID

currency

string Required

ISO 4217 three-letter currency code for the payment.

Valid values\[ "IDR", "PHP", "VND", "THB", "SGD", "MYR", "USD", "HKD", "AUD", "GBP", "EUR", "JPY", "MXN" \]

ExampleIDR

channel\_code

string Required

Channel code used to select the payment method provider.

channel\_properties

object (Payments\_API\_ChannelPropertiesWidgetPayWithToken)

Data required to initiate transaction with payment method provider. Refer to the Channel Data Finder widget in the channel\_code field above for the full list of required properties for each channel.

request\_amount

number Required

The intended payment amount to be collected from the end user.

Minimum0.0

Example10000.0

capture\_method

string

AUTOMATIC: payment capture will be processed immediately after payment request is created. MANUAL: payment capture requires merchant's trigger via payment capture endpoint before being processed

Valid values\[ "AUTOMATIC", "MANUAL" \]

Default"AUTOMATIC"

ExampleAUTOMATIC

description

string

A custom description for the Payment Request.

Min length1

Max length1000

ExamplePayment for your order #123

customer\_id

string

Xendit unique Capture ID generated as reference for the end user

Max length41

Examplecust-b98d6f63-d240-44ec-9bd5-aa42954c4f48

customer

object (Payments\_API\_XenditStandardCustomer)

type

string Required

Type of customer

Valid values\[ "INDIVIDUAL" \]

reference\_id

string Required

Merchant provided identifier for the customer. Must be unique. Alphanumeric no special characters allowed

Min length1

Max length255

email

string (email)

E-mail address of customer. Maximum length 50 characters

Min length4

Max length50

mobile\_number

string

Mobile number of customer in E.164 format +(country code)(subscriber number)

Min length1

Max length50

individual\_detail

object (Payments\_API\_XenditStandardIndividualDetail) Required

given\_names

string Required

Primary or first name/s of customer. Alphanumeric. No special characters is allowed.

Min length1

Max length50

surname

string

Last or family name of customer. Alphanumeric. No special characters is allowed.

Min length1

Max length50

nationality

string

Country code for customer nationality. ISO 3166-1 alpha-2 Country Code

Min length2

Max length2

place\_of\_birth

string

City or other relevant location for the customer birth place. Alphanumeric. No special characters is allowed.

Min length1

Max length60

date\_of\_birth

string

Date of birth of the customer. Format: YYYY-MM-DD

Min length10

Max length10

gender

Gender of customer

Valid values\[ "MALE", "FEMALE", "OTHER" \]

items

Array of object (Payments\_API\_XenditStandardItem)

Array of objects describing the item/s attached to the payment.

object

reference\_id

string Required

Merchant provided identifier for the item

Min length1

Max length255

type

Type of item

Valid values\[ "DIGITAL\_PRODUCT", "PHYSICAL\_PRODUCT", "DIGITAL\_SERVICE", "PHYSICAL\_SERVICE", "FEE" \]

name

string Required

Name of item

Min length1

Max length255

net\_unit\_amount

number Required

Net amount to be charged per unit

quantity

integer Required

Number of units of this item in the basket

Minimum1.0

url

string

URL of the item. Must be HTTPS or HTTP

image\_url

string

URL of the image of the item. Must be HTTPS or HTTP

category

string Required

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

metadata

object (Payments\_API\_MerchantMetadata)

Key-value entries for your custom data. You can specify up to 50 keys, with key names up to 40 characters and values up to 500 characters. This is for your convenience. Xendit will not use this data for any processing.

Example{ "my\_custom\_id": "merchant-123", "my\_custom\_order\_id": "order-123" }

Payments\_API\_ReusablePaymentCode

object (Payments\_API\_ReusablePaymentCode)

reference\_id

string Required

A reference ID from merchants to identify their request. For "CARDS" channel code, reference ID must be unique.

Min length1

Max length255

type

string Required

REUSABLE\_PAYMENT\_CODE: Create one payment request that is able to receive multiple payments. This is only used for repeat use payment method like static QR, static OTC payment code or a predefined Virtual Account number.

Valid values\[ "REUSABLE\_PAYMENT\_CODE" \]

country

string Required

ISO 3166-1 alpha-2 two-letter country code for the country of transaction.

Valid values\[ "ID", "PH", "VN", "TH", "SG", "MY", "HK", "MX" \]

ExampleID

currency

string Required

ISO 4217 three-letter currency code for the payment.

Valid values\[ "IDR", "PHP", "VND", "THB", "SGD", "MYR", "USD", "HKD", "AUD", "GBP", "EUR", "JPY", "MXN" \]

ExampleIDR

channel\_code

string Required

Channel code used to select the payment method provider.

channel\_properties

object (Payments\_API\_ChannelPropertiesWidgetReusablePaymentCode) Required

Data required to initiate transaction with payment method provider. Refer to the Channel Data Finder widget in the channel\_code field above for the full list of required properties for each channel.

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

string Required

Merchant provided identifier for the item

Min length1

Max length255

type

Type of item

Valid values\[ "DIGITAL\_PRODUCT", "PHYSICAL\_PRODUCT", "DIGITAL\_SERVICE", "PHYSICAL\_SERVICE", "FEE" \]

name

string Required

Name of item

Min length1

Max length255

net\_unit\_amount

number Required

Net amount to be charged per unit

quantity

integer Required

Number of units of this item in the basket

Minimum1.0

url

string

URL of the item. Must be HTTPS or HTTP

image\_url

string

URL of the image of the item. Must be HTTPS or HTTP

category

string Required

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

Responses

201

Payment Request Created

application/json

PAY\_Cards\_3DS\_Auth PAY\_Cards\_No\_3DS PAY\_Cards\_Manual\_Capture PAY\_AND\_SAVE\_Cards PAY\_WithToken\_Cards PAY\_WithToken\_TOUCHNGO PAY\_Ewallet\_Shopeepay PAY\_QR\_PromptPay PAY\_DirectDebit\_BPI REUSABLE\_PAYMENT\_CODE\_VA\_VietCapital

PAY\_Cards\_3DS\_Auth

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "order_123456_3ds",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas123",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 100000,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "masked_card_number": "4000****1091",
      "cardholder_first_name": "John",
      "cardholder_last_name": "Doe",
      "cardholder_email": "john.doe@example.com",
      "cardholder_phone_number": "+628123456789",
      "expiry_month": "12",
      "expiry_year": "2025",
      "fingerprint": "62397498595752001b9fdeba",
      "type": "DEBIT",
      "network": "VISA",
      "country": "ID",
      "issuer": "PT BANK MANDIRI"
    },
    "skip_three_ds": false,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [
    {
      "type": "REDIRECT_CUSTOMER",
      "value": "https://redirect.partner.com/3ds-auth-xyz",
      "descriptor": "WEB_URL"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Payment for Order #123456",
  "metadata": {
    "order_id": "123456",
    "customer_type": "premium"
  },
  "created": "2025-01-15T12:30:45Z",
  "updated": "2025-01-15T12:30:45Z"
}
```


PAY\_Cards\_No\_3DS

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "order_123457_no3ds",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas124",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 50000,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "masked_card_number": "5200****1096",
      "cardholder_first_name": "Jane",
      "cardholder_last_name": "Doe",
      "cardholder_email": "jane.doe@example.com",
      "cardholder_phone_number": "+6312345678901",
      "expiry_month": "06",
      "expiry_year": "2026",
      "fingerprint": "62397498595752001b9fdebc",
      "type": "CREDIT",
      "network": "MASTERCARD",
      "country": "ID",
      "issuer": "PT BANK BCA"
    },
    "skip_three_ds": true,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [],
  "status": "SUCCEEDED",
  "description": "Quick checkout for Order #123457",
  "metadata": {
    "order_id": "123457",
    "checkout_type": "express"
  },
  "created": "2025-01-15T12:31:22Z",
  "updated": "2025-01-15T12:31:25Z"
}
```


PAY\_Cards\_Manual\_Capture

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "booking_123458",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas125",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 250000,
  "capture_method": "MANUAL",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "masked_card_number": "4000****0002",
      "cardholder_first_name": "John",
      "cardholder_last_name": "Doe",
      "cardholder_email": "john.doe@example.com",
      "cardholder_phone_number": "+6212345678902",
      "expiry_month": "03",
      "expiry_year": "2027",
      "fingerprint": "62397498595752001b9fdebd",
      "type": "CREDIT",
      "network": "VISA",
      "country": "ID",
      "issuer": "PT BANK BNI"
    },
    "skip_three_ds": false,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [],
  "status": "AUTHORIZED",
  "description": "Hotel booking pre-authorization #123458",
  "metadata": {
    "booking_id": "123458",
    "booking_type": "hotel"
  },
  "created": "2025-01-15T12:32:10Z",
  "updated": "2025-01-15T12:32:12Z"
}
```


PAY\_AND\_SAVE\_Cards

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "order_123459_save",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas126",
  "customer_id": "cust-90392f42-d98a-49ef-a7f3-abcezas789",
  "type": "PAY_AND_SAVE",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 500000,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "card_details": {
      "masked_card_number": "4000****1109",
      "cardholder_first_name": "John",
      "cardholder_last_name": "Doe",
      "cardholder_email": "john.doe@example.com",
      "cardholder_phone_number": "+6212345678901",
      "expiry_month": "11",
      "expiry_year": "2026",
      "fingerprint": "62397498595752001b9fdebe",
      "type": "CREDIT",
      "network": "VISA",
      "country": "ID",
      "issuer": "PT BANK MANDIRI"
    },
    "skip_three_ds": false,
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [
    {
      "type": "REDIRECT_CUSTOMER",
      "value": "https://redirect.partner.com/3ds-auth-abc",
      "descriptor": "WEB_URL"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Initial subscription payment - Premium Plan",
  "metadata": {
    "order_id": "123459",
    "subscription_plan": "premium_monthly"
  },
  "created": "2025-01-15T12:33:15Z",
  "updated": "2025-01-15T12:33:15Z"
}
```


PAY\_WithToken\_Cards

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "recurring_payment_123460",
  "payment_token_id": "pt-90392f42-d98a-49ef-a7f3-abcezas123",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas127",
  "type": "PAY",
  "country": "ID",
  "currency": "IDR",
  "request_amount": 500000,
  "capture_method": "AUTOMATIC",
  "channel_code": "CARDS",
  "channel_properties": {
    "mid_label": "CTV_TEST",
    "skip_three_ds": true,
    "card_on_file_type": "RECURRING",
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [],
  "status": "SUCCEEDED",
  "description": "Monthly subscription - Month 2",
  "metadata": {
    "subscription_id": "sub_789",
    "billing_cycle": "2"
  },
  "created": "2025-01-15T12:34:20Z",
  "updated": "2025-01-15T12:34:22Z"
}
```


PAY\_WithToken\_TOUCHNGO

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "recurring_touchngo_123461",
  "payment_token_id": "pt-tng-90392f42-d98a-49ef-a7f3",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas128",
  "type": "PAY",
  "country": "MY",
  "currency": "MYR",
  "request_amount": 150,
  "capture_method": "AUTOMATIC",
  "channel_code": "TOUCHNGO",
  "channel_properties": {
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [
    {
      "type": "REDIRECT_CUSTOMER",
      "value": "https://redirect.partner.com/touchngo-auth-xyz",
      "descriptor": "WEB_URL"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Monthly toll payment - Touch n Go",
  "metadata": {
    "account_id": "tng_123461",
    "payment_type": "toll_recurring"
  },
  "created": "2025-01-15T12:35:30Z",
  "updated": "2025-01-15T12:35:30Z"
}
```


PAY\_Ewallet\_Shopeepay

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "order_ph_987654",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas129",
  "type": "PAY",
  "country": "PH",
  "currency": "PHP",
  "request_amount": 1250,
  "capture_method": "AUTOMATIC",
  "channel_code": "SHOPEEPAY",
  "channel_properties": {
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [
    {
      "type": "REDIRECT_CUSTOMER",
      "value": "https://redirect.partner.com/shopeepay-auth-xyz",
      "descriptor": "WEB_URL"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Online purchase - Order #987654",
  "metadata": {
    "order_id": "987654",
    "store_location": "manila_mall"
  },
  "created": "2025-01-15T12:36:45Z",
  "updated": "2025-01-15T12:36:45Z"
}
```


PAY\_QR\_PromptPay

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "qr_th_321654",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas130",
  "type": "PAY",
  "country": "TH",
  "currency": "THB",
  "request_amount": 1500,
  "capture_method": "AUTOMATIC",
  "channel_code": "QRPROMPTPAY",
  "channel_properties": {
    "expires_at": "2025-01-15T23:59:59Z",
    "qr_string_type": "DYNAMIC"
  },
  "actions": [
    {
      "type": "PRESENT_TO_CUSTOMER",
      "descriptor": "QR_STRING",
      "value": "00020101021153037643011A0000000677010111013202120000000000000000052048405303764540415005802TH5925THAI QR PAYMENT (EXAMPLE)6007Bangkok62070703001"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Restaurant bill - Table 15",
  "metadata": {
    "order_id": "321654",
    "table_number": "15",
    "branch": "bangkok_central",
    "service_type": "dine_in"
  },
  "created": "2025-01-15T12:37:50Z",
  "updated": "2025-01-15T12:37:50Z"
}
```


PAY\_DirectDebit\_BPI

```
{
  "business_id": "5f27a14a9bf05c73dd040bc8",
  "reference_id": "dd_ph_456789_with_cust",
  "payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas132",
  "customer_id": "cust-e60078dd-494c-483c-a1e2-4f7bc34ff728",
  "type": "PAY",
  "country": "PH",
  "currency": "PHP",
  "request_amount": 2500,
  "capture_method": "AUTOMATIC",
  "channel_code": "BPI_DIRECT_DEBIT",
  "channel_properties": {
    "failure_return_url": "https://xendit.co/failure",
    "success_return_url": "https://xendit.co/success"
  },
  "actions": [
    {
      "type": "API_POST_REQUEST",
      "descriptor": "DEBIT_ACCOUNT",
      "value": "https://api.xendit.co/v1/direct_debits/authorize"
    }
  ],
  "status": "REQUIRES_ACTION",
  "description": "Bill payment - Account #456789",
  "metadata": {
    "bill_id": "456789",
    "customer_account": "ACC456789"
  },
  "created": "2025-01-15T12:39:10Z",
  "updated": "2025-01-15T12:39:10Z"
}
```


REUSABLE\_PAYMENT\_CODE\_VA\_VietCapital

```
{
  "created": "2025-10-21T09:14:48.966Z",
  "updated": "2025-10-21T09:14:48.966Z",
  "channel_properties": {
    "display_name": "xenCommerce Shop",
    "expires_at": "2025-10-23T09:14:48.966Z"
  },
  "business_id": "6577c85379425b82e415c673",
  "reference_id": "rpc_vn_123456",
  "payment_request_id": "pr-73762759-0c9b-8178-b48f-0480d694b643",
  "type": "REUSABLE_PAYMENT_CODE",
  "country": "VN",
  "currency": "VND",
  "request_amount": 10000,
  "capture_method": "AUTOMATIC",
  "channel_code": "VIETCAPITAL_VIRTUAL_ACCOUNT",
  "actions": [
    {
      "type": "PRESENT_TO_CUSTOMER",
      "descriptor": "VIRTUAL_ACCOUNT_NUMBER",
      "value": "8881761038089006"
    },
    {
      "type": "PRESENT_TO_CUSTOMER",
      "descriptor": "QR_STRING",
      "value": "00020101021238600010A00000072701300006970454011688817610380890060208QRIBFTTA53037045802VN5405100006304C696"
    }
  ],
  "status": "REQUIRES_ACTION"
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

Payments\_API\_Http400InvalidValueError

object (Payments\_API\_Http400InvalidValueError)

error\_code

string

Valid values\[ "INVALID\_VALUE\_ERROR" \]

message

string

Values in the payment request is not within expected range or expected configurations. Check the specific error message for debugging.

Payments\_API\_Http400ApiValidationError

object (Payments\_API\_Http400ApiValidationError)

error\_code

string

Valid values\[ "API\_VALIDATION\_ERROR" \]

message

string

Fields or values in the payment request does not comply with our API specification. Check the specific error message for debugging.

Payments\_API\_Http400CardExpired

object (Payments\_API\_Http400CardExpired)

error\_code

string

Valid values\[ "CARD\_EXPIRED" \]

message

string

Card expiry specified in the request should not be earlier than current month.

Payments\_API\_Http400InvalidPaymentDetails

object (Payments\_API\_Http400InvalidPaymentDetails)

error\_code

string

Valid values\[ "INVALID\_PAYMENT\_DETAILS" \]

message

string

The payment details entered by the end user is invalid. Check the specific error message for debugging.

Payments\_API\_Http400PaymentRequestRateLimited

object (Payments\_API\_Http400PaymentRequestRateLimited)

error\_code

string

Valid values\[ "PAYMENT\_REQUEST\_RATE\_LIMITED" \]

message

string

Maximum number of requests to this payment channel has been exceeded in a given time frame.

Payments\_API\_Http400InvalidToken

object (Payments\_API\_Http400InvalidToken)

error\_code

string

Valid values\[ "INVALID\_TOKEN" \]

message

string

Payment token ID specified in the payment request has expired or has been cancelled. Please reinitiate linking before retrying.

403

Forbidden

application/json

OneOf

Payments\_API\_Http403Skip3dsForbidden

object (Payments\_API\_Http403Skip3dsForbidden)

error\_code

string

Valid values\[ "SKIP\_3DS\_FORBIDDEN" \]

message

string

Non 3DS payment request for cards is not allowed. Please activate the feature on Xendit dashboard before proceeding.

Payments\_API\_Http403InvalidMerchantSettings

object (Payments\_API\_Http403InvalidMerchantSettings)

error\_code

string

Valid values\[ "INVALID\_MERCHANT\_SETTINGS" \]

message

string

Merchant credentials met with an error with the provider. Please contact Xendit customer support to resolve this issue.

Payments\_API\_Http403AccountAccessBlocked

object (Payments\_API\_Http403AccountAccessBlocked)

error\_code

string

Valid values\[ "ACCOUNT\_ACCESS\_BLOCKED" \]

message

string

Payment token ID specified in the request was denied access by the payment method provider.

409

Conflict

application/json

OneOf

Payments\_API\_Http409DuplicateError

object (Payments\_API\_Http409DuplicateError)

error\_code

string

Valid values\[ "DATA\_NOT\_FOUND" \]

message

string

Duplication is not allowed. Check specific error message for debugging.

Payments\_API\_Http409AccountAlreadyLinked

object (Payments\_API\_Http409AccountAlreadyLinked)

error\_code

string

Valid values\[ "ACCOUNT\_ALREADY\_LINKED" \]

message

string

The end user has already linked their account previously.

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

503

Service unavailable

application/json

OneOf

Payments\_API\_Http503ChannelUnavailable

object (Payments\_API\_Http503ChannelUnavailable)

error\_code

string

Valid values\[ "CHANNEL\_UNAVAILABLE" \]

message

string

The channel requested is currently experiencing unexpected issues. The provider will be notified to resolve this issue.

Payments\_API\_Http503IssuerUnavailable

object (Payments\_API\_Http503IssuerUnavailable)

error\_code

string

Valid values\[ "ISSUER\_UNAVAILABLE" \]

message

string

The end user's payment method provider is currently experiencing unexpected issues. The provider will be notified to resolve this issue.

Was this article helpful?

Yes No