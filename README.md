# Simflux ⚡

> A lightweight, local payment provider simulator and webhook dispatcher for developers.

Test and develop payment flows locally without depending on live sandboxes, test cards, or rate limits.

---

## Quick Start

Run instantly without installing:

```bash
npx simflux-cli serve
```

Or install globally:

```bash
npm install -g simflux-cli
```

---

## Defining Your Provider

Create a `simflux.yaml` in your project root:

```yaml
name: my-payment-provider
version: "1.0"
provider: stripe

payment:
  fields:
    transaction_id:
      type: string
      required: true
    amount:
      type: integer
      required: true

webhook:
  url: "http://localhost:3000/api/webhooks"
  signature:
    key: "Stripe-Signature"
    value: "whsec_test_secret_123"
  events:
    - payment.success
    - payment.failed
```

---

## CLI Commands

### 1. Load Definition

Load and validate your provider definition into Simflux:

```bash
simflux load

# Optionally reset all existing payments and state:
simflux load --reset
```

### 2. Start HTTP Server

Start the simulation server:

```bash
# Default port 8080
simflux serve

# Custom port
simflux serve --port 3000
```

### 3. Simulate Webhooks

Simulate different webhook events during payment processing:

```bash
# Set active event override
simflux webhook set --event payment.failed

# Inspect active event
simflux webhook get

# Reset override (disables webhook dispatch until set)
simflux webhook reset
```

### 4. Inspect Payments

List and inspect simulated payments:

```bash
# List all payments
simflux payments list

# Filter by reference
simflux payments list --reference ref_12345

# View details of a specific payment
simflux payments show 1

# Clear all payments
simflux payments clear --force
```

---

## HTTP API Reference

When `simflux serve` is running, your application interacts with standard REST endpoints:

| Method   | Endpoint                                | Description                                              |
| -------- | --------------------------------------- | -------------------------------------------------------- |
| `POST`   | `/payments`                             | Create payment (dispatches active webhook in background) |
| `GET`    | `/payments`                             | List payments (`?status=...`)                            |
| `GET`    | `/payments/:id`                         | Get payment by ID                                        |
| `GET`    | `/payments/:id/status`                  | Get payment status by ID                                 |
| `GET`    | `/payments/reference/:reference`        | Get payment by reference                                 |
| `GET`    | `/payments/reference/:reference/status` | Get payment status by reference                          |
| `PATCH`  | `/payments/reference/:reference`        | Update payment status                                    |
| `DELETE` | `/payments/:id`                         | Delete payment record                                    |

---

## License

MIT © [Faithful001](https://github.com/Faithful001/simflux)
