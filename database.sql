CREATE DATABASE InvoicingDB;

USE InvoicingDB;
-- Create the Customers table
CREATE TABLE IF NOT EXISTS Customers (
    CustomerID INT AUTO_INCREMENT,
    FirstName VARCHAR(50),
    LastName VARCHAR(50),
    Email VARCHAR(100),
    Phone VARCHAR(20),
    Address VARCHAR(255),
    BirthDate DATE,
    PRIMARY KEY (CustomerID)
);

-- Create the Invoices table
CREATE TABLE IF NOT EXISTS Invoices (
    InvoiceID INT AUTO_INCREMENT,
    CustomerID INT,
    InvoiceDate DATETIME DEFAULT CURRENT_TIMESTAMP,
    StatusId TINYINT CHECK (StatusId IN (1,  2,  3)), -- Assuming  1=Pending,  2=Paid,  3=Deleted
    UserId INT, -- Auditing user
    PRIMARY KEY (InvoiceID),
    FOREIGN KEY (CustomerID) REFERENCES Customers(CustomerID)
);

-- Create the InvoiceLines table
CREATE TABLE IF NOT EXISTS InvoiceLines (
    LineId INT AUTO_INCREMENT,
    InvoiceId INT,
    Quantity DECIMAL(9,4),
    Title VARCHAR(512),
    Comment VARCHAR(512),
    UnitPrice DECIMAL(10,2),
    PRIMARY KEY (LineId),
    FOREIGN KEY (InvoiceId) REFERENCES Invoices(InvoiceId)
);

-- Insert some initial data into the Customers table
INSERT INTO Customers (FirstName, LastName, Email, Phone, Address, BirthDate)
VALUES ('John', 'Doe', 'john.doe@example.com', '555-1234', '123 Main St', '1980-01-01'),
       ('Jane', 'Smith', 'jane.smith@example.com', '555-5678', '456 Elm St', '1985-02-15');

-- Insert some initial data into the Invoices table
INSERT INTO Invoices (CustomerID, StatusId, UserId)
VALUES (1,   1,   1), -- John Doe's pending invoice
       (2,   2,   2); -- Jane Smith's paid invoice

-- Insert some initial data into the InvoiceLines table
INSERT INTO InvoiceLines (InvoiceId, Quantity, Title, Comment, UnitPrice)
VALUES (1,   1, 'Product A', 'First product',   10.99), -- Line item for John Doe's invoice
       (1,   2, 'Product B', 'Second product',   20.99), -- Line item for John Doe's invoice
       (2,   1, 'Service X', 'Monthly service',   50.00); -- Line item for Jane Smith's invoice
