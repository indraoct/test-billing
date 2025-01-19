-- Create Customers Table
CREATE TABLE public.customers (
                           id SERIAL PRIMARY KEY,               -- Auto-incrementing primary key
                           name VARCHAR(100) NOT NULL,          -- Customer's name
                           email VARCHAR(150) UNIQUE NOT NULL,  -- Unique email address
                           phone VARCHAR(15) UNIQUE NOT NULL,   -- Unique phone number
                           created_at TIMESTAMP DEFAULT NOW(),  -- Timestamp for record creation
                           updated_at TIMESTAMP DEFAULT NOW()   -- Timestamp for last update
);

-- Create Loans Table
CREATE TABLE public.loans (
                       id SERIAL PRIMARY KEY,                    -- Auto-incrementing primary key
                       customer_id INT NOT NULL,                 -- Foreign key referencing Customers
                       principal_amount NUMERIC(12, 2) NOT NULL, -- Loan principal amount
                       interest_rate NUMERIC(5, 2) NOT NULL,     -- Interest rate as a percentage
                       weeks INT NOT NULL,                       -- Duration of the loan in weeks
                       weekly_payment NUMERIC(12, 2) NOT NULL,   -- Weekly payment amount
                       outstanding_balance NUMERIC(12, 2) NOT NULL, -- Current outstanding balance
                       status VARCHAR(20) DEFAULT 'Active',      -- Loan status: Active, Paid, Delinquent
                       created_at TIMESTAMP DEFAULT NOW(),       -- Timestamp for record creation
                       updated_at TIMESTAMP DEFAULT NOW(),       -- Timestamp for last update
                       FOREIGN KEY (customer_id) REFERENCES Customers (id) ON DELETE CASCADE
);

-- Create RepaymentSchedules Table
CREATE TABLE public.repaymentschedules (
                                    id SERIAL PRIMARY KEY,                    -- Auto-incrementing primary key
                                    loan_id INT NOT NULL,                     -- Foreign key referencing Loans
                                    week_number INT NOT NULL,                 -- Week number in the repayment schedule
                                    due_date DATE NOT NULL,                   -- Scheduled due date for the payment
                                    status VARCHAR(10) DEFAULT 'Unpaid',      -- Status: Paid, Unpaid
                                    created_at TIMESTAMP DEFAULT NOW(),       -- Timestamp for record creation
                                    updated_at TIMESTAMP DEFAULT NOW(),       -- Timestamp for last update
                                    FOREIGN KEY (loan_id) REFERENCES Loans (id) ON DELETE CASCADE
);

-- Create Payments Table
CREATE TABLE public.payments (
                          id SERIAL PRIMARY KEY,                    -- Auto-incrementing primary key
                          loan_id INT NOT NULL,                     -- Foreign key referencing Loans
                          week_number INT NOT NULL,                 -- Corresponding week number in the repayment schedule
                          amount NUMERIC(12, 2) NOT NULL,           -- Payment amount
                          payment_date TIMESTAMP NOT NULL,          -- Date of payment
                          created_at TIMESTAMP DEFAULT NOW(),       -- Timestamp for record creation
                          updated_at TIMESTAMP DEFAULT NOW(),       -- Timestamp for last update
                          FOREIGN KEY (loan_id) REFERENCES Loans (id) ON DELETE CASCADE
);

-- Insert a customer
INSERT INTO public.customers (name, email, phone)
VALUES ('John Doe', 'john.doe@example.com', '08123456789')
    RETURNING id; -- Returning ID to reference in further inserts

-- Assume the returned customer ID is 1 for this example
-- Insert a loan for the customer
INSERT INTO public.loans (customer_id, principal_amount, interest_rate, weeks, weekly_payment, outstanding_balance, status)
VALUES (1, 5000000, 10, 50, 110000, 5500000, 'Active')
    RETURNING id; -- Returning ID to reference in further inserts

-- Assume the returned loan ID is 1 for this example
-- Insert repayment schedules for the loan
DO $$
DECLARE
week INT;
    due_date DATE := CURRENT_DATE;
BEGIN
FOR week IN 1..50 LOOP
        INSERT INTO public.repaymentschedules (loan_id, week_number, due_date, status)
        VALUES (1, week, due_date, 'Unpaid');
        due_date := due_date + INTERVAL '7 days'; -- Increment due date by 7 days for each week
END LOOP;
END $$;
