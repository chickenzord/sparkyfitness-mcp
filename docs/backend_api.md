# backend api

## Authentication

Use `Authorization: Bearer <api_key>` header.

## Common

All request and response payload are in JSON format.

## Endpoints

### Search Foods

Example: `GET /foods?name=Rice&broadMatch=true&limit=10`

Query params:

- `name`: query string
- `broadMatch`: boolean
- `exactMatch`: boolean
- `limit`: integer

Notes:
- `broadMatch` and `exactMatch` are mutually exclusive
- `limit` is required

Example response:

```json
{
  "searchResults": [
    {
      "id": "330c0435-e6ab-471c-9eb9-6baf40b8499b",
      "name": "Steamed White Rice",
      "brand": null,
      "is_custom": true,
      "user_id": "01c9a380-3bfb-424f-87da-943a5e33ec51",
      "shared_with_public": false,
      "provider_external_id": null,
      "provider_type": null,
      "default_variant": {
        "id": "ed96d32a-b995-47fe-b1c8-0adacda62be3",
        "serving_size": 100,
        "serving_unit": "g",
        "calories": 130,
        "protein": 2.7,
        "carbs": 28.6,
        "fat": 0.3,
        "saturated_fat": 0.1,
        "polyunsaturated_fat": 0.1,
        "monounsaturated_fat": 0.1,
        "trans_fat": 0,
        "cholesterol": 0,
        "sodium": 1,
        "potassium": 26,
        "dietary_fiber": 0.4,
        "sugars": 0.1,
        "vitamin_a": 0,
        "vitamin_c": 0,
        "calcium": 10,
        "iron": 0.2,
        "is_default": true,
        "glycemic_index": null,
        "custom_nutrients": {}
      }
    },
    {
      "id": "4baf23e2-0c68-4b11-8101-67a091d935d1",
      "name": "Steamed brown rice",
      "brand": "",
      "is_custom": true,
      "user_id": "01c9a380-3bfb-424f-87da-943a5e33ec51",
      "shared_with_public": false,
      "provider_external_id": null,
      "provider_type": null,
      "default_variant": {
        "id": "da01b17b-94e7-44f1-a6de-22792048fe38",
        "serving_size": 100,
        "serving_unit": "g",
        "calories": 140,
        "protein": 3.5,
        "carbs": 29.8,
        "fat": 0.9,
        "saturated_fat": 0,
        "polyunsaturated_fat": 0,
        "monounsaturated_fat": 0,
        "trans_fat": 0,
        "cholesterol": 0,
        "sodium": 0,
        "potassium": 0,
        "dietary_fiber": 0,
        "sugars": 0,
        "vitamin_a": 0,
        "vitamin_c": 0,
        "calcium": 0,
        "iron": 0,
        "is_default": true,
        "glycemic_index": "None",
        "custom_nutrients": {}
      }
    }
  ]
}
```
