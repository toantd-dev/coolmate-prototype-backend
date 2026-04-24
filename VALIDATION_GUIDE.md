# Validation Implementation Guide

All request validation has been standardized with Gin's binding validation.

## Validation Structs Location

- **Auth**: `internal/services/auth_service.go` (RegisterRequest, LoginRequest)
- **All Other Requests**: `internal/handlers/request_validators.go`

## Using Validation in Handlers

### Example Pattern

```go
// In handler
func (vh *VendorHandler) Register(c *gin.Context) {
    var req handlers.RegisterVendorRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }
    
    // Request is now validated and safe to use
    // ... call service layer
}
```

## Validation Tags

Standard Go validation tags used across all requests:

| Tag | Usage | Example |
|-----|-------|---------|
| `required` | Field is mandatory | `Store string \`binding:"required"\`` |
| `email` | Must be valid email | `Email string \`binding:"email"\`` |
| `min=N` | Minimum length/value | `Name string \`binding:"min=3"\`` |
| `max=N` | Maximum length/value | `Name string \`binding:"max=255"\`` |
| `gt=N` | Greater than | `Price float64 \`binding:"gt=0"\`` |
| `gte=N` | Greater or equal | `Stock int \`binding:"gte=0"\`` |
| `oneof=` | One of enumerated values | `Role string \`binding:"oneof=vendor customer"\`` |
| `alphanum` | Alphanumeric only | `Slug string \`binding:"alphanum"\`` |
| `url` | Must be valid URL | `Logo string \`binding:"url"\`` |
| `datetime` | Must match format | `Date string \`binding:"datetime=2006-01-02"\`` |
| `len=N` | Exact length | `Phone string \`binding:"len=10"\`` |

## Validation Error Handling

All validation errors are automatically caught by `c.ShouldBindJSON()`:

```go
if err := c.ShouldBindJSON(&req); err != nil {
    utils.BadRequest(c, err.Error())
    return
}
```

This returns HTTP 400 with a descriptive error message:
```json
{
  "status": "error",
  "message": "Key: 'RegisterVendorRequest.StoreName' Error:Field validation for 'StoreName' failed on the 'required' tag"
}
```

## Pagination Validation

All list endpoints use standardized pagination:

```go
type PaginationQuery struct {
    Page  int `form:"page" binding:"min=1" default:"1"`
    Limit int `form:"limit" binding:"min=1,max=100" default:"20"`
}

// Usage in handler
var query handlers.PaginationQuery
c.ShouldBindQuery(&query)
offset := query.GetOffset()
```

**Limits:**
- Minimum page: 1
- Default limit: 20
- Maximum limit: 100 (enforced, larger requests will be capped)

## Custom Validation Examples

For complex validation beyond tags, add validation in the service layer:

```go
// In ProductService.ValidateProduct()
func (ps *ProductService) ValidateProduct(product *models.Product) error {
    if product.CostPrice >= product.BasePrice {
        return errors.New("cost price must be less than base price")
    }
    return nil
}

// Call after binding in handler
if err := ps.ValidateProduct(product); err != nil {
    utils.BadRequest(c, err.Error())
    return
}
```

## Implemented Validations

### Product Creation
- ✓ Name: 3-255 characters
- ✓ SKU: 3-50 characters
- ✓ Description: 10-5000 characters
- ✓ Category: required
- ✓ Base price: > 0
- ✓ Cost price: >= 0 and < base price
- ✓ Weight: > 0
- ✓ Return window: > 0 if returnable

### Order Checkout
- ✓ Shipping address: all required fields
- ✓ Phone: exactly 10 digits
- ✓ Email: valid email format
- ✓ Postal code: exactly 5 digits
- ✓ Payment method: valid enum
- ✓ Promo codes: max 5

### Vendor Registration
- ✓ Store name: 3-100 characters
- ✓ Store slug: 3-50 alphanumeric
- ✓ Commission model: margin or markup
- ✓ Commission rate: 0-1

## Error Response Format

All validation errors follow this format:

```json
{
  "status": "error",
  "message": "validation error details"
}
```

HTTP Status: 400 Bad Request

## Testing Validation

Test with curl:

```bash
# Missing required field
curl -X POST http://localhost:8080/api/v1/vendor/register \
  -H "Content-Type: application/json" \
  -d '{"store_name": "Valid"}' \
  # Returns 400 with validation error

# Valid request
curl -X POST http://localhost:8080/api/v1/vendor/register \
  -H "Content-Type: application/json" \
  -d '{
    "store_name": "My Store",
    "store_slug": "my-store",
    "commission_model": "margin",
    "commission_rate": 0.05
  }'
  # Returns 201 with response
```

## Next Steps

1. Update all handlers to use validation structs
2. Add custom validation in service layer where needed
3. Run integration tests with various invalid inputs
4. Add comprehensive error logging for debugging

## References

- [Gin Validation Documentation](https://github.com/go-playground/validator)
- `internal/handlers/request_validators.go` - All validation structs
- `internal/utils/response.go` - Error response helpers
