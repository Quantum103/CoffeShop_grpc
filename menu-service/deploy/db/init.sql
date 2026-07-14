CREATE TABLE IF NOT EXISTS menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO menu_items (name, description, price, is_available) VALUES
    ('Американо', 'Классический черный кофе', 150.00, true),
    ('Латте', 'Кофе с молоком и нежной пенкой', 200.00, true),
    ('Капучино', 'Идеальный баланс кофе и молока', 180.00, true),
    ('Чизкейк', 'Временно закончился', 250.00, false);