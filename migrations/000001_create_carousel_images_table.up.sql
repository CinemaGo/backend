CREATE TABLE carousel_images (
    id SERIAL PRIMARY KEY,  -- Unique identifier for each image, automatically incremented
    image_url VARCHAR(255) NOT NULL,     -- URL or path to the image file
    title VARCHAR(255),                 -- Title or caption for the image (optional)
    description TEXT,                   -- Description of the image (optional)
    order_priority INT NOT NULL,        -- The order in which images should appear in the carousel
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- When the image was added
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- When the image info was last updated
);

CREATE INDEX idx_order_priority ON carousel_images(order_priority);       
CREATE INDEX idx_created_at ON carousel_images(created_at);
