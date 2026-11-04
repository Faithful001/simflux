<p align="center">
  <h1 align="center">Simflux</h1>
  <p align="center">
    <strong>A local-first payment provider simulator and deterministic webhook dispatcher.</strong>
  </p>
  <p align="center">
    <a href="https://www.npmjs.com/package/simflux-cli"><img src="https://img.shields.io/npm/v/simflux-cli.svg?style=flat-square" alt="npm version"></a>
    <a href="https://github.com/Faithful001/simflux"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License: MIT"></a>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go" alt="Go Version"></a>
    <a href="https://github.com/Faithful001/simflux/releases"><img src="https://img.shields.io/github/v/release/Faithful001/simflux?style=flat-square" alt="Release"></a>
  </p>
</p>

---

## Table of Contents

- [Overview](#-overview)
- [Why Simflux?](#-why-simflux)
- [Installation](#-installation)
- [How It Works](#-how-it-works)
- [Configuration Guide (`simflux.yaml`)](#-configuration-guide-simfluxyaml)
- [CLI Reference](#-cli-reference)
  - [`simflux load`](#simflux-load)
  - [`simflux serve`](#simflux-serve)
  - [`simflux webhook`](#simflux-webhook)
  - [`simflux payments`](#simflux-payments)
  - [`simflux scenario`](#simflux-scenario)
- [HTTP API Reference](#-http-api-reference)
- [Automated Testing & Integration](#-automated-testing--integration)
- [License](#-license)

---

## Overview

**Simflux** is an open-source development tool designed to replace flaky live payment sandboxes (like Stripe, Paystack, Flutterwave, or PayPal) during local development and automated CI/CD testing.

Instead of registering fake credit cards, dealing with rate limits, or waiting on external sandbox webhooks, Simflux gives you a local HTTP payment gateway and an instant webhook dispatcher that you define using a simple YAML specification.

---

## 💡 Why Simflux?

| Feature                  | Live Sandboxes (Stripe, Paystack, etc.)             | Simflux                                                     |
| ------------------------ | --------------------------------------------------- | ----------------------------------------------------------- |
| **Speed**                | 1–5 seconds per request over the internet           | **Sub-millisecond** (runs locally on localhost)             |
| **Offline Testing**      | Requires active internet access                     | 100% offline                                                |
| **Failure Reproduction** | Complex (magic card numbers, specific cents values) | Single command: `simflux webhook set --event charge.failed` |
| **Custom Providers**     | Tied to one specific provider's API                 | Define **any** provider's schema with `simflux.yaml`        |
| **CI/CD Reliability**    | Flaky, rate-limited, requires live test secrets     | Zero credentials, predictable, resets instantly             |

---

## Installation

### Option 1: Via npm / npx (Zero Go Required)

Run instantly:

```bash
npx simflux-cli serve
```

Or install globally:

```bash
npm install -g simflux-cli
simflux --help
```

### Option 2: Via Go Toolchain

```bash
go install github.com/Faithful001/simflux.git/cmd/simflux@latest
```

### Option 3: From Source

```bash
git clone https://github.com/Faithful001/simflux.git
cd simflux
go build -o simflux.exe ./cmd/simflux
```

---

## How It Works

```
                        ┌──────────────────────────────────────┐
                        │            Your App / API            │
                        └───────────┬───────────────────▲──────┘
                                    │                   │
                  1. POST /payments │                   │ 3. Dispatches Webhook
                     (Creates Rec)  │                   │    (Background POST)
                                    ▼                   │
                     ┌──────────────────────────────────┴─────┐
                     │              Simflux Server            │
                     │  - Validates against simflux.yaml      │
                     │  - Stores record in SQLite DB          │
                     │  - Injects active simulated event      │
                     └────────────────────────────────────────┘
```

1. **You define your provider** in `simflux.yaml` (payment fields, types, webhook endpoints, and signatures).
2. **Load your configuration** with `simflux load`.
3. **Start the local server** with `simflux serve`.
4. Your application calls `http://localhost:8080/payments`.
5. Simflux validates the payload, stores the payment, and immediately dispatches the simulated webhook to your configured webhook URL with signatures attached.

---

## Configuration Guide (`simflux.yaml`)

Simflux reads `simflux.yaml` (or `simflux.yml`) from your working directory.

### Complete Example:

```yaml
name: stripe-sample
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
    currency:
      type: string
      required: false

webhook:
  url: "http://localhost:3000/api/webhooks/payments"
  signature:
    key: "Stripe-Signature"
    value: "whsec_test_secret_123"
  events:
    - payment.success
    - payment.failed
    - charge.refunded
```

### Field Types & Validation

Simflux enforces strict type checking against incoming payments:

| Type      | Accepted Values                                           |
| --------- | --------------------------------------------------------- |
| `string`  | Text string                                               |
| `integer` | Whole numbers (`100`, `5000`)                             |
| `number`  | Floating-point or integer numeric values (`19.99`, `100`) |
| `boolean` | `true` or `false`                                         |

If a required field is missing or the type is invalid, Simflux responds with `400 Bad Request` and details about the mismatch.

---

## CLI Reference

### `simflux load`

Parses, validates, and activates the provider definition into the SQLite database.

```bash
# Load default simflux.yaml:
simflux load

# Load from a specific file path:
simflux load -f ./configs/paystack.yaml

# Clear all existing payments and active webhook overrides when loading:
simflux load --reset
```

---

### `simflux serve`

Starts the HTTP server to receive simulation requests.

```bash
# Start on default port :8080
simflux serve

# Specify a custom port
simflux serve --port 3000
simflux serve -p 3000
```

_Supports graceful shutdown on `Ctrl+C` (`SIGINT`/`SIGTERM`) with a 10-second timeout._

---

### `simflux webhook`

Controls which webhook event will be fired when payments are processed.

```bash
# Set active event override (must be registered in your simflux.yaml)
simflux webhook set --event payment.failed
simflux webhook set -e payment.success

# Inspect currently active event override
simflux webhook get

# Reset override (disables webhook delivery until explicitly set)
simflux webhook reset
```

---

### `simflux payments`

Inspect and manage simulated payment records stored in the SQLite database.

```bash
# List all simulated payments
simflux payments list

# Filter payments by status
simflux payments list --status pending
simflux payments list -s success

# Lookup payment by reference
simflux payments list --reference ref_81726391
simflux payments list -r ref_81726391

# Inspect full payment JSON by ID
simflux payments show 1

# Purge all simulated payments
simflux payments clear --force
```

---

### `simflux scenario`

Manage simulated payment scenarios:

```bash
simflux scenario list
simflux scenario create --file ./scenario.json
simflux scenario show <id>
simflux scenario update <id> --file ./scenario.json
simflux scenario delete <id>
```

---

## HTTP API Reference

When `simflux serve` is running, your application interacts with the following endpoints:

### 1. Create Payment

- **`POST /payments`**
- Validates payment fields according to `simflux.yaml`, persists the record, and asynchronously delivers the active webhook event.

**Request Body:**

```json
{
  "reference": "ref_order_100234",
  "status": "pending",
  "fields": {
    "transaction_id": "tx_abc123",
    "amount": 2500
  }
}
```

**Response (`201 Created`):**

```json
{
  "message": "Payment created successfully"
}
```

---

### 2. Get Payment by Reference

- **`GET /payments/reference/:reference`**

**Response (`200 OK`):**

```json
{
  "id": 1,
  "reference": "ref_order_100234",
  "status": "pending",
  "created_at": "2026-11-04T20:00:00Z",
  "updated_at": "2026-11-04T20:00:00Z",
  "transaction_id": "tx_abc123",
  "amount": 2500
}
```

---

### 3. Get Payment Status by Reference

- **`GET /payments/reference/:reference/status`**

**Response (`200 OK`):**

```json
{
  "id": 1,
  "reference": "ref_order_100234",
  "status": "pending"
}
```

---

### 4. Update Payment

- **`PATCH /payments/reference/:reference`**

**Request Body:**

```json
{
  "status": "success"
}
```

**Response (`200 OK`):**

```json
{
  "message": "Payment updated successfully"
}
```

---

### 5. Delete Payment

- **`DELETE /payments/:id`**

**Response (`200 OK`):**

```json
{
  "message": "Payment deleted successfully"
}
```

---

## Automated Testing & Integration

Simflux is designed to run in your integration test pipelines.

### Example with Node.js / Jest:

```javascript
const axios = require("axios");
const { execSync } = require("child_process");

describe("Checkout Webhook Flow", () => {
  beforeAll(() => {
    // Reset state before tests
    execSync("simflux load --reset");
  });

  test("handles failed payment gracefully", async () => {
    // 1. Tell Simflux to simulate a failure
    execSync("simflux webhook set --event payment.failed");

    // 2. Trigger checkout in your application
    const res = await axios.post("http://localhost:3000/api/checkout", {
      orderId: "order_123",
      amount: 5000,
    });

    expect(res.status).toBe(200);

    // 3. Assert your application handled the failed webhook correctly!
  });
});
```

---

## License

This project is licensed under the [MIT License](LICENSE).
