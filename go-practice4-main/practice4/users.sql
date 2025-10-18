CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(100),
    balance NUMERIC(10, 2)
);

INSERT INTO users (name, email, balance) VALUES
('Samat', 'samat@kbtu.com', 1000.00),
('Arsen', 'arsen@kbtu.com', 500.00),
('Zhansaya', 'zhansaya@kbtu.com', 700.00);
