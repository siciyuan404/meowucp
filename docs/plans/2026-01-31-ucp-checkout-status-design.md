# UCP Checkout Status and Error Design

Date: 2026-01-31

## Goal

Refine checkout status computation and error messaging for UCP Checkout. The focus is on deterministic `requires_escalation` behavior and UCP-compliant message structure.

## Scope

- Add rules for status transitions based on message severities.
- Define error codes for recoverable vs buyer-input errors.
- Ensure `requires_escalation` always includes `continue_url` and at least one `requires_buyer_input` message.

## Status Rules

1) If any `requires_buyer_input` message exists, set status to `requires_escalation`.
2) Else if any `recoverable` message exists, set status to `incomplete`.
3) Else set status to `ready_for_complete`.

## Triggers for requires_escalation

- No available payment handlers.
- Business requires user input or identity (sign-in, identity linking).

## Message Structure

Use UCP Checkout Message Error format:
- `type`: "error"
- `code`: string
- `severity`: one of `recoverable`, `requires_buyer_input`
- `content`: human-readable message
- `path`: optional JSONPath when tied to a field

## Error Codes (Initial Set)

- `payment_handlers_missing` (requires_buyer_input)
- `requires_sign_in` (requires_buyer_input)
- `requires_identity_linking` (requires_buyer_input)
- `missing_field` (recoverable)
- `invalid_item` (recoverable)
- `invalid_quantity` (recoverable)

## Implementation Notes

- Introduce a helper (e.g., `resolveMessagesAndStatus`) in `internal/ucp/api/checkout_handler.go` used by Create/Update.
- Preserve or regenerate `continue_url` only when status is `requires_escalation`.
- Do not force `continue_url` for `incomplete` or `ready_for_complete`.

## Test Strategy

- Create: no payment handlers -> requires_escalation + message + continue_url.
- Create: missing required fields -> incomplete + recoverable message, no escalation.
- Update: fix recoverable errors -> ready_for_complete.
- Update: requires_sign_in -> requires_escalation.
- Validate that any requires_escalation response includes at least one `requires_buyer_input` message.
