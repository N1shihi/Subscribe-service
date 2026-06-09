CREATE TABLE IF NOT EXISTS subscriptions (
                                             user_id UUID NOT NULL,
                                             service_name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),

    start_date VARCHAR(7) NOT NULL,
    end_date VARCHAR(7),

    PRIMARY KEY (user_id, service_name)
    );