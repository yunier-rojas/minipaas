-- migrate:up

-- Queues: `example_queue` for work and `example_stream` for events.
SELECT pgque.create_queue('example_queue');
SELECT pgque.create_queue('example_stream');

-- Tick as soon as a single event is pending so the demo feels responsive.
SELECT pgque.set_queue_config('example_queue', 'ticker_max_count', '1');
SELECT pgque.set_queue_config('example_stream', 'ticker_max_count', '1');
SELECT pgque.set_queue_config('example_queue', 'ticker_max_lag', '1 second');
SELECT pgque.set_queue_config('example_stream', 'ticker_max_lag', '1 second');
SELECT pgque.set_queue_config('example_queue', 'ticker_idle_period', '5 seconds');
SELECT pgque.set_queue_config('example_stream', 'ticker_idle_period', '5 seconds');

-- Register one consumer per queue. Consumers track their own cursor.
SELECT pgque.subscribe('example_queue', 'example-worker');
SELECT pgque.subscribe('example_stream', 'example-consumer');

-- migrate:down

SELECT pgque.drop_queue('example_stream', true);
SELECT pgque.drop_queue('example_queue', true);
