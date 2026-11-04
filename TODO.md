# Simflux Roadmap & Tasks

## Completed
- [x] Remove the `data` column from the `payments` table
- [x] Create the `provider_definitions` table with dynamic validation
- [x] Implement single active provider definition model
- [x] Add `simflux load` CLI with `--reset` support
- [x] Add `simflux webhook set/get/reset` CLI commands
- [x] Webhook configuration persistence and automatic dispatch on payment creation
- [x] Add graceful shutdown to `simflux serve --port`
- [x] Fix Scenario repository CRUD operations
- [x] Expose `PATCH /payments/reference/:reference` and `DELETE /payments/:id` HTTP endpoints

## Future Enhancements
- [ ] Scenario runner: Trigger automated payment flow sequences (e.g. pending -> success / failed with delay)
- [ ] Webhook delivery history log table (recording attempts, status codes, and timestamps)
- [ ] Support dynamic HMAC-SHA256 signatures for webhook security validation
- [ ] Automated integration test suite
