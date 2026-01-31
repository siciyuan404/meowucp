# UCP Order Outbound Webhook Design

Date: 2026-01-31

## Goal

Add outbound UCP order event webhooks for the minimal event set: created, paid, shipped, cancelled. Default to async queue delivery while also supporting manual synchronous delivery.

## Scope

- Trigger outbound events automatically on order status changes.
- Provide admin endpoint to trigger synchronous delivery.
- Reuse existing webhook queue, retry, and delivery URL (`ucp.webhook.delivery_url`).

## Event Mapping

- created: order created successfully
- paid: status changed to paid
- shipped: status changed to shipped
- cancelled: status changed to cancelled

## Components

- OrderService: triggers outbound events after status changes.
- WebhookQueueService: adds enqueue + deliver entrypoints for order events.
- Admin API: manual trigger endpoint.

## Data Flow

Automatic:
1) Order status updated
2) OrderService calls EnqueueOrderEvent(order, eventType)
3) WebhookQueueService writes job
4) Worker delivers to delivery_url with retries

Manual sync:
1) Admin POST /api/v1/admin/orders/:id/webhook
2) Handler calls DeliverOrderEvent(order, eventType)
3) HTTP POST to delivery_url, return result

## Error Handling

- Async: delivery failures do not block status change; retries + alerts handle failures.
- Sync: return error to caller, optional enqueue for retry.
- Optional de-dupe: ignore duplicate event type for same order within a short window.

## Tests

- Status change triggers enqueue call.
- Manual endpoint triggers sync delivery.
- Payload builder includes event_type, order_no, status, timestamp.
- Async failures do not affect order status update.
