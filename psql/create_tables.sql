CREATE TABLE users (
    id SERIAL PRIMARY KEY,                  
    email VARCHAR(255) NOT NULL UNIQUE,    
    password VARCHAR(255) NOT NULL,         
    isim VARCHAR(100) NOT NULL,           
    soyisim VARCHAR(100) NOT NULL,         
    resim TEXT                   
);