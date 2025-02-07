import os
from time import sleep
from google.cloud import pubsub_v1

PROJECT_ID = "test-pubsub-327810"
os.environ["PUBSUB_EMULATOR_HOST"] = "localhost:8085"

topics = {
    "drivers-downloader-topic": "http://host.docker.internal:8081",
}


def delete_topic(topic_path: str) -> None:
    publisher = pubsub_v1.PublisherClient()
    publisher.delete_topic(request={"topic": topic_path})
    print(f"Topic deleted: {topic_path}")


def create_topic(topic_id: str) -> None:
    publisher = pubsub_v1.PublisherClient()
    topic_path = publisher.topic_path(PROJECT_ID, topic_id)
    topic = publisher.create_topic(request={"name": topic_path})
    print(f"Created topic: {topic.name}")


def create_push_subscription(
    topic_id: str, subscription_id: str, endpoint: str
) -> None:
    publisher = pubsub_v1.PublisherClient()
    subscriber = pubsub_v1.SubscriberClient()

    topic_path = publisher.topic_path(PROJECT_ID, topic_id)
    subscription_path = subscriber.subscription_path(PROJECT_ID, subscription_id)
    push_config = pubsub_v1.types.PushConfig(push_endpoint=endpoint)

    with subscriber:
        subscriber.create_subscription(
            request={
                "name": subscription_path,
                "topic": topic_path,
                "push_config": push_config,
            }
        )

    print(f"Push subscription created: {subscription_id}")


def delete_subscription(subscription_path: str) -> None:
    subscriber = pubsub_v1.SubscriberClient()
    with subscriber:
        subscriber.delete_subscription(request={"subscription": subscription_path})

    print(f"Subscription deleted: {subscription_path}.")


def list_subscriptions_in_project() -> None:
    subscriber = pubsub_v1.SubscriberClient()
    project_path = f"projects/{PROJECT_ID}"

    with subscriber:
        return subscriber.list_subscriptions(request={"project": project_path})


def list_topics() -> None:
    publisher = pubsub_v1.PublisherClient()
    project_path = f"projects/{PROJECT_ID}"
    return publisher.list_topics(request={"project": project_path})


def publish_message(topic_id: str, data_str: str) -> None:
    publisher = pubsub_v1.PublisherClient()

    topic_path = publisher.topic_path(PROJECT_ID, topic_id)
    data = data_str.encode("utf-8")
    future = publisher.publish(topic_path, data)
    future.result()

    print(f"Published messages to {topic_path}.")


def init():
    existing_subscriptions = list_subscriptions_in_project()
    for subscription in existing_subscriptions:
        delete_subscription(subscription.name)

    existing_topics = list_topics()
    for topic in existing_topics:
        delete_topic(topic.name)

    for topic in topics:
        create_topic(topic)
        create_push_subscription(topic, f"{topic}-sub", topics[topic])


import sys
import json


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 example.py <command> [args]")
        sys.exit(1)

    command = sys.argv[1]

    if command == "init":
        init()
    elif command == "send":
        if len(sys.argv) < 4:
            print("Usage: python3 example.py send <topic> <json_data>")
            sys.exit(1)

        topic = sys.argv[2]
        data = sys.argv[3]
        publish_message(topic, data)
    else:
        print("Unknown command.")
        sys.exit(1)


if __name__ == "__main__":
    main()
